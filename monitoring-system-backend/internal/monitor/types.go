package monitor

import "time"

type Endpoint struct {
    ID                  int       `json:"id,omitempty"`
    ServiceName         string    `json:"service_name"`
    URL                 string    `json:"url"`
    ServerName          string    `json:"server_name"`
    APIMethod           string    `json:"api_method"`
    ExpectedCode        int       `json:"expected_status_code"`
    GitlabURL           *string   `json:"gitlab_url,omitempty"`
    DockerContainerName *string   `json:"docker_container_name,omitempty"`
    KubernetesPodName   *string   `json:"kubernetes_pod_name,omitempty"`
    Tags                []string  `json:"tags,omitempty"`
    Description         *string   `json:"description,omitempty"`
    LastChangedBy       *string   `json:"last_changed_by,omitempty"`
}


type EndpointDetail struct {
	// Core endpoint info (endpoints table)
	ID           int    `db:"id" json:"id"`
	ServiceName  string `db:"service_name" json:"service_name"`
	URL          string `db:"url" json:"url"`
	ServerName   string `db:"server_name" json:"server_name"`
	APIMethod    string `db:"api_method" json:"api_method"`
	ExpectedCode int    `db:"expected_status_code" json:"expected_code"`

	// Extra info (endpoint_info table)
	Description         string   `db:"description" json:"description"`
	GitlabURL           string   `db:"gitlab_url" json:"gitlab_url"`
	DockerContainerName string   `db:"docker_container_name" json:"docker_container_name"`
	KubernetesPodName   string   `db:"kubernetes_pod_name" json:"kubernetes_pod_name"`
	Tags                []string `db:"tags" json:"tags"`
	HasBeenModified     bool     `db:"has_been_modified" json:"has_been_modified"`
	LastChangedBy       string   `db:"last_changed_by" json:"last_changed_by"`

	// Monitoring stats (endpoint_stats + checks tables)
	EndpointID       int     `db:"endpoint_id" json:"endpoint_id"`
	TotalChecks      int     `db:"total_checks" json:"total_checks"`
	AvgLatency       float64 `db:"avg_latency" json:"avg_latency"`
	SuccessfulChecks int     `db:"successful_checks" json:"successful_checks"`
	UptimePercentage float64 `db:"uptime_percentage" json:"uptime_percentage"`
	LastRunSucceeded string  `db:"last_run" json:"last_run_succeeded"` // true = succeeded, false = failed
	FailureCount     int     `db:"failure_count" json:"failure_count"`

	// Timestamps
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}



type EndpointBasics struct {
	ID					int			`db:"id"`
	ServiceName 		string 		`db:"service_name"`
	ServerName  		string 		`db:"server_name"`
	URL					string		`db:"url"`
	TotalChecks			int			`db:"total_checks"`
	UptimePercentage	float64		`db:"uptime_percentage"`
	SuccessfulChecks	int			`db:"successful_checks"`
	AverageLatency		float64		`db:"avg_latency"`
	LastRun				bool		`db:"last_run"`
	FailureCount		int			`db:"failure_count"`
}

type EndpointBasicsDTO struct {
	ID					int			`json:"id"`
	ServiceName 		string 		`json:"service_name"`
	ServerName  		string 		`json:"server_name"`
	URL					string		`json:"url"`
	TotalChecks			int			`json:"total_checks"`
	UptimePercentage	float64		`json:"uptime_percentage"`
	DownTimeCount		int			`json:"downtime_count"`
	SuccessfulChecks	int			`json:"successful_checks"`
	AverageLatency		float64		`json:"avg_latency"`
	LastRun				bool		`json:"last_run"`
	FailureCount		int			`json:"failure_count"`
}

type AggregateDTO struct {
	TotalEndpoints     int     `json:"total_endpoints"`
	TotalChecks        int     `json:"total_checks"`
	SuccessfulChecks   int     `json:"successful_checks"`
	DownTimeCount      int     `json:"down_time_count"`
	OverallUptime      float64 `json:"overall_uptime"`
	AverageLatency     float64 `json:"average_latency"`
}