package stake

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	addrDels = []sdk.Address{
		addrs[0],
		addrs[1],
	}
	addrVals = []sdk.Address{
		addrs[2],
		addrs[3],
		addrs[4],
		addrs[5],
		addrs[6],
	}
)

// This function tests GetValidator, GetValidatorsBonded, setValidator, removeValidator
func TestValidator(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 0)

	//construct the validators
	var validators [3]Validator
	amts := []int64{9, 8, 7}
	for i, amt := range amts {
		validators[i] = NewValidator(addrVals[i], pks[i], Description{})
		validators[i].BondedShares = sdk.NewRat(amt)
		validators[i].DelegatorShares = sdk.NewRat(amt)
	}

	// check the empty keeper first
	_, found := keeper.GetValidator(ctx, addrVals[0])
	assert.False(t, found)
	resCands := keeper.GetValidatorsBonded(ctx)
	assert.Zero(t, len(resCands))

	// set and retrieve a record
	keeper.setValidator(ctx, validators[0])
	resCand, found := keeper.GetValidator(ctx, addrVals[0])
	require.True(t, found)
	assert.True(t, validators[0].equal(resCand), "%v \n %v", resCand, validators[0])

	// modify a records, save, and retrieve
	validators[0].DelegatorShares = sdk.NewRat(99)
	keeper.setValidator(ctx, validators[0])
	resCand, found = keeper.GetValidator(ctx, addrVals[0])
	require.True(t, found)
	assert.True(t, validators[0].equal(resCand))

	// also test that the address has been added to address list
	resCands = keeper.GetValidatorsBonded(ctx)
	require.Equal(t, 1, len(resCands))
	assert.Equal(t, addrVals[0], resCands[0].Address)

	// add other validators
	keeper.setValidator(ctx, validators[1])
	keeper.setValidator(ctx, validators[2])
	resCand, found = keeper.GetValidator(ctx, addrVals[1])
	require.True(t, found)
	assert.True(t, validators[1].equal(resCand), "%v \n %v", resCand, validators[1])
	resCand, found = keeper.GetValidator(ctx, addrVals[2])
	require.True(t, found)
	assert.True(t, validators[2].equal(resCand), "%v \n %v", resCand, validators[2])
	resCands = keeper.GetValidatorsBonded(ctx)
	require.Equal(t, 3, len(resCands))
	assert.True(t, validators[0].equal(resCands[0]), "%v \n %v", resCands[0], validators[0])
	assert.True(t, validators[1].equal(resCands[1]), "%v \n %v", resCands[1], validators[1])
	assert.True(t, validators[2].equal(resCands[2]), "%v \n %v", resCands[2], validators[2])

	// remove a record
	keeper.removeValidator(ctx, validators[1].Address)
	_, found = keeper.GetValidator(ctx, addrVals[1])
	assert.False(t, found)
}

