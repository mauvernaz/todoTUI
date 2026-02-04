package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Application states
type appState int

const (
	browsing appState = iota
	inputting
	helping
)

// Styles using Lip Gloss for a minimalist aesthetic
var (
	// Selected item style: bold with a subtle accent color
	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")).
			Bold(true).
			PaddingLeft(2)

	// Unselected items: dimmed for visual hierarchy
	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("245")).
			PaddingLeft(4)

	// Cursor indicator for selected item
	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")).
			Bold(true)

	// Title style
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("99")).
			Bold(true).
			MarginBottom(1)

	// Help text style
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)

	// Input prompt style
	inputPromptStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("212")).
				Bold(true).
				MarginTop(1)
)

// Model holds the application state following the Elm architecture
type model struct {
	tasks    []string          // List of tasks
	cursor   int               // Currently selected task index
	state    appState          // Current application state (browsing or inputting)
	input    textinput.Model   // Text input component for adding tasks
	quitting bool              // Flag to indicate app is quitting
}

// initialModel creates and returns the initial model state
func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Enter task name..."
	ti.CharLimit = 100
	ti.Width = 40

	return model{
		tasks:  []string{},
		cursor: 0,
		state:  browsing,
		input:  ti,
	}
}

// Init implements tea.Model - called once when the program starts
func (m model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model - handles all messages and user input
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle key presses based on current state
		switch m.state {
		case browsing:
			return m.updateBrowsing(msg)
		case inputting:
			return m.updateInputting(msg)
		case helping:
			return m.updateHelping(msg)
		}
	}
	return m, nil
}

// updateBrowsing handles key input when in browse mode
func (m model) updateBrowsing(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	// Quit commands
	case "q", "esc", "ctrl+c":
		m.quitting = true
		return m, tea.Quit

	// Navigation: move up
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	// Navigation: move down
	case "down", "j":
		if m.cursor < len(m.tasks)-1 {
			m.cursor++
		}

	// Add new task: switch to input mode
	case "n", "a":
		m.state = inputting
		m.input.Focus()
		return m, textinput.Blink

	// Toggle help
	case "?", "h":
		m.state = helping

	// Delete/Complete task
	case "x", "backspace", "d":
		if len(m.tasks) > 0 {
			// Remove the selected task
			m.tasks = append(m.tasks[:m.cursor], m.tasks[m.cursor+1:]...)
			// Adjust cursor if needed
			if m.cursor >= len(m.tasks) && m.cursor > 0 {
				m.cursor--
			}
		}
	}
	return m, nil
}

// updateHelping handles key input when in help mode
func (m model) updateHelping(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q", "?", "h", "enter":
		m.state = browsing
	}
	return m, nil
}

// updateInputting handles key input when in input mode
func (m model) updateInputting(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	// Cancel input and return to browse mode
	case "esc":
		m.state = browsing
		m.input.Reset()
		return m, nil

	// Submit the new task
	case "enter":
		value := m.input.Value()
		if value != "" {
			m.tasks = append(m.tasks, value)
			m.cursor = len(m.tasks) - 1 // Move cursor to new task
		}
		m.state = browsing
		m.input.Reset()
		return m, nil
	}

	// Update the text input component
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// View implements tea.Model - renders the UI
func (m model) View() string {
	if m.quitting {
		return "Goodbye! ✨\n"
	}

	var s string

	// Title
	s += titleStyle.Render("📝 To-Do") + "\n\n"

	// Render task list
	if len(m.tasks) == 0 {
		s += normalStyle.Render("No tasks yet. Press 'n' to add one.") + "\n"
	} else {
		for i, task := range m.tasks {
			if i == m.cursor {
				// Selected item with cursor indicator
				s += cursorStyle.Render("→ ") + selectedStyle.Render(task) + "\n"
			} else {
				// Unselected items
				s += normalStyle.Render(task) + "\n"
			}
		}
	}

	// Render input field when in input mode
	if m.state == inputting {
		s += "\n" + inputPromptStyle.Render("New Task:") + "\n"
		s += "  " + m.input.View() + "\n"
	}

	// Render help text or help view
	s += "\n"
	if m.state == browsing {
		s += helpStyle.Render("↑/↓: navigate • n: add • x: delete • ?: help • q: quit")
	} else if m.state == inputting {
		s += helpStyle.Render("enter: save • esc: cancel")
	} else if m.state == helping {
		s = titleStyle.Render("📖 Help & Commands") + "\n\n"
		s += lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("Navigation:") + "\n"
		s += normalStyle.Render("↑ / k      - Move selection up") + "\n"
		s += normalStyle.Render("↓ / j      - Move selection down") + "\n\n"

		s += lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("Tasks:") + "\n"
		s += normalStyle.Render("n / a      - Add a new task (New/Add)") + "\n"
		s += normalStyle.Render("x / d / bk - Remove selected task (Delete)") + "\n"
		s += normalStyle.Render("Enter      - Confirm new task (In input mode)") + "\n\n"

		s += lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("Application:") + "\n"
		s += normalStyle.Render("? / h      - Toggle this help view") + "\n"
		s += normalStyle.Render("q / Esc    - Return to list or Quit") + "\n"
		s += normalStyle.Render("Ctrl+C     - Force quit") + "\n\n"

		s += helpStyle.Render("Press any key to return...")
	}

	return s + "\n"
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
