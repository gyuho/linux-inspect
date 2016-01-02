## ssn [![Build Status](https://img.shields.io/travis/gyuho/ssn.svg?style=flat-square)](https://travis-ci.org/gyuho/ssn) [![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/gyuho/ssn)

ss, netstat in Go. ssn is an utility to investigate sockets.

```
go get -v -u github.com/gyuho/ssn
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

But I want something easier.

