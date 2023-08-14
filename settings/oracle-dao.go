package settings

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/dao/trustednode"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

const (
	// Members
	quorumSettingPath                 = "members.quorum"
	rplBondSettingPath                = "members.rplbond"
	minipoolUnbondedMaxSettingPath    = "members.minipool.unbonded.max"
	minipoolUnbondedMinFeeSettingPath = "members.minipool.unbonded.min.fee"
	challengeCooldownSettingPath      = "members.challenge.cooldown"
	challengeWindowSettingPath        = "members.challenge.window"
	challengeCostSettingPath          = "members.challenge.cost"

	// Minipools
	scrubPeriodPath               = "minipool.scrub.period"
	promotionScrubPeriodPath      = "minipool.promotion.scrub.period"
	scrubPenaltyEnabledPath       = "minipool.scrub.penalty.enabled"
	bondReductionWindowStartPath  = "minipool.bond.reduction.window.start"
	bondReductionWindowLengthPath = "minipool.bond.reduction.window.length"

	// Proposals
	proposalCooldownTimeSettingPath = "proposal.cooldown.time"
	voteTimeSettingPath             = "proposal.vote.time"
	voteDelayTimeSettingPath        = "proposal.vote.delay.time"
	proposalExecuteTimeSettingPath  = "proposal.execute.time"
	proposalActionTimeSettingPath   = "proposal.action.time"
)

// ===============
// === Structs ===
// ===============

// Binding for Oracle DAO settings
type OracleDaoSettings struct {
	Details           OracleDaoSettingsDetails
	MembersContract   *core.Contract
	MinipoolContract  *core.Contract
	ProposalsContract *core.Contract
	RewardsContract   *core.Contract

	rp                              *rocketpool.RocketPool
	daoNodeTrustedContract          *trustednode.DaoNodeTrusted
	daoNodeTrustedProposalsContract *trustednode.DaoNodeTrustedProposals
}

