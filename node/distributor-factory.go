package node

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rocket-pool/rocketpool-go/core"
	"github.com/rocket-pool/rocketpool-go/rocketpool"
	"github.com/rocket-pool/rocketpool-go/utils/multicall"
)

// ===============
// === Structs ===
// ===============

// Binding for RocketNodeDistributorFactory
type NodeDistributorFactory struct {
	rp       *rocketpool.RocketPool
	contract *core.Contract
}

// ====================
// === Constructors ===
// ====================

// Creates a new NodeDistributorFactory contract binding
func NewNodeDistributorFactory(rp *rocketpool.RocketPool, opts *bind.CallOpts) (*NodeDistributorFactory, error) {
	// Create the contract
	contract, err := rp.GetContract("rocketNodeDistributorFactory", opts)
	if err != nil {
		return nil, fmt.Errorf("error getting node distributor factory contract: %w", err)
	}

	return &NodeDistributorFactory{
		rp:       rp,
		contract: contract,
	}, nil
}

// =============
// === Calls ===
// =============

// Gets the deterministic address for a node's reward distributor contract
func (c *NodeDistributorFactory) GetDistributorAddress(mc *multicall.MultiCaller, nodeAddress common.Address, address_Out *common.Address) {
	multicall.AddCall(mc, c.contract, address_Out, "getProxyAddress", nodeAddress)
}

// ===================
// === Sub-Getters ===
// ===================

// Get a node's distributor with details
func (c *NodeDistributorFactory) GetNodeDistributor(nodeAddress common.Address, distributorAddress common.Address, opts *bind.CallOpts) (*NodeDistributor, error) {
	// Create the distributor
	distributor, err := NewNodeDistributor(c.rp, nodeAddress, distributorAddress, opts)
	if err != nil {
		return nil, fmt.Errorf("error creating node distributor binding for node %s at %s: %w", nodeAddress.Hex(), distributorAddress.Hex(), err)
	}

	// Get details via a multicall query
	err = c.rp.Query(func(mc *multicall.MultiCaller) {
		distributor.GetAllDetails(mc)
	}, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting node distributor for node %s at %s: %w", nodeAddress.Hex(), distributorAddress.Hex(), err)
	}

	// Return
	return distributor, nil
}
