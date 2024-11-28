// Package todo provides functionality for cycling through TODO states and priorities
// in a Logseq-compatible format.
package todo

import "strings"

// States represents the valid TODO states in order of cycling
var States = []string{"TODO", "DOING", "DONE"}

// Priorities represents the valid priority levels in order of cycling
var Priorities = []string{"[#A]", "[#B]", "[#C]"}

// CycleState takes a line of text and returns the same line with the next TODO state.
// If no state exists, it adds "TODO" at the start of the line.
func CycleState(line string) string {
	trimmed := strings.TrimLeft(line, " \t")
	if trimmed == "" {
		return line
	}

	// Get indentation
	indent := strings.Repeat(" ", len(line)-len(trimmed))

	// Find current state
	currentState := ""
	for _, state := range States {
		if strings.HasPrefix(trimmed, state+" ") {
			currentState = state
			break
		}
	}

	// If no current state, add TODO
	if currentState == "" {
		return indent + "TODO " + trimmed
	}

	// Find next state
	var nextState string
	for i, state := range States {
		if state == currentState && i < len(States)-1 {
			nextState = States[i+1]
			break
		}
	}

	// Replace current state with next state
	return indent + nextState + trimmed[len(currentState):]
}

// CyclePriority takes a line of text and returns the same line with the next priority level.
// Only adds/cycles priority if the line starts with a TODO state.
func CyclePriority(line string) string {
	trimmed := strings.TrimLeft(line, " \t")
	if trimmed == "" {
		return line
	}

	// Get indentation
	indent := strings.Repeat(" ", len(line)-len(trimmed))

	// Check if line starts with a TODO state
	var (
		hasState    bool
		statePrefix string
		content     = trimmed
	)

	for _, state := range States {
		if strings.HasPrefix(trimmed, state+" ") {
			hasState = true
			statePrefix = state + " "
			content = trimmed[len(statePrefix):]
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
		if strings.HasPrefix(content, priority+" ") {
			currentPriority = priority
			break
		}
	}

	// If no current priority, add [#A]
	if currentPriority == "" {
		return indent + statePrefix + "[#A] " + content
	}

	// Find next priority
	var nextPriority string
	for i, priority := range Priorities {
		if priority == currentPriority {
			if i < len(Priorities)-1 {
				nextPriority = Priorities[i+1]
			} else {
				// Remove priority if we're at the end of the cycle
				return indent + statePrefix + content[len(currentPriority)+1:]
			}
			break
		}
	}

	// Replace current priority with next priority
	return indent + statePrefix + nextPriority + content[len(currentPriority):]
}
