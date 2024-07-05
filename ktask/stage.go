package ktask

import "errors"

type Stage string

const (
	Todo       Stage = "todo"
	InProgress Stage = "in progress"
	Done       Stage = "done"
)

func (s *Stage) Valid() error {
	if *s != Todo && *s != InProgress && *s != Done {
		return errors.New("Invalid stage provided")
	}
	return nil
}
