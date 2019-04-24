package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/narqo/test-flock/flock"
	"github.com/narqo/test-flock/slock"
)

type Locker interface {
	Lock() error
	Unlock() error
}

func main() {
	var (
		lockFile1 string
		lockFile2 string
	)

	flag.StringVar(&lockFile1, "lock-file1", "/tmp/flock.lock", "path to flock lockfile")
	flag.StringVar(&lockFile2, "lock-file2", "/tmp/slock.lock", "path to slock lockfile")

	flag.Parse()

	if lockFile1 == "" || lockFile2 == "" {
		log.Fatalf("no lockfile: (1) %q, (2) %q\n", lockFile1, lockFile2)
	}

	var err error

	_, err = runflock(lockFile1)
	if err != nil {
		log.Printf("flocker: %v", err)
	}
	//defer locker1.Unlock()

	//_, err = runslock(lockFile2)
	//if err != nil {
	//	log.Printf("slocker: %v", err)
	//}
	//defer locker2.Unlock()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	<-sigs
}

func runflock(lockFile string) (locker Locker, err error) {
	locker, err = flock.New(lockFile)
	if err != nil {
		return nil, err
	}
	if err := locker.Lock(); err != nil {
		return nil, fmt.Errorf("unable to lock %v: %v\n", locker, err)
	}

	log.Printf("flocker: lock %v\n", locker)

	return locker, nil
}

func runslock(lockFile string) (locker Locker, err error) {
	locker, err = slock.New(lockFile)
	if err != nil {
		return nil, err
	}

	if err := locker.Lock(); err != nil {
		return nil, fmt.Errorf("unable to lock %v: %v\n", locker, err)
	}

	log.Printf("slocker: lock %v\n", locker)

	return locker, nil
}
