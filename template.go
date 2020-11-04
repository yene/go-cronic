package main

func defaultMailTemplate() string {
	return `Cronic output for the command:
{{.Command}}

RESULT CODE: {{.ResultCode}}

ERROR OUTPUT:
{{.ErrorOutput}}
STANDARD OUTPUT:
{{.StandardOutput}}
`
}
