package ui

import (
	"context"
	"fmt"

	tea "charm.land/bubbletea/v2"
	"phub/internal/terminal"
)

type terminalSession interface {
	Read() ([]byte, error)
	Wait() error
	WriteOutput([]byte) error
	SendKey(tea.Key) error
	Resize(int, int) error
	Render() string
	CursorPosition() (int, int)
	Close() error
}

type terminalStarter func(context.Context, string, string, int, int) (terminalSession, error)

type terminalStartedMsg struct {
	session     terminalSession
	projectName string
	err         error
}

type terminalOutputMsg struct {
	session  terminalSession
	readErr  error
	waitErr  error
	writeErr error
}

type terminalResizeMsg struct {
	session terminalSession
	err     error
}

func defaultTerminalStarter(ctx context.Context, shell string, projectPath string, width int, height int) (terminalSession, error) {
	return terminal.Start(ctx, shell, projectPath, width, height)
}

func (m model) beginTerminal() (tea.Model, tea.Cmd) {
	if m.startingTerminal || m.terminal != nil {
		return m, nil
	}
	if m.startTerminal == nil {
		m.notice = "Embedded terminal is unavailable."
		return m, nil
	}
	if len(m.projects) == 0 {
		m.notice = "No projects are available to open."
		return m, nil
	}

	project := m.projects[m.selected]
	width, height := terminalViewportDimensions(m.width, m.height)
	start := m.startTerminal
	ctx := m.terminalContext
	m.startingTerminal = true
	m.notice = fmt.Sprintf("Opening %s...", project.Name)
	return m, func() tea.Msg {
		session, err := start(ctx, m.shell, project.Path, width, height)
		return terminalStartedMsg{session: session, projectName: project.Name, err: err}
	}
}

func readTerminal(session terminalSession) tea.Cmd {
	return func() tea.Msg {
		output, readErr := session.Read()
		var writeErr error
		if len(output) > 0 {
			writeErr = session.WriteOutput(output)
		}
		var waitErr error
		if readErr != nil {
			waitErr = session.Wait()
		}
		return terminalOutputMsg{session: session, readErr: readErr, waitErr: waitErr, writeErr: writeErr}
	}
}

func resizeTerminal(session terminalSession, width int, height int) tea.Cmd {
	return func() tea.Msg {
		width, height = terminalViewportDimensions(width, height)
		return terminalResizeMsg{session: session, err: session.Resize(width, height)}
	}
}

func terminalDimensions(width int, height int) (int, int) {
	if width < 1 {
		width = 80
	}
	if height < 1 {
		height = 24
	}
	return width, height
}

func (m model) handleTerminalOutput(message terminalOutputMsg) (tea.Model, tea.Cmd) {
	if message.session != m.terminal {
		return m, nil
	}
	if message.writeErr != nil {
		return m.finishTerminal(fmt.Sprintf("Terminal exited with error: %v", message.writeErr))
	}
	if message.readErr != nil {
		notice := fmt.Sprintf("Returned from %s.", m.terminalProject)
		if message.waitErr != nil {
			notice = fmt.Sprintf("Returned from %s (shell status: %v).", m.terminalProject, message.waitErr)
		}
		return m.finishTerminal(notice)
	}
	return m, readTerminal(m.terminal)
}

func (m model) finishTerminal(notice string) (tea.Model, tea.Cmd) {
	projectName := m.terminalProject
	closeErr := m.terminal.Close()
	m.terminal = nil
	m.terminalProject = ""
	m.startingTerminal = false

	if closeErr != nil {
		m.notice = fmt.Sprintf("Could not close terminal: %v", closeErr)
	} else if notice == "" {
		m.notice = fmt.Sprintf("Returned from %s.", projectName)
	} else {
		m.notice = notice
	}
	return m, m.metadataCommand()
}
