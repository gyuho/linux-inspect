// Package schema defines '/proc' schema.
package schema

import (
	"reflect"

	"github.com/gyuho/linux-inspect/pkg/schemautil"
)

// NetDev represents '/proc/net/dev'
// (See http://man7.org/linux/man-pages/man5/proc.5.html
// or http://www.onlamp.com/pub/a/linux/2000/11/16/LinuxAdmin.html).
var NetDev = schemautil.RawData{
	IsYAML: false,
	Columns: []schemautil.Column{
		{"interface", "network interface", reflect.String},

		{"receive_bytes", "total number of bytes of data received by the interface", reflect.Uint64},
		{"receive_packets", "total number of packets of data received by the interface", reflect.Uint64},
		{"receive_errs", "total number of receive errors detected by the device driver", reflect.Uint64},
		{"receive_drop", "total number of packets dropped by the device driver", reflect.Uint64},
		{"receive_fifo", "number of FIFO buffer errors", reflect.Uint64},
		{"receive_frame", "number of packet framing errors", reflect.Uint64},
		{"receive_compressed", "number of compressed packets received by the device driver", reflect.Uint64},
		{"receive_multicast", "number of multicast frames received by the device driver", reflect.Uint64},

		{"transmit_bytes", "total number of bytes of data transmitted by the interface", reflect.Uint64},
		{"transmit_packets", "total number of packets of data transmitted by the interface", reflect.Uint64},
		{"transmit_errs", "total number of receive errors detected by the device driver", reflect.Uint64},
		{"transmit_drop", "total number of packets dropped by the device driver", reflect.Uint64},
		{"transmit_fifo", "number of FIFO buffer errors", reflect.Uint64},
		{"transmit_colls", "number of collisions detected on the interface", reflect.Uint64},
		{"transmit_carrier", "number of carrier losses detected by the device driver", reflect.Uint64},
	},
	ColumnsToParse: map[string]schemautil.RawDataType{
		"receive_bytes":  schemautil.TypeBytes,
		"transmit_bytes": schemautil.TypeBytes,
	},
}

// NetTCP represents '/proc/net/tcp' and '/proc/net/tcp6'
// (See http://man7.org/linux/man-pages/man5/proc.5.html
// and http://www.onlamp.com/pub/a/linux/2000/11/16/LinuxAdmin.html).
var NetTCP = schemautil.RawData{
	IsYAML: false,
	Columns: []schemautil.Column{
		{"sl", "kernel hash slot", reflect.Uint64},
		{"local_address", "local-address:port", reflect.String},
		{"rem_address", "remote-address:port", reflect.String},
		{"st", "internal status of socket", reflect.String},
		{"tx_queue", "outgoing data queue in terms of kernel memory usage", reflect.String},
		{"rx_queue", "incoming data queue in terms of kernel memory usage", reflect.String},
		{"tr", "internal information of the kernel socket state", reflect.String},
		{"tm->when", "internal information of the kernel socket state", reflect.String},
		{"retrnsmt", "internal information of the kernel socket state", reflect.String},
		{"uid", "effective UID of the creator of the socket", reflect.Uint64},
		{"timeout", "timeout", reflect.Uint64},
		{"inode", "inode raw data", reflect.String},
	},
	ColumnsToParse: map[string]schemautil.RawDataType{
		"local_address": schemautil.TypeIPAddress,
		"rem_address":   schemautil.TypeIPAddress,
		"st":            schemautil.TypeStatus,
	},
}

// LoadAvg represents '/proc/loadavg'
// (See http://man7.org/linux/man-pages/man5/proc.5.html).
var LoadAvg = schemautil.RawData{
	IsYAML: false,
	Columns: []schemautil.Column{
		{"load-avg-1-minute", "total uptime in seconds", reflect.Float64},
		{"load-avg-5-minute", "total uptime in seconds", reflect.Float64},
		{"load-avg-15-minute", "total uptime in seconds", reflect.Float64},
		{"runnable-kernel-scheduling-entities", "number of currently runnable kernel scheduling entities (processes, threads)", reflect.Int64},
		{"current-kernel-scheduling-entities", "number of kernel scheduling entities that currently exist on the system", reflect.Int64},
		{"pid", "PID of the process that was most recently created on the system", reflect.Int64},
	},
	ColumnsToParse: map[string]schemautil.RawDataType{},
}

// Uptime represents '/proc/uptime'
// (See http://man7.org/linux/man-pages/man5/proc.5.html).
var Uptime = schemautil.RawData{
	IsYAML: false,
	Columns: []schemautil.Column{
		{"uptime-total", "total uptime in seconds", reflect.Float64},
		{"uptime-idle", "total amount of time in seconds spent in idle process", reflect.Float64},
	},
	ColumnsToParse: map[string]schemautil.RawDataType{
		"uptime-total": schemautil.TypeTimeSeconds,
		"uptime-idle":  schemautil.TypeTimeSeconds,
	},
}

