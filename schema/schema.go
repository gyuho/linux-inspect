// Package schema defines proc schema.
package schema

import "reflect"

// RawDataType defines how the raw data bytes are defined.
type RawDataType int

const (
	TypeBytes RawDataType = iota
	TypeInt64
	TypeFloat64
	TypeTimeMicroseconds
	TypeTimeSeconds
	TypeIPAddress
	TypeStatus
)

// RawData defines 'proc' raw data.
type RawData struct {
	// IsYAML is true if raw data is parsable in YAML.
	IsYAML bool

	Columns        []Column
	ColumnsToParse map[string]RawDataType
}

// Column represents the schema column.
type Column struct {
	Name  string
	Godoc string
	Kind  reflect.Kind
}

// NetDev represents '/proc/net/dev'
// (See http://man7.org/linux/man-pages/man5/proc.5.html
// or http://www.onlamp.com/pub/a/linux/2000/11/16/LinuxAdmin.html).
var NetDev = RawData{
	IsYAML: false,
	Columns: []Column{
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
	ColumnsToParse: map[string]RawDataType{
		"receive_bytes":  TypeBytes,
		"transmit_bytes": TypeBytes,
	},
}

// NetTCP represents '/proc/net/tcp' and '/proc/net/tcp6'
// (See http://man7.org/linux/man-pages/man5/proc.5.html
// and http://www.onlamp.com/pub/a/linux/2000/11/16/LinuxAdmin.html).
var NetTCP = RawData{
	IsYAML: false,
	Columns: []Column{
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
	ColumnsToParse: map[string]RawDataType{
		"local_address": TypeIPAddress,
		"rem_address":   TypeIPAddress,
		"st":            TypeStatus,
	},
}

// TopCommandRow represents a row in 'top' command output.
// (See http://man7.org/linux/man-pages/man1/top.1.html).
var TopCommandRow = RawData{
	IsYAML: false,
	Columns: []Column{
		{"PID", "pid of the process", reflect.Int64},
		{"USER", "user name", reflect.String},
		{"PR", "priority", reflect.String},
		{"NI", "nice value of the task", reflect.String},
		{"VIRT", "total amount  of virtual memory used by the task (in KiB)", reflect.String},
		{"RES", "non-swapped physical memory a task is using (in KiB)", reflect.String},
		{"SHR", "amount of shared memory available to a task, not all of which is typically resident (in KiB)", reflect.String},
		{"S", "process status", reflect.String},
		{"CPUPercent", "%CPU", reflect.Float64},
		{"MEMPercent", "%MEM", reflect.Float64},
		{"TIME", "CPU time (TIME+)", reflect.String},
		{"COMMAND", "command", reflect.String},
	},
	ColumnsToParse: map[string]RawDataType{
		"S":    TypeStatus,
		"VIRT": TypeBytes,
		"RES":  TypeBytes,
		"SHR":  TypeBytes,
	},
}

// LoadAvg represents '/proc/loadavg'
// (See http://man7.org/linux/man-pages/man5/proc.5.html).
var LoadAvg = RawData{
	IsYAML: false,
	Columns: []Column{
		{"load-avg-1-minute", "total uptime in seconds", reflect.Float64},
		{"load-avg-5-minute", "total uptime in seconds", reflect.Float64},
		{"load-avg-15-minute", "total uptime in seconds", reflect.Float64},
		{"runnable-kernel-scheduling-entities", "number of currently runnable kernel scheduling entities (processes, threads)", reflect.Int64},
		{"current-kernel-scheduling-entities", "number of kernel scheduling entities that currently exist on the system", reflect.Int64},
		{"pid", "PID of the process that was most recently created on the system", reflect.Int64},
	},
	ColumnsToParse: map[string]RawDataType{},
}

// Uptime represents '/proc/uptime'
// (See http://man7.org/linux/man-pages/man5/proc.5.html).
var Uptime = RawData{
	IsYAML: false,
	Columns: []Column{
		{"uptime-total", "total uptime in seconds", reflect.Float64},
		{"uptime-idle", "total amount of time in seconds spent in idle process", reflect.Float64},
	},
	ColumnsToParse: map[string]RawDataType{
		"uptime-total": TypeTimeSeconds,
		"uptime-idle":  TypeTimeSeconds,
	},
}

// Mtab represents '/etc/mtab'
// (See https://en.wikipedia.org/wiki/Fstab
// and https://en.wikipedia.org/wiki/Mtab).
var Mtab = RawData{
	IsYAML: false,
	Columns: []Column{
		{"file-system", "file system", reflect.String},
		{"mounted-on", "'mounted on'", reflect.String},
		{"file-system-type", "file system type", reflect.String},
		{"options", "file system type", reflect.String},
		{"dump", "number indicating whether and how often the file system should be backed up by the dump program; a zero indicates the file system will never be automatically backed up", reflect.Int},
		{"pass", "number indicating the order in which the fsck program will check the devices for errors at boot time; this is 1 for the root file system and either 2 (meaning check after root) or 0 (do not check) for all other devices", reflect.Int},
	},
	ColumnsToParse: map[string]RawDataType{},
}

// DfCommandRow represents 'df' command output row
// (See https://en.wikipedia.org/wiki/Df_(Unix)
// and https://www.gnu.org/software/coreutils/manual/html_node/df-invocation.html
// and 'df --all --sync --block-size=1024 --output=source,target,fstype,file,itotal,iavail,iused,ipcent,size,avail,used,pcent'
// and the output unit is kilobytes).
var DfCommandRow = RawData{
	IsYAML: false,
	Columns: []Column{
		{"file-system", "file system ('source')", reflect.String},
		{"mounted-on", "'mounted on' ('target')", reflect.String},
		{"file-system-type", "file system type ('fstype')", reflect.String},
		{"file", "file name if specified on the command line ('file')", reflect.String},

		{"inodes", "total number of inodes ('itotal')", reflect.Int64},
		{"ifree", "number of available inodes ('iavail')", reflect.Int64},
		{"iused", "number of used inodes ('iused')", reflect.Int64},
		{"iused-percent", "percentage of iused divided by itotal ('ipcent')", reflect.String},

		{"total-blocks", "total number of blocks ('size')", reflect.Int64},
		{"available-blocks", "number of available blocks ('avail')", reflect.Int64},
		{"used-blocks", "number of used blocks ('used')", reflect.Int64},
		{"used-percent", "percentage of iused divided by itotal ('pcent')", reflect.String},
	},
	ColumnsToParse: map[string]RawDataType{
		"total-blocks":     TypeStatus,
		"available-blocks": TypeBytes,
		"used-blocks":      TypeBytes,
	},
}

// DiskStat represents '/proc/diskstats'
// (See https://www.kernel.org/doc/Documentation/ABI/testing/procfs-diskstats
// and https://www.kernel.org/doc/Documentation/iostats.txt).
var DiskStat = RawData{
	IsYAML: false,
	Columns: []Column{
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

		{"I/O-in-progress", "only field that should go to zero (incremented as requests are on request_queue)", reflect.Uint64},
		{"time-spent-on-I/O-ms", "milliseconds spent doing I/Os", reflect.Uint64},
		{"weighted-time-spent-on-I/O-ms", "weighted milliseconds spent doing I/Os (incremented at each I/O start, I/O completion, I/O merge)", reflect.Uint64},
	},
	ColumnsToParse: map[string]RawDataType{
		"time-spent-on-reading-ms":      TypeTimeMicroseconds,
		"time-spent-on-writing-ms":      TypeTimeMicroseconds,
		"time-spent-on-I/O-ms":          TypeTimeMicroseconds,
		"weighted-time-spent-on-I/O-ms": TypeTimeMicroseconds,
	},
}

// IO represents 'proc/$PID/io'
// (See http://man7.org/linux/man-pages/man5/proc.5.html).
var IO = RawData{
	IsYAML: true,
	Columns: []Column{
		{"rchar", "number of bytes which this task has caused to be read from storage (sum of bytes which this process passed to read)", reflect.Uint64},
		{"wchar", "number of bytes which this task has caused, or shall cause to be written to disk", reflect.Uint64},
		{"syscr", "number of read I/O operations", reflect.Uint64},
		{"syscw", "number of write I/O operations", reflect.Uint64},
		{"read_bytes", "number of bytes which this process really did cause to be fetched from the storage layer", reflect.Uint64},
		{"write_bytes", "number of bytes which this process caused to be sent to the storage layer", reflect.Uint64},
		{"cancelled_write_bytes", "number of bytes which this process caused to not happen by truncating pagecache", reflect.Uint64},
	},
	ColumnsToParse: map[string]RawDataType{
		"rchar":                 TypeBytes,
		"wchar":                 TypeBytes,
		"read_bytes":            TypeBytes,
		"write_bytes":           TypeBytes,
		"cancelled_write_bytes": TypeBytes,
	},
}

// Stat represents '/proc/$PID/stat'
// (See http://man7.org/linux/man-pages/man5/proc.5.html).
var Stat = RawData{
	IsYAML: false,
	Columns: []Column{
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
	ColumnsToParse: map[string]RawDataType{
		"state":  TypeStatus,
		"vsize":  TypeBytes,
		"rss":    TypeBytes,
		"rsslim": TypeBytes,
	},
}

// Status represents 'proc/$PID/status'
// (See http://man7.org/linux/man-pages/man5/proc.5.html).
var Status = RawData{
	IsYAML: true,
	Columns: []Column{
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
	ColumnsToParse: map[string]RawDataType{
		"State":        TypeStatus,
		"VmPeak":       TypeBytes,
		"VmSize":       TypeBytes,
		"VmLck":        TypeBytes,
		"VmPin":        TypeBytes,
		"VmHWM":        TypeBytes,
		"VmRSS":        TypeBytes,
		"VmData":       TypeBytes,
		"VmStk":        TypeBytes,
		"VmExe":        TypeBytes,
		"VmLib":        TypeBytes,
		"VmPTE":        TypeBytes,
		"VmPMD":        TypeBytes,
		"VmSwap":       TypeBytes,
		"HugetlbPages": TypeBytes,
	},
}
