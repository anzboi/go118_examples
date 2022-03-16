package chans

import (
	"reflect"
)

func Drain[T any](ch <-chan T) {
	for range ch {
	}
}

// Sets up a goroutine that reads all messages from inputs and sends them down a single channel
//
// Returns the new output channel, Out channel will close when all input channels are closed.
func Merge[T any](inputs ...<-chan T) <-chan T {
	r := make(chan T)
	// parent goroutine listens for messages and sends them down r
	go func(r chan<- T, inputs ...<-chan T) {
		defer close(r)

		// We cannot use a normal select statement because we don't have a fixed number of channels
		cases := make([]reflect.SelectCase, len(inputs))
		for i, ch := range inputs {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
		}

		// select loop
		for {
			if len(cases) == 0 {
				break
			}
			chosen, value, ok := reflect.Select(cases)
			if !ok {
				// remove the closed channel from cases
				cases = append(cases[:chosen], cases[chosen+1:]...)
			} else {
				r <- value.Interface().(T)
			}
		}

	}(r, inputs...)

	return r
}

// Sets up a goroutine to send input messages down all outputs.
func Split[T any](in <-chan T, out ...chan<- T) {
	go func(in <-chan T, out ...chan<- T) {
		for _, ch := range out {
			defer close(ch)
		}

		notifyAll := func(item T) {
			for _, ch := range out {
				ch := ch
				// Sets up a single goroutine per send so we don't block
				// all channels because 1 has an inactive reciever
				go func() {
					ch <- item
				}()
			}
		}

		for item := range in {
			notifyAll(item)
		}
	}(in, out...)
}
