package main

// A simple program demonstrating the text input and spinner components from the Bubbles
// component library.

import (
	"log"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

type model struct {
	textInput textinput.Model
	spinner   spinner.Model
	loading   bool
	err       error
	quitting  bool
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Charizard"
	ti.SetVirtualCursor(true)
	ti.Focus()
	ti.CharLimit = 156
	ti.SetWidth(20)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{textInput: ti, spinner: s, loading: true}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var tiCmd, spinCmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			m.quitting = true

			return m, tea.Quit, tea.QuitMsg("Nice, I like %s")

		case "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		}
	}

	// Update text input state and save command
	m.textInput, tiCmd = m.textInput.Update(msg)

	// Update spinner state and save command so animation continues
	m.spinner, spinCmd = m.spinner.Update(msg)

	// Batch both commands together so neither is dropped
	return m, tea.Batch(tiCmd, spinCmd)
}

func (m model) View() tea.View {
	var c *tea.Cursor
	if !m.textInput.VirtualCursor() {
		c = m.textInput.Cursor()
		c.Y += lipgloss.Height(m.headerView())
	}

	// Place the spinner right next to the text input box
	inputLine := lipgloss.JoinHorizontal(lipgloss.Center, m.textInput.View(), " ", m.spinner.View())

	str := lipgloss.JoinVertical(lipgloss.Top, m.headerView(), inputLine, m.footerView())
	if m.quitting {
		str += "\n"
	}

	v := tea.NewView(str)
	v.Cursor = c
	return v
}

func (m model) headerView() string { return "What’s your favorite Pokémon?\n" }
func (m model) footerView() string { return "\n(esc to quit)" }