// Details for Oracle DAO settings
type OracleDaoSettingsDetails struct {
	// Members
	Members struct {
		Quorum                 core.Parameter[float64] `json:"quorum"`
		RplBond                *big.Int                `json:"rplBond"`
		UnbondedMinipoolMax    core.Parameter[uint64]  `json:"unbondedMinipoolMax"`
		UnbondedMinipoolMinFee core.Parameter[float64] `json:"unbondedMinipoolMinFee"`
		ChallengeCooldown      core.Parameter[uint64]  `json:"challengeCooldown"`
		ChallengeWindow        core.Parameter[uint64]  `json:"challengeWindow"`
		ChallengeCost          *big.Int                `json:"challengeCost"`
	} `json:"members"`

	// Minipools
	Minipools struct {
		ScrubPeriod               core.Parameter[time.Duration] `json:"scrubPeriod"`
		PromotionScrubPeriod      core.Parameter[time.Duration] `json:"promotionScrubPeriod"`
		IsScrubPenaltyEnabled     bool                          `json:"isScrubPenaltyEnabled"`
		BondReductionWindowStart  core.Parameter[time.Duration] `json:"bondReductionWindowStart"`
		BondReductionWindowLength core.Parameter[time.Duration] `json:"bondReductionWindowLength"`
	} `json:"minipools"`

	// Proposals
	Proposals struct {
		CooldownTime  core.Parameter[time.Duration] `json:"cooldownTime"`
		VoteTime      core.Parameter[time.Duration] `json:"voteTime"`
		VoteDelayTime core.Parameter[time.Duration] `json:"voteDelayTime"`
		ExecuteTime   core.Parameter[time.Duration] `json:"executeTime"`
		ActionTime    core.Parameter[time.Duration] `json:"actionTime"`
	} `json:"proposals"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new Oracle DAO settings binding
func NewOracleDaoSettings(rp *rocketpool.RocketPool) (*OracleDaoSettings, error) {
	daoNodeTrustedContract, err := trustednode.NewDaoNodeTrusted(rp)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO node trusted binding: %w", err)
	}
	daoNodeTrustedProposalsContract, err := trustednode.NewDaoNodeTrustedProposals(rp)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO node trusted proposals binding: %w", err)
	}

	// Get the contracts
	contracts, err := rp.GetContracts([]rocketpool.ContractName{
		rocketpool.ContractName_RocketDAONodeTrustedSettingsMembers,
		rocketpool.ContractName_RocketDAONodeTrustedSettingsMinipool,
		rocketpool.ContractName_RocketDAONodeTrustedSettingsProposals,
		rocketpool.ContractName_RocketDAONodeTrustedSettingsRewards,
	}...)
	if err != nil {
		return nil, fmt.Errorf("error getting Oracle DAO settings contracts: %w", err)
	}

	return &OracleDaoSettings{
		Details:                         OracleDaoSettingsDetails{},
		rp:                              rp,
		daoNodeTrustedContract:          daoNodeTrustedContract,
		daoNodeTrustedProposalsContract: daoNodeTrustedProposalsContract,

		MembersContract:   contracts[0],
		MinipoolContract:  contracts[1],
		ProposalsContract: contracts[2],
		RewardsContract:   contracts[3],
	}, nil
}

// =============
// === Calls ===
// =============

// === RocketDAONodeTrustedSettingsMembers ===

// Get the member proposal quorum threshold
func (c *OracleDaoSettings) GetQuorum(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MembersContract, &c.Details.Members.Quorum.RawValue, "getQuorum")
}

// Get the RPL bond required for a member
func (c *OracleDaoSettings) GetRplBond(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MembersContract, &c.Details.Members.RplBond, "getRPLBond")
}

// Get the maximum number of unbonded minipools a member can run
func (c *OracleDaoSettings) GetUnbondedMinipoolMax(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MembersContract, &c.Details.Members.UnbondedMinipoolMax.RawValue, "getMinipoolUnbondedMax")
}

// Get the minimum commission rate before unbonded minipools are allowed
func (c *OracleDaoSettings) GetUnbondedMinipoolMinFee(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MembersContract, &c.Details.Members.UnbondedMinipoolMinFee.RawValue, "getMinipoolUnbondedMinFee")
}

// Get the period a member must wait for before submitting another challenge, in blocks
func (c *OracleDaoSettings) GetChallengeCooldown(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MembersContract, &c.Details.Members.ChallengeCooldown.RawValue, "getChallengeCooldown")
}

// Get the period during which a member can respond to a challenge, in blocks
func (c *OracleDaoSettings) GetChallengeWindow(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MembersContract, &c.Details.Members.ChallengeWindow.RawValue, "getChallengeWindow")
}

// Get the fee for a non-member to challenge a member, in wei
func (c *OracleDaoSettings) GetChallengeCost(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MembersContract, &c.Details.Members.ChallengeCost, "getChallengeCost")
}

// === RocketDAONodeTrustedSettingsMinipool ===

// Get the amount of time, in seconds, the scrub check lasts before a minipool can move from prelaunch to staking
func (c *OracleDaoSettings) GetScrubPeriod(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MinipoolContract, &c.Details.Minipools.ScrubPeriod.RawValue, "getScrubPeriod")
}

// Get the amount of time, in seconds, the promotion scrub check lasts before a vacant minipool can be promoted
func (c *OracleDaoSettings) GetPromotionScrubPeriod(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MinipoolContract, &c.Details.Minipools.PromotionScrubPeriod.RawValue, "getPromotionScrubPeriod")
}

// Check if the RPL slashing penalty is applied to scrubbed minipools
func (c *OracleDaoSettings) GetScrubPenaltyEnabled(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MinipoolContract, &c.Details.Minipools.IsScrubPenaltyEnabled, "getScrubPenaltyEnabled")
}

// Get the amount of time, in seconds, a minipool must wait after beginning a bond reduction before it can apply the bond reduction (how long the Oracle DAO has to cancel the reduction if required)
func (c *OracleDaoSettings) GetBondReductionWindowStart(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MinipoolContract, &c.Details.Minipools.BondReductionWindowStart.RawValue, "getBondReductionWindowStart")
}

// Get the amount of time, in seconds, a minipool has to reduce its bond once it has passed the check window
func (c *OracleDaoSettings) GetBondReductionWindowLength(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.MinipoolContract, &c.Details.Minipools.BondReductionWindowLength.RawValue, "getBondReductionWindowLength")
}

// === RocketDAONodeTrustedSettingsProposals ===

// Get the cooldown period a member must wait, in seconds, after making a proposal before making another
func (c *OracleDaoSettings) GetProposalCooldownTime(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.ProposalsContract, &c.Details.Proposals.CooldownTime.RawValue, "getCooldownTime")
}

// Get the period, in seconds, a proposal can be voted on
func (c *OracleDaoSettings) GetVoteTime(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.ProposalsContract, &c.Details.Proposals.VoteTime.RawValue, "getVoteTime")
}

// Get the delay, in seconds, after creation before a proposal can be voted on
func (c *OracleDaoSettings) GetVoteDelayTime(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.ProposalsContract, &c.Details.Proposals.VoteDelayTime.RawValue, "getVoteDelayTime")
}

// Get the period, in seconds, during which a passed proposal can be executed
func (c *OracleDaoSettings) GetProposalExecuteTime(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.ProposalsContract, &c.Details.Proposals.ExecuteTime.RawValue, "getExecuteTime")
}

// Get the period, in seconds, during which an action can be performed on an executed proposal
func (c *OracleDaoSettings) GetProposalActionTime(mc *multicall.MultiCaller) {
	multicall.AddCall(mc, c.ProposalsContract, &c.Details.Proposals.ActionTime.RawValue, "getActionTime")
}

// === RocketDAONodeTrustedSettingsRewards ===

// Get whether or not the provided rewards network is enabled
func (c *OracleDaoSettings) GetNetworkEnabled(mc *multicall.MultiCaller, enabled_Out *bool, network uint64) {
	multicall.AddCall(mc, c.RewardsContract, enabled_Out, "getNetworkEnabled", big.NewInt(0).SetUint64(network))
}

// == Meta ==

// Get all basic details
func (c *OracleDaoSettings) GetAllDetails(mc *multicall.MultiCaller) {
	// Members
	c.GetQuorum(mc)
	c.GetRplBond(mc)
	c.GetUnbondedMinipoolMax(mc)
	c.GetUnbondedMinipoolMinFee(mc)
	c.GetChallengeCooldown(mc)
	c.GetChallengeWindow(mc)
	c.GetChallengeCost(mc)

	// Minipools
	c.GetScrubPeriod(mc)
	c.GetPromotionScrubPeriod(mc)
	c.GetScrubPenaltyEnabled(mc)
	c.GetBondReductionWindowStart(mc)
	c.GetBondReductionWindowLength(mc)

	// Proposals
	c.GetProposalCooldownTime(mc)
	c.GetVoteTime(mc)
	c.GetVoteDelayTime(mc)
	c.GetProposalExecuteTime(mc)
	c.GetProposalActionTime(mc)
}

// ====================
// === Transactions ===
// ====================

// === RocketDAONodeTrustedSettingsMembers ===

// Get info for setting the member proposal quorum threshold
func (c *OracleDaoSettings) BootstrapQuorum(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoNodeTrustedContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsMembers, quorumSettingPath, value, opts)
}

// Get info for setting the RPL bond required for a member
func (c *OracleDaoSettings) BootstrapRplBond(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoNodeTrustedContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsMembers, rplBondSettingPath, value, opts)
}

// Get info for setting the maximum number of unbonded minipools a member can run
func (c *OracleDaoSettings) BootstrapUnbondedMinipoolMax(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoNodeTrustedContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsMembers, minipoolUnbondedMaxSettingPath, value, opts)
}

// Get info for setting the minimum commission rate before unbonded minipools are allowed
func (c *OracleDaoSettings) BootstrapUnbondedMinipoolMinFee(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoNodeTrustedContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsMembers, minipoolUnbondedMinFeeSettingPath, value, opts)
}

// Get info for setting the period a member must wait for before submitting another challenge, in blocks
func (c *OracleDaoSettings) BootstrapChallengeCooldown(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoNodeTrustedContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsMembers, challengeCooldownSettingPath, value, opts)
}

// Get info for setting the period during which a member can respond to a challenge, in blocks
func (c *OracleDaoSettings) BootstrapChallengeWindow(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoNodeTrustedContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsMembers, challengeWindowSettingPath, value, opts)
}

// Get info for setting the fee for a non-member to challenge a member, in wei
func (c *OracleDaoSettings) BootstrapChallengeCost(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoNodeTrustedContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsMembers, challengeCostSettingPath, value, opts)
}

// Get info for setting the member proposal quorum threshold
func (c *OracleDaoSettings) ProposeQuorum(value float64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return proposeSetValue(c.daoNodeTrustedProposalsContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsMembers, quorumSettingPath, value, opts)
}

// Get info for setting the RPL bond required for a member
func (c *OracleDaoSettings) ProposeRplBond(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return proposeSetValue(c.daoNodeTrustedProposalsContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsMembers, rplBondSettingPath, value, opts)
}

// Get info for setting the maximum number of unbonded minipools a member can run
func (c *OracleDaoSettings) ProposeUnbondedMinipoolMax(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return proposeSetValue(c.daoNodeTrustedProposalsContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsMembers, minipoolUnbondedMaxSettingPath, value, opts)
}

// Get info for setting the minimum commission rate before unbonded minipools are allowed
func (c *OracleDaoSettings) ProposeUnbondedMinipoolMinFee(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return proposeSetValue(c.daoNodeTrustedProposalsContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsMembers, minipoolUnbondedMinFeeSettingPath, value, opts)
}

// Get info for setting the period a member must wait for before submitting another challenge, in blocks
func (c *OracleDaoSettings) ProposeChallengeCooldown(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return proposeSetValue(c.daoNodeTrustedProposalsContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsMembers, challengeCooldownSettingPath, value, opts)
}

// Get info for setting the period during which a member can respond to a challenge, in blocks
func (c *OracleDaoSettings) ProposeChallengeWindow(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return proposeSetValue(c.daoNodeTrustedProposalsContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsMembers, challengeWindowSettingPath, value, opts)
}

// Get info for setting the fee for a non-member to challenge a member, in wei
func (c *OracleDaoSettings) ProposeChallengeCost(value *big.Int, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return proposeSetValue(c.daoNodeTrustedProposalsContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsMembers, challengeCostSettingPath, value, opts)
}

// === RocketDAONodeTrustedSettingsMinipool ===

// Get info for setting the amount of time, in seconds, the scrub check lasts before a minipool can move from prelaunch to staking
func (c *OracleDaoSettings) BootstrapScrubPeriod(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoNodeTrustedContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsMinipool, scrubPeriodPath, value, opts)
}

// Get info for setting the amount of time, in seconds, the promotion scrub check lasts before a vacant minipool can be promoted
func (c *OracleDaoSettings) BootstrapPromotionScrubPeriod(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoNodeTrustedContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsMinipool, promotionScrubPeriodPath, value, opts)
}

// Get info for setting the flag for the RPL slashing penalty on scrubbed minipools
func (c *OracleDaoSettings) BootstrapScrubPenaltyEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoNodeTrustedContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsMinipool, scrubPenaltyEnabledPath, value, opts)
}

// Get info for setting the amount of time, in seconds, a minipool must wait after beginning a bond reduction before it can apply the bond reduction (how long the Oracle DAO has to cancel the reduction if required)
func (c *OracleDaoSettings) BootstrapBondReductionWindowStart(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoNodeTrustedContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsMinipool, bondReductionWindowStartPath, value, opts)
}

// Get info for setting the amount of time, in seconds, a minipool has to reduce its bond once it has passed the check window
func (c *OracleDaoSettings) BootstrapBondReductionWindowLength(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoNodeTrustedContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsMinipool, bondReductionWindowLengthPath, value, opts)
}

// Get info for setting the amount of time, in seconds, the scrub check lasts before a minipool can move from prelaunch to staking
func (c *OracleDaoSettings) ProposeScrubPeriod(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return proposeSetValue(c.daoNodeTrustedProposalsContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsMinipool, scrubPeriodPath, value, opts)
}

// Get info for setting the amount of time, in seconds, the promotion scrub check lasts before a vacant minipool can be promoted
func (c *OracleDaoSettings) ProposePromotionScrubPeriod(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return proposeSetValue(c.daoNodeTrustedProposalsContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsMinipool, promotionScrubPeriodPath, value, opts)
}

// Get info for setting the flag for the RPL slashing penalty on scrubbed minipools
func (c *OracleDaoSettings) ProposeScrubPenaltyEnabled(value bool, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return proposeSetValue(c.daoNodeTrustedProposalsContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsMinipool, scrubPenaltyEnabledPath, value, opts)
}

// Get info for setting the amount of time, in seconds, a minipool must wait after beginning a bond reduction before it can apply the bond reduction (how long the Oracle DAO has to cancel the reduction if required)
func (c *OracleDaoSettings) ProposeBondReductionWindowStart(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return proposeSetValue(c.daoNodeTrustedProposalsContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsMinipool, bondReductionWindowStartPath, value, opts)
}

// Get info for setting the amount of time, in seconds, a minipool has to reduce its bond once it has passed the check window
func (c *OracleDaoSettings) ProposeBondReductionWindowLength(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return proposeSetValue(c.daoNodeTrustedProposalsContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsMinipool, bondReductionWindowLengthPath, value, opts)
}

// === RocketDAONodeTrustedSettingsProposals ===

// Get info for setting the cooldown period a member must wait, in seconds, after making a proposal before making another
func (c *OracleDaoSettings) BootstrapProposalCooldownTime(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoNodeTrustedContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsProposals, proposalCooldownTimeSettingPath, value, opts)
}

// Get info for setting the period, in seconds, a proposal can be voted on
func (c *OracleDaoSettings) BootstrapVoteTime(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoNodeTrustedContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsProposals, voteTimeSettingPath, value, opts)
}

// Get info for setting the delay, in seconds, after creation before a proposal can be voted on
func (c *OracleDaoSettings) BootstrapVoteDelayTime(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoNodeTrustedContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsProposals, voteDelayTimeSettingPath, value, opts)
}

// Get info for setting the period, in seconds, during which a passed proposal can be executed
func (c *OracleDaoSettings) BootstrapProposalExecuteTime(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoNodeTrustedContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsProposals, proposalExecuteTimeSettingPath, value, opts)
}

// Get info for setting the period, in seconds, during which an action can be performed on an executed proposal
func (c *OracleDaoSettings) BootstrapProposalActionTime(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return bootstrapValue(c.daoNodeTrustedContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsProposals, proposalActionTimeSettingPath, value, opts)
}

// Get info for setting the cooldown period a member must wait, in seconds, after making a proposal before making another
func (c *OracleDaoSettings) ProposeProposalCooldownTime(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return proposeSetValue(c.daoNodeTrustedProposalsContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsProposals, proposalCooldownTimeSettingPath, value, opts)
}

// Get info for setting the period, in seconds, a proposal can be voted on
func (c *OracleDaoSettings) ProposeVoteTime(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return proposeSetValue(c.daoNodeTrustedProposalsContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsProposals, voteTimeSettingPath, value, opts)
}

// Get info for setting the delay, in seconds, after creation before a proposal can be voted on
func (c *OracleDaoSettings) ProposeVoteDelayTime(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return proposeSetValue(c.daoNodeTrustedProposalsContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsProposals, voteDelayTimeSettingPath, value, opts)
}

// Get info for setting the period, in seconds, during which a passed proposal can be executed
func (c *OracleDaoSettings) ProposeProposalExecuteTime(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return proposeSetValue(c.daoNodeTrustedProposalsContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsProposals, proposalExecuteTimeSettingPath, value, opts)
}

// Get info for setting the period, in seconds, during which an action can be performed on an executed proposal
func (c *OracleDaoSettings) ProposeProposalActionTime(value uint64, opts *bind.TransactOpts) (*core.TransactionInfo, error) {
	return proposeSetValue(c.daoNodeTrustedProposalsContract, rocketpool.ContractName_RocketDAONodeTrustedSettingsProposals, proposalActionTimeSettingPath, value, opts)
}
