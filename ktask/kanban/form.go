package kanban

import (
	"ktask/ktask"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Form struct {
	help        help.Model
	title       textinput.Model
	description textarea.Model
	col         Column
	index       int
	totalWidth  int
	totalHeight int
}

func newDefaultForm() *Form {
	return NewForm("task name", "")
}

func NewForm(title, description string) *Form {
	form := Form{
		help:        help.New(),
		title:       textinput.New(),
		description: textarea.New(),
	}

	form.title.Width = 10
	form.description.SetWidth(10)
	form.description.SetHeight(10)

	form.title.Placeholder = title
	form.description.Placeholder = description
	form.title.Focus()
	return &form
}

func (f Form) Init() tea.Cmd {
	return nil
}

func (f Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case Column:
		f.col = msg
		f.col.List.Index()
	case tea.WindowSizeMsg:
		f.setSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Back):
			return f.col.board, nil
		case key.Matches(msg, keys.Enter):
			if f.title.Focused() {
				f.title.Blur()
				f.description.Focus()
				return f, textarea.Blink
			}
			// Return the completed form as a message.
			if f.title.Value() != "" {
				tag := f.description.Value()
				sep := ""
				if tag != "" {
					if !strings.HasPrefix(tag, "#") {
						tag = "#" + tag
					}
					sep = " "
				}
				item := ktask.NewEntry(ktask.Name(strings.Split(tag+sep+f.title.Value(), "\n")), time.Now(), time.Now(), f.index)
				return f.col.board.Update(item)
			}
			return f.col.board, nil
		}
	}
	if f.title.Focused() {
		f.title, cmd = f.title.Update(msg)
		return f, cmd
	}
	f.description, cmd = f.description.Update(msg)
	return f, cmd
}

func (f Form) View() string {
	return lipgloss.Place(
		f.totalWidth, f.totalHeight, 0.5, 0.75,
		lipgloss.NewStyle().
			Padding(0, 0).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Render(
				lipgloss.JoinVertical(
					lipgloss.Left,
					"Create a new task",
					f.title.View(),
					f.description.View(),
					f.help.View(keys),
				),
			),
	)
}

func (f *Form) setSize(width, height int) {
	f.totalWidth, f.totalHeight = width, height
	f.title.Width = 80
	f.description.SetWidth(80)
	f.description.SetHeight(4)
}
