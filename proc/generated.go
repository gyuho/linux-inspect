package proc

// updated at 2017-12-21 12:15:57.243686 -0800 PST

// NetDev is '/proc/net/dev' in Linux.
// The dev pseudo-file contains network device status information.
type NetDev struct {
	// Interface is network interface.
	Interface string `column:"interface"`
	// ReceiveBytes is total number of bytes of data received by the interface.
	ReceiveBytes            uint64 `column:"receive_bytes"`
	ReceiveBytesBytesN      uint64 `column:"receive_bytes_bytes_n"`
	ReceiveBytesParsedBytes string `column:"receive_bytes_parsed_bytes"`
	// ReceivePackets is total number of packets of data received by the interface.
	ReceivePackets uint64 `column:"receive_packets"`
	// ReceiveErrs is total number of receive errors detected by the device driver.
	ReceiveErrs uint64 `column:"receive_errs"`
	// ReceiveDrop is total number of packets dropped by the device driver.
	ReceiveDrop uint64 `column:"receive_drop"`
	// ReceiveFifo is number of FIFO buffer errors.
	ReceiveFifo uint64 `column:"receive_fifo"`
	// ReceiveFrame is number of packet framing errors.
	ReceiveFrame uint64 `column:"receive_frame"`
	// ReceiveCompressed is number of compressed packets received by the device driver.
	ReceiveCompressed uint64 `column:"receive_compressed"`
	// ReceiveMulticast is number of multicast frames received by the device driver.
	ReceiveMulticast uint64 `column:"receive_multicast"`
	// TransmitBytes is total number of bytes of data transmitted by the interface.
	TransmitBytes            uint64 `column:"transmit_bytes"`
	TransmitBytesBytesN      uint64 `column:"transmit_bytes_bytes_n"`
	TransmitBytesParsedBytes string `column:"transmit_bytes_parsed_bytes"`
	// TransmitPackets is total number of packets of data transmitted by the interface.
	TransmitPackets uint64 `column:"transmit_packets"`
	// TransmitErrs is total number of receive errors detected by the device driver.
	TransmitErrs uint64 `column:"transmit_errs"`
	// TransmitDrop is total number of packets dropped by the device driver.
	TransmitDrop uint64 `column:"transmit_drop"`
	// TransmitFifo is number of FIFO buffer errors.
	TransmitFifo uint64 `column:"transmit_fifo"`
	// TransmitColls is number of collisions detected on the interface.
	TransmitColls uint64 `column:"transmit_colls"`
	// TransmitCarrier is number of carrier losses detected by the device driver.
	TransmitCarrier uint64 `column:"transmit_carrier"`
}

// NetTCP is '/proc/net/tcp', '/proc/net/tcp6' in Linux.
// Holds a dump of the TCP socket table.
type NetTCP struct {
	Type string `column:"type"`
	// Sl is kernel hash slot.
	Sl uint64 `column:"sl"`
	// LocalAddress is local-address:port.
	LocalAddress             string `column:"local_address"`
	LocalAddressParsedIPHost string `column:"local_address_parsed_ip_host"`
	LocalAddressParsedIPPort int64  `column:"local_address_parsed_ip_port"`
	// RemAddress is remote-address:port.
	RemAddress             string `column:"rem_address"`
	RemAddressParsedIPHost string `column:"rem_address_parsed_ip_host"`
	RemAddressParsedIPPort int64  `column:"rem_address_parsed_ip_port"`
	// St is internal status of socket.
	St             string `column:"st"`
	StParsedStatus string `column:"st_parsed_status"`
	// TxQueue is outgoing data queue in terms of kernel memory usage.
	TxQueue string `column:"tx_queue"`
	// RxQueue is incoming data queue in terms of kernel memory usage.
	RxQueue string `column:"rx_queue"`
	// Tr is internal information of the kernel socket state.
	Tr string `column:"tr"`
	// TmWhen is internal information of the kernel socket state.
	TmWhen string `column:"tm_when"`
	// Retrnsmt is internal information of the kernel socket state.
	Retrnsmt string `column:"retrnsmt"`
	// Uid is effective UID of the creator of the socket.
	Uid uint64 `column:"uid"`
	// Timeout is timeout.
	Timeout uint64 `column:"timeout"`
	// Inode is inode raw data.
	Inode string `column:"inode"`
}

