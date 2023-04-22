package main

// UploadQueue 项目启动初始化
func UploadQueue() {
	upload_task_queue := gqueue.New()

	Set(TASK_UPLOAD_FILE_PROCESS_QUERY_NAME, upload_task_queue)

	go ProcessFile(upload_task_queue)

}
