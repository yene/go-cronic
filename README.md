# Go Chronic
From `0 1 * * * backup >/dev/null 2>&1`
To `0 1 * * * cronic backup`

![screenshot](shot.png)

## Configuration
```bash
mkdir -p ~/.config/chronic/
cp chronic.conf ~/.config/chronic/chronic.conf
vim ~/.config/chronic/chronic.conf

# or set new path with -c
./cronic -c ./path/chronic.conf echo "hello world"
```

> To change mail template or subject modify the code.


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
