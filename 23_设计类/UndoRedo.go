package main

import "errors"

var (
	ErrNothingToUndo = errors.New("nothing to undo")
	ErrNothingToRedo = errors.New("nothing to redo")
)

type Command interface {
	Apply() error
	Invert() error
}

// A simple undo/redo stack implementation.
type UndoRedo struct {
	pos  int
	cmds []Command
}

func NewUndoRedo() *UndoRedo {
	return &UndoRedo{pos: -1}
}

func (ur *UndoRedo) Do(cmd Command) error {
	// Clear the redo stack if we are adding a new command after undoing.
	if ur.pos < len(ur.cmds)-1 {
		ur.cmds = ur.cmds[:ur.pos+1]
	}

	err := cmd.Apply()
	if err != nil {
		return err
	}

	ur.cmds = append(ur.cmds, cmd)
	ur.pos++
	return nil
}

func (ur *UndoRedo) Undo() error {
	if !ur.CanUndo() {
		return ErrNothingToUndo
	}

	err := ur.cmds[ur.pos].Invert()
	if err != nil {
		return err
	}

	ur.pos--
	return nil
}

func (ur *UndoRedo) Redo() error {
	if !ur.CanRedo() {
		return ErrNothingToRedo
	}

	ur.pos++
	err := ur.cmds[ur.pos].Apply()
	if err != nil {
		return err
	}

	return nil
}

func (ur *UndoRedo) CanUndo() bool {
	return ur.pos >= 0
}

func (ur *UndoRedo) CanRedo() bool {
	return ur.pos < len(ur.cmds)-1
}

func (ur *UndoRedo) Clear() {
	ur.pos = -1
	ur.cmds = nil
}
