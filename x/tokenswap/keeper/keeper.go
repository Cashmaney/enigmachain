package keeper

import (
	"fmt"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/subspace"
	"github.com/enigmampc/EnigmaBlockchain/x/tokenswap/types"
	"math"
	"strconv"
	"strings"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	cdc          *codec.Codec // The wire codec for binary encoding/decoding.
	storeKey     sdk.StoreKey // Unexposed key to access store from sdk.Context
	paramSpace   paramtypes.Subspace
	supplyKeeper types.SupplyKeeper
}

// NewKeeper creates new instances of the Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, params paramtypes.Subspace, supplyKeeper types.SupplyKeeper) Keeper {
	if !params.HasKeyTable() {
		params = params.WithKeyTable(types.ParamKeyTable())
	}
	return Keeper{
		cdc:          cdc,
		storeKey:     storeKey,
		paramSpace:   params,
		supplyKeeper: supplyKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) SwapIsEnabled(ctx sdk.Context) error {
	if !k.GetMintingEnabled(ctx) {
		return fmt.Errorf("Token swap is disabled. Requires a parameter change proposal to enable")
	}
	return nil
}

func (k Keeper) ValidateTokenSwapSigner(ctx sdk.Context, signer sdk.AccAddress) error {
	if signer.String() != k.GetMultisigApproveAddress(ctx).String() {
		return fmt.Errorf("invalid signer address")
	}
	return nil
}

// ProcessTokenSwapRequest processes a claim that has just completed successfully with consensus
func (k Keeper) ProcessTokenSwapRequest(ctx sdk.Context, ethereumTxHash string, ethereumSender string, receiver sdk.AccAddress, amountENG string) error {

	// Convert ENG to uSCRT
	engFloat, err := strconv.ParseFloat(amountENG, 64)
	if err != nil {
		return err
	}

	mintMul, err := strconv.ParseFloat(k.GetMintingMultiple(ctx).String(), 64)
	if err != nil {
		return err
	}

	engToMint := engFloat * mintMul

	amountUscrt := int64(math.Ceil(engToMint * 1e6))
	amountUscrtCoins := sdk.NewCoins(sdk.NewCoin("uscrt", sdk.NewInt(amountUscrt)))

	// Lowercase ethereumTxHash as this is our indexed field
	ethereumTxHashLowercase := strings.ToLower(ethereumTxHash)
	tokenSwap := types.NewTokenSwap(ethereumTxHashLowercase, ethereumSender, receiver, amountUscrtCoins)

	// Mint new uSCRTs
	err = k.supplyKeeper.MintCoins(
		ctx,
		types.ModuleName,
		tokenSwap.AmountUSCRT,
	)
	if err != nil {
		return err
	}

	// Transfer new funds to receiver
	err = k.supplyKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		tokenSwap.Receiver,
		tokenSwap.AmountUSCRT,
	)
	if err != nil {
		return err
	}

	// Store the token swap request in our state
	// We need this to verify we process each request only once
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(tokenSwap.EthereumTxHash), k.cdc.MustMarshalBinaryBare(tokenSwap))

	return nil
}

// GetPastTokenSwapRequest retrives a past token swap request
func (k Keeper) GetPastTokenSwapRequest(ctx sdk.Context, ethereumTxHash string) (types.TokenSwap, error) {
	store := ctx.KVStore(k.storeKey)

	// Lowercase ethereumTxHash as this is our indexed field
	ethereumTxHashLowercase := strings.ToLower(ethereumTxHash)

	if !store.Has([]byte(ethereumTxHashLowercase)) {
		return types.TokenSwap{}, sdkerrors.Wrap(
			sdkerrors.ErrUnknownRequest,
			"Unknown Ethereum tx hash "+ethereumTxHash)

	}

	bz := store.Get([]byte(ethereumTxHashLowercase))
	var tokenSwap types.TokenSwap
	k.cdc.MustUnmarshalBinaryBare(bz, &tokenSwap)

	return tokenSwap, nil
}

// GetTokenSwapRecordsIterator get an iterator over tokenswap records
func (k Keeper) GetTokenSwapRecordsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}
