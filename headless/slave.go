package headless

func Register() <-chan interface{} {
	c := make(chan interface{}, 1)
	if on {
		c <- struct{}{}
	} else {
		listeners = append(listeners, c)
	}
	return c
}
