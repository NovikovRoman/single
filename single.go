package single

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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

// Lock tries to obtain an exclude lock on a lockfile and exits the program if an error occurs
func (s *Single) Lock() {
	if _, err := s.CheckLock(); err != nil {
		log.Fatal(err)
	}
}

// Unlock releases the lock, closes and removes the lockfile. All errors will be reported directly.
func (s *Single) Unlock() {
	if err := s.TryUnlock(); err != nil {
		log.Print(err)
	}
}

// Filename returns an absolute filename, appropriate for the operating system
func (s *Single) Filename() string {
	return filepath.Join("/var/lock", fmt.Sprintf("%s.lock", s.name))
}

// CheckLock tries to obtain an exclude lock on a lockfile and returns an error if one occurs
func (s *Single) CheckLock() (bool, error) {
	// open/create lock file
	f, err := os.OpenFile(s.Filename(), os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return false, err
	}

	s.file = f
	// set the lock type to F_WRLCK, therefor the file has to be opened writable
	flock := syscall.Flock_t{
		Type: syscall.F_WRLCK,
		Pid:  int32(os.Getpid()),
	}

	// try to obtain an exclusive lock - FcntlFlock seems to be the portable *ix way
	return syscall.FcntlFlock(s.file.Fd(), syscall.F_SETLK, &flock) != nil, nil
}

// TryUnlock unlocks, closes and removes the lockfile
func (s *Single) TryUnlock() error {
	// set the lock type to F_UNLCK
	flock := syscall.Flock_t{
		Type: syscall.F_UNLCK,
		Pid:  int32(os.Getpid()),
	}
	if err := syscall.FcntlFlock(s.file.Fd(), syscall.F_SETLK, &flock); err != nil {
		return fmt.Errorf("failed to unlock the lock file: %v", err)
	}
	if err := s.file.Close(); err != nil {
		return fmt.Errorf("failed to close the lock file: %v", err)
	}
	if err := os.Remove(s.Filename()); err != nil {
		return fmt.Errorf("failed to remove the lock file: %v", err)
	}
	return nil
}
