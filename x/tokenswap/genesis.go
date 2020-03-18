package tokenswap

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
)

func InitGenesis(ctx sdk.Context, supplyKeeper SupplyKeeper, keeper SwapKeeper, data GenesisState) []abci.ValidatorUpdate {
	tokenSwapAccount := supply.NewEmptyModuleAccount(ModuleName, supply.Burner, supply.Minter)
	supplyKeeper.SetModuleAccount(ctx, tokenSwapAccount)
	keeper.SetParams(ctx, data.Params)
	return nil
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper SwapKeeper) GenesisState {
	params := keeper.GetParams(ctx)
	return NewGenesisState(params)
}
