// Package schema defines proc schema.
package schema

import "reflect"

// Column represents the schema column.
type Column struct {
	Name string
	Kind reflect.Kind

	// true to humanize float64
	HumanizedSeconds bool

	// true to humanize bytes string
	HumanizedBytes bool
	YAMLTag        bool
}

// Uptime represents '/proc/uptime'
// (See http://man7.org/linux/man-pages/man5/proc.5.html).
var Uptime = []Column{
	{"uptime-total", reflect.Float64, true, false, false},
	{"uptime-idle", reflect.Float64, true, false, false},
}

// DiskStats represents '/proc/diskstats'
// (See https://www.kernel.org/doc/Documentation/ABI/testing/procfs-diskstats
// and https://www.kernel.org/doc/Documentation/iostats.txt).
var DiskStats = []Column{
	{"major-number", reflect.Uint64, false, false, false}, // 1
	{"minor-number", reflect.Uint64, false, false, false}, // 2

	{"device-name", reflect.String, false, false, false}, // 3

	{"reads-completed", reflect.Uint64, false, false, false},          // 4
	{"reads-merged", reflect.Uint64, false, false, false},             // 5
	{"sectors-read", reflect.Uint64, false, false, false},             // 6
	{"time-spent-on-reading-ms", reflect.Uint64, false, false, false}, // 7

	{"writes-completed", reflect.Uint64, false, false, false},         // 8
	{"writes-merged", reflect.Uint64, false, false, false},            // 9
	{"sectors-written", reflect.Uint64, false, false, false},          // 10
	{"time-spent-on-writing-ms", reflect.Uint64, false, false, false}, // 11

	{"i/o-in-progress", reflect.Uint64, false, false, false},               // 12
	{"time-spent-on-i/o-ms", reflect.Uint64, false, false, false},          // 13
	{"weighted-time-spent-on-i/o-ms", reflect.Uint64, false, false, false}, // 14
}

// Stat represents '/proc/$PID/stat'
// (See http://man7.org/linux/man-pages/man5/proc.5.html).
var Stat = []Column{
	{"pid", reflect.Int64, false, false, false},    // the process ID
	{"comm", reflect.String, false, false, false},  // filename of the executable (originally in parentheses, automatically removed by this package)
	{"state", reflect.String, false, false, false}, // One character that represents the state of the process
	{"ppid", reflect.Int64, false, false, false},   // PID of the parent of this process
	{"pgrp", reflect.Int64, false, false, false},   // process group ID of the process
	{"session", reflect.Int64, false, false, false},
	{"tty_nr", reflect.Int64, false, false, false},
	{"tpgid", reflect.Int64, false, false, false}, // ID of the foreground process group of the controlling terminal of the process
	{"flags", reflect.Uint64, false, false, false},
	{"minflt", reflect.Uint64, false, false, false},  // number of minor faults the process has made which have not required loading a memory page from disk.
	{"cminflt", reflect.Uint64, false, false, false}, // number of minor faults that the process's waited-for children have made.
	{"majflt", reflect.Uint64, false, false, false},  // number of major faults the process has made which have required loading a memory page from disk.
	{"cmajflt", reflect.Uint64, false, false, false}, // number of major faults that the process's waited-for children have made.
	{"utime", reflect.Uint64, false, false, false},   // Amount of time that this process has been scheduled in user mode, measured in clock ticks.
	{"stime", reflect.Uint64, false, false, false},   // Amount of time that this process has been scheduled in kernel mode, measured in clock ticks.
	{"cutime", reflect.Uint64, false, false, false},  // Amount of time that this process's waited-for children have been scheduled in user mode.
	{"cstime", reflect.Uint64, false, false, false},  // Amount of time that this process's waited-for children have been scheduled in kernel mode.
	{"priority", reflect.Int64, false, false, false},
	{"nice", reflect.Int64, false, false, false},
	{"num_threads", reflect.Int64, false, false, false},
	{"itrealvalue", reflect.Int64, false, false, false},
	{"starttime", reflect.Uint64, false, false, false}, // time the process started after system boot.
	{"vsize", reflect.Uint64, false, true, false},      // Virtual memory size in bytes.
	{"rss", reflect.Int64, false, true, false},         // Resident Set Size: number of pages the process has in real memory.
	{"rsslim", reflect.Uint64, false, true, false},
	{"startcode", reflect.Uint64, false, false, false},
	{"endcode", reflect.Uint64, false, false, false},
	{"startstack", reflect.Uint64, false, false, false},
	{"kstkesp", reflect.Uint64, false, false, false},
	{"kstkeip", reflect.Uint64, false, false, false},
	{"signal", reflect.Uint64, false, false, false},
	{"blocked", reflect.Uint64, false, false, false},
	{"sigignore", reflect.Uint64, false, false, false},
	{"sigcatch", reflect.Uint64, false, false, false},
	{"wchan", reflect.Uint64, false, false, false},
	{"nswap", reflect.Uint64, false, false, false},
	{"cnswap", reflect.Uint64, false, false, false},
	{"exit_signal", reflect.Int64, false, false, false},
	{"processor", reflect.Int64, false, false, false}, // CPU number last executed on.
	{"rt_priority", reflect.Uint64, false, false, false},
	{"policy", reflect.Uint64, false, false, false},
	{"delayacct_blkio_ticks", reflect.Uint64, false, false, false},
	{"guest_time", reflect.Uint64, false, false, false},
	{"cguest_time", reflect.Int64, false, false, false},
	{"start_data", reflect.Uint64, false, false, false},
	{"end_data", reflect.Uint64, false, false, false},
	{"start_brk", reflect.Uint64, false, false, false},
	{"arg_start", reflect.Uint64, false, false, false},
	{"arg_end", reflect.Uint64, false, false, false},
	{"env_start", reflect.Uint64, false, false, false},
	{"env_end", reflect.Uint64, false, false, false},
	{"exit_code", reflect.Int64, false, false, false},
}

