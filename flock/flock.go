package flock

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"syscall"
)

var (
	ErrLocked   = errors.New("already locked")
	ErrUnlocked = errors.New("already unlocked")
)

type flocker struct {
	mu     sync.Mutex
	file   *os.File
	flockT *syscall.Flock_t
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
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.flock()
}

func (l *flocker) Unlock() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.funlock()
}

func (l *flocker) flock() error {
	if l.flockT != nil {
		return ErrLocked
	}

	l.flockT = &syscall.Flock_t{
		Type:   syscall.F_WRLCK, // F_WRLCK is for write lock on a file
		Whence: io.SeekStart,
		Start:  0,
		Len:    0,
	}
	return syscallFcntlFlock(l)
}

func (l *flocker) funlock() error {
	if l.flockT == nil {
		return ErrUnlocked
	}

	// F_UNLCK is for release lock
	l.flockT.Type = syscall.F_UNLCK

	err := syscallFcntlFlock(l)
	if err == nil {
		l.flockT = nil
	}
	return err
}

func (l *flocker) String() string {
	return fmt.Sprintf("<flocker(%s)>", l.file.Name())
}

func syscallFcntlFlock(l *flocker) error {
	return syscall.FcntlFlock(l.file.Fd(), syscall.F_SETLK, l.flockT)
}