// tests GetDelegation, GetDelegations, SetDelegation, removeDelegation, GetBonds
func TestBond(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 0)

	//construct the validators
	amts := []int64{9, 8, 7}
	var validators [3]Validator
	for i, amt := range amts {
		validators[i] = NewValidator(addrVals[i], pks[i], Description{})
		validators[i].BondedShares = sdk.NewRat(amt)
		validators[i].DelegatorShares = sdk.NewRat(amt)
	}

	// first add a validators[0] to delegate too
	keeper.setValidator(ctx, validators[0])

	bond1to1 := Delegation{
		DelegatorAddr: addrDels[0],
		ValidatorAddr: addrVals[0],
		Shares:        sdk.NewRat(9),
	}

	// check the empty keeper first
	_, found := keeper.GetDelegation(ctx, addrDels[0], addrVals[0])
	assert.False(t, found)

	// set and retrieve a record
	keeper.setDelegation(ctx, bond1to1)
	resBond, found := keeper.GetDelegation(ctx, addrDels[0], addrVals[0])
	assert.True(t, found)
	assert.True(t, bond1to1.equal(resBond))

	// modify a records, save, and retrieve
	bond1to1.Shares = sdk.NewRat(99)
	keeper.setDelegation(ctx, bond1to1)
	resBond, found = keeper.GetDelegation(ctx, addrDels[0], addrVals[0])
	assert.True(t, found)
	assert.True(t, bond1to1.equal(resBond))

	// add some more records
	keeper.setValidator(ctx, validators[1])
	keeper.setValidator(ctx, validators[2])
	bond1to2 := Delegation{addrDels[0], addrVals[1], sdk.NewRat(9), 0}
	bond1to3 := Delegation{addrDels[0], addrVals[2], sdk.NewRat(9), 1}
	bond2to1 := Delegation{addrDels[1], addrVals[0], sdk.NewRat(9), 2}
	bond2to2 := Delegation{addrDels[1], addrVals[1], sdk.NewRat(9), 3}
	bond2to3 := Delegation{addrDels[1], addrVals[2], sdk.NewRat(9), 4}
	keeper.setDelegation(ctx, bond1to2)
	keeper.setDelegation(ctx, bond1to3)
	keeper.setDelegation(ctx, bond2to1)
	keeper.setDelegation(ctx, bond2to2)
	keeper.setDelegation(ctx, bond2to3)

	// test all bond retrieve capabilities
	resBonds := keeper.GetDelegations(ctx, addrDels[0], 5)
	require.Equal(t, 3, len(resBonds))
	assert.True(t, bond1to1.equal(resBonds[0]))
	assert.True(t, bond1to2.equal(resBonds[1]))
	assert.True(t, bond1to3.equal(resBonds[2]))
	resBonds = keeper.GetDelegations(ctx, addrDels[0], 3)
	require.Equal(t, 3, len(resBonds))
	resBonds = keeper.GetDelegations(ctx, addrDels[0], 2)
	require.Equal(t, 2, len(resBonds))
	resBonds = keeper.GetDelegations(ctx, addrDels[1], 5)
	require.Equal(t, 3, len(resBonds))
	assert.True(t, bond2to1.equal(resBonds[0]))
	assert.True(t, bond2to2.equal(resBonds[1]))
	assert.True(t, bond2to3.equal(resBonds[2]))
	allBonds := keeper.getBonds(ctx, 1000)
	require.Equal(t, 6, len(allBonds))
	assert.True(t, bond1to1.equal(allBonds[0]))
	assert.True(t, bond1to2.equal(allBonds[1]))
	assert.True(t, bond1to3.equal(allBonds[2]))
	assert.True(t, bond2to1.equal(allBonds[3]))
	assert.True(t, bond2to2.equal(allBonds[4]))
	assert.True(t, bond2to3.equal(allBonds[5]))

	// delete a record
	keeper.removeDelegation(ctx, bond2to3)
	_, found = keeper.GetDelegation(ctx, addrDels[1], addrVals[2])
	assert.False(t, found)
	resBonds = keeper.GetDelegations(ctx, addrDels[1], 5)
	require.Equal(t, 2, len(resBonds))
	assert.True(t, bond2to1.equal(resBonds[0]))
	assert.True(t, bond2to2.equal(resBonds[1]))

	// delete all the records from delegator 2
	keeper.removeDelegation(ctx, bond2to1)
	keeper.removeDelegation(ctx, bond2to2)
	_, found = keeper.GetDelegation(ctx, addrDels[1], addrVals[0])
	assert.False(t, found)
	_, found = keeper.GetDelegation(ctx, addrDels[1], addrVals[1])
	assert.False(t, found)
	resBonds = keeper.GetDelegations(ctx, addrDels[1], 5)
	require.Equal(t, 0, len(resBonds))
}

