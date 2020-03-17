package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"
)

const (
	// ModuleName is the name of the module
	ModuleName = "tokenswap"

	// StoreKey is used to register the module's store
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the tokenswap module
	QuerierRoute = ModuleName
)

// TokenSwap struct containing the data of the TokenSwap. json and yaml tags are used to specify field names when marshalled
type TokenSwap struct {
	EthereumTxHash string         `json:"ethereum_tx_hash" yaml:"ethereum_tx_hash"`
	EthereumSender string         `json:"ethereum_sender" yaml:"ethereum_sender"`
	Receiver       sdk.AccAddress `json:"receiver" yaml:"receiver"`
	AmountUSCRT    sdk.Coins      `json:"amount_uscrt" yaml:"amount_uscrt"`
}

// TokenSwap struct containing the data of the TokenSwap. json and yaml tags are used to specify field names when marshalled
type Params struct {
	MultisigApproveAddress sdk.AccAddress `json:"minting_approver_address" yaml:"minting_approver_address"`
	MintingMultiple        sdk.Dec        `json:"minting_multiple" yaml:"minting_multiple"`
	MintingEnabled         bool           `json:"minting_enabled" yaml:"minting_enabled"`
}

// NewTokenSwap Returns a new TokenSwap
func NewTokenSwap(ethereumTxHash string, ethereumSender string, receiver sdk.AccAddress, AmountUSCRT sdk.Coins) TokenSwap {
	return TokenSwap{
		EthereumTxHash: ethereumTxHash,
		EthereumSender: ethereumSender,
		Receiver:       receiver,
		AmountUSCRT:    AmountUSCRT,
	}
}

// String implement fmt.Stringer
func (s TokenSwap) String() string {
	return strings.TrimSpace(
		fmt.Sprintf(`EthereumTxHash=%s EthereumSender=%s Receiver=%s Amount=%s`,
			s.EthereumTxHash,
			s.EthereumSender,
			s.Receiver.String(),
			s.AmountUSCRT.String(),
		),
	)
}

// AccountKeeper defines the expected account keeper
type AccountKeeper interface {
	GetAccount(sdk.Context, sdk.AccAddress) authexported.Account
}

// SupplyKeeper defines the expected supply keeper
type SupplyKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	SetModuleAccount(sdk.Context, supplyexported.ModuleAccountI)
}
