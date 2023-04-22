package util

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/gogf/gf/util/grand"

	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gtime"

	"github.com/360EntSecGroup-Skylar/excelize"
)

//func DbToExcel(rows *sql.Rows) (headers []string, records []map[string]interface{}) {
func DbToExcel(rows *sql.Rows) (export_path string, err error) {

	column_head_list, _ := rows.Columns()
	//column_type_list, _ := rows.ColumnTypes()

	row_tmpl := make([]interface{}, len(column_head_list))
	//构建每行空结构 ，就是有个名字key做对应
	for hi, _ := range column_head_list {
		var np interface{}
		row_tmpl[hi] = &np
	}

	row_data_list := make([]map[string]interface{}, 0)

	//然后开始遍历
	for rows.Next() {
		rows.Scan(row_tmpl...)
		//然后整理每一个数据对应列名
		tem_row := map[string]interface{}{}
		for tem_cell_index, tem_cell_val := range row_tmpl {
			//	转类型
			//每一列的数据转成对应类型 在封装进新的struct
			tem_row[column_head_list[tem_cell_index]] = *tem_cell_val.(*interface{})
		}
		row_data_list = append(row_data_list, tem_row)
	}

	//headers := column_head_list
	//records := row_data_list

	//开始处理 excel 文件

	f := excelize.NewFile()
	const sheet_1_name = "Sheet1"

	sheet_index := f.NewSheet(sheet_1_name)
	f.SetActiveSheet(sheet_index)

	//先构建第一行输出列表头
	for tem_head_index, tem_head_label := range column_head_list {
		y := 1
		x := convertToTitle(tem_head_index + 1)
		h := fmt.Sprintf("%s%d", x, y)
		f.SetCellValue(sheet_1_name, h, tem_head_label)
	}

	//遍历每行每个单元格
	for data_row_index, data_row_dct := range row_data_list {
		y := data_row_index + 2
		//按结果集的顺序列开始读取
		for tem_head_index, tem_head_val := range column_head_list {
			//c_x:=
			x := convertToTitle(tem_head_index + 1)
			l := fmt.Sprintf("%s%d", x, y)
			cell_val := data_row_dct[tem_head_val]
			f.SetCellValue(sheet_1_name, l, cell_val)
		}
	}

	f.SetActiveSheet(sheet_index)

	base_path := "export"
	excel_path := "/excel/" + gtime.Now().Format("Ymd") + "/" //服务器访问路径
	if !gfile.Exists(base_path + excel_path) {
		gfile.Mkdir(base_path + excel_path)
	}

	excel_name := strings.ToLower(strconv.FormatInt(gtime.TimestampNano(), 36)) + grand.S(6) + ".xlsx"
	excel_name = base_path + excel_path + excel_name

	if err := f.SaveAs(excel_name); err != nil {
		fmt.Println(err)
		return "", err
	}
	export_path = excel_name

	return
}

func ExportToExcel() {
	//然后开始处理excel
	//f := excelize.NewFile()
	//
	//const sheet_1_name = "Sheet1"
	//
	//sheet_index := f.NewSheet(sheet_1_name)
	//f.SetActiveSheet(sheet_index)
	//
	////先构建第一行输出列表头
	//for tem_head_index, tem_head_label := range column_head_list {
	//	y := 1
	//	x := convertToTitle(tem_head_index + 1)
	//	h := fmt.Sprintf("%s%d", x, y)
	//	f.SetCellValue(sheet_1_name, h, tem_head_label)
	//}
	//
	////遍历每行每个单元格
	//for data_row_index, data_row_dct := range row_data_list {
	//	y := data_row_index + 2
	//	//按结果集的顺序列开始读取
	//	for tem_head_index, tem_head_val := range column_head_list {
	//		//c_x:=
	//		x := convertToTitle(tem_head_index + 1)
	//		l := fmt.Sprintf("%s%d", x, y)
	//		cell_val := data_row_dct[tem_head_val]
	//		f.SetCellValue(sheet_1_name, l, cell_val)
	//	}
	//}
	//
	//f.SetActiveSheet(sheet_index)
	//
	////// Create a new sheet.
	////index := f.NewSheet("Sheet2")
	////// Set value of a cell.
	////f.SetCellValue("Sheet2", "A2", "Hello world.")
	////f.SetCellValue("Sheet1", "B2", 100)
	////// Set active sheet of the workbook.
	////f.SetActiveSheet(index)
	//// Save xlsx file by the given path.
	//
	//base_path := "export"
	//excel_path := "/bill_excel/" + gtime.Now().Format("Ymd") + "/" //服务器访问的路径
	//if !gfile.Exists(base_path + excel_path) {
	//	gfile.Mkdir(base_path + excel_path)
	//}
	//
	//excel_name := strings.ToLower(strconv.FormatInt(gtime.TimestampNano(), 36)) + grand.S(6) + ".xlsx"
	//excel_name = base_path + excel_path + excel_name
	//
	//if err := f.SaveAs(excel_name); err != nil {
	//	fmt.Println(err)
	//	return "", "", err
	//}
	//
	//file_path = excel_name
	//bill_date = bill.BillDateTip
	//return
}

func convertToTitle(n int) string {
	// 26进制,一共26个字符，不是27进制
	// 字母与数字运算转字符串
	str := ""
	for n != 0 { //从末尾向头转换
		str = string(byte('A')+byte((n-1)%26)) + str // i=26
		n = (n - 1) / 26                             // 注意是对n-1进行的运算
	}
	return str
}
