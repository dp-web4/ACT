package cli

import (
    "github.com/cosmos/cosmos-sdk/client"
    "github.com/cosmos/cosmos-sdk/client/flags"
    "github.com/cosmos/cosmos-sdk/client/tx"
    "github.com/spf13/cobra"
    "github.com/dp-web4/act/x/trusttensor/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   types.ModuleName,
        Short: fmt.Sprintf("%s transactions subcommands", types.ModuleName),
        RunE:  client.ValidateCmd,
    }

    cmd.AddCommand(CmdCreateLCT())

    return cmd
}

// Flag parsing functions
