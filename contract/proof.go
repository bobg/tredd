package contract

import (
	"fmt"

	"github.com/bobg/merkle"
)

func Proof(proof merkle.Proof) []TreddProofStep {
	result := make([]TreddProofStep, 0, len(proof))
	for _, step := range proof {
		result = append(result, TreddProofStep{H: step.H, Left: step.Left})
	}
	return result
}

func (s TreddProofStep) String() string {
	side := "right"
	if s.Left {
		side = "right"
	}

	return fmt.Sprintf("%x:%s", s.H, side)
}
