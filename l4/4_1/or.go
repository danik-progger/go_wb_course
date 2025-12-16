package or

func Or[T any](chans ...<-chan T) <-chan T {
	switch len(chans) {
	case 0:
		return nil
	case 1:
		return chans[0]
	}

	out := make(chan T)
	go func() {
		defer close(out)

		m := len(chans) / 2
		select {
		case <-Or(chans[:m]...):
		case <-Or(chans[m:]...):
		}
	}()
	return out
}
