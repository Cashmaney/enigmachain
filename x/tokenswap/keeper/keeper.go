package keeper

import (
	"fmt"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/subspace"
	"github.com/enigmampc/EnigmaBlockchain/x/tokenswap/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/tendermint/tendermint/libs/log"

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

func (k Keeper) GetMintedCoin(ctx sdk.Context, amtEngDust sdk.Dec) sdk.Coin {
	// ENG has 8 decimals, and SCRT has 6, so we divide the number of ENG dust by 100
	amountUscrt := amtEngDust.Mul(sdk.NewDecWithPrec(1, 2))
	mintMultiplier := k.GetMintingMultiplier(ctx)
	uscrtToMint := mintMultiplier.MulInt64(amountUscrt.Int64())

	return sdk.NewCoin("uscrt", uscrtToMint.RoundInt())
}

// ProcessTokenSwapRequest processes a claim that has just completed successfully with consensus
// Also note that at this stage we already validated the swap request parameters
func (k Keeper) ProcessTokenSwapRequest(
	ctx sdk.Context, request types.MsgSwapRequest) error {
	// Store the token swap request in our state
	// We need this to verify we process each request only once
	store := ctx.KVStore(k.storeKey)

	uscrtCoin := k.GetMintedCoin(ctx, request.AmountENG)

	tokenSwap := types.NewTokenSwapRecord(request.BurnTxHash, request.EthereumSender, request.Receiver, uscrtCoin, false)

	// Register the swap as started
	store.Set(tokenSwap.BurnTxHash.Bytes(), k.cdc.MustMarshalBinaryBare(tokenSwap))

	// Mint new uSCRTs
	err := k.supplyKeeper.MintCoins(
		ctx,
		types.ModuleName,
		sdk.NewCoins(tokenSwap.AmountUSCRT),
	)
	if err != nil {
		return err
	}

	// Transfer new funds to receiver
	err = k.supplyKeeper.SendCoinsFromModuleToAccount(
		ctx,
		types.ModuleName,
		tokenSwap.Receiver,
		sdk.NewCoins(tokenSwap.AmountUSCRT),
	)
	if err != nil {
		return err
	}

	tokenSwap.Done = true

	//update the status of the swap as successful
	store.Set(tokenSwap.BurnTxHash.Bytes(), k.cdc.MustMarshalBinaryBare(tokenSwap))

	return nil
}

// GetPastTokenSwapRequest retrives a past token swap request
func (k Keeper) GetPastTokenSwapRequest(ctx sdk.Context, ethereumTxHash types.EthereumTxHash) (types.TokenSwapRecord, error) {
	store := ctx.KVStore(k.storeKey)

	if !store.Has(ethereumTxHash.Bytes()) {
		return types.TokenSwapRecord{}, sdkerrors.Wrap(
			sdkerrors.ErrUnknownRequest,
			"Unknown Ethereum tx hash: "+ethereumTxHash.String())

	}

	bz := store.Get(ethereumTxHash.Bytes())
	var tokenSwap types.TokenSwapRecord
	k.cdc.MustUnmarshalBinaryBare(bz, &tokenSwap)

	return tokenSwap, nil
}

// GetTokenSwapRecordsIterator get an iterator over tokenswap records
func (k Keeper) GetTokenSwapRecordsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}
