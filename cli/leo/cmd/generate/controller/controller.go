package controller

/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

type Styles struct {
	BorderColor lipgloss.Color
	InputField  lipgloss.Style
}

func DefaultStyles() *Styles {
	s := new(Styles)
	s.BorderColor = lipgloss.Color("36")
	s.InputField = lipgloss.NewStyle().BorderForeground(s.BorderColor).BorderStyle(lipgloss.NormalBorder()).Padding(1).Width(80)
	return s
}

// modelCmd represents the model command
var ControllerCmd = &cobra.Command{
	Use:   "controller",
	Short: "Generate controller.",
	Long:  `Generate controller.`,
	Run: func(cmd *cobra.Command, args []string) {
		inputField := textinput.NewModel()
		inputField.Placeholder = "Enter Controller Name"
		inputField.Focus()
		model := Controller{
			Name:       "",
			done:       false,
			inputField: inputField,
		}

		p := tea.NewProgram(model)
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}

	},
}

type Controller struct {
	Name       string
	done       bool
	height     int
	width      int
	inputField textinput.Model
}

func (m Controller) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m Controller) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			m.done = true
			m.Name = m.inputField.Value()
			return m, nil
		}
		m.inputField, cmd = m.inputField.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
	}

	// No messages are coming in, so just return the model and no command.
	return m, nil
}

func (m Controller) View() string {

	if m.height == 0 {
		return "Loading ...."
	}

	if m.done {
		generateController(m.Name)

		tea.Quit()
		os.Exit(0)

	}

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Center,
			"Enter Controller Name",
			m.inputField.View(),
		),
	)

}