// TODO seperate out into multiple tests
func TestGetValidators(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 0)

	// initialize some validators into the state
	amts := []int64{0, 100, 1, 400, 200}
	n := len(amts)
	var validators [5]Validator
	for i, amt := range amts {
		validators[i] = NewValidator(addrs[i], pks[i], Description{})
		validators[i].BondedShares = sdk.NewRat(amt)
		validators[i].DelegatorShares = sdk.NewRat(amt)
		keeper.setValidator(ctx, validators[i])
	}

	// first make sure everything made it in to the gotValidator group
	gotValidators := keeper.GetValidatorsBondedByPower(ctx)
	require.Equal(t, len(gotValidators), n)
	assert.Equal(t, sdk.NewRat(400), gotValidators[0].BondedShares, "%v", gotValidators)
	assert.Equal(t, sdk.NewRat(200), gotValidators[1].BondedShares, "%v", gotValidators)
	assert.Equal(t, sdk.NewRat(100), gotValidators[2].BondedShares, "%v", gotValidators)
	assert.Equal(t, sdk.NewRat(1), gotValidators[3].BondedShares, "%v", gotValidators)
	assert.Equal(t, sdk.NewRat(0), gotValidators[4].BondedShares, "%v", gotValidators)
	assert.Equal(t, validators[3].Address, gotValidators[0].Address, "%v", gotValidators)
	assert.Equal(t, validators[4].Address, gotValidators[1].Address, "%v", gotValidators)
	assert.Equal(t, validators[1].Address, gotValidators[2].Address, "%v", gotValidators)
	assert.Equal(t, validators[2].Address, gotValidators[3].Address, "%v", gotValidators)
	assert.Equal(t, validators[0].Address, gotValidators[4].Address, "%v", gotValidators)

	// test a basic increase in voting power
	validators[3].BondedShares = sdk.NewRat(500)
	keeper.setValidator(ctx, validators[3])
	gotValidators = keeper.GetValidatorsBondedByPower(ctx)
	require.Equal(t, len(gotValidators), n)
	assert.Equal(t, sdk.NewRat(500), gotValidators[0].BondedShares, "%v", gotValidators)
	assert.Equal(t, validators[3].Address, gotValidators[0].Address, "%v", gotValidators)

	// test a decrease in voting power
	validators[3].BondedShares = sdk.NewRat(300)
	keeper.setValidator(ctx, validators[3])
	gotValidators = keeper.GetValidatorsBondedByPower(ctx)
	require.Equal(t, len(gotValidators), n)
	assert.Equal(t, sdk.NewRat(300), gotValidators[0].BondedShares, "%v", gotValidators)
	assert.Equal(t, validators[3].Address, gotValidators[0].Address, "%v", gotValidators)
	assert.Equal(t, validators[4].Address, gotValidators[1].Address, "%v", gotValidators)

	// test equal voting power, different age
	validators[3].BondedShares = sdk.NewRat(200)
	ctx = ctx.WithBlockHeight(10)
	keeper.setValidator(ctx, validators[3])
	gotValidators = keeper.GetValidatorsBondedByPower(ctx)
	require.Equal(t, len(gotValidators), n)
	assert.Equal(t, sdk.NewRat(200), gotValidators[0].BondedShares, "%v", gotValidators)
	assert.Equal(t, sdk.NewRat(200), gotValidators[1].BondedShares, "%v", gotValidators)
	assert.Equal(t, validators[3].Address, gotValidators[0].Address, "%v", gotValidators)
	assert.Equal(t, validators[4].Address, gotValidators[1].Address, "%v", gotValidators)
	assert.Equal(t, int64(0), gotValidators[0].BondHeight, "%v", gotValidators)
	assert.Equal(t, int64(0), gotValidators[1].BondHeight, "%v", gotValidators)

	// no change in voting power - no change in sort
	ctx = ctx.WithBlockHeight(20)
	keeper.setValidator(ctx, validators[4])
	gotValidators = keeper.GetValidatorsBondedByPower(ctx)
	require.Equal(t, len(gotValidators), n)
	assert.Equal(t, validators[3].Address, gotValidators[0].Address, "%v", gotValidators)
	assert.Equal(t, validators[4].Address, gotValidators[1].Address, "%v", gotValidators)

	// change in voting power of both validators, both still in v-set, no age change
	validators[3].BondedShares = sdk.NewRat(300)
	validators[4].BondedShares = sdk.NewRat(300)
	keeper.setValidator(ctx, validators[3])
	gotValidators = keeper.GetValidatorsBondedByPower(ctx)
	require.Equal(t, len(gotValidators), n)
	ctx = ctx.WithBlockHeight(30)
	keeper.setValidator(ctx, validators[4])
	gotValidators = keeper.GetValidatorsBondedByPower(ctx)
	require.Equal(t, len(gotValidators), n, "%v", gotValidators)
	assert.Equal(t, validators[3].Address, gotValidators[0].Address, "%v", gotValidators)
	assert.Equal(t, validators[4].Address, gotValidators[1].Address, "%v", gotValidators)
}

