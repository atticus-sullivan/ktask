package kanban

import (
	"ktask/ktask"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const APPEND = -1

type Column struct {
	focus  bool
	List   list.Model
	height int
	width  int
	board  *Board
	cnt    uint
}

func (c *Column) Focus() {
	c.List.SetShowHelp(true)
	c.focus = true
}

func (c *Column) Blur() {
	c.List.SetShowHelp(false)
	c.focus = false
}

func (c *Column) Focused() bool {
	return c.focus
}

// NewColumn creates a new column from a list.
func NewColumn(l []list.Item, focus bool) Column {
	defaultList := list.New(l, list.NewDefaultDelegate(), 0, 0)
	defaultList.SetShowHelp(focus)
	return Column{focus: focus, List: defaultList}
}

func NewColumnFromRecord(r ktask.Record, focus bool) Column {
	items := TasksToItems(r.Entries())
	ret := Column{
		focus: focus,
		List:  list.New(items, list.NewDefaultDelegate(), 0, 0),
	}
	ret.List.SetShowHelp(focus)
	ret.List.Title = string(r.Stage())

	km := &ret.List.KeyMap
	km.CloseFullHelp.Unbind()
	km.ShowFullHelp.Unbind()
	km.Quit.Unbind()

	return ret
}

// Init does initial setup for the column.
func (c Column) Init() tea.Cmd {
	return nil
}

// Update handles all the I/O for columns.
func (c Column) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case ktask.Entry:
		return c, c.Set(msg.Index(), msg)
	case tea.WindowSizeMsg:
		c.setSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		if !c.List.SettingFilter() {
			switch {
			case key.Matches(msg, keys.Edit):
				if len(c.List.VisibleItems()) != 0 {
					item := c.List.SelectedItem().(ktask.Entry)
					f := NewForm(item.Title(), item.Description(), item.CreatedAt(), time.Now())
					f.title.SetValue(item.Title())
					f.description.SetValue(item.Description())
					f.index = c.List.Index()
					f.col = c
					return f.Update(c.board.lastWin) // not nice to have to store lastWin
				}
			case key.Matches(msg, keys.New):
				f := newDefaultForm()
				f.index = APPEND
				f.col = c
				return f.Update(c.board.lastWin) // not nice to have to store lastWin
			case key.Matches(msg, keys.Delete):
				return c, c.DeleteCurrent()
			case key.Matches(msg, keys.Prev):
				return c, c.MoveToPrev()
			case key.Matches(msg, keys.Next):
				return c, c.MoveToNext()
			}
		}
	}
	c.List, cmd = c.List.Update(msg)
	return c, cmd
}

func (c Column) View() string {
	return c.getStyle().Render(c.List.View())
}

func (c *Column) DeleteCurrent() tea.Cmd {
	if len(c.List.VisibleItems()) > 0 {
		c.List.RemoveItem(c.List.Index())
	}

	var cmd tea.Cmd
	c.List, cmd = c.List.Update(nil)
	return cmd
}

// Set adds an item to a column.
func (c *Column) Set(i int, item list.Item) tea.Cmd {
	if i != APPEND {
		return c.List.SetItem(i, item)
	}
	return c.List.InsertItem(APPEND, item)
}

func (c *Column) setSize(width, height int) {
	s := c.getStyle()
	hb, vb := s.GetHorizontalBorderSize(), s.GetVerticalBorderSize()
	hp, vp := s.GetHorizontalPadding(), s.GetVerticalPadding()
	help_h := lipgloss.Height(c.board.help.View(keys)) // TODO don't need board if board sends the appropriate WindowSize
	c.width, c.height = width/int(c.cnt)-hp, height-help_h-vp
	c.List.SetSize(c.width-hb, c.height-vb)
}

func (c *Column) getStyle() lipgloss.Style {
	if c.Focused() {
		return lipgloss.NewStyle().
			Padding(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Height(c.height).
			Width(c.width)
	}
	return lipgloss.NewStyle().
		Padding(1).
		Border(lipgloss.HiddenBorder()).
		Height(c.height).
		Width(c.width)
}

// MoveMsg can be handled by the lib user to update the status of their items.
type MoveMsg struct {
	direction int
	item      list.Item
}

// MoveToNext returns the new column index for the selected item.
func (c *Column) MoveToNext() tea.Cmd {
	// If nothing is selected, the SelectedItem will return Nil.
	item := c.List.SelectedItem()
	if item == nil {
		return nil
	}
	// move item
	c.List.RemoveItem(c.List.Index())

	// refresh list
	var cmd tea.Cmd
	c.List, cmd = c.List.Update(nil)

	return tea.Sequence(cmd, func() tea.Msg { return MoveMsg{+1, item} })
}

// MoveToPrev returns the new column index for the selected item.
func (c *Column) MoveToPrev() tea.Cmd {
	// If nothing is selected, the SelectedItem will return Nil.
	item := c.List.SelectedItem()
	if item == nil {
		return nil
	}
	// move item
	c.List.RemoveItem(c.List.Index())

	// refresh list
	var cmd tea.Cmd
	c.List, cmd = c.List.Update(nil)

	return tea.Sequence(cmd, func() tea.Msg { return MoveMsg{-1, item} })
}

// convert tasks to items for a list
func TasksToItems(tasks []ktask.Entry) []list.Item {
	var items []list.Item
	for _, t := range tasks {
		items = append(items, t)
	}
	return items
}

// convert tasks to items for a list
func ItemsToTasks(items []list.Item) []ktask.Entry {
	var tasks []ktask.Entry
	for _, i := range items {
		t := i.(ktask.Entry)
		tasks = append(tasks, t)
	}
	return tasks
}
