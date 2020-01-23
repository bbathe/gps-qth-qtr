package main

// our type
type NonBlockingMutex struct {
	c chan struct{}
}

// NewNonBlockingMutex initializes a new NonBlockingMutex
func NewNonBlockingMutex() *NonBlockingMutex {
	return &NonBlockingMutex{
		c: make(chan struct{}, 1),
	}
}

// Lock attempts to acquire a lock, returns true on success
func (tm NonBlockingMutex) Lock() bool {
	select {
	case tm.c <- struct{}{}:
		return true
	default:
		return false
	}
}

// Unlock releases the lock
func (tm NonBlockingMutex) Unlock() {
	select {
	case <-tm.c:
	default:
		panic("unlock of unlocked mutex")
	}
}
