package cli_options

import (
	"errors"

	"github.com/gravitl/netmaker/netclient/command"
	"github.com/gravitl/netmaker/netclient/config"
	"github.com/urfave/cli/v2"
)

// GetCommands - return commands that CLI uses
func GetCommands(cliFlags []cli.Flag) []*cli.Command {
	return []*cli.Command{
		{
			Name:  "join",
			Usage: "Join a Netmaker network.",
			Flags: cliFlags,
			Action: func(c *cli.Context) error {
				cfg, pvtKey, err := config.GetCLIConfig(c)
				if err != nil {
					return err
				}
				if cfg.Network == "all" {
					err = errors.New("no network provided")
					return err
				}
				if cfg.Server.GRPCAddress == "" {
					err = errors.New("no server address provided")
					return err
				}
				err = command.Join(&cfg, pvtKey)
				return err
			},
		},
		{
			Name:  "leave",
			Usage: "Leave a Netmaker network.",
			Flags: cliFlags,
			// the action, or code that will be executed when
			// we execute our `ns` command
			Action: func(c *cli.Context) error {
				cfg, _, err := config.GetCLIConfig(c)
				if err != nil {
					return err
				}
				err = command.Leave(&cfg, c.String("force") == "yes")
				return err
			},
		},
		{
			Name:  "pull",
			Usage: "Pull latest configuration and peers from server.",
			Flags: cliFlags,
			// the action, or code that will be executed when
			// we execute our `ns` command
			Action: func(c *cli.Context) error {
				cfg, _, err := config.GetCLIConfig(c)
				if err != nil {
					return err
				}
				err = command.Pull(&cfg)
				return err
			},
		},
		{
			Name:  "list",
			Usage: "Get list of networks.",
			Flags: cliFlags,
			// the action, or code that will be executed when
			// we execute our `ns` command
			Action: func(c *cli.Context) error {
				cfg, _, err := config.GetCLIConfig(c)
				if err != nil {
					return err
				}
				err = command.List(cfg)
				return err
			},
		},
		{
			Name:  "uninstall",
			Usage: "Uninstall the netclient system service.",
			Flags: cliFlags,
			// the action, or code that will be executed when
			// we execute our `ns` command
			Action: func(c *cli.Context) error {
				err := command.Uninstall()
				return err
			},
		},
		{
			Name:  "daemon",
			Usage: "run netclient as daemon",
			Flags: cliFlags,
			Action: func(c *cli.Context) error {
				err := command.Daemon()
				return err
			},
		},
	}
}
