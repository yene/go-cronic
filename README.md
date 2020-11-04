# Go Chronic
Chronic runs the given command and sends the error output per SMTP mail. If you configure `sendstdout=true` in chronic.conf it also sends a mail on success.

From `0 1 * * * backup_cmd >/dev/null 2>&1`

To `0 1 * * * cronic backup_cmd`



![screenshot](shot.png)

## Configuration
Uses the [TOML format](https://toml.io/en/)

```bash
mkdir -p ~/.config/chronic/
cp chronic.conf ~/.config/chronic/chronic.conf
vim ~/.config/chronic/chronic.conf

# or set new path with -c
./cronic -c ./path/chronic.conf echo "hello world"
```

## Subject and Boady Template
Example on how to change them:
```toml
subject="custom subject"
template="""
The command had a problem:
{{.Command}}
RESULT CODE: {{.ResultCode}}
ERROR OUTPUT:
{{.ErrorOutput}}
STANDARD OUTPUT:
{{.StandardOutput}}
"""

```


## Releae
`cronic` in this repo is the binary for linux

## Features and Todo
- [x] Send stderr per mail
- [x] Consume TOML config
- [x] Take path to config from flag
- [x] Option for TLS/SSL/none
- [x] Option to always send stdout
- [ ] Forward stdout/stderr to parent
- [ ] Inform if a script is not executable (chmod +x)

## Inspiration
* https://habilis.net/cronic/

## Dependencies
* github.com/xhit/go-simple-mail
* github.com/naoina/toml