// Status represents 'proc/$PID/status'
// (See http://man7.org/linux/man-pages/man5/proc.5.html).
var Status = []Column{
	{"Name", reflect.String, false, false, true}, // Name is the command run by this process.

	// State is Current state of the process.
	// One of "R (running)", "S (sleeping)", "D (disk sleep)",
	// "T (stopped)", "T (tracing stop)", "Z (zombie)", or "X (dead)".
	{"State", reflect.String, false, false, true},

	{"Tgid", reflect.Int64, false, false, true}, // Tgid is thread group ID.
	{"Ngid", reflect.Int64, false, false, true},
	{"Pid", reflect.Int64, false, false, true},       // Pid is process ID.
	{"PPid", reflect.Int64, false, false, true},      // PPid is parent process ID, which launches the Pid.
	{"TracerPid", reflect.Int64, false, false, true}, // TracerPid is PID of process tracing this process (0 if not being traced).

	{"Uid", reflect.String, false, false, true},
	{"Gid", reflect.String, false, false, true},

	{"FDSize", reflect.Uint64, false, false, true}, // FDSize is the number of file descriptor slots currently allocated.

	{"Groups", reflect.String, false, false, true}, // Groups is supplementary group list.

	{"VmPeak", reflect.String, false, true, true}, // VmPeak is peak virtual memory usage. Vm includes physical memory and swap.
	{"VmSize", reflect.String, false, true, true}, // VmSize is current virtual memory usage. VmSize is the total amount of memory required for this process.
	{"VmLck", reflect.String, false, true, true},  // VmLck is current mlocked memory.
	{"VmPin", reflect.String, false, true, true},  // VmPin is pinned memory size.
	{"VmHWM", reflect.String, false, true, true},  // VmHWM is peak resident set size ("high water mark").

	// VmRSS is resident set size. VmRSS is the actual
	// amount in memory. Some memory can be swapped out
	// to physical disk. So this is the real memory usage
	// of the process.
	{"VmRSS", reflect.String, false, true, true},

	{"VmData", reflect.String, false, true, true}, // VmRSS is size of data segment.

	{"VmStk", reflect.String, false, true, true}, // VmStk is size of stack.
	{"VmExe", reflect.String, false, true, true}, // VmExe is size of text segment.
	{"VmLib", reflect.String, false, true, true}, // VmLib is shared library usage.
	{"VmPMD", reflect.String, false, true, true},
	{"VmPTE", reflect.String, false, true, true},  // VmPTE is page table entries size.
	{"VmSwap", reflect.String, false, true, true}, // VmSwap is swap space used.

	{"Threads", reflect.Uint64, false, false, true}, // Threads is the number of threads in process containing this thread (process).

	{"SigQ", reflect.String, false, false, true},
	{"SigPnd", reflect.String, false, false, true},
	{"ShdPnd", reflect.String, false, false, true},
	{"SigBlk", reflect.String, false, false, true},
	{"SigIgn", reflect.String, false, false, true},
	{"SigCgt", reflect.String, false, false, true},
	{"CapInh", reflect.String, false, false, true},
	{"CapPrm", reflect.String, false, false, true},
	{"CapEff", reflect.String, false, false, true},
	{"CapBnd", reflect.String, false, false, true},
	{"CapAmb", reflect.String, false, false, true},
	{"Seccomp", reflect.Uint64, false, false, true},

	{"Cpus_allowed", reflect.String, false, false, true},
	{"Cpus_allowed_list", reflect.String, false, false, true},
	{"Mems_allowed", reflect.String, false, false, true},
	{"Mems_allowed_list", reflect.String, false, false, true},

	{"voluntary_ctxt_switches", reflect.Uint64, false, false, true},
	{"nonvoluntary_ctxt_switches", reflect.Uint64, false, false, true},
}

// IO represents 'proc/$PID/io'
// (See http://man7.org/linux/man-pages/man5/proc.5.html).
var IO = []Column{
	// number of bytes which this task has caused to be read from storage
	// sum of bytes which this process passed to read(2)
	{"rchar", reflect.Uint64, false, true, true},

	// number of bytes which this task has caused, or shall cause to be written to disk
	{"wchar", reflect.Uint64, false, true, true},

	// number of read I/O operations
	{"syscr", reflect.Uint64, false, false, true},

	// number of write I/O operations
	{"syscw", reflect.Uint64, false, false, true},

	// number of bytes which this process really did cause to be fetched
	// from the storage layer
	{"read_bytes", reflect.Uint64, false, true, true},

	// number of bytes which this process
	// caused to be sent to the storage layer
	{"write_bytes", reflect.Uint64, false, true, true},

	// number of bytes which this process caused to not happen
	// by truncating pagecache
	{"cancelled_write_bytes", reflect.Uint64, false, true, true},
}
