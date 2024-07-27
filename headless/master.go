package headless

var listeners []chan interface{}
var on = false

func TurnOn() {
	on = true
	for _, listener := range listeners {
		listener <- struct{}{}
	}
}
