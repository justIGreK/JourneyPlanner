package service

type TaskRepository interface{

}

type TaskSrv struct{
	TaskRepository
}

func NewTaskSrv(taskRepo TaskRepository) *TaskSrv{
	return &TaskSrv{TaskRepository: taskRepo}
}