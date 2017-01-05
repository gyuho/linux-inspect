package proc

// updated at 2017-01-04 17:09:34.981980906 -0800 PST

// Proc represents '/proc' in linux.
type Proc struct {
	DiskStats DiskStats
	Stat      Stat
	Status    Status
	IO        IO
}

// Uptime is 'proc/uptime' in linux.
type Uptime struct {
	UptimeTotal              float64 `column:"uptime-total"`
	UptimeTotalHumanizedTime string  `column:"uptime-total_humanized_time"`
	UptimeIdle               float64 `column:"uptime-idle"`
	UptimeIdleHumanizedTime  string  `column:"uptime-idle_humanized_time"`
}

// DiskStats is 'proc/diskstats' in linux.
type DiskStats struct {
	MajorNumber             uint64 `column:"major-number"`
	MinorNumber             uint64 `column:"minor-number"`
	DeviceName              string `column:"device-name"`
	ReadsCompleted          uint64 `column:"reads-completed"`
	ReadsMerged             uint64 `column:"reads-merged"`
	SectorsRead             uint64 `column:"sectors-read"`
	TimeSpentOnReadingMs    uint64 `column:"time-spent-on-reading-ms"`
	WritesCompleted         uint64 `column:"writes-completed"`
	WritesMerged            uint64 `column:"writes-merged"`
	SectorsWritten          uint64 `column:"sectors-written"`
	TimeSpentOnWritingMs    uint64 `column:"time-spent-on-writing-ms"`
	IoInProgress            uint64 `column:"i/o-in-progress"`
	TimeSpentOnIoMs         uint64 `column:"time-spent-on-i/o-ms"`
	WeightedTimeSpentOnIoMs uint64 `column:"weighted-time-spent-on-i/o-ms"`
}

// Stat is 'proc/$PID/stat' in linux.
type Stat struct {
	Pid                  int64   `column:"pid"`
	Comm                 string  `column:"comm"`
	State                string  `column:"state"`
	Ppid                 int64   `column:"ppid"`
	Pgrp                 int64   `column:"pgrp"`
	Session              int64   `column:"session"`
	TtyNr                int64   `column:"tty_nr"`
	Tpgid                int64   `column:"tpgid"`
	Flags                uint64  `column:"flags"`
	Minflt               uint64  `column:"minflt"`
	Cminflt              uint64  `column:"cminflt"`
	Majflt               uint64  `column:"majflt"`
	Cmajflt              uint64  `column:"cmajflt"`
	Utime                uint64  `column:"utime"`
	Stime                uint64  `column:"stime"`
	Cutime               uint64  `column:"cutime"`
	Cstime               uint64  `column:"cstime"`
	Priority             int64   `column:"priority"`
	Nice                 int64   `column:"nice"`
	NumThreads           int64   `column:"num_threads"`
	Itrealvalue          int64   `column:"itrealvalue"`
	Starttime            uint64  `column:"starttime"`
	Vsize                uint64  `column:"vsize"`
	VsizeHumanizedBytes  string  `column:"vsize_humanized_bytes"`
	Rss                  int64   `column:"rss"`
	RssHumanizedBytes    string  `column:"rss_humanized_bytes"`
	Rsslim               uint64  `column:"rsslim"`
	RsslimHumanizedBytes string  `column:"rsslim_humanized_bytes"`
	Startcode            uint64  `column:"startcode"`
	Endcode              uint64  `column:"endcode"`
	Startstack           uint64  `column:"startstack"`
	Kstkesp              uint64  `column:"kstkesp"`
	Kstkeip              uint64  `column:"kstkeip"`
	Signal               uint64  `column:"signal"`
	Blocked              uint64  `column:"blocked"`
	Sigignore            uint64  `column:"sigignore"`
	Sigcatch             uint64  `column:"sigcatch"`
	Wchan                uint64  `column:"wchan"`
	Nswap                uint64  `column:"nswap"`
	Cnswap               uint64  `column:"cnswap"`
	ExitSignal           int64   `column:"exit_signal"`
	Processor            int64   `column:"processor"`
	RtPriority           uint64  `column:"rt_priority"`
	Policy               uint64  `column:"policy"`
	DelayacctBlkioTicks  uint64  `column:"delayacct_blkio_ticks"`
	GuestTime            uint64  `column:"guest_time"`
	CguestTime           int64   `column:"cguest_time"`
	StartData            uint64  `column:"start_data"`
	EndData              uint64  `column:"end_data"`
	StartBrk             uint64  `column:"start_brk"`
	ArgStart             uint64  `column:"arg_start"`
	ArgEnd               uint64  `column:"arg_end"`
	EnvStart             uint64  `column:"env_start"`
	EnvEnd               uint64  `column:"env_end"`
	ExitCode             int64   `column:"exit_code"`
	CpuUsage             float64 `column:"cpu_usage"`
}

