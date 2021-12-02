package single

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

// Single represents the name and the open file descriptor
type Single struct {
	name string
	file *os.File
}

// New creates a Single instance
func New(name string) *Single {
	return &Single{name: name}
}

// filenameLock returns an absolute filename, appropriate for the operating system
func (s *Single) filenameLock() string {
	return filepath.Join("/var/lock", fmt.Sprintf("%s.lock", s.name))
}

// filenamePID returns an absolute filename, appropriate for the operating system
func (s *Single) filenamePID() string {
	return filepath.Join("/var/run", fmt.Sprintf("%s.pid", s.name))
}

// Lock tries to obtain an excluded lock on a lockfile and returns an error if one occurs
func (s *Single) Lock() (busy bool, err error) {
	// open/create lock file
	s.file, err = os.OpenFile(s.filenameLock(), os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return
	}

	// set the lock type to F_WRLCK, therefor the file has to be opened writable
	flock := syscall.Flock_t{
		Type: syscall.F_WRLCK,
		Pid:  int32(os.Getpid()),
	}

	// error if not root
	_ = ioutil.WriteFile(s.filenamePID(), []byte(strconv.Itoa(int(flock.Pid))), 0644)

	// try to obtain an exclusive lock - FcntlFlock seems to be the portable *ix way
	busy = syscall.FcntlFlock(s.file.Fd(), syscall.F_SETLK, &flock) != nil
	return
}

// Unlock unlocks, closes and removes the lockfile
func (s *Single) Unlock() (err error) {
	// set the lock type to F_UNLCK
	flock := syscall.Flock_t{
		Type: syscall.F_UNLCK,
		Pid:  int32(os.Getpid()),
	}

	if err = syscall.FcntlFlock(s.file.Fd(), syscall.F_SETLK, &flock); err != nil {
		return fmt.Errorf("failed to unlock the lock file: %v", err)
	}

	if err = s.file.Close(); err != nil {
		return fmt.Errorf("failed to close the lock file: %v", err)
	}

	if err = os.Remove(s.filenameLock()); err != nil {
		return fmt.Errorf("failed to remove the lock file: %v", err)
	}

	return
}
