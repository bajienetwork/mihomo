package headless

var listeners []chan interface{}
var On = false

func TurnOn() {
	On = true
	for _, listener := range listeners {
		listener <- struct{}{}
	}
}
