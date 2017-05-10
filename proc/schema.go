package proc

import (
	"reflect"

	"github.com/gyuho/linux-inspect/schema"
)

// NetDevSchema represents '/proc/net/dev'.
// Reference http://man7.org/linux/man-pages/man5/proc.5.html
// and http://www.onlamp.com/pub/a/linux/2000/11/16/LinuxAdmin.html.
var NetDevSchema = schema.RawData{
	IsYAML: false,
	Columns: []schema.Column{
		{Name: "interface", Godoc: "network interface", Kind: reflect.String},

		{Name: "receive_bytes", Godoc: "total number of bytes of data received by the interface", Kind: reflect.Uint64},
		{Name: "receive_packets", Godoc: "total number of packets of data received by the interface", Kind: reflect.Uint64},
		{Name: "receive_errs", Godoc: "total number of receive errors detected by the device driver", Kind: reflect.Uint64},
		{Name: "receive_drop", Godoc: "total number of packets dropped by the device driver", Kind: reflect.Uint64},
		{Name: "receive_fifo", Godoc: "number of FIFO buffer errors", Kind: reflect.Uint64},
		{Name: "receive_frame", Godoc: "number of packet framing errors", Kind: reflect.Uint64},
		{Name: "receive_compressed", Godoc: "number of compressed packets received by the device driver", Kind: reflect.Uint64},
		{Name: "receive_multicast", Godoc: "number of multicast frames received by the device driver", Kind: reflect.Uint64},

		{Name: "transmit_bytes", Godoc: "total number of bytes of data transmitted by the interface", Kind: reflect.Uint64},
		{Name: "transmit_packets", Godoc: "total number of packets of data transmitted by the interface", Kind: reflect.Uint64},
		{Name: "transmit_errs", Godoc: "total number of receive errors detected by the device driver", Kind: reflect.Uint64},
		{Name: "transmit_drop", Godoc: "total number of packets dropped by the device driver", Kind: reflect.Uint64},
		{Name: "transmit_fifo", Godoc: "number of FIFO buffer errors", Kind: reflect.Uint64},
		{Name: "transmit_colls", Godoc: "number of collisions detected on the interface", Kind: reflect.Uint64},
		{Name: "transmit_carrier", Godoc: "number of carrier losses detected by the device driver", Kind: reflect.Uint64},
	},
	ColumnsToParse: map[string]schema.RawDataType{
		"receive_bytes":  schema.TypeBytes,
		"transmit_bytes": schema.TypeBytes,
	},
}

// NetTCPSchema represents '/proc/net/tcp' and '/proc/net/tcp6'.
// Reference http://man7.org/linux/man-pages/man5/proc.5.html
// and http://www.onlamp.com/pub/a/linux/2000/11/16/LinuxAdmin.html.
var NetTCPSchema = schema.RawData{
	IsYAML: false,
	Columns: []schema.Column{
		{Name: "sl", Godoc: "kernel hash slot", Kind: reflect.Uint64},
		{Name: "local_address", Godoc: "local-address:port", Kind: reflect.String},
		{Name: "rem_address", Godoc: "remote-address:port", Kind: reflect.String},
		{Name: "st", Godoc: "internal status of socket", Kind: reflect.String},
		{Name: "tx_queue", Godoc: "outgoing data queue in terms of kernel memory usage", Kind: reflect.String},
		{Name: "rx_queue", Godoc: "incoming data queue in terms of kernel memory usage", Kind: reflect.String},
		{Name: "tr", Godoc: "internal information of the kernel socket state", Kind: reflect.String},
		{Name: "tm->when", Godoc: "internal information of the kernel socket state", Kind: reflect.String},
		{Name: "retrnsmt", Godoc: "internal information of the kernel socket state", Kind: reflect.String},
		{Name: "uid", Godoc: "effective UID of the creator of the socket", Kind: reflect.Uint64},
		{Name: "timeout", Godoc: "timeout", Kind: reflect.Uint64},
		{Name: "inode", Godoc: "inode raw data", Kind: reflect.String},
	},
	ColumnsToParse: map[string]schema.RawDataType{
		"local_address": schema.TypeIPAddress,
		"rem_address":   schema.TypeIPAddress,
		"st":            schema.TypeStatus,
	},
}

