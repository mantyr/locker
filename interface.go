package locker

type Mutex interface {
	Lock() error
	TryLock() (bool, error)
	Unlock() error
}
