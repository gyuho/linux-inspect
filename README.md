## psn [![Build Status](https://img.shields.io/travis/gyuho/psn.svg?style=flat-square)](https://travis-ci.org/gyuho/psn) [![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/gyuho/psn)

Process, socket utilities in Go. It provides features in ps, ss, netstat.
This is still in active development and only supports Linux system.

```
go get -v -u github.com/gyuho/psn/...
```

```
psn provides utilities to investigate OS processes and sockets.

Usage:
  psn [flags]
  psn [command]

Available Commands:
  ss          ss investigates sockets.
  kill        kill kills programs using syscall. Make sure to specify the flags to find the program.

Flags:
  -l, --local-ip="": Specify the local IP. Empty lists all local IPs.
  -p, --local-port="": Specify the local port. Empty lists all local ports.
  -s, --program="": Specify the program. Empty lists all programs.
  -t, --protocol="": 'tcp' or 'tcp6'. Empty lists all protocols.
  -r, --remote-ip="": Specify the remote IP. Empty lists all remote IPs.
  -m, --remote-port="": Specify the remote port. Empty lists all remote ports.
  -a, --state="": Specify the state. Empty lists all states.
  -u, --username="": Specify the user name. Empty lists all user names.

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
tcp        0      0 127.0.0.1:22379         0.0.0.0:*               LISTEN      21526/bin/etcd
tcp        0      0 127.0.0.1:22380         0.0.0.0:*               LISTEN      21526/bin/etcd
tcp        0      0 127.0.0.1:32379         0.0.0.0:*               LISTEN      21528/bin/etcd
tcp        0      0 127.0.0.1:12379         0.0.0.0:*               LISTEN      21529/bin/etcd
tcp        0      0 127.0.0.1:32380         0.0.0.0:*               LISTEN      21528/bin/etcd
tcp        0      0 127.0.0.1:12380         0.0.0.0:*               LISTEN      21529/bin/etcd
tcp        0      0 127.0.0.1:53697         0.0.0.0:*               LISTEN      2608/python2
tcp6       0      0 :::8555                 :::*                    LISTEN      21516/goreman
```

And stop those PIDs or:

```
sudo kill $(sudo netstat -tlpn | perl -ne 'my @a = split /[ \/]+/; print "$a[6]\n" if m/:12379/gio');
sudo kill $(sudo netstat -tlpn | perl -ne 'my @a = split /[ \/]+/; print "$a[6]\n" if m/:22379/gio');
sudo kill $(sudo netstat -tlpn | perl -ne 'my @a = split /[ \/]+/; print "$a[6]\n" if m/:32379/gio');
sudo kill $(sudo netstat -tlpn | perl -ne 'my @a = split /[ \/]+/; print "$a[6]\n" if m/:2379/gio');
sudo kill $(sudo netstat -tlpn | perl -ne 'my @a = split /[ \/]+/; print "$a[6]\n" if m/:8080/gio');
```

But I want something easier. Here's sample output:

