package main

import (
	"errors"
	"fmt"
	"ktask/ktask"
	"ktask/ktask/kanban"
	"ktask/ktask/parser"
	"os"
	"path/filepath"
	"slices"
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

func readData(source string) ([]ktask.Record, ktask.Error) {
	lock := source + ".lock"
	if exists(lock) {
		return nil, ktask.NewError(
			"Lock file exists",
			lock,
			errors.New("file exists"),
		)
	}
	content, err := os.ReadFile(source)
	if err != nil {
		return nil, ktask.NewErrorWithCode(
			ktask.NO_INPUT_ERROR,
			"Error reading file",
			"Location: "+source,
			err,
		)
	}
	f, err := os.Create(lock)
	if err == nil {
		f.Close()
	}

	records, _, errs := parser.NewSerialParser().Parse(string(content))
	if errs != nil {
		return nil, ktask.NewParserErrors(errs)
	}
	return records, nil
}

func writeData(destination string, data []ktask.Record) error {
	var err error
	lock := destination + ".lock"

	ser := parser.NewSerialiser(tf.NewStyler(tf.COLOUR_THEME_NO_COLOUR), false)
	lines := parser.SerialiseRecords(ser, data...)

	content := strings.Builder{}
	for _, l := range lines {
		content.WriteString(l.Text)
		content.WriteRune('\n')
	}

	if exists(destination) {
		err = os.Rename(destination, destination+".bak")
		if err != nil {
			return fmt.Errorf("making backup before writing failed: %w", err)
		}
	}

	err = os.WriteFile(destination, []byte(content.String()), 0777)
	if err != nil {
		return fmt.Errorf("writing output file failed. If needed you can find a backup of the original file located at the same place suffixed with .bak\n%w", err)
	}
	return os.Remove(lock)
}

type rootCmd struct {
	Kanban *argKanban `arg:"subcommand:kanban"`
}

type argKanban struct {
	File   string   `arg:"positional" help:"specify the file that should be read from / written to"`
	Tags   []string `arg:"--tags,-t,separate" help:"if set, only entries with this/these tags will be shown, may be specified multiple times"`
	NoTags []string `arg:"--no-tags,-T,separate" help:"if set, entries with this/these tags will NOT be shown, may be specified multiple times"`
}

func main() {
	var args rootCmd
	var err error
	var errK ktask.Error
	arg.MustParse(&args)

	switch {
	case args.Kanban != nil:
		var data []ktask.Record
		path := args.Kanban.File
		if path == "" {
			path = filepath.Join(setupPath(), "tasks.ktask")
		}
		data, errK = readData(path)
		if errK != nil {
			switch errK := errK.(type) {
			case ktask.ParserErrors:
				panic(ktask.PrettifyParsingError(errK, tf.NewStyler(tf.COLOUR_THEME_NO_COLOUR)))
			case ktask.Error:
				panic(ktask.PrettifyAppError(errK, false))
			default:
				if err != nil {
					panic(err)
				}
			}
		}

		var data_shown []ktask.Record
		var data_hidden []ktask.Record
		if len(args.Kanban.Tags) == 0 && len(args.Kanban.NoTags) == 0 {
			data_shown = data
		} else {
			for _, i := range data {
				r1, r2 := i.SplitOnFunc(func(e *ktask.Entry) bool {
					return (slices.ContainsFunc(args.Kanban.Tags, func(s string) bool {
						t, _ := ktask.NewTagFromString(s)
						return e.Name().Tags().Contains(t)
					}) || len(args.Kanban.Tags) == 0) && !slices.ContainsFunc(args.Kanban.NoTags, func(s string) bool {
						t, _ := ktask.NewTagFromString(s)
						return e.Name().Tags().Contains(t)
					})
				})
				data_shown = append(data_shown, r1)
				data_hidden = append(data_hidden, r2)
			}
		}

		var cols []kanban.Column
		for i, r := range data_shown {
			cols = append(cols, kanban.NewColumnFromRecord(r, i == 0))
		}
		board := kanban.NewDefaultBoard(cols)

		var rboard tea.Model
		p := tea.NewProgram(board)
		rboard, err = p.Run()
		if err != nil {
			panic(err)
		}

		nboard, ok := rboard.(*kanban.Board)
		if !ok {
			panic("tea returned something else than a board")
		}

		data = nil
		for i, c := range nboard.Cols {
			r := ktask.NewRecord(ktask.Stage(c.List.Title))
			r.SetEntries(kanban.ItemsToTasks(c.List.Items()))
			if i < len(data_hidden) {
				r.Merge(data_hidden[i])
			}
			data = append(data, r)
		}

		err = writeData(path, data)
		if err != nil {
			panic(err)
		}
	}
}
