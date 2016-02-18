## psn [![Build Status](https://img.shields.io/travis/gyuho/psn.svg?style=flat-square)](https://travis-ci.org/gyuho/psn) [![Godoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/gyuho/psn)

Process, socket utilities in Go. It provides features in ps, ss, netstat.
This is still in active development and only supports Linux system.

```
go get -v -u -f github.com/gyuho/psn
```

```
psn provides utilities to investigate OS processes and sockets.

Usage:
  psn [command]

Available Commands:
  ss          ss investigates sockets.
  kill        kill kills programs using syscall. Make sure to specify the flags to find the program.
  monitor     monitor monitors programs.

Flags:
  -h, --help[=false]: help for psn

Use "psn [command] --help" for more information about a command.

```


<br>

## Motivation

Programmatically find ports and PIDs of web servers, otherwise done
by `ss` or `netstat`. For example, when stopping a web server, one can:

```
netstat -tlpn

(Not all processes could be identified, non-owned process info
	will not be shown, you would have to be root to see it all.)
Active Internet connections (only servers)
Proto Recv-Q Send-Q Local Address           Foreign Address         State       PID/Program name
tcp        0      0 127.0.0.1:2379          0.0.0.0:*               LISTEN      21524/bin/etcd
tcp6       0      0 :::8555                 :::*                    LISTEN      21516/goreman
```

And stop those PIDs or:

```
sudo kill $(sudo netstat -tlpn | perl -ne 'my @a = split /[ \/]+/; print "$a[6]\n" if m/:12379/gio');
sudo kill $(sudo netstat -tlpn | perl -ne 'my @a = split /[ \/]+/; print "$a[6]\n" if m/:2379/gio');
sudo kill $(sudo netstat -tlpn | perl -ne 'my @a = split /[ \/]+/; print "$a[6]\n" if m/:8080/gio');
```

But I want something easier. Here's sample output:

```
+----------+-------------------------------------------------------------+-------+------------------+-------------------+-------+
| PROTOCOL |                           PROGRAM                           |  PID  |    LOCAL ADDR    |    REMOTE ADDR    | USER  |
+----------+-------------------------------------------------------------+-------+------------------+-------------------+-------+
| tcp      |                                                             |     0 | 127.0.0.1:631    | 0.0.0.0:0         | root  |
| tcp      | /usr/bin/vim.gnome                                          | 11076 | 127.0.0.1:37829  | 127.0.0.1:1380    | root  |
| tcp      | /usr/lib/gvfs/gvfsd-http                                    |  4722 | 10.0.0.122:58791 | 54.230.141.221:80 | gyuho |
| tcp      | /usr/lib/x86_64-linux-gnu/ubuntu-geoip-provider             |  2218 | 10.0.0.122:42453 | 91.189.94.25:80   | gyuho |
| tcp      | bin/etcd                                                    | 11430 | 127.0.0.1:1278   | 0.0.0.0:0         | gyuho |
| tcp      | bin/etcd                                                    | 11432 | 127.0.0.1:37814  | 127.0.0.1:1380    | gyuho |
| tcp      | bin/etcd                                                    | 11432 | 127.0.0.1:57527  | 127.0.0.1:1280    | gyuho |
| tcp      | bin/etcd                                                    | 11432 | 127.0.0.1:57528  | 127.0.0.1:1280    | gyuho |
+----------+-------------------------------------------------------------+-------+------------------+-------------------+-------+

```