// LoadAvg is '/proc/loadavg' in Linux.
type LoadAvg struct {
	// LoadAvg1Minute is total uptime in seconds.
	LoadAvg1Minute float64 `column:"load_avg_1_minute"`
	// LoadAvg5Minute is total uptime in seconds.
	LoadAvg5Minute float64 `column:"load_avg_5_minute"`
	// LoadAvg15Minute is total uptime in seconds.
	LoadAvg15Minute float64 `column:"load_avg_15_minute"`
	// RunnableKernelSchedulingEntities is number of currently runnable kernel scheduling entities (processes, threads).
	RunnableKernelSchedulingEntities int64 `column:"runnable_kernel_scheduling_entities"`
	// CurrentKernelSchedulingEntities is number of kernel scheduling entities that currently exist on the system.
	CurrentKernelSchedulingEntities int64 `column:"current_kernel_scheduling_entities"`
	// Pid is PID of the process that was most recently created on the system.
	Pid int64 `column:"pid"`
}

// Uptime is '/proc/uptime' in Linux.
type Uptime struct {
	// UptimeTotal is total uptime in seconds.
	UptimeTotal           float64 `column:"uptime_total"`
	UptimeTotalParsedTime string  `column:"uptime_total_parsed_time"`
	// UptimeIdle is total amount of time in seconds spent in idle process.
	UptimeIdle           float64 `column:"uptime_idle"`
	UptimeIdleParsedTime string  `column:"uptime_idle_parsed_time"`
}

// DiskStat is '/proc/diskstats' in Linux.
type DiskStat struct {
	// MajorNumber is major device number.
	MajorNumber uint64 `column:"major_number"`
	// MinorNumber is minor device number.
	MinorNumber uint64 `column:"minor_number"`
	// DeviceName is device name.
	DeviceName string `column:"device_name"`
	// ReadsCompleted is total number of reads completed successfully.
	ReadsCompleted uint64 `column:"reads_completed"`
	// ReadsMerged is total number of reads merged when adjacent to each other.
	ReadsMerged uint64 `column:"reads_merged"`
	// SectorsRead is total number of sectors read successfully.
	SectorsRead uint64 `column:"sectors_read"`
	// TimeSpentOnReadingMs is total number of milliseconds spent by all reads.
	TimeSpentOnReadingMs           uint64 `column:"time_spent_on_reading_ms"`
	TimeSpentOnReadingMsParsedTime string `column:"time_spent_on_reading_ms_parsed_time"`
	// WritesCompleted is total number of writes completed successfully.
	WritesCompleted uint64 `column:"writes_completed"`
	// WritesMerged is total number of writes merged when adjacent to each other.
	WritesMerged uint64 `column:"writes_merged"`
	// SectorsWritten is total number of sectors written successfully.
	SectorsWritten uint64 `column:"sectors_written"`
	// TimeSpentOnWritingMs is total number of milliseconds spent by all writes.
	TimeSpentOnWritingMs           uint64 `column:"time_spent_on_writing_ms"`
	TimeSpentOnWritingMsParsedTime string `column:"time_spent_on_writing_ms_parsed_time"`
	// IOsInProgress is only field that should go to zero (incremented as requests are on request_queue).
	IOsInProgress uint64 `column:"ios_in_progress"`
	// TimeSpentOnIOsMs is milliseconds spent doing I/Os.
	TimeSpentOnIOsMs           uint64 `column:"time_spent_on_ios_ms"`
	TimeSpentOnIOsMsParsedTime string `column:"time_spent_on_ios_ms_parsed_time"`
	// WeightedTimeSpentOnIOsMs is weighted milliseconds spent doing I/Os (incremented at each I/O start, I/O completion, I/O merge).
	WeightedTimeSpentOnIOsMs           uint64 `column:"weighted_time_spent_on_ios_ms"`
	WeightedTimeSpentOnIOsMsParsedTime string `column:"weighted_time_spent_on_ios_ms_parsed_time"`
}

