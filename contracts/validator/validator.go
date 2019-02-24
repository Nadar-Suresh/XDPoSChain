// Copyright (c) 2018 XDCchain
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package validator

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/validator/contract"
)

type Validator struct {
	*contract.XDCValidatorSession
	contractBackend bind.ContractBackend
}

func NewValidator(transactOpts *bind.TransactOpts, contractAddr common.Address, contractBackend bind.ContractBackend) (*Validator, error) {
	validator, err := contract.NewXDCValidator(contractAddr, contractBackend)
	if err != nil {
		return nil, err
	}

	return &Validator{
		&contract.XDCValidatorSession{
			Contract:     validator,
			TransactOpts: *transactOpts,
		},
		contractBackend,
	}, nil
}

func DeployValidator(transactOpts *bind.TransactOpts, contractBackend bind.ContractBackend, validatorAddress []common.Address, caps []*big.Int, ownerAddress common.Address) (common.Address, *Validator, error) {
	minDeposit := new(big.Int)
	minDeposit.SetString("10000000000000000000000000", 10)
	minVoterCap := new(big.Int)
	minVoterCap.SetString("25000000000000000000000", 10)
	// Deposit 10Mil XDC
	// Min Voter Cap 25K XDC
	// 21 masternodes
	// Candidate Delay Withdraw 30 days = 1296000 blocks
	// Voter Delay Withdraw 2 days = 86400 blocks
	validatorAddr, _, _, err := contract.DeployXDCValidator(transactOpts, contractBackend, validatorAddress, caps, ownerAddress, minDeposit, minVoterCap, big.NewInt(150), big.NewInt(1296000), big.NewInt(86400))
	if err != nil {
		return validatorAddr, nil, err
	}

	validator, err := NewValidator(transactOpts, validatorAddr, contractBackend)
	if err != nil {
		return validatorAddr, nil, err
	}

	return validatorAddr, validator, nil
}
