# Locker

[![Build Status](https://travis-ci.org/mantyr/locker.svg?branch=master)](https://travis-ci.org/mantyr/locker)
[![GoDoc](https://godoc.org/github.com/mantyr/locker?status.png)](http://godoc.org/github.com/mantyr/locker)
[![Go Report Card](https://goreportcard.com/badge/github.com/mantyr/locker?v=1)][goreport]
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg)](LICENSE.md)

This is the stable version.

## Description

- `FileMutex` - To lock via operating system file
    - Supports functions:
        - [x] Lock
        - [x] TryLock
        - [x] Unlock
    - Important:
        - [ ] thread safe - You can take several locks from one mutex within a single running application.

### Supports platforms

- [x] unix (linux, mac)
- [ ] windows
- [ ] other

## Installation

    $ go get github.com/mantyr/locker

## Example

> **Note**
>
> You will probably be protecting the entire application, rather than a single object.

```go

type Storage struct {
	mutex locker.Mutex
}

func New(fileLock string) (*Storage, error) {
	mutex, err := locker.NewFileMutex(fileLock)
	if err != nil {
		return nil, err
	}
	return &Storage{
		mutex: mutex,
	}, nil
}

func (s *Storage) Action() (err error) {
	defer s.mutex.Unlock()
	err = s.mutex.Lock()
	if err != nil {
		return err
	}
	fmt.Println("action")
	return nil
}

```

```bash
go run ./example/main.go lock ./testdata/file.lock
go run ./example/main.go try-lock ./testdata/file.lock
go run ./example/main.go multi-lock ./testdata/file.lock
go run ./example/main.go unlock ./testdata/file.lock
```

## Author

[Oleg Shevelev][mantyr]

[mantyr]: https://github.com/mantyr

[build_status]: https://travis-ci.org/mantyr/locker
[godoc]:        http://godoc.org/github.com/mantyr/locker
[goreport]:     https://goreportcard.com/report/github.com/mantyr/locker
