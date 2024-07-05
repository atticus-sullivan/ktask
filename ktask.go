package main

import (
	"errors"
	"ktask/ktask"
	"ktask/ktask/kanban"
	"ktask/ktask/parser"
	"os"
	"path/filepath"
	"strings"

	arg "github.com/alexflint/go-arg"
	tea "github.com/charmbracelet/bubbletea"
	tf "github.com/jotaen/klog/klog/app/cli/terminalformat"
	gap "github.com/muesli/go-app-paths"
)

func initTaskDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return os.Mkdir(path, 0o770)
		}
		return err
	}
	return nil
}

// setupPath uses XDG to create the necessary data dirs for the program.
func setupPath() string {
	// get XDG paths
	scope := gap.NewScope(gap.User, "ktask")
	dirs, err := scope.DataDirs()
	if err != nil {
		panic(err)
	}
	// create the app base dir, if it doesn't exist
	var taskDir string
	if len(dirs) > 0 {
		taskDir = dirs[0]
	} else {
		taskDir, _ = os.UserHomeDir()
	}
	if err := initTaskDir(taskDir); err != nil {
		panic(err)
	}
	return taskDir
}

func exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func readData(path string) ([]ktask.Record, Error) {
	if exists(filepath.Join(path, ".tasks.ktask.lock")) {
		return nil, NewError(
			"Lock file exists",
			filepath.Join(path, ".tasks.ktask.lock"),
			errors.New("file exists"),
		)
	}
	source := filepath.Join(path, "tasks.ktask")
	content, err := os.ReadFile(source)
	if err != nil {
		return nil, NewErrorWithCode(
			NO_INPUT_ERROR,
			"Error reading file",
			"Location: "+source,
			err,
		)
	}
	f, err := os.Create(filepath.Join(path, ".tasks.ktask.lock"))
	if err == nil {
		f.Close()
	}

	records, _, errs := parser.NewSerialParser().Parse(string(content))
	if len(errs) > 0 {
		panic(NewParserErrors(errs))
	}
	return records, nil
}

func writeData(path string, data []ktask.Record) Error {
	var err error

	ser := NewSerialiser(tf.NewStyler(tf.COLOUR_THEME_NO_COLOUR), false)
	lines := parser.SerialiseRecords(ser, data...)

	content := strings.Builder{}
	for _, l := range lines {
		content.WriteString(l.Text)
		content.WriteRune('\n')
	}

	destination := filepath.Join(path, "tasks.ktask")
	err = os.WriteFile(destination, []byte(content.String()), 0777)
	if err != nil {
		return NewErrorWithCode(
			NO_INPUT_ERROR,
			"Error writing file",
			"Location: "+destination,
			err,
		)
	}
	os.Remove(filepath.Join(path, ".tasks.ktask.lock"))
	return nil
}

type rootCmd struct {
	Kanban *argKanban `arg:"subcommand:kanban"`
}

type argKanban struct{}

func main() {

	var args rootCmd
	arg.MustParse(&args)

	data, err := readData(setupPath())
	if err != nil {
		panic(err)
	}

	switch {
	case args.Kanban != nil:
		var cols []kanban.Column
		for i, r := range data {
			cols = append(cols, kanban.NewColumnFromRecord(r, i == 0))
		}
		board := kanban.NewDefaultBoard(cols)

		p := tea.NewProgram(board)
		rboard, err := p.Run()
		if err != nil {
			panic(err)
		}

		nboard, ok := rboard.(*kanban.Board)
		if !ok {
			panic("tea returned something else than a board")
		}

		data = nil
		for _, c := range nboard.Cols {
			r := ktask.NewRecord(ktask.Stage(c.List.Title))
			r.SetEntries(kanban.ItemsToTasks(c.List.Items()))
			data = append(data, r)
		}
	}

	err = writeData(setupPath(), data)
	if err != nil {
		panic(err)
	}
}
