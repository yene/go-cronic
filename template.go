package main

func mailTemplate() string {
	return `Cronic detected failure or error output for the command:
{{.Command}}

RESULT CODE: {{.ResultCode}}

ERROR OUTPUT:
{{.ErrorOutput}}
STANDARD OUTPUT:
{{.StandardOutput}}
`
}