// Status is 'proc/$PID/status' in linux.
type Status struct {
	Name                     string `yaml:"Name"`
	State                    string `yaml:"State"`
	Tgid                     int64  `yaml:"Tgid"`
	Ngid                     int64  `yaml:"Ngid"`
	Pid                      int64  `yaml:"Pid"`
	PPid                     int64  `yaml:"PPid"`
	TracerPid                int64  `yaml:"TracerPid"`
	Uid                      string `yaml:"Uid"`
	Gid                      string `yaml:"Gid"`
	FDSize                   uint64 `yaml:"FDSize"`
	Groups                   string `yaml:"Groups"`
	VmPeak                   string `yaml:"VmPeak"`
	VmPeakBytesN             uint64 `yaml:"VmPeak_bytes_n"`
	VmPeakHumanizedBytes     string `yaml:"VmPeak_humanized_bytes"`
	VmSize                   string `yaml:"VmSize"`
	VmSizeBytesN             uint64 `yaml:"VmSize_bytes_n"`
	VmSizeHumanizedBytes     string `yaml:"VmSize_humanized_bytes"`
	VmLck                    string `yaml:"VmLck"`
	VmLckBytesN              uint64 `yaml:"VmLck_bytes_n"`
	VmLckHumanizedBytes      string `yaml:"VmLck_humanized_bytes"`
	VmPin                    string `yaml:"VmPin"`
	VmPinBytesN              uint64 `yaml:"VmPin_bytes_n"`
	VmPinHumanizedBytes      string `yaml:"VmPin_humanized_bytes"`
	VmHWM                    string `yaml:"VmHWM"`
	VmHWMBytesN              uint64 `yaml:"VmHWM_bytes_n"`
	VmHWMHumanizedBytes      string `yaml:"VmHWM_humanized_bytes"`
	VmRSS                    string `yaml:"VmRSS"`
	VmRSSBytesN              uint64 `yaml:"VmRSS_bytes_n"`
	VmRSSHumanizedBytes      string `yaml:"VmRSS_humanized_bytes"`
	VmData                   string `yaml:"VmData"`
	VmDataBytesN             uint64 `yaml:"VmData_bytes_n"`
	VmDataHumanizedBytes     string `yaml:"VmData_humanized_bytes"`
	VmStk                    string `yaml:"VmStk"`
	VmStkBytesN              uint64 `yaml:"VmStk_bytes_n"`
	VmStkHumanizedBytes      string `yaml:"VmStk_humanized_bytes"`
	VmExe                    string `yaml:"VmExe"`
	VmExeBytesN              uint64 `yaml:"VmExe_bytes_n"`
	VmExeHumanizedBytes      string `yaml:"VmExe_humanized_bytes"`
	VmLib                    string `yaml:"VmLib"`
	VmLibBytesN              uint64 `yaml:"VmLib_bytes_n"`
	VmLibHumanizedBytes      string `yaml:"VmLib_humanized_bytes"`
	VmPMD                    string `yaml:"VmPMD"`
	VmPMDBytesN              uint64 `yaml:"VmPMD_bytes_n"`
	VmPMDHumanizedBytes      string `yaml:"VmPMD_humanized_bytes"`
	VmPTE                    string `yaml:"VmPTE"`
	VmPTEBytesN              uint64 `yaml:"VmPTE_bytes_n"`
	VmPTEHumanizedBytes      string `yaml:"VmPTE_humanized_bytes"`
	VmSwap                   string `yaml:"VmSwap"`
	VmSwapBytesN             uint64 `yaml:"VmSwap_bytes_n"`
	VmSwapHumanizedBytes     string `yaml:"VmSwap_humanized_bytes"`
	Threads                  uint64 `yaml:"Threads"`
	SigQ                     string `yaml:"SigQ"`
	SigPnd                   string `yaml:"SigPnd"`
	ShdPnd                   string `yaml:"ShdPnd"`
	SigBlk                   string `yaml:"SigBlk"`
	SigIgn                   string `yaml:"SigIgn"`
	SigCgt                   string `yaml:"SigCgt"`
	CapInh                   string `yaml:"CapInh"`
	CapPrm                   string `yaml:"CapPrm"`
	CapEff                   string `yaml:"CapEff"`
	CapBnd                   string `yaml:"CapBnd"`
	CapAmb                   string `yaml:"CapAmb"`
	Seccomp                  uint64 `yaml:"Seccomp"`
	CpusAllowed              string `yaml:"Cpus_allowed"`
	CpusAllowedList          string `yaml:"Cpus_allowed_list"`
	MemsAllowed              string `yaml:"Mems_allowed"`
	MemsAllowedList          string `yaml:"Mems_allowed_list"`
	VoluntaryCtxtSwitches    uint64 `yaml:"voluntary_ctxt_switches"`
	NonvoluntaryCtxtSwitches uint64 `yaml:"nonvoluntary_ctxt_switches"`
}

// IO is 'proc/$PID/io' in linux.
type IO struct {
	Rchar                             uint64 `yaml:"rchar"`
	RcharHumanizedBytes               string `yaml:"rchar_humanized_bytes"`
	Wchar                             uint64 `yaml:"wchar"`
	WcharHumanizedBytes               string `yaml:"wchar_humanized_bytes"`
	Syscr                             uint64 `yaml:"syscr"`
	Syscw                             uint64 `yaml:"syscw"`
	ReadBytes                         uint64 `yaml:"read_bytes"`
	ReadBytesHumanizedBytes           string `yaml:"read_bytes_humanized_bytes"`
	WriteBytes                        uint64 `yaml:"write_bytes"`
	WriteBytesHumanizedBytes          string `yaml:"write_bytes_humanized_bytes"`
	CancelledWriteBytes               uint64 `yaml:"cancelled_write_bytes"`
	CancelledWriteBytesHumanizedBytes string `yaml:"cancelled_write_bytes_humanized_bytes"`
}