// LoadAvgSchema represents '/proc/loadavg'.
// Reference http://man7.org/linux/man-pages/man5/proc.5.html.
var LoadAvgSchema = schema.RawData{
	IsYAML: false,
	Columns: []schema.Column{
		{Name: "load-avg-1-minute", Godoc: "total uptime in seconds", Kind: reflect.Float64},
		{Name: "load-avg-5-minute", Godoc: "total uptime in seconds", Kind: reflect.Float64},
		{Name: "load-avg-15-minute", Godoc: "total uptime in seconds", Kind: reflect.Float64},
		{Name: "runnable-kernel-scheduling-entities", Godoc: "number of currently runnable kernel scheduling entities (processes, threads)", Kind: reflect.Int64},
		{Name: "current-kernel-scheduling-entities", Godoc: "number of kernel scheduling entities that currently exist on the system", Kind: reflect.Int64},
		{Name: "pid", Godoc: "PID of the process that was most recently created on the system", Kind: reflect.Int64},
	},
	ColumnsToParse: map[string]schema.RawDataType{},
}

// UptimeSchema represents '/proc/uptime'.
// Reference http://man7.org/linux/man-pages/man5/proc.5.html.
var UptimeSchema = schema.RawData{
	IsYAML: false,
	Columns: []schema.Column{
		{Name: "uptime-total", Godoc: "total uptime in seconds", Kind: reflect.Float64},
		{Name: "uptime-idle", Godoc: "total amount of time in seconds spent in idle process", Kind: reflect.Float64},
	},
	ColumnsToParse: map[string]schema.RawDataType{
		"uptime-total": schema.TypeTimeSeconds,
		"uptime-idle":  schema.TypeTimeSeconds,
	},
}

// DiskStatSchema represents '/proc/diskstats'.
// Reference https://www.kernel.org/doc/Documentation/ABI/testing/procfs-diskstats
// and https://www.kernel.org/doc/Documentation/iostats.txt.
var DiskStatSchema = schema.RawData{
	IsYAML: false,
	Columns: []schema.Column{
		{Name: "major-number", Godoc: "major device number", Kind: reflect.Uint64},
		{Name: "minor-number", Godoc: "minor device number", Kind: reflect.Uint64},
		{Name: "device-name", Godoc: "device name", Kind: reflect.String},

		{Name: "reads-completed", Godoc: "total number of reads completed successfully", Kind: reflect.Uint64},
		{Name: "reads-merged", Godoc: "total number of reads merged when adjacent to each other", Kind: reflect.Uint64},
		{Name: "sectors-read", Godoc: "total number of sectors read successfully", Kind: reflect.Uint64},
		{Name: "time-spent-on-reading-ms", Godoc: "total number of milliseconds spent by all reads", Kind: reflect.Uint64},

		{Name: "writes-completed", Godoc: "total number of writes completed successfully", Kind: reflect.Uint64},
		{Name: "writes-merged", Godoc: "total number of writes merged when adjacent to each other", Kind: reflect.Uint64},
		{Name: "sectors-written", Godoc: "total number of sectors written successfully", Kind: reflect.Uint64},
		{Name: "time-spent-on-writing-ms", Godoc: "total number of milliseconds spent by all writes", Kind: reflect.Uint64},

		{Name: "I/Os-in-progress", Godoc: "only field that should go to zero (incremented as requests are on request_queue)", Kind: reflect.Uint64},
		{Name: "time-spent-on-I/Os-ms", Godoc: "milliseconds spent doing I/Os", Kind: reflect.Uint64},
		{Name: "weighted-time-spent-on-I/Os-ms", Godoc: "weighted milliseconds spent doing I/Os (incremented at each I/O start, I/O completion, I/O merge)", Kind: reflect.Uint64},
	},
	ColumnsToParse: map[string]schema.RawDataType{
		"time-spent-on-reading-ms":       schema.TypeTimeMicroseconds,
		"time-spent-on-writing-ms":       schema.TypeTimeMicroseconds,
		"time-spent-on-I/Os-ms":          schema.TypeTimeMicroseconds,
		"weighted-time-spent-on-I/Os-ms": schema.TypeTimeMicroseconds,
	},
}

