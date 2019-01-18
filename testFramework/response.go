package testFramework

import (
	"fmt"
	"strconv"
	"time"

	"github.com/logrusorgru/aurora"
)

type TestResponse struct {
	Success  bool    `json:"success"`  // Pass/Fail of the test
	Sent     int     `json:"sent"`     // Number of tests started
	Received int     `json:"received"` // Number of responses received
	Rate     float32 `json:"rate"`     // Requests pr second
	Raw      string  `json:"raw"`      // Raw output from test
}

type TestResult struct {
	Name              string        `json:"name"`     // Name of test
	Duration          time.Duration `json:"duration"` // Time to execute test
	SuccessPercentage float32       // Percentage success
	TestResponse
}

func (r TestResult) StringSlice() [][]string {
	var success string
	if r.Success {
		success = fmt.Sprintf("[%s]", aurora.BgGreen(" PASS "))
	} else {
		success = fmt.Sprintf("[%s]", aurora.BgRed(" FAIL "))
	}
	return [][]string{
		[]string{"Name", r.Name},
		[]string{"Result", success},
		[]string{"Sent", strconv.Itoa(r.Sent)},
		[]string{"Received", strconv.Itoa(r.Received)},
		[]string{"Percentage", fmt.Sprintf("%.1f%%", r.SuccessPercentage)},
		[]string{"Rate", fmt.Sprintf("%.2f", r.Rate)},
		[]string{"Duration", r.Duration.Truncate(1 * time.Second).String()},
	}
}
