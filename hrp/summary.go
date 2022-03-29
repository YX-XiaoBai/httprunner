package hrp

import (
	"bufio"
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/httprunner/httprunner/hrp/internal/builtin"
	"github.com/httprunner/httprunner/hrp/internal/version"
	"github.com/rs/zerolog/log"
)

func newOutSummary() *Summary {
	platForm := &platform{
		HttprunnerVersion: version.VERSION,
		GoVersion:         runtime.Version(),
		Platform:          fmt.Sprintf("%v-%v", runtime.GOOS, runtime.GOARCH),
	}
	return &Summary{
		Success: true,
		Stat:    &stat{},
		Time: &testCaseTime{
			StartAt: time.Now(),
		},
		Platform: platForm,
	}
}

// Summary stores tests summary for current task execution, maybe include one or multiple testcases
type Summary struct {
	Success  bool               `json:"success" yaml:"success"`
	Stat     *stat              `json:"stat" yaml:"stat"`
	Time     *testCaseTime      `json:"time" yaml:"time"`
	Platform *platform          `json:"platform" yaml:"platform"`
	Details  []*testCaseSummary `json:"details" yaml:"details"`
}

func (s *Summary) appendCaseSummary(caseSummary *testCaseSummary) {
	s.Success = s.Success && caseSummary.Success
	s.Stat.TestCases.Total += 1
	s.Stat.TestSteps.Total += len(caseSummary.Records)
	if caseSummary.Success {
		s.Stat.TestCases.Success += 1
	} else {
		s.Stat.TestCases.Fail += 1
	}
	s.Stat.TestSteps.Successes += caseSummary.Stat.Successes
	s.Stat.TestSteps.Failures += caseSummary.Stat.Failures
	s.Details = append(s.Details, caseSummary)
	s.Success = s.Success && caseSummary.Success
}

type stat struct {
	TestCases testCaseStat `json:"testcases" yaml:"test_cases"`
	TestSteps testStepStat `json:"teststeps" yaml:"test_steps"`
}

type testCaseStat struct {
	Total   int `json:"total" yaml:"total"`
	Success int `json:"success" yaml:"success"`
	Fail    int `json:"fail" yaml:"fail"`
}

type testStepStat struct {
	Total     int `json:"total" yaml:"total"`
	Successes int `json:"successes" yaml:"successes"`
	Failures  int `json:"failures" yaml:"failures"`
}

type testCaseTime struct {
	StartAt  time.Time `json:"start_at,omitempty" yaml:"start_at,omitempty"`
	Duration float64   `json:"duration,omitempty" yaml:"duration,omitempty"`
}

type platform struct {
	HttprunnerVersion string `json:"httprunner_version" yaml:"httprunner_version"`
	GoVersion         string `json:"go_version" yaml:"go_version"`
	Platform          string `json:"platform" yaml:"platform"`
}

// testCaseSummary stores tests summary for one testcase
type testCaseSummary struct {
	Name    string         `json:"name" yaml:"name"`
	Success bool           `json:"success" yaml:"success"`
	CaseId  string         `json:"case_id,omitempty" yaml:"case_id,omitempty"` // TODO
	Stat    *testStepStat  `json:"stat" yaml:"stat"`
	Time    *testCaseTime  `json:"time" yaml:"time"`
	InOut   *testCaseInOut `json:"in_out" yaml:"in_out"`
	Log     string         `json:"log,omitempty" yaml:"log,omitempty"` // TODO
	Records []*stepData    `json:"records" yaml:"records"`
}

type stepData struct {
	Name        string                 `json:"name" yaml:"name"`                                   // step name
	StepType    stepType               `json:"step_type" yaml:"step_type"`                         // step type, testcase/request/transaction/rendezvous
	Success     bool                   `json:"success" yaml:"success"`                             // step execution result
	Elapsed     int64                  `json:"elapsed_ms" yaml:"elapsed_ms"`                       // step execution time in millisecond(ms)
	Data        interface{}            `json:"data,omitempty" yaml:"data,omitempty"`               // session data or slice of step data
	ContentSize int64                  `json:"content_size" yaml:"content_size"`                   // response body length
	ExportVars  map[string]interface{} `json:"export_vars,omitempty" yaml:"export_vars,omitempty"` // extract variables
	Attachment  string                 `json:"attachment,omitempty" yaml:"attachment,omitempty"`   // step error information
}

type testCaseInOut struct {
	ConfigVars map[string]interface{} `json:"config_vars" yaml:"config_vars"`
	ExportVars map[string]interface{} `json:"export_vars" yaml:"export_vars"`
}

func newSessionData() *SessionData {
	return &SessionData{
		Success:  false,
		ReqResps: &reqResps{},
	}
}

type SessionData struct {
	Success    bool                `json:"success" yaml:"success"`
	ReqResps   *reqResps           `json:"req_resps" yaml:"req_resps"`
	Address    *address            `json:"address,omitempty" yaml:"address,omitempty"` // TODO
	Validators []*validationResult `json:"validators,omitempty" yaml:"validators,omitempty"`
}

type reqResps struct {
	Request  interface{} `json:"request" yaml:"request"`
	Response interface{} `json:"response" yaml:"response"`
}

type address struct {
	ClientIP   string `json:"client_ip,omitempty" yaml:"client_ip,omitempty"`
	ClientPort string `json:"client_port,omitempty" yaml:"client_port,omitempty"`
	ServerIP   string `json:"server_ip,omitempty" yaml:"server_ip,omitempty"`
	ServerPort string `json:"server_port,omitempty" yaml:"server_port,omitempty"`
}

type validationResult struct {
	Validator
	CheckValue  interface{} `json:"check_value" yaml:"check_value"`
	CheckResult string      `json:"check_result" yaml:"check_result"`
}

func newSummary() *testCaseSummary {
	return &testCaseSummary{
		Success: true,
		Stat:    &testStepStat{},
		Time:    &testCaseTime{},
		InOut:   &testCaseInOut{},
	}
}

func (r *caseRunner) getSummary() *testCaseSummary {
	caseSummary := r.summary
	caseSummary.Time.StartAt = r.startTime
	caseSummary.Time.Duration = time.Since(r.startTime).Seconds()
	exportVars := make(map[string]interface{})
	for _, value := range r.Config.Export {
		exportVars[value] = r.sessionVariables[value]
	}
	caseSummary.InOut.ExportVars = exportVars
	caseSummary.InOut.ConfigVars = r.Config.Variables
	return caseSummary
}

//go:embed internal/scaffold/templates/report/template.html
var reportTemplate string

func (s *Summary) genHTMLReport() error {
	dir, _ := filepath.Split(reportPath)
	err := builtin.EnsureFolderExists(dir)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(fmt.Sprintf(reportPath, s.Time.StartAt.Unix()), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Error().Err(err).Msg("open file failed")
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	tmpl := template.Must(template.New("report").Parse(reportTemplate))
	err = tmpl.Execute(writer, s)
	if err != nil {
		log.Error().Err(err).Msg("execute applies a parsed template to the specified data object failed")
		return err
	}
	err = writer.Flush()
	return err
}
