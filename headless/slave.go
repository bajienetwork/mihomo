package headless

import "sync"

var regMutex sync.Mutex

func Register() <-chan interface{} {
	regMutex.Lock()
	defer regMutex.Unlock()

	c := make(chan interface{}, 1)
	if On {
		c <- struct{}{}
	} else {
		listeners = append(listeners, c)
	}
	return c
}
