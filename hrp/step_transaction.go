package hrp

import "fmt"

// StepTransaction implements IStep interface.
type StepTransaction struct {
	step *TStep
}

func (s *StepTransaction) Name() string {
	if s.step.Name != "" {
		return s.step.Name
	}
	return fmt.Sprintf("transaction %s %s", s.step.Transaction.Name, s.step.Transaction.Type)
}

func (s *StepTransaction) Type() string {
	return "transaction"
}

func (s *StepTransaction) ToStruct() *TStep {
	return s.step
}
