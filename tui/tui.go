package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/jrswab/lsq/config"
	"github.com/jrswab/lsq/editor"
	"github.com/jrswab/lsq/todo"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type tuiModel struct {
	//viewport viewport.Model
	textarea  textarea.Model
	config    *config.Config
	filepath  string
	statusMsg string
}

func InitialModel(cfg *config.Config, fp string) tuiModel {
	// Read file content for TUI
	content, err := os.ReadFile(fp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading journal file: %v\n", err)
		os.Exit(1)
	}

	var ta = textarea.New()
	ta.SetValue(string(content))
	ta.Focus()
	ta.CharLimit = -1

	return tuiModel{
		textarea: ta,
		config:   cfg,
		filepath: fp,
	}
}

func (m tuiModel) Init() tea.Cmd {
	return textarea.Blink
}

type statusMsg struct{}

func (m tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case statusMsg:
		m.statusMsg = ""
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyCtrlS:
			content := m.textarea.Value()
			err := os.WriteFile(m.filepath, []byte(content), 0644)
			if err != nil {
				m.statusMsg = "Error saving file!"
				return m, tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
					return statusMsg{}
				})
			}

			m.statusMsg = "File saved successfully!"
			return m, tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
				return statusMsg{}
			})

		case tea.KeyCtrlE:
			return m, tea.ExecProcess(
				exec.Command(editor.Select("EDITOR"), m.filepath),
				nil,
			)
		// Cycle through TODO states:
		case tea.KeyCtrlT:
			// Get current content and line number
			content := m.textarea.Value()
			lineNum := m.textarea.Line()

			// Split content into lines
			lines := strings.Split(content, "\n")

			// Make sure we're within bounds
			if lineNum < len(lines) {
				// Update the specific line
				lines[lineNum] = todo.CycleState(lines[lineNum])

				// Join lines back together
				newContent := strings.Join(lines, "\n")

				// Update textarea
				m.textarea.SetValue(newContent)
			}
		case tea.KeyCtrlP:
			content := m.textarea.Value()
			lineNum := m.textarea.Line()

			lines := strings.Split(content, "\n")

			if lineNum < len(lines) {
				lines[lineNum] = todo.CyclePriority(lines[lineNum])
				newContent := strings.Join(lines, "\n")
				m.textarea.SetValue(newContent)
			}
		}

	case tea.WindowSizeMsg:
		m.textarea.SetWidth(msg.Width)
		m.textarea.SetHeight(msg.Height - 2)
	}

	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m tuiModel) View() string {
	var footer string
	if m.statusMsg != "" {
		footer = m.statusMsg
	} else {
		footer = "^S save, ^E external editor, ^C quit"
	}

	return fmt.Sprintf(
		"LSQ TUI Mode - %s\n%s\n%s",
		filepath.Base(m.filepath),
		m.textarea.View(),
		footer,
	)
}
