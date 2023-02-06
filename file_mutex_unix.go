package locker

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/sys/unix"
)

type fileMutex struct {
	mutex      sync.Mutex
	descriptor int
}

const (
	fileMutexFileExt              = ".lock"
	fileMutexFileMode os.FileMode = 0666
	fileMutexDirMode  os.FileMode = 0755
)

func NewFileMutex(address string) (Mutex, error) {
	if filepath.Ext(address) != fileMutexFileExt {
		return nil, fmt.Errorf(
			`expected lock file path should have "%s" extension`,
			fileMutexFileExt,
		)
	}
	err := os.MkdirAll(filepath.Dir(address), fileMutexDirMode)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(address, os.O_CREATE|os.O_RDONLY, fileMutexFileMode)
	if err != nil {
		return nil, err
	}
	// defer file.Close() // if you close the file descriptor, then the mutex will not work
	m := fileMutex{
		descriptor: int(file.Fd()),
	}
	return &m, nil
}

func (m *fileMutex) Lock() error {
	m.mutex.Lock()
	err := unix.Flock(m.descriptor, unix.LOCK_EX)
	if err != nil {
		m.mutex.Unlock()
	}
	return err
}

func (m *fileMutex) TryLock() (bool, error) {
	ok := m.mutex.TryLock()
	if !ok {
		return false, nil
	}
	err := unix.Flock(m.descriptor, unix.LOCK_EX|unix.LOCK_NB)
	if err != nil {
		m.mutex.Unlock()
		if err == unix.EWOULDBLOCK {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (m *fileMutex) Unlock() error {
	m.mutex.Unlock()
	return unix.Flock(m.descriptor, unix.LOCK_UN)
}
