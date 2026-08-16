package ui

import (
	"context"
	"errors"
	"fmt"

	tea "charm.land/bubbletea/v2"
	"phub/internal/git"
	"phub/internal/metadata"
)

type metadataLoader func(ctx context.Context, projects []Project) (map[string]metadata.Info, error)

type projectsMetadataMsg struct {
	infos map[string]metadata.Info
	err   error
}

func (m model) metadataCommand() tea.Cmd {
	loader := m.loadMetadata
	if loader == nil {
		loader = defaultMetadataLoader
	}
	projects := make([]Project, len(m.projects))
	copy(projects, m.projects)
	return func() tea.Msg {
		infos, err := loader(context.Background(), projects)
		return projectsMetadataMsg{infos: infos, err: err}
	}
}

func defaultMetadataLoader(ctx context.Context, projects []Project) (map[string]metadata.Info, error) {
	paths := make([]string, len(projects))
	for index, current := range projects {
		paths[index] = current.Path
	}
	infos, err := metadata.LoadAll(ctx, paths)
	if err != nil {
		return nil, err
	}
	for _, info := range infos {
		if errors.Is(info.GitErr, git.ErrNotFound) {
			return infos, git.ErrNotFound
		}
	}
	return infos, nil
}

func gitStatusText(status git.Status) string {
	switch {
	case status.Modified > 0 && status.Untracked > 0:
		return fmt.Sprintf("● %d modified ? %d untracked", status.Modified, status.Untracked)
	case status.Modified > 0:
		return fmt.Sprintf("● %d modified", status.Modified)
	case status.Untracked > 0:
		return fmt.Sprintf("? %d untracked", status.Untracked)
	default:
		return "✓ clean"
	}
}
