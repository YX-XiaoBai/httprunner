package hrp

import (
	"sync"
)

type stepType string

const (
	stepTypeRequest     stepType = "request"
	stepTypeAPI         stepType = "api"
	stepTypeTestCase    stepType = "testcase"
	stepTypeTransaction stepType = "transaction"
	stepTypeRendezvous  stepType = "rendezvous"
	stepTypeThinkTime   stepType = "thinktime"
)

// TStep represents teststep data structure.
// Each step maybe three different types: make one request or reference another api/testcase.
type TStep struct {
	Name          string                 `json:"name" yaml:"name"` // required
	Request       *Request               `json:"request,omitempty" yaml:"request,omitempty"`
	API           interface{}            `json:"api,omitempty" yaml:"api,omitempty"`           // *APIPath or *API
	TestCase      interface{}            `json:"testcase,omitempty" yaml:"testcase,omitempty"` // *TestCasePath or *TestCase
	Transaction   *Transaction           `json:"transaction,omitempty" yaml:"transaction,omitempty"`
	Rendezvous    *Rendezvous            `json:"rendezvous,omitempty" yaml:"rendezvous,omitempty"`
	ThinkTime     *ThinkTime             `json:"think_time,omitempty" yaml:"think_time,omitempty"`
	Variables     map[string]interface{} `json:"variables,omitempty" yaml:"variables,omitempty"`
	SetupHooks    []string               `json:"setup_hooks,omitempty" yaml:"setup_hooks,omitempty"`
	TeardownHooks []string               `json:"teardown_hooks,omitempty" yaml:"teardown_hooks,omitempty"`
	Extract       map[string]string      `json:"extract,omitempty" yaml:"extract,omitempty"`
	Validators    []interface{}          `json:"validate,omitempty" yaml:"validate,omitempty"`
	Export        []string               `json:"export,omitempty" yaml:"export,omitempty"`
}

type Transaction struct {
	Name string          `json:"name" yaml:"name"`
	Type transactionType `json:"type" yaml:"type"`
}

type transactionType string

const (
	transactionStart transactionType = "start"
	transactionEnd   transactionType = "end"
)

const (
	defaultRendezvousTimeout int64   = 5000
	defaultRendezvousPercent float32 = 1.0
)

type Rendezvous struct {
	Name           string  `json:"name" yaml:"name"`                           // required
	Percent        float32 `json:"percent,omitempty" yaml:"percent,omitempty"` // default to 1(100%)
	Number         int64   `json:"number,omitempty" yaml:"number,omitempty"`
	Timeout        int64   `json:"timeout,omitempty" yaml:"timeout,omitempty"` // milliseconds
	cnt            int64
	releasedFlag   uint32
	spawnDoneFlag  uint32
	wg             sync.WaitGroup
	timerResetChan chan struct{}
	activateChan   chan struct{}
	releaseChan    chan struct{}
	once           *sync.Once
	lock           sync.Mutex
}

type ThinkTime struct {
	Time float64 `json:"time" yaml:"time"`
}

// IStep represents interface for all types for teststeps, includes:
// StepRequest, StepRequestWithOptionalArgs, StepRequestValidation, StepRequestExtraction,
// StepTestCaseWithOptionalArgs,
// StepTransaction, StepRendezvous.
type IStep interface {
	Name() string
	Type() string
	ToStruct() *TStep
}

// NewStep returns a new constructed teststep with specified step name.
func NewStep(name string) *StepRequest {
	return &StepRequest{
		step: &TStep{
			Name:      name,
			Variables: make(map[string]interface{}),
		},
	}
}