// IO is '/proc/$PID/io' in Linux.
type IO struct {
	// Rchar is number of bytes which this task has caused to be read from storage (sum of bytes which this process passed to read).
	Rchar            uint64 `yaml:"rchar"`
	RcharBytesN      uint64 `yaml:"rchar_bytes_n"`
	RcharParsedBytes string `yaml:"rchar_parsed_bytes"`
	// Wchar is number of bytes which this task has caused, or shall cause to be written to disk.
	Wchar            uint64 `yaml:"wchar"`
	WcharBytesN      uint64 `yaml:"wchar_bytes_n"`
	WcharParsedBytes string `yaml:"wchar_parsed_bytes"`
	// Syscr is number of read I/O operations.
	Syscr uint64 `yaml:"syscr"`
	// Syscw is number of write I/O operations.
	Syscw uint64 `yaml:"syscw"`
	// ReadBytes is number of bytes which this process really did cause to be fetched from the storage layer.
	ReadBytes            uint64 `yaml:"read_bytes"`
	ReadBytesBytesN      uint64 `yaml:"read_bytes_bytes_n"`
	ReadBytesParsedBytes string `yaml:"read_bytes_parsed_bytes"`
	// WriteBytes is number of bytes which this process caused to be sent to the storage layer.
	WriteBytes            uint64 `yaml:"write_bytes"`
	WriteBytesBytesN      uint64 `yaml:"write_bytes_bytes_n"`
	WriteBytesParsedBytes string `yaml:"write_bytes_parsed_bytes"`
	// CancelledWriteBytes is number of bytes which this process caused to not happen by truncating pagecache.
	CancelledWriteBytes            uint64 `yaml:"cancelled_write_bytes"`
	CancelledWriteBytesBytesN      uint64 `yaml:"cancelled_write_bytes_bytes_n"`
	CancelledWriteBytesParsedBytes string `yaml:"cancelled_write_bytes_parsed_bytes"`
}

