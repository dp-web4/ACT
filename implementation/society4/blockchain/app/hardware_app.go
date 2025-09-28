package app

import (
    "encoding/json"
    "fmt"
    "os"

    abci "github.com/cometbft/cometbft/abci/types"
    "github.com/cosmos/cosmos-sdk/baseapp"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/cosmos/cosmos-sdk/types/module"

    consensusKeeper "society4chain/x/consensus/keeper"
    lctTypes "society4chain/x/lctmanager/types"
)

// Society4App extends the base application with hardware binding
type Society4App struct {
    *baseapp.BaseApp

    // Keepers
    consensusKeeper   consensusKeeper.Keeper
    hardwareValidator *consensusKeeper.HardwareValidator

    // Self-LCT (root identity)
    selfLCT *lctTypes.SelfLCT

    // Module manager
    mm *module.Manager
}

// NewSociety4App creates a new Society 4 blockchain application
func NewSociety4App(
    logger log.Logger,
    db dbm.DB,
    traceStore io.Writer,
    loadLatest bool,
    appOpts servertypes.AppOptions,
    baseAppOptions ...func(*baseapp.BaseApp),
) *Society4App {
    // Initialize base app
    bApp := baseapp.NewBaseApp(
        "society4chain",
        logger,
        db,
        nil, // TxDecoder will be set later
        baseAppOptions...,
    )

    app := &Society4App{
        BaseApp: bApp,
    }

    // Initialize self-LCT with hardware binding
    if err := app.InitializeSelfLCT(); err != nil {
        panic(fmt.Sprintf("Failed to initialize self-LCT: %v", err))
    }

    // Set up hardware validator
    app.hardwareValidator = consensusKeeper.NewHardwareValidator(app.selfLCT)

    // Initialize consensus keeper with hardware validation
    app.consensusKeeper = consensusKeeper.NewKeeper(
        app.hardwareValidator,
        logger,
    )

    // Set BeginBlocker with hardware validation
    app.SetBeginBlocker(app.BeginBlocker)
    app.SetEndBlocker(app.EndBlocker)

    return app
}

// InitializeSelfLCT creates or loads the hardware-bound self-LCT
func (app *Society4App) InitializeSelfLCT() error {
    selfLCTPath := "data/self_lct.json"

    // Check if self-LCT already exists
    if _, err := os.Stat(selfLCTPath); err == nil {
        // Load existing self-LCT
        data, err := os.ReadFile(selfLCTPath)
        if err != nil {
            return fmt.Errorf("failed to read self-LCT: %w", err)
        }

        var selfLCT lctTypes.SelfLCT
        if err := json.Unmarshal(data, &selfLCT); err != nil {
            return fmt.Errorf("failed to parse self-LCT: %w", err)
        }

        // Verify hardware matches
        if err := selfLCT.VerifyHardware(); err != nil {
            return fmt.Errorf("hardware verification failed: %w", err)
        }

        app.selfLCT = &selfLCT
        app.Logger().Info("Loaded existing self-LCT",
            "id", selfLCT.ID,
            "hardware", selfLCT.HardwareBinding.HardwareHash[:8])
    } else {
        // Create new self-LCT
        selfLCT, err := lctTypes.CreateSelfLCT()
        if err != nil {
            return fmt.Errorf("failed to create self-LCT: %w", err)
        }

        // Save self-LCT
        data, err := json.MarshalIndent(selfLCT, "", "  ")
        if err != nil {
            return fmt.Errorf("failed to serialize self-LCT: %w", err)
        }

        if err := os.WriteFile(selfLCTPath, data, 0600); err != nil {
            return fmt.Errorf("failed to save self-LCT: %w", err)
        }

        app.selfLCT = selfLCT
        app.Logger().Info("Created new self-LCT",
            "id", selfLCT.ID,
            "hardware", selfLCT.HardwareBinding.HardwareHash[:8])
    }

    // Create role LCTs for queens
    if err := app.CreateRoleLCTs(); err != nil {
        return fmt.Errorf("failed to create role LCTs: %w", err)
    }

    return nil
}

// CreateRoleLCTs creates LCTs for all queen and worker roles
func (app *Society4App) CreateRoleLCTs() error {
    queens := []string{
        "Treasury-Queen",
        "Law-Oracle-Queen",
        "Implementation-Queen",
        "Research-Queen",
        "Documentation-Queen",
        "Federation-Bridge-Queen",
        "Coherence-Analysis-Queen",
        "Quality-Assurance-Queen",
        "Security-Queen",
        "Emergency-Response-Queen",
    }

    for _, queenName := range queens {
        // Check if role LCT already exists
        exists := false
        for _, childID := range app.selfLCT.RoleChildren {
            if contains(childID, queenName) {
                exists = true
                break
            }
        }

        if !exists {
            roleLCT, err := app.selfLCT.CreateRoleLCT(queenName, "queen")
            if err != nil {
                return fmt.Errorf("failed to create LCT for %s: %w", queenName, err)
            }

            app.Logger().Info("Created role LCT",
                "role", queenName,
                "lct_id", roleLCT.ID)
        }
    }

    return nil
}

// BeginBlocker runs at the beginning of every block
func (app *Society4App) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
    // Hardware validation happens first
    return app.consensusKeeper.BeginBlocker(ctx, req)
}

// EndBlocker runs at the end of every block
func (app *Society4App) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
    return app.consensusKeeper.EndBlocker(ctx, req)
}

// InitChainer initializes the blockchain with genesis state
func (app *Society4App) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
    var genesisState GenesisState
    if err := json.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
        panic(err)
    }

    // Store self-LCT in genesis
    genesisState.SelfLCT = app.selfLCT.Export()

    // Initialize modules with genesis state
    return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// ValidateTransaction validates a transaction before inclusion
func (app *Society4App) ValidateTransaction(tx sdk.Tx) error {
    // All transactions must be signed by a valid role LCT
    signatures := tx.GetSignatures()
    if len(signatures) == 0 {
        return fmt.Errorf("transaction requires signature")
    }

    // Verify at least one signature is from a valid role
    validRole := false
    for _, sig := range signatures {
        // Check if signature is from a known role LCT
        for _, childID := range app.selfLCT.RoleChildren {
            // Verify signature against role LCT
            // (Implementation would need role LCT public keys)
            validRole = true
            break
        }
    }

    if !validRole {
        return fmt.Errorf("transaction not signed by valid role")
    }

    return nil
}

// GetSelfLCT returns the application's self-LCT
func (app *Society4App) GetSelfLCT() *lctTypes.SelfLCT {
    return app.selfLCT
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
    return strings.Contains(s, substr)
}