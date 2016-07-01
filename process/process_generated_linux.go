package process

// updated at 2016-03-19 11:09:30.73210258 -0700 PDT

// Process represents '/proc' in linux.
type Process struct {
	Stat   Stat
	Status Status
}

// Stat is 'proc/$PID/stat' in linux.
type Stat struct {
	Pid                 int64   `column:"pid"`
	Comm                string  `column:"comm"`
	State               string  `column:"state"`
	Ppid                int64   `column:"ppid"`
	Pgrp                int64   `column:"pgrp"`
	Session             int64   `column:"session"`
	TtyNr               int64   `column:"tty_nr"`
	Tpgid               int64   `column:"tpgid"`
	Flags               uint64  `column:"flags"`
	Minflt              uint64  `column:"minflt"`
	Cminflt             uint64  `column:"cminflt"`
	Majflt              uint64  `column:"majflt"`
	Cmajflt             uint64  `column:"cmajflt"`
	Utime               uint64  `column:"utime"`
	Stime               uint64  `column:"stime"`
	Cutime              uint64  `column:"cutime"`
	Cstime              uint64  `column:"cstime"`
	Priority            int64   `column:"priority"`
	Nice                int64   `column:"nice"`
	NumThreads          int64   `column:"num_threads"`
	Itrealvalue         int64   `column:"itrealvalue"`
	Starttime           uint64  `column:"starttime"`
	Vsize               uint64  `column:"vsize"`
	VsizeHumanize       string  `column:"vsize_humanized"`
	Rss                 int64   `column:"rss"`
	RssHumanize         string  `column:"rss_humanized"`
	Rsslim              uint64  `column:"rsslim"`
	RsslimHumanize      string  `column:"rsslim_humanized"`
	Startcode           uint64  `column:"startcode"`
	Endcode             uint64  `column:"endcode"`
	Startstack          uint64  `column:"startstack"`
	Kstkesp             uint64  `column:"kstkesp"`
	Kstkeip             uint64  `column:"kstkeip"`
	Signal              uint64  `column:"signal"`
	Blocked             uint64  `column:"blocked"`
	Sigignore           uint64  `column:"sigignore"`
	Sigcatch            uint64  `column:"sigcatch"`
	Wchan               uint64  `column:"wchan"`
	Nswap               uint64  `column:"nswap"`
	Cnswap              uint64  `column:"cnswap"`
	ExitSignal          int64   `column:"exit_signal"`
	Processor           int64   `column:"processor"`
	RtPriority          uint64  `column:"rt_priority"`
	Policy              uint64  `column:"policy"`
	DelayacctBlkioTicks uint64  `column:"delayacct_blkio_ticks"`
	GuestTime           uint64  `column:"guest_time"`
	CguestTime          int64   `column:"cguest_time"`
	StartData           uint64  `column:"start_data"`
	EndData             uint64  `column:"end_data"`
	StartBrk            uint64  `column:"start_brk"`
	ArgStart            uint64  `column:"arg_start"`
	ArgEnd              uint64  `column:"arg_end"`
	EnvStart            uint64  `column:"env_start"`
	EnvEnd              uint64  `column:"env_end"`
	ExitCode            int64   `column:"exit_code"`
	CpuUsage            float64 `column:"cpu_usage"`
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
	VmPeakBytes              uint64 `yaml:"VmPeak_bytes"`
	VmSize                   string `yaml:"VmSize"`
	VmSizeBytes              uint64 `yaml:"VmSize_bytes"`
	VmLck                    string `yaml:"VmLck"`
	VmLckBytes               uint64 `yaml:"VmLck_bytes"`
	VmPin                    string `yaml:"VmPin"`
	VmPinBytes               uint64 `yaml:"VmPin_bytes"`
	VmHWM                    string `yaml:"VmHWM"`
	VmHWMBytes               uint64 `yaml:"VmHWM_bytes"`
	VmRSS                    string `yaml:"VmRSS"`
	VmRSSBytes               uint64 `yaml:"VmRSS_bytes"`
	VmData                   string `yaml:"VmData"`
	VmDataBytes              uint64 `yaml:"VmData_bytes"`
	VmStk                    string `yaml:"VmStk"`
	VmStkBytes               uint64 `yaml:"VmStk_bytes"`
	VmExe                    string `yaml:"VmExe"`
	VmExeBytes               uint64 `yaml:"VmExe_bytes"`
	VmLib                    string `yaml:"VmLib"`
	VmLibBytes               uint64 `yaml:"VmLib_bytes"`
	VmPMD                    string `yaml:"VmPMD"`
	VmPMDBytes               uint64 `yaml:"VmPMD_bytes"`
	VmPTE                    string `yaml:"VmPTE"`
	VmPTEBytes               uint64 `yaml:"VmPTE_bytes"`
	VmSwap                   string `yaml:"VmSwap"`
	VmSwapBytes              uint64 `yaml:"VmSwap_bytes"`
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
