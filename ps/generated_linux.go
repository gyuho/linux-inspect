package ps

// updated at 2016-03-16 23:40:56.693947193 -0700 PDT

// Stat is metrics in linux proc/$PID/stat.
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
