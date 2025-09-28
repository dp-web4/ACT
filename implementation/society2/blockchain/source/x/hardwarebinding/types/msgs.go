package types

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgUpdateHardware = "update_hardware"

var _ sdk.Msg = &MsgUpdateHardware{}

// MsgUpdateHardware is used for testing/admin purposes
type MsgUpdateHardware struct {
    Creator string `json:"creator,omitempty"`
    Enable  bool   `json:"enable,omitempty"`
}

// NewMsgUpdateHardware creates a new MsgUpdateHardware instance
func NewMsgUpdateHardware(creator string, enable bool) *MsgUpdateHardware {
    return &MsgUpdateHardware{
        Creator: creator,
        Enable:  enable,
    }
}

// ValidateBasic performs basic validation
func (msg *MsgUpdateHardware) ValidateBasic() error {
    _, err := sdk.AccAddressFromBech32(msg.Creator)
    if err != nil {
        return err
    }
    return nil
}

// GetSigners returns the signers of the message
func (msg *MsgUpdateHardware) GetSigners() []sdk.AccAddress {
    creator, _ := sdk.AccAddressFromBech32(msg.Creator)
    return []sdk.AccAddress{creator}
}

// Type returns the message type
func (msg *MsgUpdateHardware) Type() string {
    return TypeMsgUpdateHardware
}

// Route returns the message route
func (msg *MsgUpdateHardware) Route() string {
    return RouterKey
}

// MsgUpdateHardwareResponse defines the response
type MsgUpdateHardwareResponse struct{}

// Query messages
type QueryHardwareRequest struct{}

type QueryHardwareResponse struct {
    HardwareBinding *HardwareBinding `json:"hardware_binding,omitempty"`
    CurrentMatch    bool             `json:"current_match"`
    ValidationEnabled bool           `json:"validation_enabled"`
}

// Service descriptors (simplified stubs)
var _Msg_serviceDesc = sdk.ServiceDescriptor{
    ServiceName: "society4chain.hardwarebinding.Msg",
    Methods: []sdk.MethodDescriptor{
        {
            Name: "UpdateHardware",
        },
    },
}

var _Query_serviceDesc = sdk.ServiceDescriptor{
    ServiceName: "society4chain.hardwarebinding.Query",
    Methods: []sdk.MethodDescriptor{
        {
            Name: "Hardware",
        },
    },
}

// MsgServer interface
type MsgServer interface {
    UpdateHardware(ctx sdk.Context, msg *MsgUpdateHardware) (*MsgUpdateHardwareResponse, error)
}

// QueryServer interface
type QueryServer interface {
    Hardware(ctx sdk.Context, req *QueryHardwareRequest) (*QueryHardwareResponse, error)
}

// RegisterMsgServer registers the msg server
func RegisterMsgServer(router sdk.Router, srv MsgServer) {
    // Implementation would go here
}

// RegisterQueryServer registers the query server
func RegisterQueryServer(router sdk.Router, srv QueryServer) {
    // Implementation would go here
}