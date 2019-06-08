# wp

Rotate wallpaper/background.

## Install

```
go get -u github.com/shadowfax-chc/wallpaper/wp
```

Note: This requires [feh](https://feh.finalrewind.org/) for actually setting
the background image.

## Configuration

By default `wp` looks for a YAML config file in `~/.wp.yaml`. The path to the
config file can be changes using the `--config` option.

Options can also be set via command line args, see `wp --help` for details.

### Example config file

```yaml
directory: ~/.wallpaper
shuffle: true
# See the BACKGROUND SETTINGS section of the feh man page: https://man.finalrewind.org/1/feh/
# This supports any of the `--bg-*` options, specify it without the `--bg-` prefix
mode: fill
# How often to change the background image. This must be specified in a format
# that is parsable by https://golang.org/pkg/time/#ParseDuration.
update-frequency: 5m
# Where to send any logs, default is stdout. This also supports `syslog`.
log-handler: stdout
# Verbosity of logs. Valid values are: DEBUG, INFO, WARN, ERROR
log-level: WARN
```

## Signals

A running instance of `wp` will response to the following signals:

| Signal  | Action                                                         |
| ------- | -------------------------------------------------------------- |
| SIGHUP  | Reload config file, re-scan the directory, and set a new image |
| SIGUSR1 | Switch to the next image, and reset the timer                  |
