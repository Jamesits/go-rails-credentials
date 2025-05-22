package main

import (
	"errors"
	"fmt"
	"github.com/jamesits/go-rails-credentials/pkg/credentials"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Edit struct {
	EditorCommand string `name:"editor" env:"VISUAL,EDITOR" default:"vi" help:"Your editor program."`
}

const (
	masterKeyCreateTemplate = `Adding config/master.key to store the encryption key: %s

Save this in a password manager your team can access.

If you lose the key, no one, including you, can access anything encrypted with it.

      create  %s
`
	editorStartTemplate   = "Editing %s...\n"
	decryptFailedTemplate = "Couldn't decrypt %s. Perhaps you passed the wrong key?\n"
	savedTemplate         = "File encrypted and saved.\n"
)

func (cmd *Edit) Run(cli *Cli) error {
	var err error

	// if creation of a new master key is needed
	if cli.masterKeyGenerated {
		_, _ = fmt.Fprintf(os.Stderr, masterKeyCreateTemplate, cli.MasterKey, cli.MasterKeyFile)
		err = atomicWrite(cli.MasterKeyFile, []byte(cli.MasterKey), 0o600, 0o777)
		if err != nil {
			return fmt.Errorf("write master key file failed: %w", err)
		}
	}

	var rawCredentialsFileContent string
	// read and decrypt the file
	e, err := os.ReadFile(cli.EncryptedCredentialsFile)
	if err == nil {
		obj, err := credentials.Decrypt(cli.MasterKey, string(e))
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, decryptFailedTemplate, cli.EncryptedCredentialsFile)
			return fmt.Errorf("decrypt failed: %w", err)
		}
		rawCredentialsFileContent, err = credentials.UnmarshalSingleString(obj)
		if err != nil {
			return fmt.Errorf("unmarshal failed: %w", err)
		}
	} else if errors.Is(err, os.ErrNotExist) {
		rawCredentialsFileContent, err = credentials.NewCredentialsFileContent()
		if err != nil {
			return fmt.Errorf("render credentials.yml template failed: %w", err)
		}
	}

	// write temp file
	editorTempFile, err := os.CreateTemp("", "*-credentials.yml")
	if err != nil {
		return fmt.Errorf("unable to create temporary file for editing: %w", err)
	}
	err = editorTempFile.Close()
	if err != nil {
		return fmt.Errorf("unable to close temporary file for editing: %w", err)
	}
	defer os.Remove(editorTempFile.Name())

	err = os.WriteFile(editorTempFile.Name(), []byte(rawCredentialsFileContent), 0o600)
	if err != nil {
		return fmt.Errorf("unable to write temporary file for editing: %w", err)
	}

	// start the editor
	_, _ = fmt.Fprintf(os.Stderr, editorStartTemplate, cli.EncryptedCredentialsFile)
	editorCommandArgs := strings.Fields(cmd.EditorCommand)
	editorCommandPath, err := exec.LookPath(editorCommandArgs[0])
	if err != nil {
		return fmt.Errorf("unable to find editor executable %s: %w", editorCommandPath, err)
	}
	editorCmd := exec.Cmd{
		Path:   editorCommandPath,
		Args:   append(editorCommandArgs, editorTempFile.Name()),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	err = editorCmd.Run()
	if err != nil || !editorCmd.ProcessState.Success() {
		return fmt.Errorf("editor failed with code %d: %w", editorCmd.ProcessState.ExitCode(), err)
	}

	// read back
	newRawCredentialsFileContent, err := os.ReadFile(editorTempFile.Name())
	if err != nil {
		return fmt.Errorf("unable to read temporary file: %w", err)
	}
	if string(newRawCredentialsFileContent) == rawCredentialsFileContent {
		return nil
	}

	// encrypt the file
	newObject, err := credentials.MarshalSingleString(string(newRawCredentialsFileContent))
	if err != nil {
		return fmt.Errorf("unable to marshal object: %w", err)
	}
	newEncryptedCredentialsFileContent, err := credentials.Encrypt(cli.MasterKey, newObject)
	if err != nil {
		return fmt.Errorf("unable to encrypt: %w", err)
	}

	err = atomicWrite(cli.EncryptedCredentialsFile, []byte(newEncryptedCredentialsFileContent), 0o666, 0o777)
	if err != nil {
		return fmt.Errorf("unable to save encrypted file: %w", err)
	}
	_, _ = fmt.Fprintf(os.Stderr, savedTemplate)
	return nil
}

func atomicWrite(path string, content []byte, filePerm os.FileMode, dirPerm os.FileMode) error {
	var err error

	temp := path + ".tmp"
	_ = os.Remove(temp)

	err = os.MkdirAll(filepath.Dir(path), dirPerm)
	if err != nil {
		return fmt.Errorf("unable to create directory: %w", err)
	}

	err = os.WriteFile(temp, content, filePerm)
	if err != nil {
		return fmt.Errorf("unable to write temporary file: %w", err)
	}
	defer os.Remove(temp)

	err = os.Rename(temp, path)
	if err != nil {
		return fmt.Errorf("unable to overwrite destination file: %w", err)
	}

	return nil
}