// DiskStat represents '/proc/diskstats'
// (See https://www.kernel.org/doc/Documentation/ABI/testing/procfs-diskstats
// and https://www.kernel.org/doc/Documentation/iostats.txt).
var DiskStat = schemautil.RawData{
	IsYAML: false,
	Columns: []schemautil.Column{
		{"major-number", "major device number", reflect.Uint64},
		{"minor-number", "minor device number", reflect.Uint64},
		{"device-name", "device name", reflect.String},

		{"reads-completed", "total number of reads completed successfully", reflect.Uint64},
		{"reads-merged", "total number of reads merged when adjacent to each other", reflect.Uint64},
		{"sectors-read", "total number of sectors read successfully", reflect.Uint64},
		{"time-spent-on-reading-ms", "total number of milliseconds spent by all reads", reflect.Uint64},

		{"writes-completed", "total number of writes completed successfully", reflect.Uint64},
		{"writes-merged", "total number of writes merged when adjacent to each other", reflect.Uint64},
		{"sectors-written", "total number of sectors written successfully", reflect.Uint64},
		{"time-spent-on-writing-ms", "total number of milliseconds spent by all writes", reflect.Uint64},

		{"I/Os-in-progress", "only field that should go to zero (incremented as requests are on request_queue)", reflect.Uint64},
		{"time-spent-on-I/Os-ms", "milliseconds spent doing I/Os", reflect.Uint64},
		{"weighted-time-spent-on-I/Os-ms", "weighted milliseconds spent doing I/Os (incremented at each I/O start, I/O completion, I/O merge)", reflect.Uint64},
	},
	ColumnsToParse: map[string]schemautil.RawDataType{
		"time-spent-on-reading-ms":       schemautil.TypeTimeMicroseconds,
		"time-spent-on-writing-ms":       schemautil.TypeTimeMicroseconds,
		"time-spent-on-I/Os-ms":          schemautil.TypeTimeMicroseconds,
		"weighted-time-spent-on-I/Os-ms": schemautil.TypeTimeMicroseconds,
	},
}

// IO represents 'proc/$PID/io'
// (See http://man7.org/linux/man-pages/man5/proc.5.html).
var IO = schemautil.RawData{
	IsYAML: true,
	Columns: []schemautil.Column{
		{"rchar", "number of bytes which this task has caused to be read from storage (sum of bytes which this process passed to read)", reflect.Uint64},
		{"wchar", "number of bytes which this task has caused, or shall cause to be written to disk", reflect.Uint64},
		{"syscr", "number of read I/O operations", reflect.Uint64},
		{"syscw", "number of write I/O operations", reflect.Uint64},
		{"read_bytes", "number of bytes which this process really did cause to be fetched from the storage layer", reflect.Uint64},
		{"write_bytes", "number of bytes which this process caused to be sent to the storage layer", reflect.Uint64},
		{"cancelled_write_bytes", "number of bytes which this process caused to not happen by truncating pagecache", reflect.Uint64},
	},
	ColumnsToParse: map[string]schemautil.RawDataType{
		"rchar":                 schemautil.TypeBytes,
		"wchar":                 schemautil.TypeBytes,
		"read_bytes":            schemautil.TypeBytes,
		"write_bytes":           schemautil.TypeBytes,
		"cancelled_write_bytes": schemautil.TypeBytes,
	},
}

