package generate

import tea "github.com/charmbracelet/bubbletea"

type TeaModel struct{}

func (m TeaModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m TeaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	// No messages are coming in, so just return the model and no command.
	return m, nil
}

func (m TeaModel) View() string {
	return "Hello, Bubble Tea!"
}
