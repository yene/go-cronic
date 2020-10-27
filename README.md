# Go Chronic
From `0 1 * * * backup >/dev/null 2>&1`
To `0 1 * * * cronic backup`

![screenshot](shot.png)

## Configuration
```
mkdir -p ~/.config/chronic/
cp chronic.conf ~/.config/chronic/chronic.conf
vim ~/.config/chronic/chronic.conf
```

## Releae
`cronic` in this repo is the binary for linux

## TODO
* Scripts that are non executable are skipped

## Inspiration
* https://habilis.net/cronic/