// IOSchema represents 'proc/$PID/io'.
// Reference http://man7.org/linux/man-pages/man5/proc.5.html.
var IOSchema = schema.RawData{
	IsYAML: true,
	Columns: []schema.Column{
		{Name: "rchar", Godoc: "number of bytes which this task has caused to be read from storage (sum of bytes which this process passed to read)", Kind: reflect.Uint64},
		{Name: "wchar", Godoc: "number of bytes which this task has caused, or shall cause to be written to disk", Kind: reflect.Uint64},
		{Name: "syscr", Godoc: "number of read I/O operations", Kind: reflect.Uint64},
		{Name: "syscw", Godoc: "number of write I/O operations", Kind: reflect.Uint64},
		{Name: "read_bytes", Godoc: "number of bytes which this process really did cause to be fetched from the storage layer", Kind: reflect.Uint64},
		{Name: "write_bytes", Godoc: "number of bytes which this process caused to be sent to the storage layer", Kind: reflect.Uint64},
		{Name: "cancelled_write_bytes", Godoc: "number of bytes which this process caused to not happen by truncating pagecache", Kind: reflect.Uint64},
	},
	ColumnsToParse: map[string]schema.RawDataType{
		"rchar":                 schema.TypeBytes,
		"wchar":                 schema.TypeBytes,
		"read_bytes":            schema.TypeBytes,
		"write_bytes":           schema.TypeBytes,
		"cancelled_write_bytes": schema.TypeBytes,
	},
}

