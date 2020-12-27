# Go Chonic
Cronic runs the given command and sends the error output per SMTP mail. If you configure `sendstdout=true` in cronic.conf it also sends a mail on success.

From `0 1 * * * backup_cmd >/dev/null 2>&1`

To `0 1 * * * cronic backup_cmd`


![screenshot](shot.png)

## Configuration
You can use Config file or `.env` file or ENV variables (in this order).
Default config location `~/.config/cronic/cronic.conf` can be changed with `-c`.

```bash
mkdir -p ~/.config/cronic/
cp cronic.conf ~/.config/cronic/cronic.conf
vim ~/.config/cronic/cronic.conf

# or set new path with -c
./cronic -c ./path/cronic.conf echo "hello world"
```

Environment Variables example:
```bash
CRONIC_SMTP_HOST=example.com
CRONIC_SMTP_PORT=465
CRONIC_SMTP_USERNAME=usernamess
CRONIC_SMTP_PASSWORD=Passw0rd!
CRONIC_SMTP_ENCRYPTION=SSL
CRONIC_MAIL_SENDER=sender@example.org
CRONIC_MAIL_RECEIVER=receiver@example.org
CRONIC_MAIL_SENDSTDOUT=true
CRONIC_MAIL_SUBJECT=...
CRONIC_MAIL_TEMPLATE=...
```


## Subject and Boady Template
Adjust the template with what the problem could be and possible next steps.

Example on how to change them:
```toml
subject="custom subject"
template="""
Output for the command:
{{.Command}}
RESULT CODE: {{.ResultCode}}
ERROR OUTPUT:
{{.ErrorOutput}}
STANDARD OUTPUT:
{{.StandardOutput}}
"""

```

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