```
+----------+-------------------------------------------------------------+-------+------------------+-------------------+-------+
| PROTOCOL |                           PROGRAM                           |  PID  |    LOCAL ADDR    |    REMOTE ADDR    | USER  |
+----------+-------------------------------------------------------------+-------+------------------+-------------------+-------+
| tcp      |                                                             |     0 | 127.0.0.1:631    | 0.0.0.0:0         | root  |
| tcp      |                                                             |     0 | 127.0.1.1:53     | 0.0.0.0:0         | root  |
| tcp      | /opt/google/chrome/chrome                                   |  7011 | 10.0.0.122:57134 | 192.30.252.88:443 | gyuho |
| tcp      | /usr/bin/vim.gnome                                          | 11076 | 127.0.0.1:37829  | 127.0.0.1:1380    | root  |
| tcp      | /usr/lib/gvfs/gvfsd-http                                    |  4722 | 10.0.0.122:58791 | 54.230.141.221:80 | gyuho |
| tcp      | /usr/lib/x86_64-linux-gnu/ubuntu-geoip-provider             |  2218 | 10.0.0.122:42453 | 91.189.94.25:80   | gyuho |
| tcp      | /usr/lib/x86_64-linux-gnu/unity-scope-home/unity-scope-home |  4684 | 10.0.0.122:56710 | 91.189.92.10:443  | gyuho |
| tcp      | bin/etcd                                                    | 11430 | 127.0.0.1:1278   | 0.0.0.0:0         | gyuho |
| tcp      | bin/etcd                                                    | 11430 | 127.0.0.1:1279   | 0.0.0.0:0         | gyuho |
| tcp      | bin/etcd                                                    | 11430 | 127.0.0.1:1280   | 127.0.0.1:57533   | gyuho |
| tcp      | bin/etcd                                                    | 11430 | 127.0.0.1:1280   | 0.0.0.0:0         | gyuho |
| tcp      | bin/etcd                                                    | 11430 | 127.0.0.1:1280   | 127.0.0.1:57534   | gyuho |
| tcp      | bin/etcd                                                    | 11430 | 127.0.0.1:1280   | 127.0.0.1:57527   | gyuho |
| tcp      | bin/etcd                                                    | 11430 | 127.0.0.1:1280   | 127.0.0.1:57528   | gyuho |
| tcp      | bin/etcd                                                    | 11430 | 127.0.0.1:37821  | 127.0.0.1:1380    | gyuho |
| tcp      | bin/etcd                                                    | 11430 | 127.0.0.1:37822  | 127.0.0.1:1380    | gyuho |
| tcp      | bin/etcd                                                    | 11430 | 127.0.0.1:42808  | 127.0.0.1:1180    | gyuho |
| tcp      | bin/etcd                                                    | 11430 | 127.0.0.1:42809  | 127.0.0.1:1180    | gyuho |
| tcp      | bin/etcd                                                    | 11431 | 127.0.0.1:1378   | 0.0.0.0:0         | gyuho |
| tcp      | bin/etcd                                                    | 11431 | 127.0.0.1:1379   | 0.0.0.0:0         | gyuho |
| tcp      | bin/etcd                                                    | 11431 | 127.0.0.1:1380   | 127.0.0.1:37822   | gyuho |
| tcp      | bin/etcd                                                    | 11431 | 127.0.0.1:1380   | 127.0.0.1:37813   | gyuho |
| tcp      | bin/etcd                                                    | 11431 | 127.0.0.1:1380   | 127.0.0.1:37814   | gyuho |
| tcp      | bin/etcd                                                    | 11431 | 127.0.0.1:1380   | 0.0.0.0:0         | gyuho |
| tcp      | bin/etcd                                                    | 11431 | 127.0.0.1:1380   | 127.0.0.1:37821   | gyuho |
| tcp      | bin/etcd                                                    | 11431 | 127.0.0.1:42812  | 127.0.0.1:1180    | gyuho |
| tcp      | bin/etcd                                                    | 11431 | 127.0.0.1:42813  | 127.0.0.1:1180    | gyuho |
| tcp      | bin/etcd                                                    | 11431 | 127.0.0.1:57533  | 127.0.0.1:1280    | gyuho |
| tcp      | bin/etcd                                                    | 11431 | 127.0.0.1:57534  | 127.0.0.1:1280    | gyuho |
| tcp      | bin/etcd                                                    | 11432 | 127.0.0.1:1178   | 0.0.0.0:0         | gyuho |
| tcp      | bin/etcd                                                    | 11432 | 127.0.0.1:1179   | 0.0.0.0:0         | gyuho |
| tcp      | bin/etcd                                                    | 11432 | 127.0.0.1:1180   | 0.0.0.0:0         | gyuho |
| tcp      | bin/etcd                                                    | 11432 | 127.0.0.1:1180   | 127.0.0.1:42808   | gyuho |
| tcp      | bin/etcd                                                    | 11432 | 127.0.0.1:1180   | 127.0.0.1:42809   | gyuho |
| tcp      | bin/etcd                                                    | 11432 | 127.0.0.1:1180   | 127.0.0.1:42812   | gyuho |
| tcp      | bin/etcd                                                    | 11432 | 127.0.0.1:1180   | 127.0.0.1:42813   | gyuho |
| tcp      | bin/etcd                                                    | 11432 | 127.0.0.1:37813  | 127.0.0.1:1380    | gyuho |
| tcp      | bin/etcd                                                    | 11432 | 127.0.0.1:37814  | 127.0.0.1:1380    | gyuho |
| tcp      | bin/etcd                                                    | 11432 | 127.0.0.1:57527  | 127.0.0.1:1280    | gyuho |
| tcp      | bin/etcd                                                    | 11432 | 127.0.0.1:57528  | 127.0.0.1:1280    | gyuho |
+----------+-------------------------------------------------------------+-------+------------------+-------------------+-------+

```

