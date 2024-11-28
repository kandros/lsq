// Package todo provides functionality for cycling through TODO states and priorities
// in a Logseq-compatible format.
package todo

import (
	"fmt"
	"strings"
)

// States represents the valid TODO states in order of cycling
var States = []string{"TODO", "DOING", "DONE"}

// Priorities represents the valid priority levels in order of cycling
var Priorities = []string{"[#A]", "[#B]", "[#C]"}

// CycleState takes a line of text and returns the same line with the next TODO state.
// If no state exists, it adds "TODO" at the start of the line.
func CycleState(line string) string {
	var cutset = " \t-"

	trimmed := strings.TrimLeft(line, cutset)
	if trimmed == "" {
		return line
	}

	// Get indentation count.
	// The -2 ensures we don't count the bullet point formatting as part of the indentation level.
	indent := strings.Repeat("\t", (len(line)-2)-len(trimmed))

	// Find current state
	var currentState string
	for _, state := range States {
		if strings.HasPrefix(trimmed, fmt.Sprintf("%s ", state)) {
			currentState = state
			break
		}
	}

	switch currentState {
	case States[0]: // TODO
		trimmed = strings.TrimPrefix(trimmed, fmt.Sprintf("%s ", currentState))
		return fmt.Sprintf("%s- DOING %s", indent, trimmed)
	case States[1]: // DOING
		trimmed = strings.TrimPrefix(trimmed, fmt.Sprintf("%s ", currentState))
		return fmt.Sprintf("%s- DONE %s", indent, trimmed)
	case States[2]: // DONE
		// take out "DONE" from the trimmed string before returning
		trimmed = strings.TrimPrefix(trimmed, fmt.Sprintf("%s ", currentState))
		return fmt.Sprintf("%s- TODO %s", indent, trimmed)
	default:
		return fmt.Sprintf("%s- TODO %s", indent, trimmed)
	}
}

// CyclePriority takes a line of text and returns the same line with the next priority level.
// Only adds/cycles priority if the line starts with a TODO state.
func CyclePriority(line string) string {
	var cutset = " \t-"

	trimmed := strings.TrimLeft(line, cutset)
	if trimmed == "" {
		return line
	}

	// Get indentation count.
	// The -2 ensures we don't count the bullet point formatting as part of the indentation level.
	indent := strings.Repeat("\t", (len(line)-2)-len(trimmed))

	// Check if line starts with a TODO state
	var (
		hasState    bool
		statePrefix string
		content     = trimmed
	)

	for _, state := range States {
		if strings.HasPrefix(trimmed, fmt.Sprintf("%s ", state)) {
			hasState = true
			statePrefix = state
			content = trimmed[len(statePrefix)+1:] // the "+1" takes out the leading space
			break
		}
	}

	// If no TODO state, return original line unchanged
	if !hasState {
		return line
	}

	// Find current priority
	var currentPriority string
	for _, priority := range Priorities {
		if strings.HasPrefix(content, fmt.Sprintf("%s ", priority)) {
			currentPriority = priority
			content = content[len(currentPriority)+1:]
			break
		}
	}

	switch currentPriority {
	case Priorities[0]: // [#A]
		return fmt.Sprintf("%s- %s %s %s", indent, statePrefix, Priorities[1], content)
	case Priorities[1]: // [#B]
		return fmt.Sprintf("%s- %s %s %s", indent, statePrefix, Priorities[2], content)
	case Priorities[2]: // [#C]
		return fmt.Sprintf("%s- %s %s", indent, statePrefix, content)
	default:
		// When no TODO state is found the function returns the original line
		// before checking for a current priority.
		return fmt.Sprintf("%s- %s [#A] %s", indent, statePrefix, content)
	}
}
