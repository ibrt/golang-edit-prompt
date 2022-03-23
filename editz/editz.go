package editz

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ibrt/golang-bites/filez"
	"github.com/ibrt/golang-errors/errorz"
	"github.com/ibrt/golang-shell/shellz"
)

// Editor opens a text editor.
type Editor interface {
	Edit(filePath string)
}

// ShellEditor describes the command and parameters to open a text editor.
type ShellEditor struct {
	Command string
	Params  []string
}

// Edit runs the editor.
func (e *ShellEditor) Edit(filePath string) {
	shellz.NewCommand(e.Command).
		AddParamsString(e.Params...).
		AddParamsString(filePath).
		SetLogf(nil).
		SetStdin(os.Stdin).
		SetStdout(os.Stdout).
		SetStderr(os.Stderr).
		MustRun()
}

var (
	// DefaultEditor is the default Editor.
	DefaultEditor Editor = getDefaultEditor()
)

func getDefaultEditor() *ShellEditor {
	if osEditor := os.Getenv("EDITOR"); osEditor != "" {
		parts := strings.Split(osEditor, " ")
		goodParts := make([]string, 0, len(parts))

		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				goodParts = append(goodParts, part)
			}
		}

		if len(goodParts) > 0 {
			return &ShellEditor{
				Command: goodParts[0],
				Params:  goodParts[1:],
			}
		}
	}

	return &ShellEditor{
		Command: "vi",
		Params:  []string{},
	}
}

// Edit implements a visudo-like text editing prompt.
//
// Sequence of actions:
//
// 1. Copy the file to a temporary location.
// 2. Open it using the default text editor (or a fallback editor, depending on the system).
// 3. Wait for the editor to exit.
// 4. If no changes were made, return immediately, otherwise validate the changes using the given callback.
// 5. If validation succeeds, overwrite the original file with the changed one, otherwise abort.
//
// The function returns the contents of the file after editing and a flag indicating whether there were changes.
func Edit(filePath string, validateFunc func([]byte) error) (contents []byte, isChanged bool, err error) {
	defer func() {
		if rErr := errorz.MaybeWrapRecover(recover(), errorz.SkipPackage()); rErr != nil {
			contents = nil
			isChanged = false
			err = errorz.Unwrap(rErr)
		}
	}()

	stat, err := os.Stat(filePath)
	errorz.MaybeMustWrap(err, errorz.SkipPackage())
	origBuf := filez.MustReadFile(filePath)

	filez.WithMustWriteTempFile("golang-edit-prompt", origBuf, func(tmpFilePath string) {
		DefaultEditor.Edit(tmpFilePath)
		newBuf := filez.MustReadFile(tmpFilePath)

		if bytes.Equal(origBuf, newBuf) {
			contents = newBuf
			return
		}

		errorz.MaybeMustWrap(validateFunc(newBuf), errorz.SkipPackage())
		errorz.MaybeMustWrap(ioutil.WriteFile(filePath, newBuf, stat.Mode()), errorz.SkipPackage())

		contents = newBuf
		isChanged = true
	})

	return
}
