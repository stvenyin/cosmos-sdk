package types

import (
	"bytes"
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	tmtypes "github.com/tendermint/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

// Validator defines the total amount of bond shares and their exchange rate to
// coins. Accumulation of interest is modelled as an in increase in the
// exchange rate, and slashing as a decrease.  When coins are delegated to this
// validator, the validator is credited with a Delegation whose number of
// bond shares is based on the amount of coins delegated divided by the current
// exchange rate. Voting power can be calculated as total bonds multiplied by
// exchange rate.
type Validator struct {
	Operator sdk.AccAddress `json:"operator"` // sender of BondTx - UnbondTx returns here
	PubKey   crypto.PubKey  `json:"pub_key"`  // pubkey of validator
	Revoked  bool           `json:"revoked"`  // has the validator been revoked from bonded status?

	Status          sdk.BondStatus `json:"status"`           // validator status (bonded/unbonding/unbonded)
	Tokens          sdk.Dec        `json:"tokens"`           // delegated tokens (incl. self-delegation)
	DelegatorShares sdk.Dec        `json:"delegator_shares"` // total shares issued to a validator's delegators

	Description        Description `json:"description"`           // description terms for the validator
	BondHeight         int64       `json:"bond_height"`           // earliest height as a bonded validator
	BondIntraTxCounter int16       `json:"bond_intra_tx_counter"` // block-local tx index of validator change
	ProposerRewardPool sdk.Coins   `json:"proposer_reward_pool"`  // XXX reward pool collected from being the proposer

	Commission            sdk.Dec `json:"commission"`              // XXX the commission rate of fees charged to any delegators
	CommissionMax         sdk.Dec `json:"commission_max"`          // XXX maximum commission rate which this validator can ever charge
	CommissionChangeRate  sdk.Dec `json:"commission_change_rate"`  // XXX maximum daily increase of the validator commission
	CommissionChangeToday sdk.Dec `json:"commission_change_today"` // XXX commission rate change today, reset each day (UTC time)

	// fee related
	LastBondedTokens sdk.Dec `json:"prev_bonded_tokens"` // Previous bonded tokens held
}

// NewValidator - initialize a new validator
func NewValidator(operator sdk.AccAddress, pubKey crypto.PubKey, description Description) Validator {
	return Validator{
		Operator:              operator,
		PubKey:                pubKey,
		Revoked:               false,
		Status:                sdk.Unbonded,
		Tokens:                sdk.ZeroDec(),
		DelegatorShares:       sdk.ZeroDec(),
		Description:           description,
		BondHeight:            int64(0),
		BondIntraTxCounter:    int16(0),
		ProposerRewardPool:    sdk.Coins{},
		Commission:            sdk.ZeroDec(),
		CommissionMax:         sdk.ZeroDec(),
		CommissionChangeRate:  sdk.ZeroDec(),
		CommissionChangeToday: sdk.ZeroDec(),
		LastBondedTokens:      sdk.ZeroDec(),
	}
}

// what's kept in the store value
type validatorValue struct {
	PubKey                crypto.PubKey
	Revoked               bool
	Status                sdk.BondStatus
	Tokens                sdk.Dec
	DelegatorShares       sdk.Dec
	Description           Description
	BondHeight            int64
	BondIntraTxCounter    int16
	ProposerRewardPool    sdk.Coins
	Commission            sdk.Dec
	CommissionMax         sdk.Dec
	CommissionChangeRate  sdk.Dec
	CommissionChangeToday sdk.Dec
	LastBondedTokens      sdk.Dec
}

// return the redelegation without fields contained within the key for the store
func MustMarshalValidator(cdc *wire.Codec, validator Validator) []byte {
	val := validatorValue{
		PubKey:                validator.PubKey,
		Revoked:               validator.Revoked,
		Status:                validator.Status,
		Tokens:                validator.Tokens,
		DelegatorShares:       validator.DelegatorShares,
		Description:           validator.Description,
		BondHeight:            validator.BondHeight,
		BondIntraTxCounter:    validator.BondIntraTxCounter,
		ProposerRewardPool:    validator.ProposerRewardPool,
		Commission:            validator.Commission,
		CommissionMax:         validator.CommissionMax,
		CommissionChangeRate:  validator.CommissionChangeRate,
		CommissionChangeToday: validator.CommissionChangeToday,
		LastBondedTokens:      validator.LastBondedTokens,
	}
	return cdc.MustMarshalBinary(val)
}

// unmarshal a redelegation from a store key and value
func MustUnmarshalValidator(cdc *wire.Codec, operatorAddr, value []byte) Validator {
	validator, err := UnmarshalValidator(cdc, operatorAddr, value)
	if err != nil {
		panic(err)
	}

	return validator
}

// unmarshal a redelegation from a store key and value
func UnmarshalValidator(cdc *wire.Codec, operatorAddr, value []byte) (validator Validator, err error) {
	if len(operatorAddr) != sdk.AddrLen {
		err = fmt.Errorf("%v", ErrBadValidatorAddr(DefaultCodespace).Data())
		return
	}
	var storeValue validatorValue
	err = cdc.UnmarshalBinary(value, &storeValue)
	if err != nil {
		return
	}

	return Validator{
		Operator:              operatorAddr,
		PubKey:                storeValue.PubKey,
		Revoked:               storeValue.Revoked,
		Tokens:                storeValue.Tokens,
		Status:                storeValue.Status,
		DelegatorShares:       storeValue.DelegatorShares,
		Description:           storeValue.Description,
		BondHeight:            storeValue.BondHeight,
		BondIntraTxCounter:    storeValue.BondIntraTxCounter,
		ProposerRewardPool:    storeValue.ProposerRewardPool,
		Commission:            storeValue.Commission,
		CommissionMax:         storeValue.CommissionMax,
		CommissionChangeRate:  storeValue.CommissionChangeRate,
		CommissionChangeToday: storeValue.CommissionChangeToday,
		LastBondedTokens:      storeValue.LastBondedTokens,
	}, nil
}

// HumanReadableString returns a human readable string representation of a
// validator. An error is returned if the operator or the operator's public key
// cannot be converted to Bech32 format.
func (v Validator) HumanReadableString() (string, error) {
	bechVal, err := sdk.Bech32ifyValPub(v.PubKey)
	if err != nil {
		return "", err
	}

	resp := "Validator \n"
	resp += fmt.Sprintf("Operator: %s\n", v.Operator)
	resp += fmt.Sprintf("Validator: %s\n", bechVal)
	resp += fmt.Sprintf("Revoked: %v\n", v.Revoked)
	resp += fmt.Sprintf("Status: %s\n", sdk.BondStatusToString(v.Status))
	resp += fmt.Sprintf("Tokens: %s\n", v.Tokens.String())
	resp += fmt.Sprintf("Delegator Shares: %s\n", v.DelegatorShares.String())
	resp += fmt.Sprintf("Description: %s\n", v.Description)
	resp += fmt.Sprintf("Bond Height: %d\n", v.BondHeight)
	resp += fmt.Sprintf("Proposer Reward Pool: %s\n", v.ProposerRewardPool.String())
	resp += fmt.Sprintf("Commission: %s\n", v.Commission.String())
	resp += fmt.Sprintf("Max Commission Rate: %s\n", v.CommissionMax.String())
	resp += fmt.Sprintf("Commission Change Rate: %s\n", v.CommissionChangeRate.String())
	resp += fmt.Sprintf("Commission Change Today: %s\n", v.CommissionChangeToday.String())
	resp += fmt.Sprintf("Previous Bonded Tokens: %s\n", v.LastBondedTokens.String())

	return resp, nil
}

//___________________________________________________________________

// validator struct for bech output
type BechValidator struct {
	Operator sdk.AccAddress `json:"operator"` // in bech32
	PubKey   string         `json:"pub_key"`  // in bech32
	Revoked  bool           `json:"revoked"`  // has the validator been revoked from bonded status?

	Status          sdk.BondStatus `json:"status"`           // validator status (bonded/unbonding/unbonded)
	Tokens          sdk.Dec        `json:"tokens"`           // delegated tokens (incl. self-delegation)
	DelegatorShares sdk.Dec        `json:"delegator_shares"` // total shares issued to a validator's delegators

	Description        Description `json:"description"`           // description terms for the validator
	BondHeight         int64       `json:"bond_height"`           // earliest height as a bonded validator
	BondIntraTxCounter int16       `json:"bond_intra_tx_counter"` // block-local tx index of validator change
	ProposerRewardPool sdk.Coins   `json:"proposer_reward_pool"`  // XXX reward pool collected from being the proposer

	Commission            sdk.Dec `json:"commission"`              // XXX the commission rate of fees charged to any delegators
	CommissionMax         sdk.Dec `json:"commission_max"`          // XXX maximum commission rate which this validator can ever charge
	CommissionChangeRate  sdk.Dec `json:"commission_change_rate"`  // XXX maximum daily increase of the validator commission
	CommissionChangeToday sdk.Dec `json:"commission_change_today"` // XXX commission rate change today, reset each day (UTC time)

	// fee related
	LastBondedTokens sdk.Dec `json:"prev_bonded_shares"` // last bonded token amount
}

// get the bech validator from the the regular validator
func (v Validator) Bech32Validator() (BechValidator, error) {
	bechValPubkey, err := sdk.Bech32ifyValPub(v.PubKey)
	if err != nil {
		return BechValidator{}, err
	}

	return BechValidator{
		Operator: v.Operator,
		PubKey:   bechValPubkey,
		Revoked:  v.Revoked,

		Status:          v.Status,
		Tokens:          v.Tokens,
		DelegatorShares: v.DelegatorShares,

		Description:        v.Description,
		BondHeight:         v.BondHeight,
		BondIntraTxCounter: v.BondIntraTxCounter,
		ProposerRewardPool: v.ProposerRewardPool,

		Commission:            v.Commission,
		CommissionMax:         v.CommissionMax,
		CommissionChangeRate:  v.CommissionChangeRate,
		CommissionChangeToday: v.CommissionChangeToday,

		LastBondedTokens: v.LastBondedTokens,
	}, nil
}

//___________________________________________________________________

// only the vitals - does not check bond height of IntraTxCounter
// nolint gocyclo - why dis fail?
func (v Validator) Equal(c2 Validator) bool {
	return v.PubKey.Equals(c2.PubKey) &&
		bytes.Equal(v.Operator, c2.Operator) &&
		v.Status.Equal(c2.Status) &&
		v.Tokens.Equal(c2.Tokens) &&
		v.DelegatorShares.Equal(c2.DelegatorShares) &&
		v.Description == c2.Description &&
		v.ProposerRewardPool.IsEqual(c2.ProposerRewardPool) &&
		v.Commission.Equal(c2.Commission) &&
		v.CommissionMax.Equal(c2.CommissionMax) &&
		v.CommissionChangeRate.Equal(c2.CommissionChangeRate) &&
		v.CommissionChangeToday.Equal(c2.CommissionChangeToday) &&
		v.LastBondedTokens.Equal(c2.LastBondedTokens)
}

// constant used in flags to indicate that description field should not be updated
const DoNotModifyDesc = "[do-not-modify]"

// Description - description fields for a validator
type Description struct {
	Moniker  string `json:"moniker"`  // name
	Identity string `json:"identity"` // optional identity signature (ex. UPort or Keybase)
	Website  string `json:"website"`  // optional website link
	Details  string `json:"details"`  // optional details
}

// NewDescription returns a new Description with the provided values.
func NewDescription(moniker, identity, website, details string) Description {
	return Description{
		Moniker:  moniker,
		Identity: identity,
		Website:  website,
		Details:  details,
	}
}

// UpdateDescription updates the fields of a given description. An error is
// returned if the resulting description contains an invalid length.
func (d Description) UpdateDescription(d2 Description) (Description, sdk.Error) {
	if d2.Moniker == DoNotModifyDesc {
		d2.Moniker = d.Moniker
	}
	if d2.Identity == DoNotModifyDesc {
		d2.Identity = d.Identity
	}
	if d2.Website == DoNotModifyDesc {
		d2.Website = d.Website
	}
	if d2.Details == DoNotModifyDesc {
		d2.Details = d.Details
	}

	return Description{
		Moniker:  d2.Moniker,
		Identity: d2.Identity,
		Website:  d2.Website,
		Details:  d2.Details,
	}.EnsureLength()
}

// EnsureLength ensures the length of a validator's description.
func (d Description) EnsureLength() (Description, sdk.Error) {
	if len(d.Moniker) > 70 {
		return d, ErrDescriptionLength(DefaultCodespace, "moniker", len(d.Moniker), 70)
	}
	if len(d.Identity) > 3000 {
		return d, ErrDescriptionLength(DefaultCodespace, "identity", len(d.Identity), 3000)
	}
	if len(d.Website) > 140 {
		return d, ErrDescriptionLength(DefaultCodespace, "website", len(d.Website), 140)
	}
	if len(d.Details) > 280 {
		return d, ErrDescriptionLength(DefaultCodespace, "details", len(d.Details), 280)
	}

	return d, nil
}

// ABCIValidator returns an abci.Validator from a staked validator type.
func (v Validator) ABCIValidator() abci.Validator {
	return abci.Validator{
		PubKey:  tmtypes.TM2PB.PubKey(v.PubKey),
		Address: v.PubKey.Address(),
		Power:   v.BondedTokens().RoundInt64(),
	}
}

// ABCIValidatorZero returns an abci.Validator from a staked validator type
// with with zero power used for validator updates.
func (v Validator) ABCIValidatorZero() abci.Validator {
	return abci.Validator{
		PubKey:  tmtypes.TM2PB.PubKey(v.PubKey),
		Address: v.PubKey.Address(),
		Power:   0,
	}
}

// UpdateStatus updates the location of the shares within a validator
// to reflect the new status
func (v Validator) UpdateStatus(pool Pool, NewStatus sdk.BondStatus) (Validator, Pool) {

	switch v.Status {
	case sdk.Unbonded:

		switch NewStatus {
		case sdk.Unbonded:
			return v, pool
		case sdk.Bonded:
			pool = pool.looseTokensToBonded(v.Tokens)
		}
	case sdk.Unbonding:

		switch NewStatus {
		case sdk.Unbonding:
			return v, pool
		case sdk.Bonded:
			pool = pool.looseTokensToBonded(v.Tokens)
		}
	case sdk.Bonded:

		switch NewStatus {
		case sdk.Bonded:
			return v, pool
		default:
			pool = pool.bondedTokensToLoose(v.Tokens)
		}
	}

	v.Status = NewStatus
	return v, pool
}

// removes tokens from a validator
func (v Validator) RemoveTokens(pool Pool, tokens sdk.Dec) (Validator, Pool) {
	if v.Status == sdk.Bonded {
		pool = pool.bondedTokensToLoose(tokens)
	}

	v.Tokens = v.Tokens.Sub(tokens)
	return v, pool
}

//_________________________________________________________________________________________________________

// AddTokensFromDel adds tokens to a validator
func (v Validator) AddTokensFromDel(pool Pool, amount int64) (Validator, Pool, sdk.Dec) {

	// bondedShare/delegatedShare
	exRate := v.DelegatorShareExRate()
	amountDec := sdk.NewDec(amount)

	if v.Status == sdk.Bonded {
		pool = pool.looseTokensToBonded(amountDec)
	}

	v.Tokens = v.Tokens.Add(amountDec)
	issuedShares := amountDec.Quo(exRate)
	v.DelegatorShares = v.DelegatorShares.Add(issuedShares)

	return v, pool, issuedShares
}

// RemoveDelShares removes delegator shares from a validator.
func (v Validator) RemoveDelShares(pool Pool, delShares sdk.Dec) (Validator, Pool, sdk.Dec) {
	issuedTokens := v.DelegatorShareExRate().Mul(delShares)
	v.Tokens = v.Tokens.Sub(issuedTokens)
	v.DelegatorShares = v.DelegatorShares.Sub(delShares)

	if v.Status == sdk.Bonded {
		pool = pool.bondedTokensToLoose(issuedTokens)
	}

	return v, pool, issuedTokens
}

// DelegatorShareExRate gets the exchange rate of tokens over delegator shares.
// UNITS: tokens/delegator-shares
func (v Validator) DelegatorShareExRate() sdk.Dec {
	if v.DelegatorShares.IsZero() {
		return sdk.OneDec()
	}
	return v.Tokens.Quo(v.DelegatorShares)
}

// Get the bonded tokens which the validator holds
func (v Validator) BondedTokens() sdk.Dec {
	if v.Status == sdk.Bonded {
		return v.Tokens
	}
	return sdk.ZeroDec()
}

//______________________________________________________________________

// ensure fulfills the sdk validator types
var _ sdk.Validator = Validator{}

// nolint - for sdk.Validator
func (v Validator) GetRevoked() bool            { return v.Revoked }
func (v Validator) GetMoniker() string          { return v.Description.Moniker }
func (v Validator) GetStatus() sdk.BondStatus   { return v.Status }
func (v Validator) GetOperator() sdk.AccAddress { return v.Operator }
func (v Validator) GetPubKey() crypto.PubKey    { return v.PubKey }
func (v Validator) GetPower() sdk.Dec           { return v.BondedTokens() }
func (v Validator) GetTokens() sdk.Dec          { return v.Tokens }
func (v Validator) GetDelegatorShares() sdk.Dec { return v.DelegatorShares }
func (v Validator) GetBondHeight() int64        { return v.BondHeight }
