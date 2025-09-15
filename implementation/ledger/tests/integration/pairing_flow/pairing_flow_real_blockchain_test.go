package pairing_flow_test

import (
	"context"
	"testing"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	componentregistrykeeper "racecar-web/x/componentregistry/keeper"
	componentregistrytypes "racecar-web/x/componentregistry/types"
	energycyclekeeper "racecar-web/x/energycycle/keeper"
	energycycletypes "racecar-web/x/energycycle/types"
	lctmanagerkeeper "racecar-web/x/lctmanager/keeper"
	lctmanagertypes "racecar-web/x/lctmanager/types"
	pairingkeeper "racecar-web/x/pairing/keeper"
	pairing "racecar-web/x/pairing/module"
	pairingtypes "racecar-web/x/pairing/types"
	pairingqueuekeeper "racecar-web/x/pairingqueue/keeper"
	pairingqueuetypes "racecar-web/x/pairingqueue/types"
	trusttensorkeeper "racecar-web/x/trusttensor/keeper"
	trusttensortypes "racecar-web/x/trusttensor/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PairingFlowRealBlockchainTestSuite struct {
	suite.Suite
	ctx context.Context

	pairingKeeper      pairingkeeper.Keeper
	compRegKeeper      componentregistrykeeper.Keeper
	lctKeeper          lctmanagerkeeper.Keeper
	trustTensorKeeper  trusttensorkeeper.Keeper
	energyCycleKeeper  energycyclekeeper.Keeper
	pairingQueueKeeper pairingqueuekeeper.Keeper
}

func (suite *PairingFlowRealBlockchainTestSuite) SetupTest() {
	encCfg := moduletestutil.MakeTestEncodingConfig(pairing.AppModule{})
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())

	pairingStoreKey := storetypes.NewKVStoreKey(pairingtypes.StoreKey)
	componentRegistryStoreKey := storetypes.NewKVStoreKey(componentregistrytypes.StoreKey)
	lctManagerStoreKey := storetypes.NewKVStoreKey(lctmanagertypes.StoreKey)
	trustTensorStoreKey := storetypes.NewKVStoreKey(trusttensortypes.StoreKey)
	energyCycleStoreKey := storetypes.NewKVStoreKey(energycycletypes.StoreKey)
	pairingQueueStoreKey := storetypes.NewKVStoreKey(pairingqueuetypes.StoreKey)

	pairingStoreService := runtime.NewKVStoreService(pairingStoreKey)
	componentRegistryStoreService := runtime.NewKVStoreService(componentRegistryStoreKey)
	lctManagerStoreService := runtime.NewKVStoreService(lctManagerStoreKey)
	trustTensorStoreService := runtime.NewKVStoreService(trustTensorStoreKey)
	energyCycleStoreService := runtime.NewKVStoreService(energyCycleStoreKey)
	pairingQueueStoreService := runtime.NewKVStoreService(pairingQueueStoreKey)

	ctx := testutil.DefaultContextWithDB(
		suite.T(),
		pairingStoreKey,
		storetypes.NewTransientStoreKey("transient_test"),
	).Ctx

	authority := authtypes.NewModuleAddress(pairingtypes.GovModuleName)

	compRegKeeper := componentregistrykeeper.NewKeeper(
		componentRegistryStoreService,
		encCfg.Codec,
		addressCodec,
		authority,
		nil, nil, nil,
	)
	lctKeeper := lctmanagerkeeper.NewKeeper(
		lctManagerStoreService,
		encCfg.Codec,
		addressCodec,
		authority,
		nil, nil, nil, log.NewNopLogger(),
	)
	trustTensorKeeper := trusttensorkeeper.NewKeeper(
		trustTensorStoreService,
		encCfg.Codec,
		addressCodec,
		authority,
		nil, nil,
	)
	energyCycleKeeper := energycyclekeeper.NewKeeper(
		energyCycleStoreService,
		encCfg.Codec,
		addressCodec,
		authority,
		nil, nil, nil,
	)
	pairingQueueKeeper := pairingqueuekeeper.NewKeeper(
		encCfg.Codec,
		pairingQueueStoreService,
		nil, // authKeeper
		nil, // bankKeeper
		nil, // compKeeper
	)
	pairingKeeper := pairingkeeper.NewKeeper(
		pairingStoreService,
		encCfg.Codec,
		addressCodec,
		authority,
		nil, // bankKeeper
		nil, // componentregistryKeeper
		nil, // pairingqueueKeeper
		nil, // lctmanagerKeeper
	)

	suite.ctx = ctx
	suite.pairingKeeper = pairingKeeper
	suite.compRegKeeper = compRegKeeper
	suite.lctKeeper = lctKeeper
	suite.trustTensorKeeper = trustTensorKeeper
	suite.energyCycleKeeper = energyCycleKeeper
	suite.pairingQueueKeeper = pairingQueueKeeper
}

func (suite *PairingFlowRealBlockchainTestSuite) TestKeeperInitialization() {
	assert.NotNil(suite.T(), suite.pairingKeeper)
	assert.NotNil(suite.T(), suite.compRegKeeper)
	assert.NotNil(suite.T(), suite.lctKeeper)
	assert.NotNil(suite.T(), suite.trustTensorKeeper)
	assert.NotNil(suite.T(), suite.energyCycleKeeper)
	assert.NotNil(suite.T(), suite.pairingQueueKeeper)
}

func TestPairingFlowRealBlockchainTestSuite(t *testing.T) {
	suite.Run(t, new(PairingFlowRealBlockchainTestSuite))
}
