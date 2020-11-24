// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package api

import (
	"math/rand"
	"time"
)

// SimulatedFailure modes.
type SimulatedFailure int

// Simulated failure constants.
const (
	ExplicitFailure SimulatedFailure = iota
	SilentFailure
	NoFailure
)

type FailureDice struct {
	rand *rand.Rand

	failureRates []float64
}

// NewFailureDice creates new failure dice with specified probabilities.
func NewFailureDice(explicitFailureProbability, silentFailureProbability float64) *FailureDice {
	return &FailureDice{
		rand:         rand.New(rand.NewSource(time.Now().UnixNano())),
		failureRates: []float64{explicitFailureProbability, explicitFailureProbability + silentFailureProbability},
	}
}

// Roll the dice to get the expected failure mode.
func (dice *FailureDice) Roll() SimulatedFailure {
	val := dice.rand.Float64()

	for failure, rate := range dice.failureRates {
		if val < rate {
			return SimulatedFailure(failure)
		}
	}

	return NoFailure
}

// DefaultDice is used in the api.Client.
//
// Default value is to have no failures.
var DefaultDice = NewFailureDice(0, 0)
