package gitlab

import (
	"fmt"
	"html/template"
	// "os"
	"strings"
)

func (s *PipelineService) RenderAuthorizationRequestToHTML(req *AuthorizationRequest) (string, error) {
	// Fetch the authorization request from the repository
	fmt.Printf("\n\n--- Authorization Request ---\n%+v\n\n", *req)

	// HTML template
	tmpl := `
	<html>
	<head>
		<style>
			body { font-family: Arial, sans-serif; }
			.container { border: 1px solid #ddd; padding: 16px; border-radius: 8px; }
			.title { font-size: 20px; font-weight: bold; margin-bottom: 12px; }
			.section { margin-bottom: 8px; }
			.label { font-weight: bold; }
			.list { margin-left: 20px; }
		</style>
	</head>
	<body>
		<div class="container">
			<div class="title">Authorization Request Details</div>

			<div class="section"><span class="label">ID:</span> {{.ID}}</div>
			<div class="section"><span class="label">Pipeline Run ID:</span> {{.PipelineRunID}}</div>
			<div class="section"><span class="label">Requester:</span> {{.RequesterName}} ({{.RequesterID}})</div>
			<div class="section"><span class="label">Approver:</span> {{.ApproverName}}
				{{if .ApproverID}} ({{.ApproverID}}){{end}}
			</div>
			<div class="section"><span class="label">Status:</span> {{.Status}}</div>
			<div class="section"><span class="label">Created At:</span> {{.CreatedAt}}</div>
			<div class="section"><span class="label">Updated At:</span> {{.UpdatedAt}}</div>
			<div class="section"><span class="label">Comment:</span> {{.Comment}}</div>
			<div class="section"><span class="label">Macro Service:</span> {{.MacroServiceName}}</div>

			<div class="section"><span class="label">Micro Services:</span>
				<ul class="list">
					{{range .MicroServiceNames}}
						<li>{{.}}</li>
					{{else}}
						<li><i>No microservices listed</i></li>
					{{end}}
				</ul>
			</div>
		</div>
	</body>
	</html>`

	// Parse and execute template
	t, err := template.New("authRequest").Parse(tmpl)
	if err != nil {
		return "", err
	}

	builder := &strings.Builder{}
	if err := t.Execute(builder, req); err != nil {
		return "", err
	}

	return builder.String(), nil
}

func (s *PipelineService) RenderExecutionHistoryToHTML(h *ExecutionHistory) (string, error) {
	tmpl := `
	<html>
	<head>
		<style>
			body { font-family: Arial, sans-serif; background-color: #f9f9f9; }
			.container { background: #fff; border: 1px solid #ddd; padding: 20px; border-radius: 8px; width: 600px; margin: auto; }
			.title { font-size: 22px; font-weight: bold; margin-bottom: 16px; color: #333; }
			.section { margin-bottom: 12px; }
			.label { font-weight: bold; color: #555; }
			.list { margin-left: 20px; }
			.error { color: #c0392b; font-weight: bold; }
			.success { color: #27ae60; font-weight: bold; }
		</style>
	</head>
	<body>
		<div class="container">
			<div class="title">Pipeline Execution History</div>

			<div class="section"><span class="label">Execution ID:</span> {{.ID}}</div>
			<div class="section"><span class="label">Pipeline Run ID:</span> {{.PipelineRunID}}</div>
			<div class="section"><span class="label">Pipeline Unit ID:</span> {{.PipelineUnitID}}</div>

			<div class="section"><span class="label">Requester:</span> {{.RequesterName}} ({{.RequesterID}})</div>
			<div class="section"><span class="label">Approver:</span> {{.ApproverName}} {{if .ApproverID}} ({{.ApproverID}}){{end}}</div>

			<div class="section"><span class="label">Macro Service:</span> {{.MacroServiceName}}</div>

			<div class="section"><span class="label">Micro Services:</span>
				<ul class="list">
					{{range .MicroServiceNames}}
						<li>{{.}}</li>
					{{else}}
						<li><i>No microservices listed</i></li>
					{{end}}
				</ul>
			</div>

			<div class="section"><span class="label">Status:</span> 
				{{if eq .Status "FAILED"}}
					<span class="error">{{.Status}}</span>
				{{else if eq .Status "SUCCESS"}}
					<span class="success">{{.Status}}</span>
				{{else}}
					{{.Status}}
				{{end}}
			</div>

			<div class="section"><span class="label">Started At:</span> {{.StartedAt}}</div>
			<div class="section"><span class="label">Completed At:</span> {{.CompletedAt}}</div>

			{{if .ErrorMessage}}
				<div class="section"><span class="label">Error:</span> <span class="error">{{.ErrorMessage}}</span></div>
			{{end}}
		</div>
	</body>
	</html>`

	t, err := template.New("executionHistory").Parse(tmpl)
	if err != nil {
		return "", err
	}

	builder := &strings.Builder{}
	err = t.Execute(builder, h)
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}