// Stat represents '/proc/$PID/stat'
// (See http://man7.org/linux/man-pages/man5/proc.5.html).
var Stat = schemautil.RawData{
	IsYAML: false,
	Columns: []schemautil.Column{
		{"pid", "process ID", reflect.Int64},
		{"comm", "filename of the executable (originally in parentheses, automatically removed by this package)", reflect.String},
		{"state", "one character that represents the state of the process", reflect.String},
		{"ppid", "PID of the parent process", reflect.Int64},
		{"pgrp", "group ID of the process", reflect.Int64},
		{"session", "session ID of the process", reflect.Int64},
		{"tty_nr", "controlling terminal of the process", reflect.Int64},
		{"tpgid", "ID of the foreground process group of the controlling terminal of the process", reflect.Int64},
		{"flags", "kernel flags word of the process", reflect.Int64},
		{"minflt", "number of minor faults the process has made which have not required loading a memory page from disk", reflect.Uint64},
		{"cminflt", "number of minor faults that the process's waited-for children have made", reflect.Uint64},
		{"majflt", "number of major faults the process has made which have required loading a memory page from disk", reflect.Uint64},
		{"cmajflt", "number of major faults that the process's waited-for children have made", reflect.Uint64},
		{"utime", "number of clock ticks that this process has been scheduled in user mode (includes guest_time)", reflect.Uint64},
		{"stime", "number of clock ticks that this process has been scheduled in kernel mode", reflect.Uint64},
		{"cutime", "number of clock ticks that this process's waited-for children have been scheduled in user mode", reflect.Uint64},
		{"cstime", "number of clock ticks that this process's waited-for children have been scheduled in kernel mode", reflect.Uint64},
		{"priority", "for processes running a real-time scheduling policy, the negated scheduling priority, minus one; that is, a number in the range -2 to -100, corresponding to real-time priorities 1 to 99. For processes running under a non-real-time scheduling policy, this is the raw nice value. The kernel stores nice values as numbers in the range 0 (high) to 39 (low)", reflect.Int64},
		{"nice", "nice value, a value in the range 19 (low priority) to -20 (high priority)", reflect.Int64},
		{"num_threads", "number of threads in this process", reflect.Int64},
		{"itrealvalue", "no longer maintained", reflect.Int64},
		{"starttime", "time(number of clock ticks) the process started after system boot", reflect.Uint64},
		{"vsize", "virtual memory size in bytes", reflect.Uint64},
		{"rss", "resident set size: number of pages the process has in real memory (text, data, or stack space but does not include pages which have not been demand-loaded in, or which are swapped out)", reflect.Int64},
		{"rsslim", "current soft limit in bytes on the rss of the process", reflect.Uint64},
		{"startcode", "address above which program text can run", reflect.Uint64},
		{"endcode", "address below which program text can run", reflect.Uint64},
		{"startstack", "address of the start (i.e., bottom) of the stack", reflect.Uint64},
		{"kstkesp", "current value of ESP (stack pointer), as found in the kernel stack page for the process", reflect.Uint64},
		{"kstkeip", "current EIP (instruction pointer)", reflect.Uint64},
		{"signal", "obsolete, because it does not provide information on real-time signals (use /proc/$PID/status)", reflect.Uint64},
		{"blocked", "obsolete, because it does not provide information on real-time signals (use /proc/$PID/status)", reflect.Uint64},
		{"sigignore", "obsolete, because it does not provide information on real-time signals (use /proc/$PID/status)", reflect.Uint64},
		{"sigcatch", "obsolete, because it does not provide information on real-time signals (use /proc/$PID/status)", reflect.Uint64},
		{"wchan", "channel in which the process is waiting (address of a location in the kernel where the process is sleeping)", reflect.Uint64},
		{"nswap", "not maintained (number of pages swapped)", reflect.Uint64},
		{"cnswap", "not maintained (cumulative nswap for child processes)", reflect.Uint64},
		{"exit_signal", "signal to be sent to parent when we die", reflect.Int64},
		{"processor", "CPU number last executed on", reflect.Int64},
		{"rt_priority", "real-time scheduling priority, a number in the range 1 to 99 for processes scheduled under a real-time policy, or 0, for non-real-time processes", reflect.Uint64},
		{"policy", "scheduling policy", reflect.Uint64},
		{"delayacct_blkio_ticks", "aggregated block I/O delays, measured in clock ticks", reflect.Uint64},
		{"guest_time", "number of clock ticks spent running a virtual CPU for a guest operating system", reflect.Uint64},
		{"cguest_time", "number of clock ticks (guest_time of the process's children)", reflect.Uint64},
		{"start_data", "address above which program initialized and uninitialized (BSS) data are placed", reflect.Uint64},
		{"end_data", "address below which program initialized and uninitialized (BSS) data are placed", reflect.Uint64},
		{"start_brk", "address above which program heap can be expanded with brk", reflect.Uint64},
		{"arg_start", "address above which program command-line arguments are placed", reflect.Uint64},
		{"arg_end", "address below program command-line arguments are placed", reflect.Uint64},
		{"env_start", "address above which program environment is placed", reflect.Uint64},
		{"env_end", "address below which program environment is placed", reflect.Uint64},
		{"exit_code", "thread's exit status in the form reported by waitpid(2)", reflect.Int64},
	},
	ColumnsToParse: map[string]schemautil.RawDataType{
		"state":  schemautil.TypeStatus,
		"vsize":  schemautil.TypeBytes,
		"rss":    schemautil.TypeBytes,
		"rsslim": schemautil.TypeBytes,
	},
}

