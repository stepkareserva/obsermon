package taskpool

type semaphore struct {
	semaCh chan struct{}
}

func newSemaphore(maxReq int) *semaphore {
	return &semaphore{
		semaCh: make(chan struct{}, maxReq),
	}
}

func (s *semaphore) Acquire() {
	s.semaCh <- struct{}{}
}

func (s *semaphore) Release() {
	<-s.semaCh
}
