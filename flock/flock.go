package flock

import (
	"errors"
	"fmt"
	"io"
	"os"
	"syscall"
)

var (
	ErrLocked   = errors.New("already locked")
	ErrUnlocked = errors.New("already unlocked")
)

type flocker struct {
	file *os.File
}

func New(path string) (*flocker, error) {
	// in order to place a write lock, fd must be open for writing; see https://linux.die.net/man/2/fcntl
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	if err := file.Chmod(0660); err != nil {
		return nil, err
	}

	return &flocker{file: file}, nil
}

func (l *flocker) Lock() error {
	err := l.lock()
	if err != nil {
		l.unlock()
	}
	return err
}

func (l *flocker) Unlock() error {
	l.unlock()
	return l.file.Close()
}

func (l *flocker) lock() error {
	return syscallFcntlFlock(l.file, syscall.F_WRLCK)
}

func (l *flocker) unlock() error {
	return syscallFcntlFlock(l.file, syscall.F_UNLCK)
}

func (l *flocker) String() string {
	return fmt.Sprintf("<flocker(%s)>", l.file.Name())
}

// Note, POSIX locks apply per inode and process.
// The lock for an inode is released when *any* descriptor for that inode is closed.
// - http://0pointer.de/blog/projects/locking.html
// - https://github.com/golang/go/blob/release-branch.go1.12/src/cmd/go/internal/lockedfile/internal/filelock/filelock_fcntl.go
// For the use case described in adjust/backend#8158 this is fine.
func syscallFcntlFlock(file *os.File, lt int16) error {
	err := syscall.FcntlFlock(
		file.Fd(),
		//syscall.F_SETLK, // non-blocking, e.g. return EAGAIN error if lock is already held
		syscall.F_SETLKW, // blocking
		&syscall.Flock_t{
			Type:   lt,
			Whence: io.SeekStart,
			Start:  0,
			Len:    0,
		},
	)
	if err == syscall.EAGAIN {
		return ErrLocked
	}
	return err
}
