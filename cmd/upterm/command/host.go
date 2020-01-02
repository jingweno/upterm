package command

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/shlex"
	"github.com/jingweno/upterm/host"
	"github.com/rs/xid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var (
	flagServer         string
	flagForceCommand   string
	flagPrivateKeys    []string
	flagAuthorizedKeys string
)

func hostCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "host",
		Short: "Host a terminal session",
		Long:  "Host a terminal session via a reverse SSH tunnel to the upterm server. By default, the command authenticates against the upterm server using the private keys located at `~/.ssh/id_dsa`, `~/.ssh/id_ecdsa`, `~/.ssh/id_ed25519`, and `~/.ssh/id_rsa`. The host can permit a list of client public keys by specifying an authorized_keys file. By default, the input/output of the host attaches to the input/output of the client's. The host can force the execution of a command after the client joins, and attach the input/output of this command to the client's.",
		Example: `  # Host a session by running $SHELL.
  # The client's input/output is attached to the host's.
  upterm host

  # Host a session with a custom command.
  # The client's input/output is attached to the host's.
  upterm host -- docker run --rm -ti ubuntu bash

  # Host a session by running 'tmux new -t pair-programming'.
  # The host runs 'tmux attach -t pair-programming' after the client joins the session.
  # The client's input/output is attached to this command's.
  upterm host --force-command 'tmux attach -t pair-programming' -- tmux new -t pair-programming`,
		PreRunE: validateShareRequiredFlags,
		RunE:    shareRunE,
	}

	cmd.PersistentFlags().StringVarP(&flagServer, "server", "", "uptermd.upterm.dev:22", "upterm server address (required)")
	cmd.PersistentFlags().StringVarP(&flagForceCommand, "force-command", "f", "", "force execution of a command and attach its input/output to client's.")
	cmd.PersistentFlags().StringSliceVarP(&flagPrivateKeys, "private-key", "i", nil, "private key for public key authentication against the upterm server (required).")
	cmd.PersistentFlags().StringVarP(&flagAuthorizedKeys, "authorized-key", "a", "", "an authorized_keys file that lists public keys that are permitted to connect.")

	return cmd
}

func validateShareRequiredFlags(c *cobra.Command, args []string) error {
	missingFlagNames := []string{}
	if flagServer == "" {
		missingFlagNames = append(missingFlagNames, "server")
	}

	if len(flagPrivateKeys) == 0 {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		flagPrivateKeys = defaultPrivateKeys(homeDir)
		if len(flagPrivateKeys) == 0 {
			missingFlagNames = append(missingFlagNames, "private-key")

		}
	}
	if len(missingFlagNames) > 0 {
		return fmt.Errorf(`required flag(s) "%s" not set`, strings.Join(missingFlagNames, ", "))
	}

	return nil
}

func shareRunE(c *cobra.Command, args []string) error {
	var err error
	if len(args) == 0 {
		args, err = shlex.Split(os.Getenv("SHELL"))
		if err != nil {
			return err
		}
	}

	var forceCommand []string
	if flagForceCommand != "" {
		forceCommand, err = shlex.Split(flagForceCommand)
		if err != nil {
			return fmt.Errorf("error parsing command %s: %w", flagForceCommand, err)
		}
	}

	var authorizedKeys []ssh.PublicKey
	if flagAuthorizedKeys != "" {
		authorizedKeys, err = host.AuthorizedKeys(flagAuthorizedKeys)
	}
	if err != nil {
		return fmt.Errorf("error reading authorized keys: %w", err)
	}

	signers, cleanup, err := host.Signers(flagPrivateKeys)
	if err != nil {
		return fmt.Errorf("error reading private keys: %w", err)
	}
	if cleanup != nil {
		defer cleanup()
	}

	h := &host.Host{
		Host:           flagServer,
		SessionID:      xid.New().String(),
		Command:        args,
		ForceCommand:   forceCommand,
		Signers:        signers,
		AuthorizedKeys: authorizedKeys,
		KeepAlive:      time.Duration(30),
		Logger:         log.New(),
	}

	return h.Run(context.Background())
}

func defaultPrivateKeys(homeDir string) []string {
	var pks []string
	for _, f := range []string{
		"id_dsa",
		"id_ecdsa",
		"id_ed25519",
		"id_rsa",
	} {
		pk := filepath.Join(homeDir, ".ssh", f)
		if _, err := os.Stat(pk); os.IsNotExist(err) {
			continue
		}

		pks = append(pks, pk)
	}

	return pks
}