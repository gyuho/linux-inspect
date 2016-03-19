package ps

import "reflect"

// StatList represents the proc/$PID/stat file
// specification as documented in http://man7.org/linux/man-pages/man5/proc.5.html.
var StatList = [...]struct {
	Col      string
	Kind     reflect.Kind
	Humanize bool // true if humanized string is needed
}{
	{"pid", reflect.Int64, false},    // the process ID
	{"comm", reflect.String, false},  // filename of the executable (originally in parentheses, automatically removed by this package)
	{"state", reflect.String, false}, // One character that represents the state of the process
	{"ppid", reflect.Int64, false},   // PID of the parent of this process
	{"pgrp", reflect.Int64, false},   // process group ID of the process
	{"session", reflect.Int64, false},
	{"tty_nr", reflect.Int64, false},
	{"tpgid", reflect.Int64, false}, // ID of the foreground process group of the controlling terminal of the process
	{"flags", reflect.Uint64, false},
	{"minflt", reflect.Uint64, false},  // number of minor faults the process has made which have not required loading a memory page from disk.
	{"cminflt", reflect.Uint64, false}, // number of minor faults that the process's waited-for children have made.
	{"majflt", reflect.Uint64, false},  // number of major faults the process has made which have required loading a memory page from disk.
	{"cmajflt", reflect.Uint64, false}, // number of major faults that the process's waited-for children have made.
	{"utime", reflect.Uint64, false},   // Amount of time that this process has been scheduled in user mode, measured in clock ticks.
	{"stime", reflect.Uint64, false},   // Amount of time that this process has been scheduled in kernel mode, measured in clock ticks.
	{"cutime", reflect.Uint64, false},  // Amount of time that this process's waited-for children have been scheduled in user mode.
	{"cstime", reflect.Uint64, false},  // Amount of time that this process's waited-for children have been scheduled in kernel mode.
	{"priority", reflect.Int64, false},
	{"nice", reflect.Int64, false},
	{"num_threads", reflect.Int64, false},
	{"itrealvalue", reflect.Int64, false},
	{"starttime", reflect.Uint64, false}, // time the process started after system boot.
	{"vsize", reflect.Uint64, true},      // Virtual memory size in bytes.
	{"rss", reflect.Int64, true},         // Resident Set Size: number of pages the process has in real memory.
	{"rsslim", reflect.Uint64, true},
	{"startcode", reflect.Uint64, false},
	{"endcode", reflect.Uint64, false},
	{"startstack", reflect.Uint64, false},
	{"kstkesp", reflect.Uint64, false},
	{"kstkeip", reflect.Uint64, false},
	{"signal", reflect.Uint64, false},
	{"blocked", reflect.Uint64, false},
	{"sigignore", reflect.Uint64, false},
	{"sigcatch", reflect.Uint64, false},
	{"wchan", reflect.Uint64, false},
	{"nswap", reflect.Uint64, false},
	{"cnswap", reflect.Uint64, false},
	{"exit_signal", reflect.Int64, false},
	{"processor", reflect.Int64, false}, // CPU number last executed on.
	{"rt_priority", reflect.Uint64, false},
	{"policy", reflect.Uint64, false},
	{"delayacct_blkio_ticks", reflect.Uint64, false},
	{"guest_time", reflect.Uint64, false},
	{"cguest_time", reflect.Int64, false},
	{"start_data", reflect.Uint64, false},
	{"end_data", reflect.Uint64, false},
	{"start_brk", reflect.Uint64, false},
	{"arg_start", reflect.Uint64, false},
	{"arg_end", reflect.Uint64, false},
	{"env_start", reflect.Uint64, false},
	{"env_end", reflect.Uint64, false},
	{"exit_code", reflect.Int64, false},
}

// StatusListYAML represents the proc/$PID/status file
// specification as documented in http://man7.org/linux/man-pages/man5/proc.5.html.
var StatusListYAML = [...]struct {
	Col   string
	Kind  reflect.Kind
	Bytes bool // true if parsing humanized string in bytes is needed
}{
	{"Name", reflect.String, false}, // Name is the command run by this process.

	// State is Current state of the process.
	// One of "R (running)", "S (sleeping)", "D (disk sleep)",
	// "T (stopped)", "T (tracing stop)", "Z (zombie)", or "X (dead)".
	{"State", reflect.String, false},

	{"Tgid", reflect.Int64, false}, // Tgid is thread group ID.
	{"Ngid", reflect.Int64, false},
	{"Pid", reflect.Int64, false},       // Pid is process ID.
	{"PPid", reflect.Int64, false},      // PPid is parent process ID, which launches the Pid.
	{"TracerPid", reflect.Int64, false}, // TracerPid is PID of process tracing this process (0 if not being traced).

	{"Uid", reflect.String, false},
	{"Gid", reflect.String, false},

	{"FDSize", reflect.Uint64, false}, // FDSize is the number of file descriptor slots currently allocated.

	{"Groups", reflect.String, false}, // Groups is supplementary group list.

	{"VmPeak", reflect.String, true}, // VmPeak is peak virtual memory usage. Vm includes physical memory and swap.
	{"VmSize", reflect.String, true}, // VmSize is current virtual memory usage. VmSize is the total amount of memory required for this process.
	{"VmLck", reflect.String, true},  // VmLck is current mlocked memory.
	{"VmPin", reflect.String, true},  // VmPin is pinned memory size.
	{"VmHWM", reflect.String, true},  // VmHWM is peak resident set size ("high water mark").

	// VmRSS is resident set size. VmRSS is the actual
	// amount in memory. Some memory can be swapped out
	// to physical disk. So this is the real memory usage
	// of the process.
	{"VmRSS", reflect.String, true},

	{"VmData", reflect.String, true}, // VmRSS is size of data segment.

	{"VmStk", reflect.String, true}, // VmStk is size of stack.
	{"VmExe", reflect.String, true}, // VmExe is size of text segment.
	{"VmLib", reflect.String, true}, // VmLib is shared library usage.
	{"VmPMD", reflect.String, true},
	{"VmPTE", reflect.String, true},  // VmPTE is page table entries size.
	{"VmSwap", reflect.String, true}, // VmSwap is swap space used.

	{"Threads", reflect.Uint64, false}, // Threads is the number of threads in process containing this thread (process).

	{"SigQ", reflect.String, false},
	{"SigPnd", reflect.String, false},
	{"ShdPnd", reflect.String, false},
	{"SigBlk", reflect.String, false},
	{"SigIgn", reflect.String, false},
	{"SigCgt", reflect.String, false},
	{"CapInh", reflect.String, false},
	{"CapPrm", reflect.String, false},
	{"CapEff", reflect.String, false},
	{"CapBnd", reflect.String, false},
	{"CapAmb", reflect.String, false},
	{"Seccomp", reflect.Uint64, false},

	{"Cpus_allowed", reflect.String, false},
	{"Cpus_allowed_list", reflect.String, false},
	{"Mems_allowed", reflect.String, false},
	{"Mems_allowed_list", reflect.String, false},

	{"voluntary_ctxt_switches", reflect.Uint64, false},
	{"nonvoluntary_ctxt_switches", reflect.Uint64, false},
}
