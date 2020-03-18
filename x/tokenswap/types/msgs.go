package types

import (
	"fmt"
	"regexp"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// RouterKey is used to route messages and queriers to the greeter module
const RouterKey = "tokenswap"

// MsgSwapRequest defines the MsgSwapRequest Message
type MsgSwapRequest struct {
	BurnTxHash     string
	EthereumSender string
	Receiver       sdk.AccAddress
	AmountENG      string
	SignerAddr     sdk.AccAddress
}

// Check in compile time that MsgSwapRequest is a sdk.Msg
var _ sdk.Msg = MsgSwapRequest{}

// NewMsgSwapRequest Returns a new MsgSwapRequest
func NewMsgSwapRequest(burnTxHash string, ethereumSender string, receiver sdk.AccAddress, signerAddr sdk.AccAddress, amountENG string) MsgSwapRequest {

	return MsgSwapRequest{
		//BurnTxHash:     HexToTxHash(burnTxHash),
		//EthereumSender: HexToAddress(ethereumSender),
		BurnTxHash:     burnTxHash,
		EthereumSender: ethereumSender,
		Receiver:       receiver,
		AmountENG:      amountENG,
		SignerAddr:     signerAddr,
	}
}

// Route should return the name of the module
func (msg MsgSwapRequest) Route() string { return RouterKey }

// Type should return the action
func (msg MsgSwapRequest) Type() string { return "tokenswap" }

var ethereumTxHashRegex = regexp.MustCompile(`^0x[0-9a-fA-F]{64}$`)
var ethereumAddressRegex = regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)

// ValidateBasic runs stateless checks on the message
func (msg MsgSwapRequest) ValidateBasic() error {
	if !ethereumTxHashRegex.MatchString(msg.BurnTxHash) {
		return sdkerrors.Wrap(
			sdkerrors.ErrUnknownRequest,
			fmt.Sprintf(
				`Invalid EthereumTxHash %s accoding to regex '%s'`,
				msg.BurnTxHash,
				ethereumTxHashRegex.String(),
			),
		)
	}
	if !ethereumAddressRegex.MatchString(msg.EthereumSender) {
		return sdkerrors.Wrap(
			sdkerrors.ErrUnknownRequest,
			fmt.Sprintf(
				`Invalid EthereumSender %s accoding to regex '%s'`,
				msg.EthereumSender,
				ethereumAddressRegex.String(),
			),
		)
	}

	if msg.Receiver.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Receiver cannot be empty")
	}

	engDust, err := strconv.ParseInt(msg.AmountENG, 10, 64)
	if err != nil {
		return err
	}
	if engDust <= 0 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Amount %d must be positive", engDust))
	}
	return nil
}

// GetSigners returns the addresses of those required to sign the message
func (msg MsgSwapRequest) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.SignerAddr}
}

// GetSignBytes encodes the message for signing
func (msg MsgSwapRequest) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}
