package custom

type WorkerPool interface {
	Run()
	AddTask(task func())
}
