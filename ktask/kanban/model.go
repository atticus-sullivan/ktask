package kanban

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const margin = 4

type Board struct {
	help     help.Model
	loaded   bool
	Focused  int
	Cols     []Column
	quitting bool
}

type focus int

// NewDefaultBoard creates a new kanban board with To Do, In Progress, and Done
// columns.
func NewDefaultBoard(cols []Column) *Board {
	help := help.New()
	help.ShowAll = false
	b := &Board{Cols: cols, help: help}
	for i, c := range cols {
		if c.Focused() {
			b.Focused = int(i)
		}
		cols[i].board = b
		cols[i].cnt = uint(len(cols))
	}

	return b
}

func (m *Board) Init() tea.Cmd {
	return nil
}

func mod(a, b int) int {
	return ((a % b) + b) % b
}

func (m *Board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.help.Width = msg.Width - margin
		msg.Height -= lipgloss.Height(m.help.View(keys))
		for i := 0; i < len(m.Cols); i++ {
			var res tea.Model
			res, cmd = m.Cols[i].Update(msg)
			m.Cols[i] = res.(Column)
			cmds = append(cmds, cmd)
		}
		m.loaded = true
		return m, tea.Batch(cmds...)
	case MoveMsg:
		cmds = append(cmds, m.Cols[mod((m.Focused+msg.direction), len(m.Cols))].Set(APPEND, msg.item))
	case tea.KeyMsg:
		if !m.Cols[m.Focused].List.SettingFilter() {
			switch {
			case key.Matches(msg, keys.Help):
				m.help.ShowAll = !m.help.ShowAll
				cmds = append(cmds, tea.WindowSize())
				return m, tea.Batch(cmds...)
			case key.Matches(msg, keys.Quit):
				m.quitting = true
				return m, tea.Quit
			case key.Matches(msg, keys.Left):
				m.Cols[m.Focused].Blur()
				m.Focused = mod((m.Focused - 1), len(m.Cols))
				m.Cols[m.Focused].Focus()
			case key.Matches(msg, keys.Right):
				m.Cols[m.Focused].Blur()
				m.Focused = mod((m.Focused + 1), len(m.Cols))
				m.Cols[m.Focused].Focus()
			}
		}
	}
	res, cmd := m.Cols[m.Focused].Update(msg)
	cmds = append(cmds, cmd)
	if _, ok := res.(Column); ok {
		m.Cols[m.Focused] = res.(Column)
	} else {
		// if it's not a column, switch to the returned model
		return res, tea.Batch(cmds...)
	}
	return m, tea.Batch(cmds...)
}

// Changing to pointer receiver to get back to this model after adding a new task via the form... Otherwise I would need to pass this model along to the form and it becomes highly coupled to the other models.
func (m *Board) View() string {
	if m.quitting {
		return ""
	}
	if !m.loaded {
		return "loading..."
	}
	var cs []string
	for _, c := range m.Cols {
		cs = append(cs, c.View())
	}
	board := lipgloss.JoinHorizontal(
		lipgloss.Left,
		cs...,
	)
	return lipgloss.JoinVertical(lipgloss.Left, board, m.help.View(keys))
}
