## ssn [![Build Status](https://img.shields.io/travis/gyuho/ssn.svg?style=flat-square)](https://travis-ci.org/gyuho/ssn) [![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/gyuho/ssn)

ss, netstat in Go. ssn is an utility to investigate sockets.

```
go get -v -u github.com/gyuho/ssn/...
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
+----------+------------------------------------------------------------------+-------+------------------+--------------------+-------+
| PROTOCOL |                             PROGRAM                              |  PID  |    LOCAL ADDR    |    REMOTE ADDR     | USER  |
+----------+------------------------------------------------------------------+-------+------------------+--------------------+-------+
| tcp      |                                                                  |     0 | 127.0.0.1:631    | 0.0.0.0:0          | root  |
| tcp      |                                                                  |     0 | 127.0.1.1:53     | 0.0.0.0:0          | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 10.0.0.122:39584 | 192.30.252.86:443  | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 10.0.0.122:49993 | 192.30.252.90:443  | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:1278   | 0.0.0.0:0          | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:1279   | 0.0.0.0:0          | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:1280   | 127.0.0.1:56820    | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:1280   | 127.0.0.1:56821    | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:1280   | 127.0.0.1:56812    | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:1280   | 0.0.0.0:0          | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:1280   | 127.0.0.1:56813    | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:37102  | 127.0.0.1:1380     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:37103  | 127.0.0.1:1380     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:37110  | 127.0.0.1:1380     | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:37111  | 127.0.0.1:1380     | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:37117  | 127.0.0.1:1380     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:37119  | 127.0.0.1:1380     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:37125  | 127.0.0.1:1380     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:37127  | 127.0.0.1:1380     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:37151  | 127.0.0.1:1380     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:37154  | 127.0.0.1:1380     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:37163  | 127.0.0.1:1380     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:37165  | 127.0.0.1:1380     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:42094  | 127.0.0.1:1180     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:42101  | 127.0.0.1:1180     | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:42102  | 127.0.0.1:1180     | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:42106  | 127.0.0.1:1180     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:42108  | 127.0.0.1:1180     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:42111  | 127.0.0.1:1180     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:42114  | 127.0.0.1:1180     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:42116  | 127.0.0.1:1180     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:42119  | 127.0.0.1:1180     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:42122  | 127.0.0.1:1180     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:42140  | 127.0.0.1:1180     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:42145  | 127.0.0.1:1180     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:42152  | 127.0.0.1:1180     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:42154  | 127.0.0.1:1180     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:42157  | 127.0.0.1:1180     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:42160  | 127.0.0.1:1180     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:56816  | 127.0.0.1:1280     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:56828  | 127.0.0.1:1280     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:56830  | 127.0.0.1:1280     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:56832  | 127.0.0.1:1280     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:56836  | 127.0.0.1:1280     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:56838  | 127.0.0.1:1280     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:56840  | 127.0.0.1:1280     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:56844  | 127.0.0.1:1280     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:56862  | 127.0.0.1:1280     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:56866  | 127.0.0.1:1280     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:56874  | 127.0.0.1:1280     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:56876  | 127.0.0.1:1280     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:56878  | 127.0.0.1:1280     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10684 | 127.0.0.1:56882  | 127.0.0.1:1280     | root  |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10685 | 127.0.0.1:1178   | 0.0.0.0:0          | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10685 | 127.0.0.1:1179   | 0.0.0.0:0          | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10685 | 127.0.0.1:1180   | 127.0.0.1:42091    | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10685 | 127.0.0.1:1180   | 0.0.0.0:0          | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10685 | 127.0.0.1:1180   | 127.0.0.1:42102    | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10685 | 127.0.0.1:1180   | 127.0.0.1:42092    | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10685 | 127.0.0.1:1180   | 127.0.0.1:42101    | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10685 | 127.0.0.1:37106  | 127.0.0.1:1380     | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10685 | 127.0.0.1:37107  | 127.0.0.1:1380     | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10685 | 127.0.0.1:56820  | 127.0.0.1:1280     | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10685 | 127.0.0.1:56821  | 127.0.0.1:1280     | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10686 | 127.0.0.1:1378   | 0.0.0.0:0          | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10686 | 127.0.0.1:1379   | 0.0.0.0:0          | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10686 | 127.0.0.1:1380   | 127.0.0.1:37111    | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10686 | 127.0.0.1:1380   | 127.0.0.1:37107    | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10686 | 127.0.0.1:1380   | 127.0.0.1:37110    | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10686 | 127.0.0.1:1380   | 0.0.0.0:0          | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10686 | 127.0.0.1:1380   | 127.0.0.1:37106    | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10686 | 127.0.0.1:42091  | 127.0.0.1:1180     | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10686 | 127.0.0.1:42092  | 127.0.0.1:1180     | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10686 | 127.0.0.1:56812  | 127.0.0.1:1280     | gyuho |
| tcp      | /home/gyuho/go/src/github.com/gyuho/runetcd/run_example/bin/etcd | 10686 | 127.0.0.1:56813  | 127.0.0.1:1280     | gyuho |
| tcp      | /opt/google/chrome/chrome                                        |  7011 | 10.0.0.122:34735 | 192.30.252.125:443 | gyuho |
| tcp      | /opt/google/chrome/chrome                                        |  7011 | 10.0.0.122:34736 | 192.30.252.125:443 | gyuho |
| tcp      | /opt/google/chrome/chrome                                        |  7011 | 10.0.0.122:34751 | 192.30.252.125:443 | gyuho |
| tcp      | /opt/google/chrome/chrome                                        |  7011 | 10.0.0.122:37338 | 216.58.192.2:443   | gyuho |
| tcp      | /opt/google/chrome/chrome                                        |  7011 | 10.0.0.122:39456 | 192.30.252.86:443  | gyuho |
| tcp      | /opt/google/chrome/chrome                                        |  7011 | 10.0.0.122:42807 | 192.30.252.87:443  | gyuho |
| tcp      | /opt/google/chrome/chrome                                        |  7011 | 10.0.0.122:45255 | 192.30.252.128:443 | gyuho |
| tcp      | /opt/google/chrome/chrome                                        |  7011 | 10.0.0.122:45256 | 192.30.252.128:443 | gyuho |
| tcp      | /opt/google/chrome/chrome                                        |  7011 | 10.0.0.122:45257 | 192.30.252.128:443 | gyuho |
| tcp      | /opt/google/chrome/chrome                                        |  7011 | 10.0.0.122:46802 | 52.3.76.20:443     | gyuho |
| tcp      | /opt/google/chrome/chrome                                        |  7011 | 10.0.0.122:46803 | 52.3.76.20:443     | gyuho |
| tcp      | /opt/google/chrome/chrome                                        |  7011 | 10.0.0.122:53588 | 199.27.79.133:443  | gyuho |
| tcp      | /opt/google/chrome/chrome                                        |  7011 | 10.0.0.122:53589 | 199.27.79.133:443  | gyuho |
| tcp      | /opt/google/chrome/chrome                                        |  7011 | 10.0.0.122:53590 | 199.27.79.133:443  | gyuho |
| tcp      | /opt/google/chrome/chrome                                        |  7011 | 10.0.0.122:53591 | 199.27.79.133:443  | gyuho |
| tcp      | /opt/google/chrome/chrome                                        |  7011 | 10.0.0.122:53592 | 199.27.79.133:443  | gyuho |
| tcp      | /opt/google/chrome/chrome                                        |  7011 | 10.0.0.122:53593 | 199.27.79.133:443  | gyuho |
| tcp      | /opt/google/chrome/chrome                                        |  7011 | 10.0.0.122:53604 | 199.27.79.133:443  | gyuho |
| tcp      | /opt/google/chrome/chrome                                        |  7011 | 10.0.0.122:53853 | 23.235.47.133:443  | gyuho |
| tcp      | /opt/google/chrome/chrome                                        |  7011 | 10.0.0.122:53854 | 23.235.47.133:443  | gyuho |
| tcp      | /opt/google/chrome/chrome                                        |  7011 | 10.0.0.122:53855 | 23.235.47.133:443  | gyuho |
| tcp      | /opt/google/chrome/chrome                                        |  7011 | 10.0.0.122:53860 | 23.235.47.133:443  | gyuho |
| tcp      | /opt/google/chrome/chrome                                        |  7011 | 10.0.0.122:53868 | 23.235.47.133:443  | gyuho |
| tcp      | /opt/google/chrome/chrome                                        |  7011 | 10.0.0.122:56703 | 192.30.252.88:443  | gyuho |
| tcp      | /opt/google/chrome/chrome                                        |  7011 | 10.0.0.122:57134 | 192.30.252.88:443  | gyuho |
| tcp      | /usr/lib/gvfs/gvfsd-http                                         |  4722 | 10.0.0.122:34857 | 54.230.147.216:80  | gyuho |
| tcp      | /usr/lib/gvfs/gvfsd-http                                         |  4722 | 10.0.0.122:35589 | 184.25.56.173:80   | gyuho |
| tcp      | /usr/lib/gvfs/gvfsd-http                                         |  4722 | 10.0.0.122:46279 | 192.229.163.33:80  | gyuho |
| tcp      | /usr/lib/gvfs/gvfsd-http                                         |  4722 | 10.0.0.122:49533 | 54.230.141.244:80  | gyuho |
| tcp      | /usr/lib/gvfs/gvfsd-http                                         |  4722 | 10.0.0.122:49539 | 54.230.141.244:80  | gyuho |
| tcp      | /usr/lib/gvfs/gvfsd-http                                         |  4722 | 10.0.0.122:49540 | 54.230.141.244:80  | gyuho |
| tcp      | /usr/lib/gvfs/gvfsd-http                                         |  4722 | 10.0.0.122:50177 | 54.230.147.240:80  | gyuho |
| tcp      | /usr/lib/gvfs/gvfsd-http                                         |  4722 | 10.0.0.122:52413 | 54.230.141.223:80  | gyuho |
| tcp      | /usr/lib/gvfs/gvfsd-http                                         |  4722 | 10.0.0.122:55675 | 54.230.141.107:80  | gyuho |
| tcp      | /usr/lib/gvfs/gvfsd-http                                         |  4722 | 10.0.0.122:58358 | 54.230.141.225:80  | gyuho |
| tcp      | /usr/lib/gvfs/gvfsd-http                                         |  4722 | 10.0.0.122:58791 | 54.230.141.221:80  | gyuho |
| tcp      | /usr/lib/x86_64-linux-gnu/ubuntu-geoip-provider                  |  2218 | 10.0.0.122:42453 | 91.189.94.25:80    | gyuho |
| tcp      | /usr/lib/x86_64-linux-gnu/unity-scope-home/unity-scope-home      |  4684 | 10.0.0.122:56710 | 91.189.92.10:443   | gyuho |
+----------+------------------------------------------------------------------+-------+------------------+--------------------+-------+

+----------+---------------------------+------+-----------------------------------------------+----------------------------------------------+-------+
| PROTOCOL |          PROGRAM          | PID  |                  LOCAL ADDR                   |                 REMOTE ADDR                  | USER  |
+----------+---------------------------+------+-----------------------------------------------+----------------------------------------------+-------+
| tcp6     |                           |    0 | 0000:0001:0000:0000:0000:0000:0000:0000:59438 | 0000:0001:0000:0000:0000:0000:0000:0000:631  | root  |
| tcp6     |                           |    0 | 0000:0001:0000:0000:0000:0000:0000:0000:631   | 0000:0000:0000:0000:0000:0000:0000:0000:0    | root  |
| tcp6     | /opt/google/chrome/chrome | 7011 | 3161:2906:98DE:7E16:C100:A791:2601:0645:35027 | 0000:00BC:0000:0000:400E:0C00:2607:F8B0:5228 | gyuho |
| tcp6     | /opt/google/chrome/chrome | 7011 | 3161:2906:98DE:7E16:C100:A791:2601:0645:35291 | 0000:100C:0000:0000:4005:0802:2607:F8B0:443  | gyuho |
| tcp6     | /opt/google/chrome/chrome | 7011 | 3161:2906:98DE:7E16:C100:A791:2601:0645:36204 | 0000:100F:0000:0000:4010:0801:2607:F8B0:443  | gyuho |
| tcp6     | /opt/google/chrome/chrome | 7011 | 3161:2906:98DE:7E16:C100:A791:2601:0645:36672 | 0000:100F:0000:0000:4005:0800:2607:F8B0:443  | gyuho |
| tcp6     | /opt/google/chrome/chrome | 7011 | 3161:2906:98DE:7E16:C100:A791:2601:0645:37012 | 0000:200E:0000:0000:4005:0803:2607:F8B0:443  | gyuho |
| tcp6     | /opt/google/chrome/chrome | 7011 | 3161:2906:98DE:7E16:C100:A791:2601:0645:37887 | 0000:1000:0000:0000:4010:0801:2607:F8B0:443  | gyuho |
| tcp6     | /opt/google/chrome/chrome | 7011 | 3161:2906:98DE:7E16:C100:A791:2601:0645:37893 | 0000:1000:0000:0000:4010:0801:2607:F8B0:443  | gyuho |
| tcp6     | /opt/google/chrome/chrome | 7011 | 3161:2906:98DE:7E16:C100:A791:2601:0645:37898 | 0000:1000:0000:0000:4010:0801:2607:F8B0:443  | gyuho |
| tcp6     | /opt/google/chrome/chrome | 7011 | 3161:2906:98DE:7E16:C100:A791:2601:0645:39088 | 0000:1002:0000:0000:4005:0800:2607:F8B0:443  | gyuho |
| tcp6     | /opt/google/chrome/chrome | 7011 | 3161:2906:98DE:7E16:C100:A791:2601:0645:40842 | 0000:005E:0000:0000:400C:0C06:2607:F8B0:443  | gyuho |
| tcp6     | /opt/google/chrome/chrome | 7011 | 3161:2906:98DE:7E16:C100:A791:2601:0645:40926 | 0000:1016:0000:0000:4005:0802:2607:F8B0:443  | gyuho |
| tcp6     | /opt/google/chrome/chrome | 7011 | 3161:2906:98DE:7E16:C100:A791:2601:0645:41209 | 0000:006A:0000:0000:400E:0C03:2607:F8B0:443  | gyuho |
| tcp6     | /opt/google/chrome/chrome | 7011 | 3161:2906:98DE:7E16:C100:A791:2601:0645:43664 | 0000:1005:0000:0000:4010:0801:2607:F8B0:443  | gyuho |
| tcp6     | /opt/google/chrome/chrome | 7011 | 3161:2906:98DE:7E16:C100:A791:2601:0645:43680 | 0000:1005:0000:0000:4010:0801:2607:F8B0:443  | gyuho |
| tcp6     | /opt/google/chrome/chrome | 7011 | 3161:2906:98DE:7E16:C100:A791:2601:0645:44247 | 0000:1007:0000:0000:4005:0800:2607:F8B0:443  | gyuho |
| tcp6     | /opt/google/chrome/chrome | 7011 | 3161:2906:98DE:7E16:C100:A791:2601:0645:46763 | 0000:200D:0000:0000:4005:0803:2607:F8B0:443  | gyuho |
| tcp6     | /opt/google/chrome/chrome | 7011 | 3161:2906:98DE:7E16:C100:A791:2601:0645:47353 | 0000:1007:0000:0000:4005:0802:2607:F8B0:443  | gyuho |
| tcp6     | /opt/google/chrome/chrome | 7011 | 3161:2906:98DE:7E16:C100:A791:2601:0645:47372 | 0000:1007:0000:0000:4005:0802:2607:F8B0:443  | gyuho |
| tcp6     | /opt/google/chrome/chrome | 7011 | 3161:2906:98DE:7E16:C100:A791:2601:0645:47578 | 0000:100A:0000:0000:4005:0800:2607:F8B0:443  | gyuho |
| tcp6     | /opt/google/chrome/chrome | 7011 | 3161:2906:98DE:7E16:C100:A791:2601:0645:49144 | 0000:1000:0000:0000:4005:0802:2607:F8B0:443  | gyuho |
| tcp6     | /opt/google/chrome/chrome | 7011 | 3161:2906:98DE:7E16:C100:A791:2601:0645:52139 | 0000:200E:0000:0000:4005:0801:2607:F8B0:443  | gyuho |
| tcp6     | /opt/google/chrome/chrome | 7011 | 3161:2906:98DE:7E16:C100:A791:2601:0645:52142 | 0000:200E:0000:0000:4005:0801:2607:F8B0:443  | gyuho |
| tcp6     | /opt/google/chrome/chrome | 7011 | 3161:2906:98DE:7E16:C100:A791:2601:0645:52277 | 0000:00BD:0000:0000:400B:0C01:2001:4860:443  | gyuho |
| tcp6     | /opt/google/chrome/chrome | 7011 | 3161:2906:98DE:7E16:C100:A791:2601:0645:54963 | 0000:100A:0000:0000:4005:0802:2607:F8B0:443  | gyuho |
| tcp6     | /opt/google/chrome/chrome | 7011 | 3161:2906:98DE:7E16:C100:A791:2601:0645:58795 | 0000:1000:0000:0000:4005:0800:2607:F8B0:443  | gyuho |
| tcp6     | /opt/google/chrome/chrome | 7011 | 3161:2906:98DE:7E16:C100:A791:2601:0645:59091 | 0000:100A:0000:0000:4010:0801:2607:F8B0:443  | gyuho |
+----------+---------------------------+------+-----------------------------------------------+----------------------------------------------+-------+

```

