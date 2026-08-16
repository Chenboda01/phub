package ui

import (
	"slices"
	"strings"
	"testing"
)

func TestProjectDisplayNames_keepsUniqueNames(t *testing.T) {
	// Given
	projects := []Project{
		{Name: "Alpha", Path: "/home/user/Alpha"},
		{Name: "Beta", Path: "/home/user/Beta"},
	}

	// When
	names := projectDisplayNames(projects)

	// Then
	if !slices.Equal(names, []string{"Alpha", "Beta"}) {
		t.Fatalf("display names = %q, want unique names unchanged", names)
	}
}

func TestProjectDisplayNames_addsParentDirectory_whenNamesCollide(t *testing.T) {
	// Given
	projects := []Project{
		{Name: "phub", Path: "/home/user/phub"},
		{Name: "phub", Path: "/home/user/Projects/phub"},
	}

	// When
	names := projectDisplayNames(projects)

	// Then
	if want := []string{"phub (user)", "phub (Projects)"}; !slices.Equal(names, want) {
		t.Fatalf("display names = %q, want %q", names, want)
	}
}

func TestProjectDisplayNames_fallsBackToFullPath_whenParentsAlsoCollide(t *testing.T) {
	// Given
	projects := []Project{
		{Name: "phub", Path: "/a/shared/phub"},
		{Name: "phub", Path: "/b/shared/phub"},
	}

	// When
	names := projectDisplayNames(projects)

	// Then
	if want := []string{"phub (/a/shared/phub)", "phub (/b/shared/phub)"}; !slices.Equal(names, want) {
		t.Fatalf("display names = %q, want %q", names, want)
	}
}

func TestModelRendersDisambiguatedNames_whenProjectsShareNames(t *testing.T) {
	// Given
	current := newModel([]Project{
		{Name: "phub", Path: "/home/user/phub"},
		{Name: "phub", Path: "/home/user/Projects/phub"},
	}, "/bin/sh")

	// When
	content := current.render()

	// Then
	for _, expected := range []string{"phub (user)", "phub (Projects)"} {
		if !strings.Contains(content, expected) {
			t.Fatalf("rendered list missing %q:\n%s", expected, content)
		}
	}
}