// StatSchema represents '/proc/$PID/stat'.
// Reference http://man7.org/linux/man-pages/man5/proc.5.html.
var StatSchema = schema.RawData{
	IsYAML: false,
	Columns: []schema.Column{
		{Name: "pid", Godoc: "process ID", Kind: reflect.Int64},
		{Name: "comm", Godoc: "filename of the executable (originally in parentheses, automatically removed by this package)", Kind: reflect.String},
		{Name: "state", Godoc: "one character that represents the state of the process", Kind: reflect.String},
		{Name: "ppid", Godoc: "PID of the parent process", Kind: reflect.Int64},
		{Name: "pgrp", Godoc: "group ID of the process", Kind: reflect.Int64},
		{Name: "session", Godoc: "session ID of the process", Kind: reflect.Int64},
		{Name: "tty_nr", Godoc: "controlling terminal of the process", Kind: reflect.Int64},
		{Name: "tpgid", Godoc: "ID of the foreground process group of the controlling terminal of the process", Kind: reflect.Int64},
		{Name: "flags", Godoc: "kernel flags word of the process", Kind: reflect.Int64},
		{Name: "minflt", Godoc: "number of minor faults the process has made which have not required loading a memory page from disk", Kind: reflect.Uint64},
		{Name: "cminflt", Godoc: "number of minor faults that the process's waited-for children have made", Kind: reflect.Uint64},
		{Name: "majflt", Godoc: "number of major faults the process has made which have required loading a memory page from disk", Kind: reflect.Uint64},
		{Name: "cmajflt", Godoc: "number of major faults that the process's waited-for children have made", Kind: reflect.Uint64},
		{Name: "utime", Godoc: "number of clock ticks that this process has been scheduled in user mode (includes guest_time)", Kind: reflect.Uint64},
		{Name: "stime", Godoc: "number of clock ticks that this process has been scheduled in kernel mode", Kind: reflect.Uint64},
		{Name: "cutime", Godoc: "number of clock ticks that this process's waited-for children have been scheduled in user mode", Kind: reflect.Uint64},
		{Name: "cstime", Godoc: "number of clock ticks that this process's waited-for children have been scheduled in kernel mode", Kind: reflect.Uint64},
		{Name: "priority", Godoc: "for processes running a real-time scheduling policy, the negated scheduling priority, minus one; that is, a number in the range -2 to -100, corresponding to real-time priorities 1 to 99. For processes running under a non-real-time scheduling policy, this is the raw nice value. The kernel stores nice values as numbers in the range 0 (high) to 39 (low)", Kind: reflect.Int64},
		{Name: "nice", Godoc: "nice value, a value in the range 19 (low priority) to -20 (high priority)", Kind: reflect.Int64},
		{Name: "num_threads", Godoc: "number of threads in this process", Kind: reflect.Int64},
		{Name: "itrealvalue", Godoc: "no longer maintained", Kind: reflect.Int64},
		{Name: "starttime", Godoc: "time(number of clock ticks) the process started after system boot", Kind: reflect.Uint64},
		{Name: "vsize", Godoc: "virtual memory size in bytes", Kind: reflect.Uint64},
		{Name: "rss", Godoc: "resident set size: number of pages the process has in real memory (text, data, or stack space but does not include pages which have not been demand-loaded in, or which are swapped out)", Kind: reflect.Int64},
		{Name: "rsslim", Godoc: "current soft limit in bytes on the rss of the process", Kind: reflect.Uint64},
		{Name: "startcode", Godoc: "address above which program text can run", Kind: reflect.Uint64},
		{Name: "endcode", Godoc: "address below which program text can run", Kind: reflect.Uint64},
		{Name: "startstack", Godoc: "address of the start (i.e., bottom) of the stack", Kind: reflect.Uint64},
		{Name: "kstkesp", Godoc: "current value of ESP (stack pointer), as found in the kernel stack page for the process", Kind: reflect.Uint64},
		{Name: "kstkeip", Godoc: "current EIP (instruction pointer)", Kind: reflect.Uint64},
		{Name: "signal", Godoc: "obsolete, because it does not provide information on real-time signals (use /proc/$PID/status)", Kind: reflect.Uint64},
		{Name: "blocked", Godoc: "obsolete, because it does not provide information on real-time signals (use /proc/$PID/status)", Kind: reflect.Uint64},
		{Name: "sigignore", Godoc: "obsolete, because it does not provide information on real-time signals (use /proc/$PID/status)", Kind: reflect.Uint64},
		{Name: "sigcatch", Godoc: "obsolete, because it does not provide information on real-time signals (use /proc/$PID/status)", Kind: reflect.Uint64},
		{Name: "wchan", Godoc: "channel in which the process is waiting (address of a location in the kernel where the process is sleeping)", Kind: reflect.Uint64},
		{Name: "nswap", Godoc: "not maintained (number of pages swapped)", Kind: reflect.Uint64},
		{Name: "cnswap", Godoc: "not maintained (cumulative nswap for child processes)", Kind: reflect.Uint64},
		{Name: "exit_signal", Godoc: "signal to be sent to parent when we die", Kind: reflect.Int64},
		{Name: "processor", Godoc: "CPU number last executed on", Kind: reflect.Int64},
		{Name: "rt_priority", Godoc: "real-time scheduling priority, a number in the range 1 to 99 for processes scheduled under a real-time policy, or 0, for non-real-time processes", Kind: reflect.Uint64},
		{Name: "policy", Godoc: "scheduling policy", Kind: reflect.Uint64},
		{Name: "delayacct_blkio_ticks", Godoc: "aggregated block I/O delays, measured in clock ticks", Kind: reflect.Uint64},
		{Name: "guest_time", Godoc: "number of clock ticks spent running a virtual CPU for a guest operating system", Kind: reflect.Uint64},
		{Name: "cguest_time", Godoc: "number of clock ticks (guest_time of the process's children)", Kind: reflect.Uint64},
		{Name: "start_data", Godoc: "address above which program initialized and uninitialized (BSS) data are placed", Kind: reflect.Uint64},
		{Name: "end_data", Godoc: "address below which program initialized and uninitialized (BSS) data are placed", Kind: reflect.Uint64},
		{Name: "start_brk", Godoc: "address above which program heap can be expanded with brk", Kind: reflect.Uint64},
		{Name: "arg_start", Godoc: "address above which program command-line arguments are placed", Kind: reflect.Uint64},
		{Name: "arg_end", Godoc: "address below program command-line arguments are placed", Kind: reflect.Uint64},
		{Name: "env_start", Godoc: "address above which program environment is placed", Kind: reflect.Uint64},
		{Name: "env_end", Godoc: "address below which program environment is placed", Kind: reflect.Uint64},
		{Name: "exit_code", Godoc: "thread's exit status in the form reported by waitpid(2)", Kind: reflect.Int64},
	},
	ColumnsToParse: map[string]schema.RawDataType{
		"state":  schema.TypeStatus,
		"vsize":  schema.TypeBytes,
		"rss":    schema.TypeBytes,
		"rsslim": schema.TypeBytes,
	},
}

