package tokenswap

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	// sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns a handler for "tokenswap" type messages.
func NewHandler(keeper SwapKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case MsgSwapRequest:
			res, err := handleMsgTokenSwap(ctx, keeper, msg)
			if err != nil {
				return err.Result()
			}
			return *res

		default:
			errMsg := fmt.Sprintf("unrecognized tokenswap message type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle a message to create a token swap
func handleMsgTokenSwap(
	ctx sdk.Context, keeper SwapKeeper, msg MsgSwapRequest,
) (*sdk.Result, sdk.Error) {

	err := keeper.SwapIsEnabled(ctx)
	if err != nil {
		return nil, err
	}

	// validate signer
	err = keeper.ValidateTokenSwapSigner(ctx, msg.SignerAddr)
	if err != nil {
		return nil, err
	}

	// Check if the this tokeswap request was alread processed
	swapRecord, err := keeper.GetPastTokenSwapRequest(ctx, msg.BurnTxHash)
	if err == nil {
		// msg.EthereumTxHash already exists in db
		// So this request was already processed
		// Check if we might have failed processing the transaction
		if swapRecord.Done {
			return nil, sdk.ErrUnauthorized(fmt.Sprintf(
				"TokenSwap with EthereumTxHash %s was already processed",
				msg.BurnTxHash))
		}
	}

	err = keeper.ProcessTokenSwapRequest(
		ctx,
		msg,
	)
	return &sdk.Result{}, err

}
