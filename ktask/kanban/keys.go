package kanban

import "github.com/charmbracelet/bubbles/key"

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the help.KeyMap interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// help.KeyMap interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},    // first column
		{k.Left, k.Right}, // second column
		{k.New, k.Delete}, // third column
		{k.Edit},          // third column
		{k.Next, k.Prev},  // third column
		{k.Help, k.Quit},  // fourth column
	}
}

type keyMap struct {
	New    key.Binding
	Edit   key.Binding
	Delete key.Binding
	Up     key.Binding
	Down   key.Binding
	Right  key.Binding
	Left   key.Binding
	Next   key.Binding
	Prev   key.Binding
	Enter  key.Binding
	Help   key.Binding
	Quit   key.Binding
	Back   key.Binding
	Esc    key.Binding
}

var keys = keyMap{
	New: key.NewBinding(
		key.WithKeys("n", "a"),
		key.WithHelp("n/a", "new"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d", "x"),
		key.WithHelp("d/x", "delete"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l", "tab"),
		key.WithHelp("→/l", "move right"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h", "shift+tab"),
		key.WithHelp("←/l", "move left"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "enter"),
	),
	Next: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "next stage"),
	),
	Prev: key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("backspace", "prev stage"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "quit"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "back"),
	),
}
