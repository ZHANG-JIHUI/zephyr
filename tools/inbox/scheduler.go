package inbox

const (
	defaultThroughput = 300
)

type Scheduler interface {
	Schedule(fn func())
	Throughput() int
}

type scheduler int

func (slf scheduler) Schedule(fn func()) {
	go fn()
}

func (slf scheduler) Throughput() int {
	return int(slf)
}

func NewScheduler(throughput int) Scheduler {
	return scheduler(throughput)
}
