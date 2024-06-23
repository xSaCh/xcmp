package util

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

type TerminalSize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func GetTermSize() TerminalSize {
	ws := &TerminalSize{}
	retCode, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		panic(errno)
	}
	return *ws
}

func GetCousorPos() (uint32, uint32, error) {
	err := makeTerminalRaw()
	if err != nil {
		return 0, 0, err
	}
	// Write the query cursor position escape sequence
	print("\033[6n")

	// Read the response from the terminal
	var buf [16]byte
	n, err := os.Stdin.Read(buf[:])
	if err != nil {
		return 0, 0, err
	}

	// The response is expected to be in the format \033[{ROW};{COLUMN}R
	response := string(buf[:n])
	if !strings.HasPrefix(response, "\033[") || !strings.HasSuffix(response, "R") {
		return 0, 0, fmt.Errorf("unexpected response: %s", response)
	}

	// Strip the prefix \033[ and suffix R
	response = response[2 : len(response)-1]

	// Split the response into row and column
	parts := strings.Split(response, ";")

	row, err1 := strconv.Atoi(parts[0])
	column, err2 := strconv.Atoi(parts[1])

	if err1 != nil || err2 != nil {
		return 0, 0, err1
	}

	return uint32(column), uint32(row), nil
}

func makeTerminalRaw() error {
	fd := int(os.Stdin.Fd())
	var termios syscall.Termios
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&termios))); err != 0 {
		return err
	}

	termios.Lflag &^= syscall.ICANON | syscall.ECHO
	if _, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&termios))); err != 0 {
		return err
	}

	return nil
}
