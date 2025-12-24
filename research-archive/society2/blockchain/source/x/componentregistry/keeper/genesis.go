package keeper

import (
	"context"

	"cosmossdk.io/errors"

	"racecar-web/x/componentregistry/types"
)

// InitGenesis initializes the module's state from a provided genesis state
func (k Keeper) InitGenesis(ctx context.Context, genState types.GenesisState) error {
	// Set params
	if err := k.SetParams(ctx, genState.Params); err != nil {
		return errors.Wrap(err, "failed to set params")
	}

	// Set components
	for _, component := range genState.Components {
		if err := k.Components.Set(ctx, component.ComponentId, component); err != nil {
			return errors.Wrapf(err, "failed to initialize component %s", component.ComponentId)
		}
	}

	// Set verifications
	for _, verification := range genState.ComponentVerifications {
		if err := k.ComponentVerifications.Set(ctx, verification.ComponentId, verification); err != nil {
			return errors.Wrapf(err, "failed to initialize verification for component %s", verification.ComponentId)
		}
	}

	// Set pairing rules
	for _, rule := range genState.PairingRules {
		ruleKey := rule.SourceTypeHash + "-" + rule.TargetTypeHash
		if err := k.ComponentPairingRules.Set(ctx, ruleKey, rule); err != nil {
			return errors.Wrapf(err, "failed to set pairing rule for %s-%s", rule.SourceTypeHash, rule.TargetTypeHash)
		}
	}

	return nil
}

// ExportGenesis returns the module's exported genesis state
func (k Keeper) ExportGenesis(ctx context.Context) (*types.GenesisState, error) {
	// Get params
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get params")
	}

	// Get all components
	var components []types.Component
	err = k.Components.Walk(ctx, nil, func(key string, component types.Component) (bool, error) {
		components = append(components, component)
		return false, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list components")
	}

	// Get all verifications
	var verifications []types.ComponentVerification
	err = k.ComponentVerifications.Walk(ctx, nil, func(key string, verification types.ComponentVerification) (bool, error) {
		verifications = append(verifications, verification)
		return false, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list verifications")
	}

	// Get all pairing rules
	var pairingRules []types.ComponentPairingRule
	err = k.ComponentPairingRules.Walk(ctx, nil, func(key string, rule types.ComponentPairingRule) (bool, error) {
		pairingRules = append(pairingRules, rule)
		return false, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to list pairing rules")
	}

	return &types.GenesisState{
		Params:                 params,
		Components:             components,
		ComponentVerifications: verifications,
		PairingRules:           pairingRules,
	}, nil
}
