package types

// GenesisState - all distribution state that must be provided at genesis
type GenesisState struct {
	Params Params `json:"params" yaml:"params"`
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

func NewGenesisState(params Params) GenesisState {
	return GenesisState{
		Params: params,
	}
}

// ValidateGenesis validates the genesis state of distribution genesis input
func ValidateGenesis(gs GenesisState) error {
	if err := gs.Params.ValidateBasic(); err != nil {
		return err
	}
	return nil
}
