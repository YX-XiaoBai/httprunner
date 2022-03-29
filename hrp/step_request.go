package hrp

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

type HTTPMethod string

const (
	httpGET     HTTPMethod = "GET"
	httpHEAD    HTTPMethod = "HEAD"
	httpPOST    HTTPMethod = "POST"
	httpPUT     HTTPMethod = "PUT"
	httpDELETE  HTTPMethod = "DELETE"
	httpOPTIONS HTTPMethod = "OPTIONS"
	httpPATCH   HTTPMethod = "PATCH"
)

// Request represents HTTP request data structure.
// This is used for teststep.
type Request struct {
	Method         HTTPMethod             `json:"method" yaml:"method"` // required
	URL            string                 `json:"url" yaml:"url"`       // required
	Params         map[string]interface{} `json:"params,omitempty" yaml:"params,omitempty"`
	Headers        map[string]string      `json:"headers,omitempty" yaml:"headers,omitempty"`
	Cookies        map[string]string      `json:"cookies,omitempty" yaml:"cookies,omitempty"`
	Body           interface{}            `json:"body,omitempty" yaml:"body,omitempty"`
	Json           interface{}            `json:"json,omitempty" yaml:"json,omitempty"`
	Data           interface{}            `json:"data,omitempty" yaml:"data,omitempty"`
	Timeout        float32                `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	AllowRedirects bool                   `json:"allow_redirects,omitempty" yaml:"allow_redirects,omitempty"`
	Verify         bool                   `json:"verify,omitempty" yaml:"verify,omitempty"`
}

type StepRequest struct {
	step *TStep
}

// WithVariables sets variables for current teststep.
func (s *StepRequest) WithVariables(variables map[string]interface{}) *StepRequest {
	s.step.Variables = variables
	return s
}

// SetupHook adds a setup hook for current teststep.
func (s *StepRequest) SetupHook(hook string) *StepRequest {
	s.step.SetupHooks = append(s.step.SetupHooks, hook)
	return s
}

// GET makes a HTTP GET request.
func (s *StepRequest) GET(url string) *StepRequestWithOptionalArgs {
	s.step.Request = &Request{
		Method: httpGET,
		URL:    url,
	}
	return &StepRequestWithOptionalArgs{
		step: s.step,
	}
}

// HEAD makes a HTTP HEAD request.
func (s *StepRequest) HEAD(url string) *StepRequestWithOptionalArgs {
	s.step.Request = &Request{
		Method: httpHEAD,
		URL:    url,
	}
	return &StepRequestWithOptionalArgs{
		step: s.step,
	}
}

// POST makes a HTTP POST request.
func (s *StepRequest) POST(url string) *StepRequestWithOptionalArgs {
	s.step.Request = &Request{
		Method: httpPOST,
		URL:    url,
	}
	return &StepRequestWithOptionalArgs{
		step: s.step,
	}
}

// PUT makes a HTTP PUT request.
func (s *StepRequest) PUT(url string) *StepRequestWithOptionalArgs {
	s.step.Request = &Request{
		Method: httpPUT,
		URL:    url,
	}
	return &StepRequestWithOptionalArgs{
		step: s.step,
	}
}

// DELETE makes a HTTP DELETE request.
func (s *StepRequest) DELETE(url string) *StepRequestWithOptionalArgs {
	s.step.Request = &Request{
		Method: httpDELETE,
		URL:    url,
	}
	return &StepRequestWithOptionalArgs{
		step: s.step,
	}
}

// OPTIONS makes a HTTP OPTIONS request.
func (s *StepRequest) OPTIONS(url string) *StepRequestWithOptionalArgs {
	s.step.Request = &Request{
		Method: httpOPTIONS,
		URL:    url,
	}
	return &StepRequestWithOptionalArgs{
		step: s.step,
	}
}

// PATCH makes a HTTP PATCH request.
func (s *StepRequest) PATCH(url string) *StepRequestWithOptionalArgs {
	s.step.Request = &Request{
		Method: httpPATCH,
		URL:    url,
	}
	return &StepRequestWithOptionalArgs{
		step: s.step,
	}
}

// CallRefCase calls a referenced testcase.
func (s *StepRequest) CallRefCase(tc ITestCase) *StepTestCaseWithOptionalArgs {
	var err error
	s.step.TestCase, err = tc.ToTestCase()
	if err != nil {
		log.Error().Err(err).Msg("failed to load testcase")
		os.Exit(1)
	}
	return &StepTestCaseWithOptionalArgs{
		step: s.step,
	}
}

// CallRefAPI calls a referenced api.
func (s *StepRequest) CallRefAPI(api IAPI) *StepAPIWithOptionalArgs {
	var err error
	s.step.API, err = api.ToAPI()
	if err != nil {
		log.Error().Err(err).Msg("failed to load api")
		os.Exit(1)
	}
	return &StepAPIWithOptionalArgs{
		step: s.step,
	}
}

// StartTransaction starts a transaction.
func (s *StepRequest) StartTransaction(name string) *StepTransaction {
	s.step.Transaction = &Transaction{
		Name: name,
		Type: transactionStart,
	}
	return &StepTransaction{
		step: s.step,
	}
}

// EndTransaction ends a transaction.
func (s *StepRequest) EndTransaction(name string) *StepTransaction {
	s.step.Transaction = &Transaction{
		Name: name,
		Type: transactionEnd,
	}
	return &StepTransaction{
		step: s.step,
	}
}

// SetThinkTime sets think time.
func (s *StepRequest) SetThinkTime(time float64) *StepThinkTime {
	s.step.ThinkTime = &ThinkTime{
		Time: time,
	}
	return &StepThinkTime{
		step: s.step,
	}
}

// StepRequestWithOptionalArgs implements IStep interface.
type StepRequestWithOptionalArgs struct {
	step *TStep
}

// SetVerify sets whether to verify SSL for current HTTP request.
func (s *StepRequestWithOptionalArgs) SetVerify(verify bool) *StepRequestWithOptionalArgs {
	s.step.Request.Verify = verify
	return s
}

// SetTimeout sets timeout for current HTTP request.
func (s *StepRequestWithOptionalArgs) SetTimeout(timeout float32) *StepRequestWithOptionalArgs {
	s.step.Request.Timeout = timeout
	return s
}

// SetProxies sets proxies for current HTTP request.
func (s *StepRequestWithOptionalArgs) SetProxies(proxies map[string]string) *StepRequestWithOptionalArgs {
	// TODO
	return s
}

// SetAllowRedirects sets whether to allow redirects for current HTTP request.
func (s *StepRequestWithOptionalArgs) SetAllowRedirects(allowRedirects bool) *StepRequestWithOptionalArgs {
	s.step.Request.AllowRedirects = allowRedirects
	return s
}

// SetAuth sets auth for current HTTP request.
func (s *StepRequestWithOptionalArgs) SetAuth(auth map[string]string) *StepRequestWithOptionalArgs {
	// TODO
	return s
}

// WithParams sets HTTP request params for current step.
func (s *StepRequestWithOptionalArgs) WithParams(params map[string]interface{}) *StepRequestWithOptionalArgs {
	s.step.Request.Params = params
	return s
}

// WithHeaders sets HTTP request headers for current step.
func (s *StepRequestWithOptionalArgs) WithHeaders(headers map[string]string) *StepRequestWithOptionalArgs {
	s.step.Request.Headers = headers
	return s
}

// WithCookies sets HTTP request cookies for current step.
func (s *StepRequestWithOptionalArgs) WithCookies(cookies map[string]string) *StepRequestWithOptionalArgs {
	s.step.Request.Cookies = cookies
	return s
}

// WithBody sets HTTP request body for current step.
func (s *StepRequestWithOptionalArgs) WithBody(body interface{}) *StepRequestWithOptionalArgs {
	s.step.Request.Body = body
	return s
}

// TeardownHook adds a teardown hook for current teststep.
func (s *StepRequestWithOptionalArgs) TeardownHook(hook string) *StepRequestWithOptionalArgs {
	s.step.TeardownHooks = append(s.step.TeardownHooks, hook)
	return s
}

// Validate switches to step validation.
func (s *StepRequestWithOptionalArgs) Validate() *StepRequestValidation {
	return &StepRequestValidation{
		step: s.step,
	}
}

// Extract switches to step extraction.
func (s *StepRequestWithOptionalArgs) Extract() *StepRequestExtraction {
	s.step.Extract = make(map[string]string)
	return &StepRequestExtraction{
		step: s.step,
	}
}

func (s *StepRequestWithOptionalArgs) Name() string {
	if s.step.Name != "" {
		return s.step.Name
	}
	return fmt.Sprintf("%v %s", s.step.Request.Method, s.step.Request.URL)
}

func (s *StepRequestWithOptionalArgs) Type() string {
	return fmt.Sprintf("request-%v", s.step.Request.Method)
}

func (s *StepRequestWithOptionalArgs) ToStruct() *TStep {
	return s.step
}

// StepRequestExtraction implements IStep interface.
type StepRequestExtraction struct {
	step *TStep
}

// WithJmesPath sets the JMESPath expression to extract from the response.
func (s *StepRequestExtraction) WithJmesPath(jmesPath string, varName string) *StepRequestExtraction {
	s.step.Extract[varName] = jmesPath
	return s
}

// Validate switches to step validation.
func (s *StepRequestExtraction) Validate() *StepRequestValidation {
	return &StepRequestValidation{
		step: s.step,
	}
}

func (s *StepRequestExtraction) Name() string {
	return s.step.Name
}

func (s *StepRequestExtraction) Type() string {
	return fmt.Sprintf("request-%v", s.step.Request.Method)
}

func (s *StepRequestExtraction) ToStruct() *TStep {
	return s.step
}

// StepRequestValidation implements IStep interface.
type StepRequestValidation struct {
	step *TStep
}

func (s *StepRequestValidation) Name() string {
	if s.step.Name != "" {
		return s.step.Name
	}
	return fmt.Sprintf("%s %s", s.step.Request.Method, s.step.Request.URL)
}

func (s *StepRequestValidation) Type() string {
	return fmt.Sprintf("request-%v", s.step.Request.Method)
}

func (s *StepRequestValidation) ToStruct() *TStep {
	return s.step
}

func (s *StepRequestValidation) AssertEqual(jmesPath string, expected interface{}, msg string) *StepRequestValidation {
	v := Validator{
		Check:   jmesPath,
		Assert:  "equals",
		Expect:  expected,
		Message: msg,
	}
	s.step.Validators = append(s.step.Validators, v)
	return s
}

func (s *StepRequestValidation) AssertGreater(jmesPath string, expected interface{}, msg string) *StepRequestValidation {
	v := Validator{
		Check:   jmesPath,
		Assert:  "greater_than",
		Expect:  expected,
		Message: msg,
	}
	s.step.Validators = append(s.step.Validators, v)
	return s
}

func (s *StepRequestValidation) AssertLess(jmesPath string, expected interface{}, msg string) *StepRequestValidation {
	v := Validator{
		Check:   jmesPath,
		Assert:  "less_than",
		Expect:  expected,
		Message: msg,
	}
	s.step.Validators = append(s.step.Validators, v)
	return s
}

func (s *StepRequestValidation) AssertGreaterOrEqual(jmesPath string, expected interface{}, msg string) *StepRequestValidation {
	v := Validator{
		Check:   jmesPath,
		Assert:  "greater_or_equals",
		Expect:  expected,
		Message: msg,
	}
	s.step.Validators = append(s.step.Validators, v)
	return s
}

func (s *StepRequestValidation) AssertLessOrEqual(jmesPath string, expected interface{}, msg string) *StepRequestValidation {
	v := Validator{
		Check:   jmesPath,
		Assert:  "less_or_equals",
		Expect:  expected,
		Message: msg,
	}
	s.step.Validators = append(s.step.Validators, v)
	return s
}

func (s *StepRequestValidation) AssertNotEqual(jmesPath string, expected interface{}, msg string) *StepRequestValidation {
	v := Validator{
		Check:   jmesPath,
		Assert:  "not_equal",
		Expect:  expected,
		Message: msg,
	}
	s.step.Validators = append(s.step.Validators, v)
	return s
}

func (s *StepRequestValidation) AssertContains(jmesPath string, expected interface{}, msg string) *StepRequestValidation {
	v := Validator{
		Check:   jmesPath,
		Assert:  "contains",
		Expect:  expected,
		Message: msg,
	}
	s.step.Validators = append(s.step.Validators, v)
	return s
}

func (s *StepRequestValidation) AssertTypeMatch(jmesPath string, expected interface{}, msg string) *StepRequestValidation {
	v := Validator{
		Check:   jmesPath,
		Assert:  "type_match",
		Expect:  expected,
		Message: msg,
	}
	s.step.Validators = append(s.step.Validators, v)
	return s
}

func (s *StepRequestValidation) AssertRegexp(jmesPath string, expected interface{}, msg string) *StepRequestValidation {
	v := Validator{
		Check:   jmesPath,
		Assert:  "regex_match",
		Expect:  expected,
		Message: msg,
	}
	s.step.Validators = append(s.step.Validators, v)
	return s
}

func (s *StepRequestValidation) AssertStartsWith(jmesPath string, expected interface{}, msg string) *StepRequestValidation {
	v := Validator{
		Check:   jmesPath,
		Assert:  "startswith",
		Expect:  expected,
		Message: msg,
	}
	s.step.Validators = append(s.step.Validators, v)
	return s
}

func (s *StepRequestValidation) AssertEndsWith(jmesPath string, expected interface{}, msg string) *StepRequestValidation {
	v := Validator{
		Check:   jmesPath,
		Assert:  "endswith",
		Expect:  expected,
		Message: msg,
	}
	s.step.Validators = append(s.step.Validators, v)
	return s
}

func (s *StepRequestValidation) AssertLengthEqual(jmesPath string, expected interface{}, msg string) *StepRequestValidation {
	v := Validator{
		Check:   jmesPath,
		Assert:  "length_equals",
		Expect:  expected,
		Message: msg,
	}
	s.step.Validators = append(s.step.Validators, v)
	return s
}

func (s *StepRequestValidation) AssertContainedBy(jmesPath string, expected interface{}, msg string) *StepRequestValidation {
	v := Validator{
		Check:   jmesPath,
		Assert:  "contained_by",
		Expect:  expected,
		Message: msg,
	}
	s.step.Validators = append(s.step.Validators, v)
	return s
}

func (s *StepRequestValidation) AssertLengthLessThan(jmesPath string, expected interface{}, msg string) *StepRequestValidation {
	v := Validator{
		Check:   jmesPath,
		Assert:  "length_less_than",
		Expect:  expected,
		Message: msg,
	}
	s.step.Validators = append(s.step.Validators, v)
	return s
}

func (s *StepRequestValidation) AssertStringEqual(jmesPath string, expected interface{}, msg string) *StepRequestValidation {
	v := Validator{
		Check:   jmesPath,
		Assert:  "string_equals",
		Expect:  expected,
		Message: msg,
	}
	s.step.Validators = append(s.step.Validators, v)
	return s
}

func (s *StepRequestValidation) AssertLengthLessOrEquals(jmesPath string, expected interface{}, msg string) *StepRequestValidation {
	v := Validator{
		Check:   jmesPath,
		Assert:  "length_less_or_equals",
		Expect:  expected,
		Message: msg,
	}
	s.step.Validators = append(s.step.Validators, v)
	return s
}

func (s *StepRequestValidation) AssertLengthGreaterThan(jmesPath string, expected interface{}, msg string) *StepRequestValidation {
	v := Validator{
		Check:   jmesPath,
		Assert:  "length_greater_than",
		Expect:  expected,
		Message: msg,
	}
	s.step.Validators = append(s.step.Validators, v)
	return s
}

func (s *StepRequestValidation) AssertLengthGreaterOrEquals(jmesPath string, expected interface{}, msg string) *StepRequestValidation {
	v := Validator{
		Check:   jmesPath,
		Assert:  "length_greater_or_equals",
		Expect:  expected,
		Message: msg,
	}
	s.step.Validators = append(s.step.Validators, v)
	return s
}

// Validator represents validator for one HTTP response.
type Validator struct {
	Check   string      `json:"check" yaml:"check"` // get value with jmespath
	Assert  string      `json:"assert" yaml:"assert"`
	Expect  interface{} `json:"expect" yaml:"expect"`
	Message string      `json:"msg,omitempty" yaml:"msg,omitempty"` // optional
}