// StatusSchema represents 'proc/$PID/status'.
// Reference http://man7.org/linux/man-pages/man5/proc.5.html.
var StatusSchema = schema.RawData{
	IsYAML: true,
	Columns: []schema.Column{
		{Name: "Name", Godoc: "command run by this process", Kind: reflect.String},
		{Name: "Umask", Godoc: "process umask, expressed in octal with a leading", Kind: reflect.String},
		{Name: "State", Godoc: "current state of the process: R (running), S (sleeping), D (disk sleep), T (stopped), T (tracing stop), Z (zombie), or X (dead)", Kind: reflect.String},

		{Name: "Tgid", Godoc: "thread group ID", Kind: reflect.Int64},
		{Name: "Ngid", Godoc: "NUMA group ID", Kind: reflect.Int64},
		{Name: "Pid", Godoc: "process ID", Kind: reflect.Int64},
		{Name: "PPid", Godoc: "parent process ID, which launches the Pid", Kind: reflect.Int64},
		{Name: "TracerPid", Godoc: "PID of process tracing this process (0 if not being traced)", Kind: reflect.Int64},
		{Name: "Uid", Godoc: "real, effective, saved set, and filesystem UIDs", Kind: reflect.String},
		{Name: "Gid", Godoc: "real, effective, saved set, and filesystem UIDs", Kind: reflect.String},

		{Name: "FDSize", Godoc: "number of file descriptor slots currently allocated", Kind: reflect.Uint64},

		{Name: "Groups", Godoc: "supplementary group list", Kind: reflect.String},

		{Name: "NStgid", Godoc: "thread group ID (i.e., PID) in each of the PID namespaces of which [pid] is a member", Kind: reflect.String},
		{Name: "NSpid", Godoc: "thread ID (i.e., PID) in each of the PID namespaces of which [pid] is a member", Kind: reflect.String},
		{Name: "NSpgid", Godoc: "process group ID (i.e., PID) in each of the PID namespaces of which [pid] is a member", Kind: reflect.String},
		{Name: "NSsid", Godoc: "descendant namespace session ID hierarchy Session ID in each of the PID namespaces of which [pid] is a member", Kind: reflect.String},

		{Name: "VmPeak", Godoc: "peak virtual memory usage. Vm includes physical memory and swap", Kind: reflect.String},
		{Name: "VmSize", Godoc: "current virtual memory usage. VmSize is the total amount of memory required for this process", Kind: reflect.String},
		{Name: "VmLck", Godoc: "locked memory size", Kind: reflect.String},
		{Name: "VmPin", Godoc: "pinned memory size (pages can't be moved, requires direct-access to physical memory)", Kind: reflect.String},
		{Name: "VmHWM", Godoc: `peak resident set size ("high water mark")`, Kind: reflect.String},
		{Name: "VmRSS", Godoc: "resident set size. VmRSS is the actual amount in memory. Some memory can be swapped out to physical disk. So this is the real memory usage of the process", Kind: reflect.String},
		{Name: "VmData", Godoc: "size of data segment", Kind: reflect.String},
		{Name: "VmStk", Godoc: "size of stack", Kind: reflect.String},
		{Name: "VmExe", Godoc: "size of text segments", Kind: reflect.String},
		{Name: "VmLib", Godoc: "shared library code size", Kind: reflect.String},
		{Name: "VmPTE", Godoc: "page table entries size", Kind: reflect.String},
		{Name: "VmPMD", Godoc: "size of second-level page tables", Kind: reflect.String},
		{Name: "VmSwap", Godoc: "swapped-out virtual memory size by anonymous private", Kind: reflect.String},
		{Name: "HugetlbPages", Godoc: "size of hugetlb memory portions", Kind: reflect.String},

		{Name: "Threads", Godoc: "number of threads in process containing this thread (process)", Kind: reflect.Uint64},

		{Name: "SigQ", Godoc: "queued signals for the real user ID of this process (queued signals / limits)", Kind: reflect.String},
		{Name: "SigPnd", Godoc: "number of signals pending for thread", Kind: reflect.String},
		{Name: "ShdPnd", Godoc: "number of signals pending for process as a whole", Kind: reflect.String},

		{Name: "SigBlk", Godoc: "masks indicating signals being blocked", Kind: reflect.String},
		{Name: "SigIgn", Godoc: "masks indicating signals being ignored", Kind: reflect.String},
		{Name: "SigCgt", Godoc: "masks indicating signals being caught", Kind: reflect.String},

		{Name: "CapInh", Godoc: "masks of capabilities enabled in inheritable sets", Kind: reflect.String},
		{Name: "CapPrm", Godoc: "masks of capabilities enabled in permitted sets", Kind: reflect.String},
		{Name: "CapEff", Godoc: "masks of capabilities enabled in effective sets", Kind: reflect.String},
		{Name: "CapBnd", Godoc: "capability Bounding set", Kind: reflect.String},
		{Name: "CapAmb", Godoc: "ambient capability set", Kind: reflect.String},

		{Name: "Seccomp", Godoc: "seccomp mode of the process (0 means SECCOMP_MODE_DISABLED; 1 means SECCOMP_MODE_STRICT; 2 means SECCOMP_MODE_FILTER)", Kind: reflect.Uint64},

		{Name: "Cpus_allowed", Godoc: "mask of CPUs on which this process may run", Kind: reflect.String},
		{Name: "Cpus_allowed_list", Godoc: "list of CPUs on which this process may run", Kind: reflect.String},
		{Name: "Mems_allowed", Godoc: "mask of memory nodes allowed to this process", Kind: reflect.String},
		{Name: "Mems_allowed_list", Godoc: "list of memory nodes allowed to this process", Kind: reflect.String},

		{Name: "voluntary_ctxt_switches", Godoc: "number of voluntary context switches", Kind: reflect.Uint64},
		{Name: "nonvoluntary_ctxt_switches", Godoc: "number of involuntary context switches", Kind: reflect.Uint64},
	},
	ColumnsToParse: map[string]schema.RawDataType{
		"State":        schema.TypeStatus,
		"VmPeak":       schema.TypeBytes,
		"VmSize":       schema.TypeBytes,
		"VmLck":        schema.TypeBytes,
		"VmPin":        schema.TypeBytes,
		"VmHWM":        schema.TypeBytes,
		"VmRSS":        schema.TypeBytes,
		"VmData":       schema.TypeBytes,
		"VmStk":        schema.TypeBytes,
		"VmExe":        schema.TypeBytes,
		"VmLib":        schema.TypeBytes,
		"VmPTE":        schema.TypeBytes,
		"VmPMD":        schema.TypeBytes,
		"VmSwap":       schema.TypeBytes,
		"HugetlbPages": schema.TypeBytes,
	},
}
