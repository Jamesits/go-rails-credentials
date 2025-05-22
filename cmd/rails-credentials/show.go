package main

import (
	"fmt"
	"github.com/jamesits/go-rails-credentials/pkg/credentials"
	"os"
)

type Show struct{}

func (cmd *Show) Run(cli *Cli) error {
	var err error

	var rawObject []byte

	e, err := os.ReadFile(cli.EncryptedCredentialsFile)
	if err != nil {
		return fmt.Errorf("read encrypted file failed: %w", err)
	}

	rawObject, err = credentials.Decrypt(cli.MasterKey, string(e))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, decryptFailedTemplate, cli.EncryptedCredentialsFile)
		return fmt.Errorf("decrypt failed: %w", err)
	}

	rawString, err := credentials.UnmarshalSingleString(rawObject)
	if err != nil {
		return fmt.Errorf("unmarshal failed: %w", err)
	}

	_, _ = fmt.Fprintf(os.Stdout, rawString)
	return nil
}