// TODO seperate out into multiple tests
func TestGetValidatorsEdgeCases(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 0)

	// now 2 max gotValidators
	params := keeper.GetParams(ctx)
	params.MaxValidators = 2
	keeper.setParams(ctx, params)

	// initialize some validators into the state
	amts := []int64{0, 100, 400, 400, 200}
	n := len(amts)
	var validators [5]Validator
	for i, amt := range amts {
		validators[i] = NewValidator(addrs[i], pks[i], Description{})
		validators[i].BondedShares = sdk.NewRat(amt)
		validators[i].DelegatorShares = sdk.NewRat(amt)
		keeper.setValidator(ctx, validators[i])
	}

	validators[0].BondedShares = sdk.NewRat(500)
	keeper.setValidator(ctx, validators[0])
	gotValidators := keeper.GetValidatorsBondedByPower(ctx)
	require.Equal(t, uint16(len(gotValidators)), params.MaxValidators)
	require.Equal(t, validators[0].Address, gotValidators[0].Address, "%v", gotValidators)
	// validator 3 was set before validator 4
	require.Equal(t, validators[2].Address, gotValidators[1].Address, "%v", gotValidators)

	// A validator which leaves the gotValidator set due to a decrease in voting power,
	// then increases to the original voting power, does not get its spot back in the
	// case of a tie.
	// ref https://github.com/cosmos/cosmos-sdk/issues/582#issuecomment-380757108
	validators[3].BondedShares = sdk.NewRat(401)
	keeper.setValidator(ctx, validators[3])
	gotValidators = keeper.GetValidatorsBondedByPower(ctx)
	require.Equal(t, uint16(len(gotValidators)), params.MaxValidators)
	require.Equal(t, validators[0].Address, gotValidators[0].Address, "%v", gotValidators)
	require.Equal(t, validators[3].Address, gotValidators[1].Address, "%v", gotValidators)
	ctx = ctx.WithBlockHeight(40)
	// validator 3 kicked out temporarily
	validators[3].BondedShares = sdk.NewRat(200)
	keeper.setValidator(ctx, validators[3])
	gotValidators = keeper.GetValidatorsBondedByPower(ctx)
	require.Equal(t, uint16(len(gotValidators)), params.MaxValidators)
	require.Equal(t, validators[0].Address, gotValidators[0].Address, "%v", gotValidators)
	require.Equal(t, validators[2].Address, gotValidators[1].Address, "%v", gotValidators)
	// validator 4 does not get spot back
	validators[3].BondedShares = sdk.NewRat(400)
	keeper.setValidator(ctx, validators[3])
	gotValidators = keeper.GetValidatorsBondedByPower(ctx)
	require.Equal(t, uint16(len(gotValidators)), params.MaxValidators)
	require.Equal(t, validators[0].Address, gotValidators[0].Address, "%v", gotValidators)
	require.Equal(t, validators[2].Address, gotValidators[1].Address, "%v", gotValidators)
	validator, exists := keeper.GetValidator(ctx, validators[3].Address)
	require.Equal(t, exists, true)
	require.Equal(t, validator.BondHeight, int64(40))

	// If two validators both increase to the same voting power in the same block,
	// the one with the first transaction should take precedence (become a gotValidator).
	// ref https://github.com/cosmos/cosmos-sdk/issues/582#issuecomment-381250392
	validators[0].BondedShares = sdk.NewRat(2000)
	keeper.setValidator(ctx, validators[0])
	validators[1].BondedShares = sdk.NewRat(1000)
	validators[2].BondedShares = sdk.NewRat(1000)
	keeper.setValidator(ctx, validators[1])
	keeper.setValidator(ctx, validators[2])
	gotValidators = keeper.GetValidatorsBondedByPower(ctx)
	require.Equal(t, uint16(len(gotValidators)), params.MaxValidators)
	require.Equal(t, validators[0].Address, gotValidators[0].Address, "%v", gotValidators)
	require.Equal(t, validators[1].Address, gotValidators[1].Address, "%v", gotValidators)
	validators[1].BondedShares = sdk.NewRat(1100)
	validators[2].BondedShares = sdk.NewRat(1100)
	keeper.setValidator(ctx, validators[2])
	keeper.setValidator(ctx, validators[1])
	gotValidators = keeper.GetValidatorsBondedByPower(ctx)
	require.Equal(t, uint16(len(gotValidators)), params.MaxValidators)
	require.Equal(t, validators[0].Address, gotValidators[0].Address, "%v", gotValidators)
	require.Equal(t, validators[2].Address, gotValidators[1].Address, "%v", gotValidators)

	// reset assets / heights
	params.MaxValidators = 100
	keeper.setParams(ctx, params)
	validators[0].BondedShares = sdk.NewRat(0)
	validators[1].BondedShares = sdk.NewRat(100)
	validators[2].BondedShares = sdk.NewRat(1)
	validators[3].BondedShares = sdk.NewRat(300)
	validators[4].BondedShares = sdk.NewRat(200)
	ctx = ctx.WithBlockHeight(0)
	keeper.setValidator(ctx, validators[0])
	keeper.setValidator(ctx, validators[1])
	keeper.setValidator(ctx, validators[2])
	keeper.setValidator(ctx, validators[3])
	keeper.setValidator(ctx, validators[4])

	// test a swap in voting power
	validators[0].BondedShares = sdk.NewRat(600)
	keeper.setValidator(ctx, validators[0])
	gotValidators = keeper.GetValidatorsBondedByPower(ctx)
	require.Equal(t, len(gotValidators), n)
	assert.Equal(t, sdk.NewRat(600), gotValidators[0].BondedShares, "%v", gotValidators)
	assert.Equal(t, validators[0].Address, gotValidators[0].Address, "%v", gotValidators)
	assert.Equal(t, sdk.NewRat(300), gotValidators[1].BondedShares, "%v", gotValidators)
	assert.Equal(t, validators[3].Address, gotValidators[1].Address, "%v", gotValidators)

	// test the max gotValidators term
	params = keeper.GetParams(ctx)
	n = 2
	params.MaxValidators = uint16(n)
	keeper.setParams(ctx, params)
	gotValidators = keeper.GetValidatorsBondedByPower(ctx)
	require.Equal(t, len(gotValidators), n)
	assert.Equal(t, sdk.NewRat(600), gotValidators[0].BondedShares, "%v", gotValidators)
	assert.Equal(t, validators[0].Address, gotValidators[0].Address, "%v", gotValidators)
	assert.Equal(t, sdk.NewRat(300), gotValidators[1].BondedShares, "%v", gotValidators)
	assert.Equal(t, validators[3].Address, gotValidators[1].Address, "%v", gotValidators)
}

