package types

import (
	"strconv"

	"cosmossdk.io/errors"
)

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params:                 DefaultParams(),
		Components:             []Component{},
		ComponentVerifications: []ComponentVerification{},
		PairingRules:           []ComponentPairingRule{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Validate params
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	// Validate components
	componentIDs := make(map[string]bool)
	for _, component := range gs.Components {
		// Check for duplicate component IDs
		if componentIDs[component.ComponentId] {
			return errors.Wrapf(ErrInvalidComponentID, "duplicate component ID: %s", component.ComponentId)
		}
		componentIDs[component.ComponentId] = true

		// Validate component
		if err := component.Validate(); err != nil {
			return errors.Wrapf(err, "invalid component %s", component.ComponentId)
		}
	}

	// Validate verifications
	for _, verification := range gs.ComponentVerifications {
		// Check that verified component exists
		if !componentIDs[verification.ComponentId] {
			return errors.Wrapf(ErrComponentNotFound, "verification for non-existent component: %s", verification.ComponentId)
		}

		// Validate verification
		if err := verification.Validate(); err != nil {
			return errors.Wrapf(err, "invalid verification for component %s", verification.ComponentId)
		}
	}

	// Validate pairing rules
	for _, rule := range gs.PairingRules {
		if err := rule.Validate(); err != nil {
			return errors.Wrapf(err, "invalid pairing rule")
		}
	}

	return nil
}

// Validate performs basic validation on a Component
func (c *Component) Validate() error {
	if c.ComponentId == "" {
		return errors.Wrap(ErrInvalidComponentID, "component ID cannot be empty")
	}
	// Use new hash-based fields for privacy-focused implementation
	if c.CategoryHash == "" {
		return errors.Wrap(ErrInvalidComponentType, "category hash cannot be empty")
	}
	if c.ManufacturerHash == "" {
		return errors.Wrap(ErrInvalidAuthority, "manufacturer hash cannot be empty")
	}
	return nil
}

// Validate performs basic validation on a ComponentVerification
func (v *ComponentVerification) Validate() error {
	if v.ComponentId == "" {
		return errors.Wrap(ErrInvalidComponentID, "component ID cannot be empty")
	}
	if v.Status == "" {
		return errors.Wrap(ErrInvalidComponentType, "verification status cannot be empty")
	}
	if v.VerifiedAt.IsZero() {
		return errors.Wrap(ErrInvalidComponentType, "verification timestamp cannot be zero")
	}
	return nil
}

// Validate performs basic validation on a ComponentPairingRule
func (r *ComponentPairingRule) Validate() error {
	if r.TargetTypeHash == "" {
		return errors.Wrap(ErrInvalidComponentType, "target component type hash cannot be empty")
	}

	// Parse MinQualityScore as float64 for validation
	if r.MinQualityScore != "" {
		score, err := strconv.ParseFloat(r.MinQualityScore, 64)
		if err != nil {
			return errors.Wrap(ErrInvalidComponentType, "invalid quality score format")
		}
		if score < 0 || score > 100 {
			return errors.Wrap(ErrInvalidComponentType, "minimum quality score must be between 0 and 100")
		}
	}

	return nil
}
