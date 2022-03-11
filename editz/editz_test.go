package editz

import (
	"fmt"
	"os"
	"testing"

	"github.com/ibrt/golang-bites/filez"
	"github.com/ibrt/golang-errors/errorz"
	"github.com/stretchr/testify/require"
)

const (
	testContents        = "Hello!"
	testChangedContents = "Hello world!"
)

var (
	_ Editor = TestEditor(func(string) {})
)

// TestEditor implements Editor for test.
type TestEditor func(filePath string)

// Edit implements the Editor interface.
func (e TestEditor) Edit(filePath string) {
	e(filePath)
}

func TestEdit_Ok_Unchanged(t *testing.T) {
	filez.WithMustWriteTempFile("golang-edit-prompt", []byte(testContents), func(filePath string) {
		isEdited := false
		isValidated := false

		DefaultEditor = TestEditor(func(filePath string) {
			require.Equal(t, testContents, string(filez.MustReadFile(filePath)))
			isEdited = true
		})

		contents, isModified, err := Edit(filePath, func(contents []byte) error {
			require.Equal(t, testContents, string(contents))
			isValidated = true
			return nil
		})

		require.NoError(t, err)
		require.Equal(t, testContents, string(contents))
		require.Equal(t, testContents, string(filez.MustReadFile(filePath)))
		require.False(t, isModified)
		require.True(t, isEdited)
		require.False(t, isValidated)
	})
}

func TestEdit_Ok_Changed(t *testing.T) {
	filez.WithMustWriteTempFile("golang-edit-prompt", []byte(testContents), func(filePath string) {
		isEdited := false
		isValidated := false

		DefaultEditor = TestEditor(func(filePath string) {
			require.Equal(t, testContents, string(filez.MustReadFile(filePath)))
			filez.MustWriteFile(filePath, 0777, 0666, []byte(testChangedContents))
			isEdited = true
		})

		contents, isModified, err := Edit(filePath, func(contents []byte) error {
			require.Equal(t, testChangedContents, string(contents))
			isValidated = true
			return nil
		})

		require.NoError(t, err)
		require.Equal(t, testChangedContents, string(contents))
		require.Equal(t, testChangedContents, string(filez.MustReadFile(filePath)))
		require.True(t, isModified)
		require.True(t, isEdited)
		require.True(t, isValidated)
	})
}

func TestEdit_Err_Invalid(t *testing.T) {
	filez.WithMustWriteTempFile("golang-edit-prompt", []byte(testContents), func(filePath string) {
		isEdited := false
		isValidated := false

		DefaultEditor = TestEditor(func(filePath string) {
			require.Equal(t, testContents, string(filez.MustReadFile(filePath)))
			filez.MustWriteFile(filePath, 0777, 0666, []byte(testChangedContents))
			isEdited = true
		})

		contents, isModified, err := Edit(filePath, func(contents []byte) error {
			require.Equal(t, testChangedContents, string(contents))
			isValidated = true
			return fmt.Errorf("invalid")
		})

		require.EqualError(t, err, "invalid")
		require.Nil(t, contents)
		require.Equal(t, testContents, string(filez.MustReadFile(filePath)))
		require.False(t, isModified)
		require.True(t, isEdited)
		require.True(t, isValidated)
	})
}

func TestEdit_Err_Panic(t *testing.T) {
	filez.WithMustWriteTempFile("golang-edit-prompt", []byte(testContents), func(filePath string) {
		isEdited := false
		isValidated := false

		DefaultEditor = TestEditor(func(filePath string) {
			require.Equal(t, testContents, string(filez.MustReadFile(filePath)))
			isEdited = true
			errorz.MustErrorf("edit error")
		})

		contents, isModified, err := Edit(filePath, func(contents []byte) error {
			require.Equal(t, testContents, string(contents))
			isValidated = true
			return nil
		})

		require.EqualError(t, err, "edit error")
		require.Nil(t, contents)
		require.Equal(t, testContents, string(filez.MustReadFile(filePath)))
		require.False(t, isModified)
		require.True(t, isEdited)
		require.False(t, isValidated)
	})
}

func TestShellEditor(t *testing.T) {
	filez.WithMustWriteTempFile("golang-edit-prompt", []byte(testContents), func(filePath string) {
		e := &ShellEditor{
			Command: "/usr/bin/env",
			Params: []string{
				"sed",
				"-i",
				".bak",
				fmt.Sprintf("s/%v/%v/", testContents, testChangedContents),
			},
		}
		e.Edit(filePath)
		require.Equal(t, testChangedContents, string(filez.MustReadFile(filePath)))
	})
}

func TestGetDefaultEditor(t *testing.T) {
	require.NoError(t, os.Setenv("EDITOR", "edit"))
	require.Equal(t, &ShellEditor{
		Command: "edit",
		Params:  []string{},
	}, getDefaultEditor())

	require.NoError(t, os.Setenv("EDITOR", "   edit  -p  -a   "))
	require.Equal(t, &ShellEditor{
		Command: "edit",
		Params:  []string{"-p", "-a"},
	}, getDefaultEditor())

	require.NoError(t, os.Unsetenv("EDITOR"))
	require.Equal(t, &ShellEditor{
		Command: "vi",
		Params:  []string{},
	}, getDefaultEditor())
}
