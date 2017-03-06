DO NOT USE THIS. EXPERIMENTAL!

## psn [![Build Status](https://img.shields.io/travis/gyuho/psn.svg?style=flat-square)](https://travis-ci.org/gyuho/psn) [![Godoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/gyuho/psn)

psn inspects Linux processes, sockets (ps, ss, netstat).

```
go get -v github.com/gyuho/psn/cmd/psn
```

```
Usage:
  psn [command]

Available Commands:
  ds          Inspects '/proc/diskstats'
  ns          Inspects '/proc/net/dev'
  ps          Inspects '/proc/$PID/status', 'top' command output
  ss          Inspects '/proc/net/tcp,tcp6'
```
