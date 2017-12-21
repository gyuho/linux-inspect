package top

// updated at 2017-12-21 12:15:58.06223 -0800 PST

// Row represents a row in 'top' command output.
type Row struct {
	// PID is pid of the process.
	PID int64 `column:"pid"`
	// USER is user name.
	USER string `column:"user"`
	// PR is priority.
	PR string `column:"pr"`
	// NI is nice value of the task.
	NI string `column:"ni"`
	// VIRT is total amount  of virtual memory used by the task (in KiB).
	VIRT            string `column:"virt"`
	VIRTBytesN      uint64 `column:"virt_bytes_n"`
	VIRTParsedBytes string `column:"virt_parsed_bytes"`
	// RES is non-swapped physical memory a task is using (in KiB).
	RES            string `column:"res"`
	RESBytesN      uint64 `column:"res_bytes_n"`
	RESParsedBytes string `column:"res_parsed_bytes"`
	// SHR is amount of shared memory available to a task, not all of which is typically resident (in KiB).
	SHR            string `column:"shr"`
	SHRBytesN      uint64 `column:"shr_bytes_n"`
	SHRParsedBytes string `column:"shr_parsed_bytes"`
	// S is process status.
	S             string `column:"s"`
	SParsedStatus string `column:"s_parsed_status"`
	// CPUPercent is %CPU.
	CPUPercent float64 `column:"cpupercent"`
	// MEMPercent is %MEM.
	MEMPercent float64 `column:"mempercent"`
	// TIME is CPU time (TIME+).
	TIME string `column:"time"`
	// COMMAND is command.
	COMMAND string `column:"command"`
}
