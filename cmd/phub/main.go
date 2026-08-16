package main

import (
	"context"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"phub/internal/discovery"
	"phub/internal/ui"
)

func main() {
	program := tea.NewProgram(ui.NewWithLoader(os.Getenv("SHELL"), func(scope ui.ProjectScope) ([]ui.Project, error) {
		return discoverProjectsForScope(context.Background(), scope)
	}))
	if _, err := program.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "phub could not start:", err)
		os.Exit(1)
	}
}

func discoverProjects(ctx context.Context) ([]ui.Project, error) {
	return discoverProjectsForScope(ctx, ui.ScopeAllLocal)
}

func discoverProjectsForScope(ctx context.Context, scope ui.ProjectScope) ([]ui.Project, error) {
	options, err := discovery.LoadOptions()
	if err != nil {
		return nil, fmt.Errorf("load discovery options: %w", err)
	}

	scanner, err := discovery.New(options.MaxDepth())
	if err != nil {
		return nil, fmt.Errorf("create scanner: %w", err)
	}

	projects, err := scanner.Scan(ctx, options.Roots())
	if err != nil {
		return nil, fmt.Errorf("scan projects: %w", err)
	}

	filtered, err := discovery.FilterProjects(ctx, projects, discoveryScope(scope))
	if err != nil {
		return nil, fmt.Errorf("filter projects: %w", err)
	}

	result := make([]ui.Project, len(filtered))
	for index, discovered := range filtered {
		result[index] = ui.Project{Name: discovered.Name(), Path: discovered.Path()}
	}

	return result, nil
}

func discoveryScope(scope ui.ProjectScope) discovery.ProjectScope {
	if scope == ui.ScopeAllLocal {
		return discovery.AllLocal
	}
	return discovery.GitHubOnly
}