// clear the tracked changes to the gotValidator set
func TestClearTendermintUpdates(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 0)

	amts := []int64{100, 400, 200}
	validators := make([]Validator, len(amts))
	for i, amt := range amts {
		validators[i] = NewValidator(addrs[i], pks[i], Description{})
		validators[i].BondedShares = sdk.NewRat(amt)
		validators[i].DelegatorShares = sdk.NewRat(amt)
		keeper.setValidator(ctx, validators[i])
	}

	updates := keeper.getTendermintUpdates(ctx)
	assert.Equal(t, len(amts), len(updates))
	keeper.clearTendermintUpdates(ctx)
	updates = keeper.getTendermintUpdates(ctx)
	assert.Equal(t, 0, len(updates))
}

// test the mechanism which keeps track of a gotValidator set change
func TestGetTendermintUpdates(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 0)
	params := defaultParams()
	params.MaxValidators = 4
	keeper.setParams(ctx, params)

	// TODO eliminate use of validatorsIn here
	// tests could be clearer if they just
	// created the validator at time of use
	// and were labelled by power in the comments
	// outlining in each test
	amts := []int64{10, 11, 12, 13, 1}
	var validatorsIn [5]Validator
	for i, amt := range amts {
		validatorsIn[i] = NewValidator(addrs[i], pks[i], Description{})
		validatorsIn[i].BondedShares = sdk.NewRat(amt)
		validatorsIn[i].DelegatorShares = sdk.NewRat(amt)
	}

	// test from nothing to something
	//  validator set: {} -> {c1, c3}
	//  gotValidator set: {} -> {c1, c3}
	//  tendermintUpdate set: {} -> {c1, c3}
	assert.Equal(t, 0, len(keeper.GetValidatorsBonded(ctx))) // GetValidatorsBonded(ctx, 5
	assert.Equal(t, 0, len(keeper.GetValidatorsBonded(ctx)))
	assert.Equal(t, 0, len(keeper.getTendermintUpdates(ctx)))

	keeper.setValidator(ctx, validatorsIn[1])
	keeper.setValidator(ctx, validatorsIn[3])

	vals := keeper.GetValidatorsBondedByPower(ctx) // to init recent gotValidator set
	require.Equal(t, 2, len(vals))
	updates := keeper.getTendermintUpdates(ctx)
	require.Equal(t, 2, len(updates))
	validators := keeper.GetValidatorsBonded(ctx) //GetValidatorsBonded(ctx, 5
	require.Equal(t, 2, len(validators))
	assert.Equal(t, validators[0].abciValidator(keeper.cdc), updates[0])
	assert.Equal(t, validators[1].abciValidator(keeper.cdc), updates[1])
	assert.True(t, validators[0].equal(vals[1]))
	assert.True(t, validators[1].equal(vals[0]))

	// test identical,
	//  validator set: {c1, c3} -> {c1, c3}
	//  tendermintUpdate set: {} -> {}
	keeper.clearTendermintUpdates(ctx)
	assert.Equal(t, 2, len(keeper.GetValidatorsBonded(ctx)))
	assert.Equal(t, 0, len(keeper.getTendermintUpdates(ctx)))

	keeper.setValidator(ctx, validators[0])
	keeper.setValidator(ctx, validators[1])

	require.Equal(t, 2, len(keeper.GetValidatorsBonded(ctx)))
	assert.Equal(t, 0, len(keeper.getTendermintUpdates(ctx)))

	// test single value change
	//  validator set: {c1, c3} -> {c1', c3}
	//  tendermintUpdate set: {} -> {c1'}
	keeper.clearTendermintUpdates(ctx)
	assert.Equal(t, 2, len(keeper.GetValidatorsBonded(ctx)))
	assert.Equal(t, 0, len(keeper.getTendermintUpdates(ctx)))

	validators[0].BondedShares = sdk.NewRat(600)
	keeper.setValidator(ctx, validators[0])

	validators = keeper.GetValidatorsBonded(ctx)
	require.Equal(t, 2, len(validators))
	assert.True(t, validators[0].BondedShares.Equal(sdk.NewRat(600)))
	updates = keeper.getTendermintUpdates(ctx)
	require.Equal(t, 1, len(updates))
	assert.Equal(t, validators[0].abciValidator(keeper.cdc), updates[0])

	// test multiple value change
	//  validator set: {c1, c3} -> {c1', c3'}
	//  tendermintUpdate set: {c1, c3} -> {c1', c3'}
	keeper.clearTendermintUpdates(ctx)
	assert.Equal(t, 2, len(keeper.GetValidatorsBonded(ctx)))
	assert.Equal(t, 0, len(keeper.getTendermintUpdates(ctx)))

	validators[0].BondedShares = sdk.NewRat(200)
	validators[1].BondedShares = sdk.NewRat(100)
	keeper.setValidator(ctx, validators[0])
	keeper.setValidator(ctx, validators[1])

	updates = keeper.getTendermintUpdates(ctx)
	require.Equal(t, 2, len(updates))
	validators = keeper.GetValidatorsBonded(ctx)
	require.Equal(t, 2, len(validators))
	require.Equal(t, validators[0].abciValidator(keeper.cdc), updates[0])
	require.Equal(t, validators[1].abciValidator(keeper.cdc), updates[1])

	// test validtor added at the beginning
	//  validator set: {c1, c3} -> {c0, c1, c3}
	//  tendermintUpdate set: {} -> {c0}
	keeper.clearTendermintUpdates(ctx)
	assert.Equal(t, 2, len(keeper.GetValidatorsBonded(ctx)))
	assert.Equal(t, 0, len(keeper.getTendermintUpdates(ctx)))

	keeper.setValidator(ctx, validatorsIn[0])
	updates = keeper.getTendermintUpdates(ctx)
	require.Equal(t, 1, len(updates))
	validators = keeper.GetValidatorsBonded(ctx)
	require.Equal(t, 3, len(validators))
	assert.Equal(t, validators[0].abciValidator(keeper.cdc), updates[0])

	// test gotValidator added at the middle
	//  validator set: {c0, c1, c3} -> {c0, c1, c2, c3]
	//  tendermintUpdate set: {} -> {c2}
	keeper.clearTendermintUpdates(ctx)
	assert.Equal(t, 3, len(keeper.GetValidatorsBonded(ctx)))
	assert.Equal(t, 0, len(keeper.getTendermintUpdates(ctx)))

	keeper.setValidator(ctx, validatorsIn[2])
	updates = keeper.getTendermintUpdates(ctx)
	require.Equal(t, 1, len(updates))
	validators = keeper.GetValidatorsBonded(ctx)
	require.Equal(t, 4, len(validators))
	assert.Equal(t, validators[2].abciValidator(keeper.cdc), updates[0])

	// test validator added at the end but not inserted in the valset
	//  validator set: {c0, c1, c2, c3} -> {c0, c1, c2, c3, c4}
	//  gotValidator set: {c0, c1, c2, c3} -> {c0, c1, c2, c3}
	//  tendermintUpdate set: {} -> {}
	keeper.clearTendermintUpdates(ctx)
	assert.Equal(t, 4, len(keeper.GetValidatorsBonded(ctx)))
	assert.Equal(t, 4, len(keeper.GetValidatorsBonded(ctx)))
	assert.Equal(t, 0, len(keeper.getTendermintUpdates(ctx)))

	keeper.setValidator(ctx, validatorsIn[4])

	assert.Equal(t, 5, len(keeper.GetValidatorsBonded(ctx)))
	assert.Equal(t, 4, len(keeper.GetValidatorsBonded(ctx)))
	require.Equal(t, 0, len(keeper.getTendermintUpdates(ctx))) // max gotValidator number is 4

	// test validator change its power but still not in the valset
	//  validator set: {c0, c1, c2, c3, c4} -> {c0, c1, c2, c3, c4}
	//  gotValidator set: {c0, c1, c2, c3}     -> {c0, c1, c2, c3}
	//  tendermintUpdate set: {}     -> {}
	keeper.clearTendermintUpdates(ctx)
	assert.Equal(t, 5, len(keeper.GetValidatorsBonded(ctx)))
	assert.Equal(t, 4, len(keeper.GetValidatorsBonded(ctx)))
	assert.Equal(t, 0, len(keeper.getTendermintUpdates(ctx)))

	validatorsIn[4].BondedShares = sdk.NewRat(1)
	keeper.setValidator(ctx, validatorsIn[4])

	assert.Equal(t, 5, len(keeper.GetValidatorsBonded(ctx)))
	assert.Equal(t, 4, len(keeper.GetValidatorsBonded(ctx)))
	require.Equal(t, 0, len(keeper.getTendermintUpdates(ctx))) // max gotValidator number is 4

	// test validator change its power and become a gotValidator (pushing out an existing)
	//  validator set: {c0, c1, c2, c3, c4} -> {c0, c1, c2, c3, c4}
	//  gotValidator set: {c0, c1, c2, c3}     -> {c1, c2, c3, c4}
	//  tendermintUpdate set: {}     -> {c0, c4}
	keeper.clearTendermintUpdates(ctx)
	assert.Equal(t, 5, len(keeper.GetValidatorsBonded(ctx)))
	assert.Equal(t, 4, len(keeper.GetValidatorsBonded(ctx)))
	assert.Equal(t, 0, len(keeper.getTendermintUpdates(ctx)))

	validatorsIn[4].BondedShares = sdk.NewRat(1000)
	keeper.setValidator(ctx, validatorsIn[4])

	validators = keeper.GetValidatorsBonded(ctx)
	require.Equal(t, 5, len(validators))
	vals = keeper.GetValidatorsBondedByPower(ctx)
	require.Equal(t, 4, len(vals))
	assert.Equal(t, validatorsIn[1].Address, vals[1].Address)
	assert.Equal(t, validatorsIn[2].Address, vals[3].Address)
	assert.Equal(t, validatorsIn[3].Address, vals[2].Address)
	assert.Equal(t, validatorsIn[4].Address, vals[0].Address)

	updates = keeper.getTendermintUpdates(ctx)
	require.Equal(t, 2, len(updates), "%v", updates)

	assert.Equal(t, validatorsIn[0].PubKey.Bytes(), updates[0].PubKey)
	assert.Equal(t, int64(0), updates[0].Power)
	assert.Equal(t, vals[0].abciValidator(keeper.cdc), updates[1])

	// test from something to nothing
	//  validator set: {c0, c1, c2, c3, c4} -> {}
	//  gotValidator set: {c1, c2, c3, c4}  -> {}
	//  tendermintUpdate set: {} -> {c1, c2, c3, c4}
	keeper.clearTendermintUpdates(ctx)
	assert.Equal(t, 5, len(keeper.GetValidatorsBonded(ctx)))
	assert.Equal(t, 4, len(keeper.GetValidatorsBonded(ctx)))
	assert.Equal(t, 0, len(keeper.getTendermintUpdates(ctx)))

	keeper.removeValidator(ctx, validatorsIn[0].Address)
	keeper.removeValidator(ctx, validatorsIn[1].Address)
	keeper.removeValidator(ctx, validatorsIn[2].Address)
	keeper.removeValidator(ctx, validatorsIn[3].Address)
	keeper.removeValidator(ctx, validatorsIn[4].Address)

	vals = keeper.GetValidatorsBondedByPower(ctx)
	assert.Equal(t, 0, len(vals), "%v", vals)
	validators = keeper.GetValidatorsBonded(ctx)
	require.Equal(t, 0, len(validators))
	updates = keeper.getTendermintUpdates(ctx)
	require.Equal(t, 4, len(updates))
	assert.Equal(t, validatorsIn[1].PubKey.Bytes(), updates[0].PubKey)
	assert.Equal(t, validatorsIn[2].PubKey.Bytes(), updates[1].PubKey)
	assert.Equal(t, validatorsIn[3].PubKey.Bytes(), updates[2].PubKey)
	assert.Equal(t, validatorsIn[4].PubKey.Bytes(), updates[3].PubKey)
	assert.Equal(t, int64(0), updates[0].Power)
	assert.Equal(t, int64(0), updates[1].Power)
	assert.Equal(t, int64(0), updates[2].Power)
	assert.Equal(t, int64(0), updates[3].Power)
}

func TestParams(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 0)
	expParams := defaultParams()

	//check that the empty keeper loads the default
	resParams := keeper.GetParams(ctx)
	assert.True(t, expParams.equal(resParams))

	//modify a params, save, and retrieve
	expParams.MaxValidators = 777
	keeper.setParams(ctx, expParams)
	resParams = keeper.GetParams(ctx)
	assert.True(t, expParams.equal(resParams))
}

func TestPool(t *testing.T) {
	ctx, _, keeper := createTestInput(t, false, 0)
	expPool := initialPool()

	//check that the empty keeper loads the default
	resPool := keeper.GetPool(ctx)
	assert.True(t, expPool.equal(resPool))

	//modify a params, save, and retrieve
	expPool.TotalSupply = 777
	keeper.setPool(ctx, expPool)
	resPool = keeper.GetPool(ctx)
	assert.True(t, expPool.equal(resPool))
}
