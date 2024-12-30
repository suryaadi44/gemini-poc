package service

import (
	"gemini-poc/utils/custom"

	"go.uber.org/zap"
)

type workerPool struct {
	maxWorker   int
	queuedTaskC chan func()

	log *zap.Logger
}

func NewWorkerPool(
	maxWorker int,
	maxQueue int,
	log *zap.Logger,
) custom.WorkerPool {
	wp := &workerPool{
		maxWorker:   maxWorker,
		queuedTaskC: make(chan func(), maxQueue),

		log: log,
	}

	return wp
}

func (wp *workerPool) Run() {
	wp.run()
}

func (wp *workerPool) AddTask(task func()) {
	wp.queuedTaskC <- task
}

func (wp *workerPool) GetTotalQueuedTask() int {
	return len(wp.queuedTaskC)
}

func (wp *workerPool) run() {
	for i := 0; i < wp.maxWorker; i++ {
		wID := i + 1
		wp.log.Info("[WorkerPool] Start worker", zap.Int("worker_id", wID))

		go func(workerID int) {
			for task := range wp.queuedTaskC {
				wp.log.Debug("[WorkerPool] Worker start processing task", zap.Int("worker_id", workerID))
				task()
				wp.log.Debug("[WorkerPool] Worker finish processing task", zap.Int("worker_id", workerID))
			}
		}(wID)
	}
}
