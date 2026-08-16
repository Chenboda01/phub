package ui

import (
	"context"
	"errors"
	"io"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestModelReturnsToProjects_whenEmbeddedTerminalExits(t *testing.T) {
	// Given
	session := &fakeTerminalSession{
		readData: []byte("shell ready"),
		readErr:  io.EOF,
	}
	current := newModel(testProjects(), "/bin/sh")
	current.startTerminal = func(_ context.Context, shell, path string, width, height int) (terminalSession, error) {
		if shell != "/bin/sh" || path != testProjects()[0].Path || width != 80 || height != 22 {
			t.Fatalf("terminal start args = %q, %q, %dx%d", shell, path, width, height)
		}
		return session, nil
	}
	current.width = 80
	current.height = 24

	// When
	starting, startCommand := updateModel(t, current, tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}))
	started, readCommand := updateModel(t, starting, startCommand())
	returned, _ := updateModel(t, started, readCommand())

	// Then
	if returned.terminal != nil {
		t.Fatal("embedded terminal remained active after EOF")
	}
	if returned.notice != "Returned from Alpha." {
		t.Fatalf("notice = %q, want return notice", returned.notice)
	}
	if !session.closed {
		t.Fatal("embedded terminal was not closed")
	}
	if string(session.output) != "shell ready" {
		t.Fatalf("terminal output = %q, want shell ready", session.output)
	}
}

func TestModelReturnsShellStatusAsNotice_whenEmbeddedTerminalExitsNonzero(t *testing.T) {
	// Given
	session := &fakeTerminalSession{
		readErr: io.EOF,
		waitErr: errors.New("exit status 1"),
	}
	current := newModel(testProjects(), "/bin/sh")
	current.terminal = session
	current.terminalProject = "Alpha"

	// When
	returned, _ := updateModel(t, current, terminalOutputMsg{session: session, readErr: io.EOF, waitErr: session.waitErr})

	// Then
	if returned.terminal != nil {
		t.Fatal("embedded terminal remained active after nonzero exit")
	}
	if returned.notice != "Returned from Alpha (shell status: exit status 1)." {
		t.Fatalf("notice = %q, want non-error shell status notice", returned.notice)
	}
}

func TestModelForwardsKeys_whenEmbeddedTerminalIsActive(t *testing.T) {
	// Given
	session := &fakeTerminalSession{readErr: errors.New("still running")}
	current := newModel(testProjects(), "/bin/sh")
	current.terminal = session

	// When
	_, command := updateModel(t, current, tea.KeyPressMsg(tea.Key{Text: "e", Code: 'e'}))
	if command != nil {
		t.Fatal("key forwarding returned an asynchronous command")
	}

	// Then
	if len(session.keys) != 1 || session.keys[0].Text != "e" {
		t.Fatalf("forwarded keys = %q, want e", session.keys)
	}
}

type fakeTerminalSession struct {
	readData []byte
	readErr  error
	waitErr  error
	output   []byte
	keys     []tea.Key
	closed   bool
	cursorX  int
	cursorY  int
}

func (s *fakeTerminalSession) Read() ([]byte, error) {
	return s.readData, s.readErr
}

func (s *fakeTerminalSession) Wait() error {
	return s.waitErr
}

func (s *fakeTerminalSession) Close() error {
	s.closed = true
	return nil
}

func (s *fakeTerminalSession) Resize(int, int) error {
	return nil
}

func (s *fakeTerminalSession) Render() string {
	return "terminal"
}

func (s *fakeTerminalSession) CursorPosition() (int, int) {
	return s.cursorX, s.cursorY
}

func (s *fakeTerminalSession) SendKey(key tea.Key) error {
	s.keys = append(s.keys, key)
	return nil
}

func (s *fakeTerminalSession) WriteOutput(output []byte) error {
	s.output = append(s.output, output...)
	return nil
}
