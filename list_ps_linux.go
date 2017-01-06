package psn

// PSEntry is a process entry.
// Simplied from 'Stat' and 'Status'.
type PSEntry struct {
	Program string
	State   string
	PID     int64
	PPID    int64

	CPU    string
	VMRSS  string
	VMSize string

	FD      int64
	Threads int64

	// extra fields for sorting
	CPUNum    float64
	VMRSSNum  uint64
	VMSizeNum uint64
}
