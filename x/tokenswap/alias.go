package tokenswap

import (
	"github.com/enigmampc/EnigmaBlockchain/x/tokenswap/keeper"
	"github.com/enigmampc/EnigmaBlockchain/x/tokenswap/types"
)

const (
	DefaultParamspace = types.DefaultParamspace
	ModuleName        = types.ModuleName
	RouterKey         = types.RouterKey
	StoreKey          = types.StoreKey
	QuerierRoute      = types.QuerierRoute
)

// functions aliases
var (
	ModuleCdc             = types.ModuleCdc
	RegisterCodec         = types.RegisterCodec
	NewKeeper             = keeper.NewKeeper
	NewQuerier            = keeper.NewQuerier
	NewTokenSwap          = types.NewTokenSwap
	NewMsgTokenSwap       = types.NewMsgTokenSwap
	NewGetTokenSwapParams = types.NewGetTokenSwapParams
)

type (
	Keeper                 = keeper.Keeper
	TokenSwap              = types.TokenSwap
	MsgTokenSwap           = types.MsgTokenSwap
	QueryEthProphecyParams = types.GetTokenSwapParams
)
