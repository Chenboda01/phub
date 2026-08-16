package ui

import (
	"fmt"
	"path/filepath"
)

func projectDisplayNames(projects []Project) []string {
	names := make([]string, len(projects))
	for index, current := range projects {
		names[index] = current.Name
	}

	grouped := make(map[string][]int)
	for index, current := range projects {
		grouped[current.Name] = append(grouped[current.Name], index)
	}

	for _, indices := range grouped {
		if len(indices) < 2 {
			continue
		}
		parentsUnique := true
		seenParents := make(map[string]struct{}, len(indices))
		for _, index := range indices {
			parent := filepath.Base(filepath.Dir(projects[index].Path))
			if _, exists := seenParents[parent]; exists {
				parentsUnique = false
				break
			}
			seenParents[parent] = struct{}{}
		}
		for _, index := range indices {
			if parentsUnique {
				parent := filepath.Base(filepath.Dir(projects[index].Path))
				names[index] = fmt.Sprintf("%s (%s)", projects[index].Name, parent)
			} else {
				names[index] = fmt.Sprintf("%s (%s)", projects[index].Name, projects[index].Path)
			}
		}
	}

	return names
}
