package cli

import (
    "context"
    "github.com/cosmos/cosmos-sdk/client"
    "github.com/cosmos/cosmos-sdk/client/flags"
    "github.com/spf13/cobra"
    "racecar-web/x/componentregistry/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
    cmd := &cobra.Command{
        Use:   types.ModuleName,
        Short: fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
        RunE:  client.ValidateCmd,
    }

    cmd.AddCommand(CmdQueryLCT())

    return cmd
}

// Query formatting functions
