package terminal

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"slices"
	"strings"
	"sync"

	tea "charm.land/bubbletea/v2"
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/vt"
	"github.com/charmbracelet/x/xpty"
	"phub/internal/launcher"
)

const (
	defaultWidth  = 80
	defaultHeight = 24
)

type Session struct {
	pty         xpty.Pty
	command     *exec.Cmd
	ctx         context.Context
	cancel      context.CancelFunc
	emulator    *vt.SafeEmulator
	slaveClosed bool
	closeOnce   sync.Once
	closeErr    error
}

func Start(ctx context.Context, shell string, projectPath string, width int, height int) (*Session, error) {
	width, height = normalizeSize(width, height)
	ptyDevice, err := xpty.NewPty(width, height)
	if err != nil {
		return nil, err
	}

	sessionContext, cancel := context.WithCancel(ctx)
	command, err := launcher.NewInteractiveShellCommandContext(sessionContext, shell, projectPath)
	if err != nil {
		closeErr := ptyDevice.Close()
		cancel()
		return nil, errors.Join(err, closeErr)
	}
	configurePTYCommand(command)
	command.Env = terminalEnvironment(os.Environ())
	if err := ptyDevice.Start(command); err != nil {
		closeErr := ptyDevice.Close()
		cancel()
		return nil, errors.Join(err, closeErr)
	}
	slaveClosed, err := closePTYSlave(ptyDevice)
	if err != nil {
		cancel()
		closeErr := ptyDevice.Close()
		return nil, errors.Join(err, closeErr)
	}

	session := &Session{
		pty:         ptyDevice,
		command:     command,
		ctx:         sessionContext,
		cancel:      cancel,
		emulator:    vt.NewSafeEmulator(width, height),
		slaveClosed: slaveClosed,
	}
	go session.forwardInput()
	return session, nil
}

func (s *Session) Read() ([]byte, error) {
	buffer := make([]byte, 4096)
	count, err := s.pty.Read(buffer)
	return slices.Clone(buffer[:count]), err
}

func (s *Session) Wait() error {
	return xpty.WaitProcess(s.ctx, s.command)
}

func (s *Session) WriteOutput(output []byte) error {
	_, err := s.emulator.Write(output)
	return err
}

func (s *Session) SendKey(key tea.Key) error {
	if key.Text != "" && (key.Mod == 0 || key.Mod == uv.ModShift) {
		s.emulator.SendText(key.Text)
		return nil
	}
	s.emulator.SendKey(uv.KeyPressEvent(uv.Key{
		Text:        key.Text,
		Mod:         key.Mod,
		Code:        key.Code,
		ShiftedCode: key.ShiftedCode,
		BaseCode:    key.BaseCode,
	}))
	return nil
}

func (s *Session) Resize(width int, height int) error {
	width, height = normalizeSize(width, height)
	if err := s.pty.Resize(width, height); err != nil {
		return err
	}
	s.emulator.Resize(width, height)
	return nil
}

func (s *Session) Render() string {
	return s.emulator.Render()
}

func (s *Session) CursorPosition() (int, int) {
	position := s.emulator.CursorPosition()
	return position.X, position.Y
}

func (s *Session) Close() error {
	s.closeOnce.Do(func() {
		s.cancel()
		ptyClose := s.pty.Close()
		if s.slaveClosed {
			ptyClose = closePTYMaster(s.pty)
		}
		s.closeErr = errors.Join(s.emulator.Close(), ptyClose)
	})
	return s.closeErr
}

func (s *Session) forwardInput() {
	buffer := make([]byte, 256)
	for {
		count, err := s.emulator.Read(buffer)
		if count > 0 {
			if _, writeErr := s.pty.Write(buffer[:count]); writeErr != nil {
				return
			}
		}
		if err != nil {
			return
		}
	}
}

func normalizeSize(width int, height int) (int, int) {
	if width < 1 {
		width = defaultWidth
	}
	if height < 1 {
		height = defaultHeight
	}
	return width, height
}

func terminalEnvironment(environment []string) []string {
	result := slices.Clone(environment)
	for _, value := range []string{
		"TERM=xterm-256color",
		"COLORTERM=truecolor",
		"TERM_PROGRAM=phub",
	} {
		key, _, _ := strings.Cut(value, "=")
		replaced := false
		for index, current := range result {
			currentKey, _, _ := strings.Cut(current, "=")
			if currentKey == key {
				result[index] = value
				replaced = true
				break
			}
		}
		if !replaced {
			result = append(result, value)
		}
	}
	return result
}

func closePTYSlave(ptyDevice xpty.Pty) (bool, error) {
	type slaveProvider interface {
		Slave() *os.File
	}
	provider, ok := ptyDevice.(slaveProvider)
	if !ok {
		return false, nil
	}
	return true, provider.Slave().Close()
}

func closePTYMaster(ptyDevice xpty.Pty) error {
	type masterProvider interface {
		Master() *os.File
	}
	provider, ok := ptyDevice.(masterProvider)
	if !ok {
		return ptyDevice.Close()
	}
	err := provider.Master().Close()
	if errors.Is(err, os.ErrClosed) || errors.Is(err, os.ErrInvalid) {
		return nil
	}
	return err
}
