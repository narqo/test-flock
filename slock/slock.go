package slock

import (
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
)

var (
	ErrLocked   = errors.New("already locked")
	ErrUnlocked = errors.New("already unlocked")
)

type slocker struct {
	mu   sync.Mutex
	file *os.File
	ln   net.Listener
}

func New(path string) (*slocker, error) {
	// in order to place a write lock, fd must be open for writing; see https://linux.die.net/man/2/fcntl
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	if err := file.Chmod(0660); err != nil {
		return nil, err
	}
	return &slocker{file: file}, nil
}

func (l *slocker) Lock() (err error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.ln, err = net.FileListener(l.file)
	//if err != nil {
	//	_, err = net.Dial("unix", l.path)
	//}
	return err
}

func (l *slocker) Unlock() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.ln == nil {
		return ErrUnlocked
	}

	err := l.ln.Close()
	if err == nil {
		l.ln = nil
	}
	return err
}

func (l *slocker) String() string {
	return fmt.Sprintf("<slocker(%s)>", l.file.Name())
}
