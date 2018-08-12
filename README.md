# clock

[![pipeline status](https://gitlab.com/rfaulhaber/clock/badges/master/pipeline.svg)](https://gitlab.com/rfaulhaber/clock/commits/master)

CLI tool for time tracking and record keeping.

By default, this program looks for $HOME/.clock.yaml to know where to write
to, otherwise it'll save to the `--dir` flag, otherwise $HOME/.clock.

This program is a work in progress and not amazingly written.

```
Usage:
  clock [command]

Available Commands:
  help        Help about any command
  report      Aggregates data from logs
  start       records start time
  stop        Stops a current clock and finalizes the record.

Flags:
      --config string   config file (default is $HOME/.clock.yaml)
  -d, --dir string      directory to save this log to (default "$HOME/.clock")
  -h, --help            help for clock
  -t, --tag string      specifies tag to either save, stop, or report on

Use "clock [command] --help" for more information about a command.
```