package locker

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

type fileMutex int

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
	m := fileMutex(int(file.Fd()))
	return &m, nil
}

func (m *fileMutex) Lock() error {
	return unix.Flock(int(*m), unix.LOCK_EX)
}

func (m *fileMutex) TryLock() (bool, error) {
	err := unix.Flock(int(*m), unix.LOCK_EX|unix.LOCK_NB)
	if err != nil {
		if err == unix.EWOULDBLOCK {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (m *fileMutex) Unlock() error {
	return unix.Flock(int(*m), unix.LOCK_UN)
}
