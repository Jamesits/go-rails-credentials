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

func executable(s ...string) string {
	return strings.Join(append([]string{os.Args[0]}, s...), " ")
}

type Cli struct {
	Edit Edit "cmd:\"\" help:\"Open the decrypted credentials in `$VISUAL` or `$EDITOR` for editing\""
	Show Show `cmd:"" help:"Show the decrypted credentials"`

	BaseDir                  string `name:"base-dir" default:"." type:"existingdir" help:"Root directory of your Rails project."`
	Environment              string `name:"environment" env:"RAILS_ENV"`
	MasterKey                string `name:"master-key" env:"RAILS_MASTER_KEY" help:"Your master key. For security, please do not provide this value by CLI argument; use the environment variable or a file instead."`
	MasterKeyFile            string `name:"master-key-file" help:"Path to your master.key file."`
	EncryptedCredentialsFile string `name:"credentials-file" help:"Path to your credential.yml.enc file."`

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
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("unable to read master key file %s: %w", cli.MasterKeyFile, err)
		}
	}
	// If master key file does not exist, generate a new one:
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
