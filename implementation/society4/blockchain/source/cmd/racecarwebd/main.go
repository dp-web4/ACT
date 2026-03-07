package main

import (
	"fmt"
	"io"
	"os"

	"cosmossdk.io/log"
	confixcmd "cosmossdk.io/tools/confix/cmd"
	cmtcfg "github.com/cometbft/cometbft/config"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/pruning"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/snapshot"
	"github.com/cosmos/cosmos-sdk/server"
	serverconfig "github.com/cosmos/cosmos-sdk/server/config"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/types/module"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cobra"

	"society4chain/app"
)

func main() {
	rootCmd := newRootCmd()
	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}

func newRootCmd() *cobra.Command {
	// Create a temporary app to get module basics and codecs
	tempApp := app.New(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		simtestutil.NewAppOptionsWithFlagHome(app.DefaultNodeHome),
	)

	// Build BasicManager from the module manager
	basicManager := module.BasicManager{}
	for name, mod := range tempApp.ModuleManager.Modules {
		if basicMod, ok := mod.(module.AppModuleBasic); ok {
			basicManager[name] = basicMod
		}
	}
	// Override genutil with proper AppModuleBasic (depinject registers GenesisOnlyAppModule)
	basicManager[genutiltypes.ModuleName] = genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator)

	txConfig := tempApp.TxConfig()

	rootCmd := &cobra.Command{
		Use:   "racecarwebd",
		Short: "Society4 Chain daemon (ACT blockchain)",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			initClientCtx := client.Context{}.
				WithCodec(tempApp.AppCodec()).
				WithInterfaceRegistry(tempApp.InterfaceRegistry()).
				WithTxConfig(txConfig).
				WithLegacyAmino(tempApp.LegacyAmino()).
				WithInput(os.Stdin).
				WithAccountRetriever(authtypes.AccountRetriever{}).
				WithHomeDir(app.DefaultNodeHome).
				WithViper("")

			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			return server.InterceptConfigsPreRunHandler(
				cmd,
				serverconfig.DefaultConfigTemplate,
				serverconfig.DefaultConfig(),
				cmtcfg.DefaultConfig(),
			)
		},
	}

	// Add commands
	rootCmd.AddCommand(
		genutilcli.InitCmd(basicManager, app.DefaultNodeHome),
		debug.Cmd(),
		confixcmd.MigrateCommand(),
		pruning.Cmd(newApp, app.DefaultNodeHome),
		snapshot.Cmd(newApp),
	)

	server.AddCommands(rootCmd, app.DefaultNodeHome, newApp, appExport, addModuleInitFlags)

	rootCmd.AddCommand(
		server.StatusCommand(),
		genesisCommand(txConfig, basicManager),
		queryCommand(),
		txCommand(),
		keys.Commands(),
	)

	return rootCmd
}

func addModuleInitFlags(startCmd *cobra.Command) {}

func genesisCommand(txConfig client.TxConfig, mbm module.BasicManager, cmds ...*cobra.Command) *cobra.Command {
	cmd := genutilcli.Commands(txConfig, mbm, app.DefaultNodeHome)
	for _, subCmd := range cmds {
		cmd.AddCommand(subCmd)
	}
	return cmd
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		rpc.QueryEventForTxCmd(),
		server.QueryBlockCmd(),
		authcmd.QueryTxsByEventsCmd(),
		server.QueryBlocksCmd(),
		authcmd.QueryTxCmd(),
		server.QueryBlockResultsCmd(),
	)
	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         false,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetMultiSignBatchCmd(),
		authcmd.GetValidateSignaturesCommand(),
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
		authcmd.GetSimulateCmd(),
	)
	return cmd
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer, appOpts servertypes.AppOptions) servertypes.Application {
	baseappOptions := server.DefaultBaseappOptions(appOpts)
	return app.New(logger, db, traceStore, true, appOpts, baseappOptions...)
}

func appExport(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	height int64,
	forZeroHeight bool,
	jailAllowedAddrs []string,
	appOpts servertypes.AppOptions,
	modulesToExport []string,
) (servertypes.ExportedApp, error) {
	homePath, ok := appOpts.Get(flags.FlagHome).(string)
	if !ok || homePath == "" {
		return servertypes.ExportedApp{}, fmt.Errorf("application home not set")
	}

	a := app.New(logger, db, traceStore, height == -1, appOpts)

	if height != -1 {
		if err := a.LoadHeight(height); err != nil {
			return servertypes.ExportedApp{}, err
		}
	}

	return a.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs, modulesToExport)
}
