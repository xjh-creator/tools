package main

type TaskUploadFileProcess int8

//队列信号常量
const TASK_UPLOAD_FILE_PROCESS_QUERY_NAME = "TASK_UPLOAD_FILE_PROCESS_QUERY_NAME"

const (
	Import = 20
)

type MsgProcessImportFileEntry struct {
	ImportFileType   TaskUploadFileProcess
	BehaviorImportID int64 `json:"behavior_import_id,string"`
}

func processImportFileMsg(msg MsgProcessImportFileEntry) {
	q := Get(TASK_UPLOAD_FILE_PROCESS_QUERY_NAME).(*gqueue.Queue)
	q.Push(msg)
}

// ProcessImportFileBehavior
func ProcessImportFileBehavior(behavior_import_id int64) {
	msg := MsgProcessImportFileEntry{
		ImportFileType:   Import,
		BehaviorImportID: behavior_import_id,
	}
	processImportFileMsg(msg)
}

//全局自定义 kv
var global_map = gmap.New(true)

func Set(name string, val interface{}) {
	global_map.Set(name, val)
}

func Get(name string) interface{} {
	return global_map.Get(name)
}
