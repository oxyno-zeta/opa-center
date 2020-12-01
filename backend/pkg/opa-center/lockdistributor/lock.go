package lockdistributor

import (
	"cirello.io/pglock"
)

type lock struct {
	name string
	pl   *pglock.Lock
	s    *service
}

func (l *lock) IsAlreadyTaken() (bool, error) {
	// Get lock
	lo, err := l.s.cl.Get(l.name)
	// Check error
	if err != nil {
		// Check if error is a not found error
		if err == pglock.ErrLockNotFound {
			return false, nil
		}

		return false, err
	}

	// Check if lock exists or not
	return lo != nil, nil
}

func (l *lock) Acquire() error {
	ll, err := l.s.cl.Acquire(l.name)
	// Check error
	if err != nil {
		return err
	}
	// Save lock
	l.pl = ll

	return nil
}

func (l *lock) IsReleased() (bool, error) {
	return l.pl.IsReleased(), nil
}

func (l *lock) Release() error {
	return l.pl.Close()
}
