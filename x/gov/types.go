package gov

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	YesOption        = "Yes"
	NoOption         = "No"
	NoWithVetoOption = "NoWithVeto"
	AbstainOption    = "Abstain"
)

type Vote struct {
	Voter      sdk.Address `json:"voter"`       //  address of the voter
	ProposalID int64       `json:"proposal_id"` //  proposalID of the proposal
	Option     string      `json:"option"`      //  option from OptionSet chosen by the voter
	Weight     int64       `json:"weight"`      //  weight of the Vote
}

func NewVote(voter sdk.Address, proposalID int64, option string, weight int64) Vote {
	return Vote{
		voter, proposalID, option, weight,
	}
}

//-----------------------------------------------------------

// Proposal
type Proposal struct {
	ProposalID   int64     `json:"proposal_id"`   //  ID of the proposal
	Title        string    `json:"title"`         //  Title of the proposal
	Description  string    `json:"description"`   //  Description of the proposal
	ProposalType string    `json:"proposal_type"` //  Type of proposal. Initial set {PlainTextProposal, SoftwareUpgradeProposal}
	Procedure    Procedure `json:"procedure"`     //  Governance Procedure that the proposal follows proposal

	SubmitBlock  int64     `json:"submit_block"`  //  Height of the block where TxGovSubmitProposal was included
	TotalDeposit sdk.Coins `json:"total_deposit"` //  Current deposit on this proposal. Initial value is set at InitialDeposit
	Deposits     []Deposit `json:"deposits"`      //  Current deposit on this proposal. Initial value is set at InitialDeposit

	VotingStartBlock int64 `json:"voting_start_block"` //  Height of the block where MinDeposit was reached. -1 if MinDeposit is not reached

	ValidatorGovInfos []ValidatorGovInfo `json:"validator_gov_infos"` //  Total voting power when proposal enters voting period (default 0)
	VoteList          []Vote             `json:"vote_list"`           //  Total votes for each option

	TotalVotingPower int64 `json:"total_voting_power"` //  The Total Voting Power
	YesVotes         int64 `json:"yes_votes"`          //  Weight of Yes Votes
	NoVotes          int64 `json:"no_votes"`           //  Weight of No Votes
	NoWithVetoVotes  int64 `json:"no_with_veto_votes"` //  Weight of NoWithVeto Votes
	AbstainVotes     int64 `json:"abstain_votes"`      //  Weight of Abstain Votes
}

func (proposal *Proposal) getValidatorGovInfo(validatorAddr sdk.Address) *ValidatorGovInfo {
	for i, validatorGovInfo := range proposal.ValidatorGovInfos {
		if bytes.Equal(validatorGovInfo.ValidatorAddr, validatorAddr) {
			return &proposal.ValidatorGovInfos[i]
		}
	}
	return nil
}

func (proposal *Proposal) getVote(voterAddr sdk.Address) *Vote {
	for i, vote := range proposal.VoteList {
		if bytes.Equal(vote.Voter, voterAddr) {
			return &proposal.VoteList[i]
		}
	}
	return nil
}

func (proposal Proposal) isActive() bool {
	return proposal.VotingStartBlock >= 0
}

func (proposal Proposal) isExpired(height int64) bool {
	return height > proposal.VotingStartBlock + proposal.Procedure.VotingPeriod
}

func (proposal *Proposal) updateTally(option string, amount int64) {
	switch option {
	case YesOption:
		proposal.YesVotes += amount
	case NoOption:
		proposal.NoVotes += amount
	case NoWithVetoOption:
		proposal.NoWithVetoVotes += amount
	case AbstainOption:
		proposal.AbstainVotes += amount
	}
}

// Procedure
type Procedure struct {
	VotingPeriod      int64     `json:"voting_period"`      //  Length of the voting period. Initial value: 2 weeks
	MinDeposit        sdk.Coins `json:"min_deposit"`        //  Minimum deposit for a proposal to enter voting period.
	ProposalTypes     []string  `json:"proposal_type"`      //  Types available to submitters. {PlainTextProposal, SoftwareUpgradeProposal}
	Threshold         sdk.Rat   `json:"threshold"`          //  Minimum propotion of Yes votes for proposal to pass. Initial value: 0.5
	Veto              sdk.Rat   `json:"veto"`               //  Minimum value of Veto votes to Total votes ratio for proposal to be vetoed. Initial value: 1/3
	FastPass          sdk.Rat   `json:"fast_pass"`          //  Minimum propotion of Yes votes for proposal to pass. Initial value: 0.5
	MaxDepositPeriod  int64     `json:"max_deposit_period"` //  Maximum period for Atom holders to deposit on a proposal. Initial value: 2 months
	GovernancePenalty sdk.Rat   `json:"governance_penalty"` //  Penalty if validator does not vote
}

func (procedure Procedure) validProposalType(proposalType string) bool {
	for _, p := range procedure.ProposalTypes {
		if p == proposalType {
			return true
		}
	}
	return false
}

// Deposit
type Deposit struct {
	Depositer sdk.Address `json:"depositer"` //  Address of the depositer
	Amount    sdk.Coins   `json:"amount"`    //  Deposit amount
}

type ValidatorGovInfo struct {
	ProposalID      int64       `json:"proposal_iD"`		//  Id of the Proposal this validator
	ValidatorAddr   sdk.Address `json:"validator_addr"`		//  Address of the validator
	InitVotingPower int64       `json:"init_voting_power"`	//  Voting power of validator when proposal enters voting period
	Minus           int64       `json:"minus"`				//  Minus of validator, used to compute validator's voting power
	LastVoteWeight  int64       `json:"last_vote_weight"`	//  Weight of the last vote by validator at time of casting, -1 if hasn't voted yet
}

type Delegation struct {
	Amount    int64
	Validator sdk.Address
}

type ProposalQueue []int64