// Status represents 'proc/$PID/status'
// (See http://man7.org/linux/man-pages/man5/proc.5.html).
var Status = schemautil.RawData{
	IsYAML: true,
	Columns: []schemautil.Column{
		{"Name", "command run by this process", reflect.String},
		{"Umask", "process umask, expressed in octal with a leading", reflect.String},
		{"State", "current state of the process: R (running), S (sleeping), D (disk sleep), T (stopped), T (tracing stop), Z (zombie), or X (dead)", reflect.String},

		{"Tgid", "thread group ID", reflect.Int64},
		{"Ngid", "NUMA group ID", reflect.Int64},
		{"Pid", "process ID", reflect.Int64},
		{"PPid", "parent process ID, which launches the Pid", reflect.Int64},
		{"TracerPid", "PID of process tracing this process (0 if not being traced)", reflect.Int64},
		{"Uid", "real, effective, saved set, and filesystem UIDs", reflect.String},
		{"Gid", "real, effective, saved set, and filesystem UIDs", reflect.String},

		{"FDSize", "number of file descriptor slots currently allocated", reflect.Uint64},

		{"Groups", "supplementary group list", reflect.String},

		{"NStgid", "thread group ID (i.e., PID) in each of the PID namespaces of which [pid] is a member", reflect.String},
		{"NSpid", "thread ID (i.e., PID) in each of the PID namespaces of which [pid] is a member", reflect.String},
		{"NSpgid", "process group ID (i.e., PID) in each of the PID namespaces of which [pid] is a member", reflect.String},
		{"NSsid", "descendant namespace session ID hierarchy Session ID in each of the PID namespaces of which [pid] is a member", reflect.String},

		{"VmPeak", "peak virtual memory usage. Vm includes physical memory and swap", reflect.String},
		{"VmSize", "current virtual memory usage. VmSize is the total amount of memory required for this process", reflect.String},
		{"VmLck", "locked memory size", reflect.String},
		{"VmPin", "pinned memory size (pages can't be moved, requires direct-access to physical memory)", reflect.String},
		{"VmHWM", `peak resident set size ("high water mark")`, reflect.String},
		{"VmRSS", "resident set size. VmRSS is the actual amount in memory. Some memory can be swapped out to physical disk. So this is the real memory usage of the process", reflect.String},
		{"VmData", "size of data segment", reflect.String},
		{"VmStk", "size of stack", reflect.String},
		{"VmExe", "size of text segments", reflect.String},
		{"VmLib", "shared library code size", reflect.String},
		{"VmPTE", "page table entries size", reflect.String},
		{"VmPMD", "size of second-level page tables", reflect.String},
		{"VmSwap", "swapped-out virtual memory size by anonymous private", reflect.String},
		{"HugetlbPages", "size of hugetlb memory portions", reflect.String},

		{"Threads", "number of threads in process containing this thread (process)", reflect.Uint64},

		{"SigQ", "queued signals for the real user ID of this process (queued signals / limits)", reflect.String},
		{"SigPnd", "number of signals pending for thread", reflect.String},
		{"ShdPnd", "number of signals pending for process as a whole", reflect.String},

		{"SigBlk", "masks indicating signals being blocked", reflect.String},
		{"SigIgn", "masks indicating signals being ignored", reflect.String},
		{"SigCgt", "masks indicating signals being caught", reflect.String},

		{"CapInh", "masks of capabilities enabled in inheritable sets", reflect.String},
		{"CapPrm", "masks of capabilities enabled in permitted sets", reflect.String},
		{"CapEff", "masks of capabilities enabled in effective sets", reflect.String},
		{"CapBnd", "capability Bounding set", reflect.String},
		{"CapAmb", "ambient capability set", reflect.String},

		{"Seccomp", "seccomp mode of the process (0 means SECCOMP_MODE_DISABLED; 1 means SECCOMP_MODE_STRICT; 2 means SECCOMP_MODE_FILTER)", reflect.Uint64},

		{"Cpus_allowed", "mask of CPUs on which this process may run", reflect.String},
		{"Cpus_allowed_list", "list of CPUs on which this process may run", reflect.String},
		{"Mems_allowed", "mask of memory nodes allowed to this process", reflect.String},
		{"Mems_allowed_list", "list of memory nodes allowed to this process", reflect.String},

		{"voluntary_ctxt_switches", "number of voluntary context switches", reflect.Uint64},
		{"nonvoluntary_ctxt_switches", "number of involuntary context switches", reflect.Uint64},
	},
	ColumnsToParse: map[string]schemautil.RawDataType{
		"State":        schemautil.TypeStatus,
		"VmPeak":       schemautil.TypeBytes,
		"VmSize":       schemautil.TypeBytes,
		"VmLck":        schemautil.TypeBytes,
		"VmPin":        schemautil.TypeBytes,
		"VmHWM":        schemautil.TypeBytes,
		"VmRSS":        schemautil.TypeBytes,
		"VmData":       schemautil.TypeBytes,
		"VmStk":        schemautil.TypeBytes,
		"VmExe":        schemautil.TypeBytes,
		"VmLib":        schemautil.TypeBytes,
		"VmPTE":        schemautil.TypeBytes,
		"VmPMD":        schemautil.TypeBytes,
		"VmSwap":       schemautil.TypeBytes,
		"HugetlbPages": schemautil.TypeBytes,
	},
}
