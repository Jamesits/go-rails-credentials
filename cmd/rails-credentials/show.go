package main

import (
	"fmt"
	"github.com/jamesits/go-rails-credentials/pkg/credentials"
	"os"
)

const (
	missingKeyMessageTemplate         = "Missing '%s' to decrypt credentials. See `%s`."
	missingCredentialsMessageTemplate = "File '%s' does not exist. Use \"%s` to change that.\n"
)

type Show struct{}

func (cmd *Show) Run(cli *Cli) error {
	var err error

	var rawObject []byte

	e, err := os.ReadFile(cli.EncryptedCredentialsFile)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, missingCredentialsMessageTemplate, cli.EncryptedCredentialsFile, executable("edit"))
		return fmt.Errorf("read encrypted file failed: %w", err)
	}

	rawObject, err = credentials.Decrypt(cli.MasterKey, string(e))
	if err != nil {
		if cli.masterKeyGenerated {
			_, _ = fmt.Fprintf(os.Stderr, missingKeyMessageTemplate, cli.MasterKeyFile, executable("--help"))
		} else {
			_, _ = fmt.Fprintf(os.Stderr, decryptFailedTemplate, cli.EncryptedCredentialsFile)
		}
		return fmt.Errorf("decrypt failed: %w", err)
	}

	rawString, err := credentials.UnmarshalSingleString(rawObject)
	if err != nil {
		return fmt.Errorf("unmarshal failed: %w", err)
	}

	_, _ = fmt.Fprint(os.Stdout, rawString)
	return nil
}