// Stat is '/proc/$PID/stat' in Linux.
type Stat struct {
	// Pid is process ID.
	Pid int64 `column:"pid"`
	// Comm is filename of the executable (originally in parentheses, automatically removed by this package).
	Comm string `column:"comm"`
	// State is one character that represents the state of the process.
	State             string `column:"state"`
	StateParsedStatus string `column:"state_parsed_status"`
	// Ppid is PID of the parent process.
	Ppid int64 `column:"ppid"`
	// Pgrp is group ID of the process.
	Pgrp int64 `column:"pgrp"`
	// Session is session ID of the process.
	Session int64 `column:"session"`
	// TtyNr is controlling terminal of the process.
	TtyNr int64 `column:"tty_nr"`
	// Tpgid is ID of the foreground process group of the controlling terminal of the process.
	Tpgid int64 `column:"tpgid"`
	// Flags is kernel flags word of the process.
	Flags int64 `column:"flags"`
	// Minflt is number of minor faults the process has made which have not required loading a memory page from disk.
	Minflt uint64 `column:"minflt"`
	// Cminflt is number of minor faults that the process's waited-for children have made.
	Cminflt uint64 `column:"cminflt"`
	// Majflt is number of major faults the process has made which have required loading a memory page from disk.
	Majflt uint64 `column:"majflt"`
	// Cmajflt is number of major faults that the process's waited-for children have made.
	Cmajflt uint64 `column:"cmajflt"`
	// Utime is number of clock ticks that this process has been scheduled in user mode (includes guest_time).
	Utime uint64 `column:"utime"`
	// Stime is number of clock ticks that this process has been scheduled in kernel mode.
	Stime uint64 `column:"stime"`
	// Cutime is number of clock ticks that this process's waited-for children have been scheduled in user mode.
	Cutime uint64 `column:"cutime"`
	// Cstime is number of clock ticks that this process's waited-for children have been scheduled in kernel mode.
	Cstime uint64 `column:"cstime"`
	// Priority is for processes running a real-time scheduling policy, the negated scheduling priority, minus one; that is, a number in the range -2 to -100, corresponding to real-time priorities 1 to 99. For processes running under a non-real-time scheduling policy, this is the raw nice value. The kernel stores nice values as numbers in the range 0 (high) to 39 (low).
	Priority int64 `column:"priority"`
	// Nice is nice value, a value in the range 19 (low priority) to -20 (high priority).
	Nice int64 `column:"nice"`
	// NumThreads is number of threads in this process.
	NumThreads int64 `column:"num_threads"`
	// Itrealvalue is no longer maintained.
	Itrealvalue int64 `column:"itrealvalue"`
	// Starttime is time(number of clock ticks) the process started after system boot.
	Starttime uint64 `column:"starttime"`
	// Vsize is virtual memory size in bytes.
	Vsize            uint64 `column:"vsize"`
	VsizeBytesN      uint64 `column:"vsize_bytes_n"`
	VsizeParsedBytes string `column:"vsize_parsed_bytes"`
	// Rss is resident set size: number of pages the process has in real memory (text, data, or stack space but does not include pages which have not been demand-loaded in, or which are swapped out).
	Rss            int64  `column:"rss"`
	RssBytesN      int64  `column:"rss_bytes_n"`
	RssParsedBytes string `column:"rss_parsed_bytes"`
	// Rsslim is current soft limit in bytes on the rss of the process.
	Rsslim            uint64 `column:"rsslim"`
	RsslimBytesN      uint64 `column:"rsslim_bytes_n"`
	RsslimParsedBytes string `column:"rsslim_parsed_bytes"`
	// Startcode is address above which program text can run.
	Startcode uint64 `column:"startcode"`
	// Endcode is address below which program text can run.
	Endcode uint64 `column:"endcode"`
	// Startstack is address of the start (i.e., bottom) of the stack.
	Startstack uint64 `column:"startstack"`
	// Kstkesp is current value of ESP (stack pointer), as found in the kernel stack page for the process.
	Kstkesp uint64 `column:"kstkesp"`
	// Kstkeip is current EIP (instruction pointer).
	Kstkeip uint64 `column:"kstkeip"`
	// Signal is obsolete, because it does not provide information on real-time signals (use /proc/$PID/status).
	Signal uint64 `column:"signal"`
	// Blocked is obsolete, because it does not provide information on real-time signals (use /proc/$PID/status).
	Blocked uint64 `column:"blocked"`
	// Sigignore is obsolete, because it does not provide information on real-time signals (use /proc/$PID/status).
	Sigignore uint64 `column:"sigignore"`
	// Sigcatch is obsolete, because it does not provide information on real-time signals (use /proc/$PID/status).
	Sigcatch uint64 `column:"sigcatch"`
	// Wchan is channel in which the process is waiting (address of a location in the kernel where the process is sleeping).
	Wchan uint64 `column:"wchan"`
	// Nswap is not maintained (number of pages swapped).
	Nswap uint64 `column:"nswap"`
	// Cnswap is not maintained (cumulative nswap for child processes).
	Cnswap uint64 `column:"cnswap"`
	// ExitSignal is signal to be sent to parent when we die.
	ExitSignal int64 `column:"exit_signal"`
	// Processor is CPU number last executed on.
	Processor int64 `column:"processor"`
	// RtPriority is real-time scheduling priority, a number in the range 1 to 99 for processes scheduled under a real-time policy, or 0, for non-real-time processes.
	RtPriority uint64 `column:"rt_priority"`
	// Policy is scheduling policy.
	Policy uint64 `column:"policy"`
	// DelayacctBlkioTicks is aggregated block I/O delays, measured in clock ticks.
	DelayacctBlkioTicks uint64 `column:"delayacct_blkio_ticks"`
	// GuestTime is number of clock ticks spent running a virtual CPU for a guest operating system.
	GuestTime uint64 `column:"guest_time"`
	// CguestTime is number of clock ticks (guest_time of the process's children).
	CguestTime uint64 `column:"cguest_time"`
	// StartData is address above which program initialized and uninitialized (BSS) data are placed.
	StartData uint64 `column:"start_data"`
	// EndData is address below which program initialized and uninitialized (BSS) data are placed.
	EndData uint64 `column:"end_data"`
	// StartBrk is address above which program heap can be expanded with brk.
	StartBrk uint64 `column:"start_brk"`
	// ArgStart is address above which program command-line arguments are placed.
	ArgStart uint64 `column:"arg_start"`
	// ArgEnd is address below program command-line arguments are placed.
	ArgEnd uint64 `column:"arg_end"`
	// EnvStart is address above which program environment is placed.
	EnvStart uint64 `column:"env_start"`
	// EnvEnd is address below which program environment is placed.
	EnvEnd uint64 `column:"env_end"`
	// ExitCode is thread's exit status in the form reported by waitpid(2).
	ExitCode int64 `column:"exit_code"`
}

