package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/mantyr/locker"
)

func main() {
	if len(os.Args) != 3 {
		panic(errors.New("expected go run main.go [command] [file.lock]"))
	}
	command := os.Args[1]
	lockFile := os.Args[2]

	mutex, err := locker.NewFileMutex(lockFile)
	if err != nil {
		panic(err)
	}
	fmt.Printf("start command: %s\r\n", command)
	var ok bool
	switch command {
	case "lock":
		err = mutex.Lock()
		if err != nil {
			panic(err)
		}
		fmt.Println("STATUS: Locked")
		time.Sleep(60 * time.Second)
	case "try-lock":
		ok, err = mutex.TryLock()
		if err != nil {
			panic(err)
		}
		if ok {
			fmt.Println("STATUS: Locked")
		} else {
			fmt.Println("STATUS: Non Locked")
		}
		time.Sleep(60 * time.Second)
	case "multi-lock":
		err = mutex.Lock()
		if err != nil {
			panic(err)
		}
		fmt.Println("STATUS: Locked")

		ok, err = mutex.TryLock()
		if err != nil {
			panic(err)
		}
		if ok {
			fmt.Println("STATUS1: Locked (dubl)")
		} else {
			fmt.Println("STATUS1: TryLock return false")
		}

		mutex2, err := locker.NewFileMutex(lockFile)
		if err != nil {
			panic(err)
		}
		ok, err = mutex2.TryLock()
		if err != nil {
			panic(err)
		}
		if ok {
			fmt.Println("STATUS2: Locked (dubl)")
		} else {
			fmt.Println("STATUS2: TryLock return false")
		}
		time.Sleep(60 * time.Second)
	case "unlock":
		err = mutex.Unlock()
		if err != nil {
			panic(err)
		}
		fmt.Println("STATUS: Unlocked")
	default:
		panic(errors.New("unexpected command"))
	}
}
