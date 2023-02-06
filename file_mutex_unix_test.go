package locker

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey" //nolint:golint,stylecheck
)

func TestFileMutexUnix(t *testing.T) {
	Convey("Checking locker.NewFileMutex (unix)", t, func() {
		fileLock := "./testdata/file.lock"
		os.Remove(fileLock)

		Convey("Checking the creation of the .lock file", func() {
			stat, err := os.Stat(fileLock)
			So(os.IsNotExist(err), ShouldBeTrue)
			So(stat, ShouldBeNil)

			mutex, err := NewFileMutex(fileLock)
			So(err, ShouldBeNil)
			So(mutex, ShouldNotBeNil)

			stat, err = os.Stat(fileLock)
			So(err, ShouldBeNil)
			So(stat, ShouldNotBeNil)
			Convey("Checking for mutex re-creation", func() {
				mutex2, err := NewFileMutex(fileLock)
				So(err, ShouldBeNil)
				So(mutex2, ShouldNotBeNil)
			})
			Convey("Checking TryLock", func() {
				So(mutex.Lock(), ShouldBeNil)

				lock, err := mutex.TryLock()
				So(err, ShouldBeNil)
				So(lock, ShouldBeFalse)

				So(mutex.Unlock(), ShouldBeNil)
				lock, err = mutex.TryLock()
				So(err, ShouldBeNil)
				So(lock, ShouldBeTrue)

				mutex2, err := NewFileMutex(fileLock)
				So(err, ShouldBeNil)
				So(mutex2, ShouldNotBeNil)

				lock2, err := mutex2.TryLock()
				So(err, ShouldBeNil)
				So(lock2, ShouldBeFalse)
			})
			Convey("Checking Unlock - After Lock", func() {
				So(mutex.Lock(), ShouldBeNil)
				So(mutex.Unlock(), ShouldBeNil)

				lock, err := mutex.TryLock()
				So(err, ShouldBeNil)
				So(lock, ShouldBeTrue)

				mutex2, err := NewFileMutex(fileLock)
				So(err, ShouldBeNil)
				So(mutex2, ShouldNotBeNil)

				lock2, err := mutex2.TryLock()
				So(err, ShouldBeNil)
				So(lock2, ShouldBeFalse)
			})
			Convey("Checking Unlock - After TryLock", func() {
				So(mutex.Lock(), ShouldBeNil)
				So(mutex.Unlock(), ShouldBeNil)

				lock, err := mutex.TryLock()
				So(err, ShouldBeNil)
				So(lock, ShouldBeTrue)
				So(mutex.Unlock(), ShouldBeNil)

				mutex2, err := NewFileMutex(fileLock)
				So(err, ShouldBeNil)
				So(mutex2, ShouldNotBeNil)

				lock2, err := mutex2.TryLock()
				So(err, ShouldBeNil)
				So(lock2, ShouldBeTrue)
			})
			Convey("Checking Lock and Unlock", func() {
				So(mutex.Lock(), ShouldBeNil)
				So(mutex.Unlock(), ShouldBeNil)
			})
			Convey("Checking Unlock", func() {
				So(mutex.Lock(), ShouldBeNil)
				So(mutex.Unlock(), ShouldBeNil)

				lock, err := mutex.TryLock()
				So(err, ShouldBeNil)
				So(lock, ShouldBeTrue)
				So(mutex.Unlock(), ShouldBeNil)
			})
			Convey("Checking Unlock - We check that it is impossible to unlock someone else's lock", func() {
				lock, err := mutex.TryLock()
				So(err, ShouldBeNil)
				So(lock, ShouldBeTrue)

				mutex2, err := NewFileMutex(fileLock)
				So(err, ShouldBeNil)
				So(mutex2, ShouldNotBeNil)

				lock2, err := mutex2.TryLock()
				So(err, ShouldBeNil)
				So(lock2, ShouldBeFalse)

				So(mutex.Unlock(), ShouldBeNil)
				lock2, err = mutex2.TryLock()
				So(err, ShouldBeNil)
				So(lock2, ShouldBeTrue)
			})
		})
	})
}