// Status is '/proc/$PID/status' in Linux.
type Status struct {
	// Name is command run by this process.
	Name string `yaml:"Name"`
	// Umask is process umask, expressed in octal with a leading.
	Umask string `yaml:"Umask"`
	// State is current state of the process: R (running), S (sleeping), D (disk sleep), T (stopped), T (tracing stop), Z (zombie), or X (dead).
	State             string `yaml:"State"`
	StateParsedStatus string `yaml:"State_parsed_status"`
	// Tgid is thread group ID.
	Tgid int64 `yaml:"Tgid"`
	// Ngid is NUMA group ID.
	Ngid int64 `yaml:"Ngid"`
	// Pid is process ID.
	Pid int64 `yaml:"Pid"`
	// PPid is parent process ID, which launches the Pid.
	PPid int64 `yaml:"PPid"`
	// TracerPid is PID of process tracing this process (0 if not being traced).
	TracerPid int64 `yaml:"TracerPid"`
	// Uid is real, effective, saved set, and filesystem UIDs.
	Uid string `yaml:"Uid"`
	// Gid is real, effective, saved set, and filesystem UIDs.
	Gid string `yaml:"Gid"`
	// FDSize is number of file descriptor slots currently allocated.
	FDSize uint64 `yaml:"FDSize"`
	// Groups is supplementary group list.
	Groups string `yaml:"Groups"`
	// NStgid is thread group ID (i.e., PID) in each of the PID namespaces of which [pid] is a member.
	NStgid string `yaml:"NStgid"`
	// NSpid is thread ID (i.e., PID) in each of the PID namespaces of which [pid] is a member.
	NSpid string `yaml:"NSpid"`
	// NSpgid is process group ID (i.e., PID) in each of the PID namespaces of which [pid] is a member.
	NSpgid string `yaml:"NSpgid"`
	// NSsid is descendant namespace session ID hierarchy Session ID in each of the PID namespaces of which [pid] is a member.
	NSsid string `yaml:"NSsid"`
	// VmPeak is peak virtual memory usage. Vm includes physical memory and swap.
	VmPeak            string `yaml:"VmPeak"`
	VmPeakBytesN      uint64 `yaml:"VmPeak_bytes_n"`
	VmPeakParsedBytes string `yaml:"VmPeak_parsed_bytes"`
	// VmSize is current virtual memory usage. VmSize is the total amount of memory required for this process.
	VmSize            string `yaml:"VmSize"`
	VmSizeBytesN      uint64 `yaml:"VmSize_bytes_n"`
	VmSizeParsedBytes string `yaml:"VmSize_parsed_bytes"`
	// VmLck is locked memory size.
	VmLck            string `yaml:"VmLck"`
	VmLckBytesN      uint64 `yaml:"VmLck_bytes_n"`
	VmLckParsedBytes string `yaml:"VmLck_parsed_bytes"`
	// VmPin is pinned memory size (pages can't be moved, requires direct-access to physical memory).
	VmPin            string `yaml:"VmPin"`
	VmPinBytesN      uint64 `yaml:"VmPin_bytes_n"`
	VmPinParsedBytes string `yaml:"VmPin_parsed_bytes"`
	// VmHWM is peak resident set size ("high water mark").
	VmHWM            string `yaml:"VmHWM"`
	VmHWMBytesN      uint64 `yaml:"VmHWM_bytes_n"`
	VmHWMParsedBytes string `yaml:"VmHWM_parsed_bytes"`
	// VmRSS is resident set size. VmRSS is the actual amount in memory. Some memory can be swapped out to physical disk. So this is the real memory usage of the process.
	VmRSS            string `yaml:"VmRSS"`
	VmRSSBytesN      uint64 `yaml:"VmRSS_bytes_n"`
	VmRSSParsedBytes string `yaml:"VmRSS_parsed_bytes"`
	// VmData is size of data segment.
	VmData            string `yaml:"VmData"`
	VmDataBytesN      uint64 `yaml:"VmData_bytes_n"`
	VmDataParsedBytes string `yaml:"VmData_parsed_bytes"`
	// VmStk is size of stack.
	VmStk            string `yaml:"VmStk"`
	VmStkBytesN      uint64 `yaml:"VmStk_bytes_n"`
	VmStkParsedBytes string `yaml:"VmStk_parsed_bytes"`
	// VmExe is size of text segments.
	VmExe            string `yaml:"VmExe"`
	VmExeBytesN      uint64 `yaml:"VmExe_bytes_n"`
	VmExeParsedBytes string `yaml:"VmExe_parsed_bytes"`
	// VmLib is shared library code size.
	VmLib            string `yaml:"VmLib"`
	VmLibBytesN      uint64 `yaml:"VmLib_bytes_n"`
	VmLibParsedBytes string `yaml:"VmLib_parsed_bytes"`
	// VmPTE is page table entries size.
	VmPTE            string `yaml:"VmPTE"`
	VmPTEBytesN      uint64 `yaml:"VmPTE_bytes_n"`
	VmPTEParsedBytes string `yaml:"VmPTE_parsed_bytes"`
	// VmPMD is size of second-level page tables.
	VmPMD            string `yaml:"VmPMD"`
	VmPMDBytesN      uint64 `yaml:"VmPMD_bytes_n"`
	VmPMDParsedBytes string `yaml:"VmPMD_parsed_bytes"`
	// VmSwap is swapped-out virtual memory size by anonymous private.
	VmSwap            string `yaml:"VmSwap"`
	VmSwapBytesN      uint64 `yaml:"VmSwap_bytes_n"`
	VmSwapParsedBytes string `yaml:"VmSwap_parsed_bytes"`
	// HugetlbPages is size of hugetlb memory portions.
	HugetlbPages            string `yaml:"HugetlbPages"`
	HugetlbPagesBytesN      uint64 `yaml:"HugetlbPages_bytes_n"`
	HugetlbPagesParsedBytes string `yaml:"HugetlbPages_parsed_bytes"`
	// Threads is number of threads in process containing this thread (process).
	Threads uint64 `yaml:"Threads"`
	// SigQ is queued signals for the real user ID of this process (queued signals / limits).
	SigQ string `yaml:"SigQ"`
	// SigPnd is number of signals pending for thread.
	SigPnd string `yaml:"SigPnd"`
	// ShdPnd is number of signals pending for process as a whole.
	ShdPnd string `yaml:"ShdPnd"`
	// SigBlk is masks indicating signals being blocked.
	SigBlk string `yaml:"SigBlk"`
	// SigIgn is masks indicating signals being ignored.
	SigIgn string `yaml:"SigIgn"`
	// SigCgt is masks indicating signals being caught.
	SigCgt string `yaml:"SigCgt"`
	// CapInh is masks of capabilities enabled in inheritable sets.
	CapInh string `yaml:"CapInh"`
	// CapPrm is masks of capabilities enabled in permitted sets.
	CapPrm string `yaml:"CapPrm"`
	// CapEff is masks of capabilities enabled in effective sets.
	CapEff string `yaml:"CapEff"`
	// CapBnd is capability Bounding set.
	CapBnd string `yaml:"CapBnd"`
	// CapAmb is ambient capability set.
	CapAmb string `yaml:"CapAmb"`
	// Seccomp is seccomp mode of the process (0 means SECCOMP_MODE_DISABLED; 1 means SECCOMP_MODE_STRICT; 2 means SECCOMP_MODE_FILTER).
	Seccomp uint64 `yaml:"Seccomp"`
	// CpusAllowed is mask of CPUs on which this process may run.
	CpusAllowed string `yaml:"Cpus_allowed"`
	// CpusAllowedList is list of CPUs on which this process may run.
	CpusAllowedList string `yaml:"Cpus_allowed_list"`
	// MemsAllowed is mask of memory nodes allowed to this process.
	MemsAllowed string `yaml:"Mems_allowed"`
	// MemsAllowedList is list of memory nodes allowed to this process.
	MemsAllowedList string `yaml:"Mems_allowed_list"`
	// VoluntaryCtxtSwitches is number of voluntary context switches.
	VoluntaryCtxtSwitches uint64 `yaml:"voluntary_ctxt_switches"`
	// NonvoluntaryCtxtSwitches is number of involuntary context switches.
	NonvoluntaryCtxtSwitches uint64 `yaml:"nonvoluntary_ctxt_switches"`
}
