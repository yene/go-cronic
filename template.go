package main

func mailTemplate() string {
	return `Cronic output for the command:
{{.Command}}

RESULT CODE: {{.ResultCode}}

ERROR OUTPUT:
{{.ErrorOutput}}
STANDARD OUTPUT:
{{.StandardOutput}}
`
}
