package web

import (
	"fmt"
	"path/filepath"

	"github.com/prysmaticlabs/prysm/v5/api"
	"github.com/prysmaticlabs/prysm/v5/cmd"
	"github.com/prysmaticlabs/prysm/v5/cmd/validator/flags"
	"github.com/prysmaticlabs/prysm/v5/config/features"
	"github.com/prysmaticlabs/prysm/v5/runtime/tos"
	"github.com/prysmaticlabs/prysm/v5/validator/rpc"
	"github.com/urfave/cli/v2"
)

// Commands for managing Prysm validator accounts.
var Commands = &cli.Command{
	Name:     "web",
	Category: "web",
	Usage:    "Defines commands for interacting with the Prysm web interface.",
	Subcommands: []*cli.Command{
		{
			Name:        "generate-auth-token",
			Description: `Generate an authentication token for the Prysm web interface`,
			Flags: cmd.WrapFlags([]cli.Flag{
				flags.WalletDirFlag,
				flags.GRPCGatewayHost,
				flags.GRPCGatewayPort,
				flags.AuthTokenPathFlag,
				cmd.AcceptTosFlag,
			}),
			Before: func(cliCtx *cli.Context) error {
				if err := cmd.LoadFlagsFromConfig(cliCtx, cliCtx.Command.Flags); err != nil {
					return err
				}
				return tos.VerifyTosAcceptedOrPrompt(cliCtx)
			},
			Action: func(cliCtx *cli.Context) error {
				if err := features.ConfigureValidator(cliCtx); err != nil {
					return err
				}
				walletDirPath := cliCtx.String(flags.WalletDirFlag.Name)
				if walletDirPath == "" {
					log.Fatal("--wallet-dir not specified")
				}
				gatewayHost := cliCtx.String(flags.GRPCGatewayHost.Name)
				gatewayPort := cliCtx.Int(flags.GRPCGatewayPort.Name)
				validatorWebAddr := fmt.Sprintf("%s:%d", gatewayHost, gatewayPort)
				authTokenPath := filepath.Join(walletDirPath, api.AuthTokenFileName)
				tempAuthTokenPath := cliCtx.String(flags.AuthTokenPathFlag.Name)
				if tempAuthTokenPath != "" {
					authTokenPath = tempAuthTokenPath
				}
				if err := rpc.CreateAuthToken(authTokenPath, validatorWebAddr); err != nil {
					log.WithError(err).Fatal("Could not create web auth token")
				}
				return nil
			},
		},
	},
}
