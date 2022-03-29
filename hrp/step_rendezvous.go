package hrp

import (
	"time"

	"github.com/rs/zerolog/log"
)

// StepRendezvous implements IStep interface.
type StepRendezvous struct {
	step *TStep
}

func (s *StepRendezvous) Name() string {
	if s.step.Name != "" {
		return s.step.Name
	}
	return s.step.Rendezvous.Name
}

func (s *StepRendezvous) Type() string {
	return "rendezvous"
}

func (s *StepRendezvous) ToStruct() *TStep {
	return s.step
}

// Rendezvous creates a new rendezvous
func (s *StepRequest) Rendezvous(name string) *StepRendezvous {
	s.step.Rendezvous = &Rendezvous{
		Name: name,
	}
	return &StepRendezvous{
		step: s.step,
	}
}

// WithUserNumber sets the user number needed to release the current rendezvous
func (s *StepRendezvous) WithUserNumber(number int64) *StepRendezvous {
	s.step.Rendezvous.Number = number
	return s
}

// WithUserPercent sets the user percent needed to release the current rendezvous
func (s *StepRendezvous) WithUserPercent(percent float32) *StepRendezvous {
	s.step.Rendezvous.Percent = percent
	return s
}

// WithTimeout sets the timeout of duration between each user arriving at the current rendezvous
func (s *StepRendezvous) WithTimeout(timeout int64) *StepRendezvous {
	s.step.Rendezvous.Timeout = timeout
	return s
}

func initRendezvous(testcase *TestCase, total int64) []*Rendezvous {
	tCase := testcase.ToTCase()
	var rendezvousList []*Rendezvous
	for _, step := range tCase.TestSteps {
		if step.Rendezvous == nil {
			continue
		}
		rendezvous := step.Rendezvous

		// either number or percent should be correctly put, otherwise set to default (total)
		if rendezvous.Number == 0 && rendezvous.Percent > 0 && rendezvous.Percent <= defaultRendezvousPercent {
			rendezvous.Number = int64(rendezvous.Percent * float32(total))
		} else if rendezvous.Number > 0 && rendezvous.Number <= total && rendezvous.Percent == 0 {
			rendezvous.Percent = float32(rendezvous.Number) / float32(total)
		} else {
			log.Warn().
				Str("name", rendezvous.Name).
				Int64("default number", total).
				Float32("default percent", defaultRendezvousPercent).
				Msg("rendezvous parameter not defined or error, set to default value")
			rendezvous.Number = total
			rendezvous.Percent = defaultRendezvousPercent
		}

		if rendezvous.Timeout <= 0 {
			rendezvous.Timeout = defaultRendezvousTimeout
		}

		rendezvous.reset()
		rendezvousList = append(rendezvousList, rendezvous)
	}
	return rendezvousList
}

func waitRendezvous(rendezvousList []*Rendezvous) {
	if rendezvousList != nil {
		lastRendezvous := rendezvousList[len(rendezvousList)-1]
		for _, rendezvous := range rendezvousList {
			go waitSingleRendezvous(rendezvous, rendezvousList, lastRendezvous)
		}
	}
}

func waitSingleRendezvous(rendezvous *Rendezvous, rendezvousList []*Rendezvous, lastRendezvous *Rendezvous) {
	for {
		// cycle start: block current checking until current rendezvous activated
		<-rendezvous.activateChan
		stop := make(chan struct{})
		timeout := time.Duration(rendezvous.Timeout) * time.Millisecond
		timer := time.NewTimer(timeout)
		go func() {
			defer close(stop)
			rendezvous.wg.Wait()
		}()
		for !rendezvous.isReleased() {
			select {
			case <-rendezvous.timerResetChan:
				timer.Reset(timeout)
			case <-stop:
				rendezvous.setReleased()
				close(rendezvous.releaseChan)
				log.Info().
					Str("name", rendezvous.Name).
					Float32("percent", rendezvous.Percent).
					Int64("number", rendezvous.Number).
					Int64("timeout(ms)", rendezvous.Timeout).
					Int64("cnt", rendezvous.cnt).
					Str("reason", "rendezvous release condition satisfied").
					Msg("rendezvous released")
			case <-timer.C:
				rendezvous.setReleased()
				close(rendezvous.releaseChan)
				log.Info().
					Str("name", rendezvous.Name).
					Float32("percent", rendezvous.Percent).
					Int64("number", rendezvous.Number).
					Int64("timeout(ms)", rendezvous.Timeout).
					Int64("cnt", rendezvous.cnt).
					Str("reason", "time's up").
					Msg("rendezvous released")
			}
		}
		// cycle end: reset all previous rendezvous after last rendezvous released
		// otherwise, block current checker until the last rendezvous end
		if rendezvous == lastRendezvous {
			for _, r := range rendezvousList {
				r.reset()
			}
		} else {
			<-lastRendezvous.releaseChan
		}
	}
}
