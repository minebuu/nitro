// Copyright 2021-2022, Offchain Labs, Inc.
// For license information, see https://github.com/nitro/blob/master/LICENSE

package l1pricing

import (
	"math/big"
	"testing"

	am "github.com/offchainlabs/nitro/util/arbmath"

	"github.com/ethereum/go-ethereum/common"
	"github.com/offchainlabs/nitro/arbos/burn"
	"github.com/offchainlabs/nitro/arbos/storage"
)

func TestL1PriceUpdate(t *testing.T) {
	sto := storage.NewMemoryBacked(burn.NewSystemBurner(nil, false))
	err := InitializeL1PricingState(sto, common.Address{})
	Require(t, err)
	ps := OpenL1PricingState(sto)

	tyme, err := ps.LastUpdateTime()
	Require(t, err)
	if tyme != 0 {
		Fail(t)
	}

	initialPriceEstimate := am.UintToBig(InitialPricePerUnitWei)
	priceEstimate, err := ps.PricePerUnit()
	Require(t, err)
	if priceEstimate.Cmp(initialPriceEstimate) != 0 {
		Fail(t)
	}
}

func TestL1Throttling(t *testing.T) {
	sto := storage.NewMemoryBacked(burn.NewSystemBurner(nil, false))
	err := InitializeL1PricingState(sto, common.Address{})
	Require(t, err)
	ps := OpenL1PricingState(sto)

	l1Price, err := ps.PricePerUnit()
	Require(t, err)
	if l1Price.Cmp(big.NewInt(InitialPricePerUnitWei)) != 0 {
		t.Fatal()
	}
	l1Price, err = ps.adjustPricePerUnit(l1Price)
	Require(t, err)
	if l1Price.Cmp(big.NewInt(InitialPricePerUnitWei)) != 0 {
		t.Fatal()
	}

	Require(t, ps.SetL1DataUnitsSpeedLimit(1000))
	Require(t, ps.SetL1DataUnitsThreshold(1000))
	l1Price, err = ps.PricePerUnit()
	Require(t, err)
	if l1Price.Cmp(big.NewInt(InitialPricePerUnitWei)) != 0 {
		t.Fatal()
	}
	l1Price, err = ps.adjustPricePerUnit(l1Price)
	Require(t, err)
	if l1Price.Cmp(big.NewInt(InitialPricePerUnitWei)) != 0 {
		t.Fatal()
	}

	Require(t, ps.SetL1DataUnitsBacklog(5000))
	l1Price, err = ps.PricePerUnit()
	Require(t, err)
	if l1Price.Cmp(big.NewInt(InitialPricePerUnitWei)) != 0 {
		t.Fatal()
	}
	l1Price, err = ps.adjustPricePerUnit(l1Price)
	Require(t, err)
	if l1Price.Cmp(big.NewInt(InitialPricePerUnitWei)) != 0 {
		t.Fatal()
	}

	Require(t, ps.SetL1DataUnitsBacklog(10000))
	l1Price, err = ps.PricePerUnit()
	Require(t, err)
	if l1Price.Cmp(big.NewInt(InitialPricePerUnitWei)) != 0 {
		t.Fatal()
	}
	l1Price, err = ps.adjustPricePerUnit(l1Price)
	Require(t, err)
	if l1Price.Cmp(big.NewInt(InitialPricePerUnitWei)) == 0 {
		t.Fatal()
	}
}
