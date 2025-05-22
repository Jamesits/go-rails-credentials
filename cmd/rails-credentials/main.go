// Package main tries to mock the behavior of rails credentials:*
// https://github.com/rails/rails/blob/main/railties/lib/rails/commands/credentials/credentials_command.rb
package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/jamesits/go-rails-credentials/pkg/credentials"
	"os"
	"path/filepath"
	"strings"
)

type Cli struct {
	Edit Edit "cmd:\"\" help:\"Open the decrypted credentials in `$VISUAL` or `$EDITOR` for editing\""
	Show Show `cmd:"" help:"Show the decrypted credentials"`

	BaseDir                  string `name:"base-dir" default:"." type:"existingdir"`
	Environment              string `name:"environment" env:"RAILS_ENV"`
	MasterKey                string `name:"master-key" env:"RAILS_MASTER_KEY" help:"your master key; it's not recommended to provide this argument from CLI.'"`
	MasterKeyFile            string `name:"master-key-file"`
	EncryptedCredentialsFile string `name:"credentials"`

	masterKeyGenerated bool
}

func (cli *Cli) AfterApply() error {
	err := os.Chdir(cli.BaseDir)
	if err != nil {
		return fmt.Errorf("unable to open base directory: %w", err)
	}

	// parse RAILS_ENV
	if cli.Environment == "" {
		if cli.MasterKeyFile == "" {
			cli.MasterKeyFile = filepath.Join("config", "master.key")
		}
		if cli.EncryptedCredentialsFile == "" {
			cli.EncryptedCredentialsFile = filepath.Join("config", "credentials.yml.enc")
		}
	} else {
		if cli.MasterKeyFile == "" {
			cli.MasterKeyFile = filepath.Join("config", "credentials", fmt.Sprintf("%s.key", cli.Environment))
		}
		if cli.EncryptedCredentialsFile == "" {
			cli.EncryptedCredentialsFile = filepath.Join("config", "credentials", fmt.Sprintf("%s.yml.enc", cli.Environment))
		}
	}

	// If RAILS_MASTER_KEY environment variable is set, we use it instead of the file content.
	// Otherwise, try read an existing master key.
	if cli.MasterKey == "" {
		m, err := os.ReadFile(cli.MasterKeyFile)
		if err == nil {
			cli.MasterKey = strings.Trim(string(m), "\r\n")
		}
	}
	// If for some reason we are unable to read the existing credentials file, just generate a new one:
	if cli.MasterKey == "" {
		cli.MasterKey, err = credentials.RandomMasterKey()
		cli.masterKeyGenerated = true
		if err != nil {
			return fmt.Errorf("unable to generate a master key: %w", err)
		}
	}

	return nil
}

func main() {
	cli := &Cli{}
	ctx := kong.Parse(cli)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
