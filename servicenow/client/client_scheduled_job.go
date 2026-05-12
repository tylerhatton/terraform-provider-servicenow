package client

// EndpointScheduledJob is the endpoint to manage scheduled job (sysauto_script) records.
const EndpointScheduledJob = "sysauto_script.do"

// ScheduledJob is the json response for a scheduled script job in ServiceNow.
//
// Note: sysauto_script does not expose `description` or `priority` columns on
// the base record, so those fields are not modelled here.
type ScheduledJob struct {
	BaseResult
	Name          string `json:"name"`
	Script        string `json:"script"`
	RunType       string `json:"run_type"`
	RunTime       string `json:"run_time"`
	RunDayOfWeek  string `json:"run_dayofweek"`
	RunDayOfMonth string `json:"run_dayofmonth"`
	RunPeriod     string `json:"run_period"`
	RunStart      string `json:"run_start"`
	Active        bool   `json:"active,string"`
	Conditional   bool   `json:"conditional,string"`
	Condition     string `json:"condition"`
}
