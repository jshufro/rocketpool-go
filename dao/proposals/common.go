package proposals

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	batch "github.com/rocket-pool/batch-query"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/types"
)

// ===============
// === Structs ===
// ===============

// Binding for proposals
type proposalCommon struct {
	*proposalCommonDetails
	rp       *rocketpool.RocketPool
	contract *core.Contract
}

// Details for proposals
type proposalCommonDetails struct {
	ID              core.Parameter[uint64]                   `json:"id"`
	ProposerAddress common.Address                           `json:"proposerAddress"`
	Message         string                                   `json:"message"`
	CreatedTime     core.Parameter[time.Time]                `json:"createdTime"`
	StartTime       core.Parameter[time.Time]                `json:"startTime"`
	EndTime         core.Parameter[time.Time]                `json:"endTime"`
	ExpiryTime      core.Parameter[time.Time]                `json:"expiryTime"`
	VotesRequired   core.Parameter[float64]                  `json:"votesRequired"`
	VotesFor        core.Parameter[float64]                  `json:"votesFor"`
	VotesAgainst    core.Parameter[float64]                  `json:"votesAgainst"`
	MemberVoted     bool                                     `json:"memberVoted"`
	MemberSupported bool                                     `json:"memberSupported"`
	IsCancelled     bool                                     `json:"isCancelled"`
	IsExecuted      bool                                     `json:"isExecuted"`
	Payload         []byte                                   `json:"payload"`
	State           core.Uint8Parameter[types.ProposalState] `json:"state"`
}

// ====================
// === Constructors ===
// ====================

// Creates a new ProposalCommon contract binding
func newProposalCommon(rp *rocketpool.RocketPool, id uint64) (*proposalCommon, error) {
	// Create the contract
	contract, err := rp.GetContract(rocketpool.ContractName_RocketDAOProposal)
	if err != nil {
		return nil, fmt.Errorf("error getting DAO proposal contract: %w", err)
	}

	return &proposalCommon{
		proposalCommonDetails: &proposalCommonDetails{
			ID: core.Parameter[uint64]{
				RawValue: big.NewInt(0).SetUint64(id),
			},
		},
		rp:       rp,
		contract: contract,
	}, nil
}

// =============
// === Calls ===
// =============

// Get the address of the node that created the proposal
func (c *proposalCommon) GetProposerAddress(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.ProposerAddress, "getProposer", c.ID.RawValue)
}

// Get the message provided with the proposal
func (c *proposalCommon) GetMessage(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Message, "getMessage", c.ID.RawValue)
}

// Get the time the proposal was created
func (c *proposalCommon) GetCreatedTime(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.CreatedTime.RawValue, "getCreated", c.ID.RawValue)
}

// Get the time the voting window on the proposal started
func (c *proposalCommon) GetStartTime(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.StartTime.RawValue, "getStart", c.ID.RawValue)
}

// Get the time the voting window on the proposal ended
func (c *proposalCommon) GetEndTime(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.EndTime.RawValue, "getEnd", c.ID.RawValue)
}

// Get the time the proposal expires
func (c *proposalCommon) GetExpiryTime(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.ExpiryTime.RawValue, "getExpires", c.ID.RawValue)
}

// Get the number of votes required for the proposal to pass
func (c *proposalCommon) GetVotesRequired(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.VotesRequired.RawValue, "getVotesRequired", c.ID.RawValue)
}

// Get the number of votes in favor of the proposal
func (c *proposalCommon) GetVotesFor(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.VotesFor.RawValue, "getVotesFor", c.ID.RawValue)
}

// Get the number of votes against the proposal
func (c *proposalCommon) GetVotesAgainst(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.VotesAgainst.RawValue, "getVotesAgainst", c.ID.RawValue)
}

// Check if the proposal has been cancelled
func (c *proposalCommon) GetIsCancelled(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.IsCancelled, "getCancelled", c.ID.RawValue)
}

// Check if the proposal has been executed
func (c *proposalCommon) GetIsExecuted(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.IsExecuted, "getExecuted", c.ID.RawValue)
}

// Get the proposal's payload
func (c *proposalCommon) GetPayload(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.Payload, "getPayload", c.ID.RawValue)
}

// Get the proposal's state
func (c *proposalCommon) GetState(mc *batch.MultiCaller) {
	core.AddCall(mc, c.contract, &c.State.RawValue, "getState", c.ID.RawValue)
}

// Get all of the proposal's details
func (c *proposalCommon) QueryAllDetails(mc *batch.MultiCaller) {
	c.GetProposerAddress(mc)
	c.GetMessage(mc)
	c.GetCreatedTime(mc)
	c.GetStartTime(mc)
	c.GetEndTime(mc)
	c.GetExpiryTime(mc)
	c.GetVotesRequired(mc)
	c.GetVotesFor(mc)
	c.GetVotesAgainst(mc)
	c.GetIsCancelled(mc)
	c.GetIsExecuted(mc)
	c.GetPayload(mc)
	c.GetState(mc)
}

// Check if a node has voted on the proposal
func (c *proposalCommon) GetMemberHasVoted(mc *batch.MultiCaller, out *bool, address common.Address) {
	core.AddCall(mc, c.contract, out, "getReceiptHasVoted", c.ID.RawValue, address)
}

// Check if a node has voted in favor of the proposal
func (c *proposalCommon) GetMemberSupported(mc *batch.MultiCaller, out *bool, address common.Address) {
	core.AddCall(mc, c.contract, out, "getReceiptSupported", c.ID.RawValue, address)
}

// Get which DAO the proposal is for - reserved for internal use
func (c *proposalCommon) getDAO(mc *batch.MultiCaller, dao_Out *string) {
	core.AddCall(mc, c.contract, dao_Out, "getDAO", c.ID.RawValue)
}
