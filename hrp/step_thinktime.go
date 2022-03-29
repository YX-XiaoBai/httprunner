package hrp

// StepThinkTime implements IStep interface.
type StepThinkTime struct {
	step *TStep
}

func (s *StepThinkTime) Name() string {
	return s.step.Name
}

func (s *StepThinkTime) Type() string {
	return "thinktime"
}

func (s *StepThinkTime) ToStruct() *TStep {
	return s.step
}
