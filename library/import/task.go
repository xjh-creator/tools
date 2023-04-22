package main

import (
	"context"
	"dongguanquandao_server/app/inner/dispatch"
	"dongguanquandao_server/app/inner/enum"
	"dongguanquandao_server/app/inner/service"
	"fmt"
	"time"
	"tools/library/xdb"

	"github.com/gogf/gf/os/glog"

	"github.com/gogf/gf/container/gqueue"
)

func ProcessFile(queue *gqueue.Queue) {
	for {
		data := queue.Pop().(MsgProcessImportFileEntry)
		glog.Info("queue pop data:", data)

		switch data.ImportFileType {
		case TaskUploadFileProcess_None:
			glog.Info("信号异常,无法处理:", data)
		case TaskUploadFileProcess_BehaviorImport:
			ProcessFileByBehaviorImport(data.BehaviorImportID)
		default:
			glog.Info("未知信号,无法处理:", data)
		}
		time.Sleep(time.Millisecond)
	}

}

func ProcessFileByBehaviorImport(behavior_import_id int64) {
	glog.Info("开始处理Behavior导入 behavior_import_id:", behavior_import_id)

	err := ImportFromExcelByID(context.Background(), behavior_import_id)
	if err != nil {
		glog.Error(err)
	}
}

func ImportFromExcelByID(ctx context.Context, id int64) (err error) {
	record, err := s.Detail(ctx, id)
	if err != nil {
		return
	}

	record.Status = enum.BehaviorImportStatus_Processing
	if err = s.dao.Update(ctx, record); err != nil {
		return
	}

	succ_count, fail_count, total_count, err := ImportFromExcel(ctx, record.FilePath)

	if err != nil {
		//失败处理
		record.Status = enum.BehaviorImportStatus_End
		record.Result = enum.BehaviorImportResult_Fail
		record.ErrMsg = err.Error()
		record.SuccessCount = succ_count
		record.FailCount = fail_count
		record.TotalCount = total_count
		err = s.dao.Update(ctx, record)

	} else {
		//成功处理
		record.Status = enum.BehaviorImportStatus_End
		record.Result = enum.BehaviorImportResult_Success
		record.SuccessCount = succ_count
		record.FailCount = fail_count
		record.TotalCount = total_count
		err = s.dao.Update(ctx, record)
	}

	return
}

func ImportFromExcel(ctx context.Context, file string) (suc_count, fail_count, total_count int, err error) {
	glog.Info("开始处理excel 文件:", file)

	xlsxFile, err := xlsx.OpenFile(file)

	if err != nil {
		//return err
		return
	}
	sheet := xlsxFile.Sheets[0]
	glog.Info("总行数:", len(sheet.Rows))

	//先全部处理封装 在一个个判断
	excel_row_list := make([]*behaviorImportExcelRow, 0)
	//就获取所有部门
	//dept_list, err := dao.NewDepartmentDao().All(ctx)
	//if err != nil {
	//	return
	//}

	for rdx, row := range sheet.Rows {
		if rdx == 0 {
			continue
		}

		row_struct := &behaviorImportExcelRow{}
		if err = row.ReadStruct(row_struct); err != nil {
			glog.Error(err)
			return 0, 0, 0, gerror.New(fmt.Sprintf("第%d行:读取异常", rdx+1))
		}

		excel_row_list = append(excel_row_list, row_struct) //加入数组
	}

	excel_row_list_count := len(excel_row_list)

	glog.Info("有效格式excel行数:", excel_row_list_count)

	behavior_list := make([]*model.Behavior, 0)

	mission_service := NewMissionService()

	mission_unique_id_list := make([]string, 0)
	mission_unique_id_dct := map[string]interface{}{}
	for _, row := range excel_row_list {
		s_uid := gmd5.MustEncryptString(fmt.Sprintf("%s%s", row.PERSON_NAME, row.TELEPHONE))
		if mission_unique_id_dct[s_uid] == nil {
			mission_unique_id_dct[s_uid] = g.Map{"name": row.PERSON_NAME, "telephone": row.TELEPHONE}
			mission_unique_id_list = append(mission_unique_id_list, s_uid)
		}
	}

	for idx, row := range excel_row_list {

		//先定位到对应的mission
		mission_unique_id := gmd5.MustEncryptString(fmt.Sprintf("%s%s", row.PERSON_NAME, row.TELEPHONE))

		tm, f_err := mission_service.FindByMissionUniqueID(ctx, mission_unique_id)
		if f_err != nil {
			if f_err == gorm.ErrRecordNotFound {
				return 0, 0, 0, gerror.New(fmt.Sprintf("第%d行:%s", idx+2, "找不到对应任务"))
			}
			return 0, 0, 0, gerror.New(fmt.Sprintf("第%d行:%s", idx+2, f_err.Error()))
		}

		behavior_at := gtime.NewFromStr(row.BEHAVIOR_AT)
		if behavior_at == nil {
			return 0, 0, 0, gerror.New(fmt.Sprintf("第%d行:日期格式不正确", idx+2))
		}

		b := &model.Behavior{
			Code:  row.BEHAVIOR_CODE,
			NAME:  "",
			SNAME: row.BEHAVIOR_TITLE,
			//WFGD:       "",
			//FLNR:       "",
			//SORTNUM:    0,
			BehaviorAt: behavior_at.TimestampMilli(),
			MissionID:  tm.ID,
		}

		behavior_list = append(behavior_list, b)
	}

	var funcs []func(db *gorm.DB) error
	funcs = append(funcs, func(db *gorm.DB) error {
		return dao.NewBehaviorDao(&xdb.XDb{DB: db}).Creates(ctx, behavior_list)
	})
	if err = xdb.Tx(funcs...); err != nil {
		glog.Error(err)
		return 0, 0, len(behavior_list), err
	}
	return len(behavior_list), 0, len(behavior_list), err
}
