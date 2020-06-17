// +build !windows,!solaris

package term

import (
	"fmt"
	"os"
	"syscall"

	"github.com/dalefarnsworth/term/termios"
)

func setOptions(t Term, options ...func(*Term) error) error {
	if err := termios.Tcgetattr(uintptr(t.fd), &t.orig); err != nil {
		return err
	}
	if err := t.SetOption(options...); err != nil {
		return err
	}
	return syscall.SetNonblock(t.fd, false)
}

// Open opens an asynchronous communications port.
func Open(name string, options ...func(*Term) error) (*Term, error) {
	fd, e := syscall.Open(name, syscall.O_NOCTTY|syscall.O_CLOEXEC|syscall.O_NDELAY|syscall.O_RDWR, 0666)
	if e != nil {
		return nil, &os.PathError{"open", name, e}
	}

	t := Term{name: name, fd: fd}

	return &t, setOptions(t, options...)
}

// OpenFD opens an asynchronous communications port.
func OpenFD(fd int, options ...func(*Term) error) (*Term, error) {
	t := Term{name: fmt.Sprintf("fd%d"), fd: fd}

	return &t, setOptions(t, options...)
}

// Restore restores the state of the terminal captured at the point that
// the terminal was originally opened.
func (t *Term) Restore() error {
	return termios.Tcsetattr(uintptr(t.fd), termios.TCIOFLUSH, &t.orig)
}
