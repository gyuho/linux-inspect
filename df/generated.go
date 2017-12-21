package df

// updated at 2017-12-21 12:15:54.764438 -0800 PST

// Row is 'df' command output row in Linux.
type Row struct {
	// FileSystem is file system ('source').
	FileSystem string `column:"file_system"`
	// Device is device name.
	Device string `column:"device"`
	// MountedOn is 'mounted on' ('target').
	MountedOn string `column:"mounted_on"`
	// FileSystemType is file system type ('fstype').
	FileSystemType string `column:"file_system_type"`
	// File is file name if specified on the command line ('file').
	File string `column:"file"`
	// Inodes is total number of inodes ('itotal').
	Inodes int64 `column:"inodes"`
	// Ifree is number of available inodes ('iavail').
	Ifree int64 `column:"ifree"`
	// Iused is number of used inodes ('iused').
	Iused int64 `column:"iused"`
	// IusedPercent is percentage of iused divided by itotal ('ipcent').
	IusedPercent string `column:"iused_percent"`
	// TotalBlocks is total number of 1K-blocks ('size').
	TotalBlocks            int64  `column:"total_blocks"`
	TotalBlocksBytesN      int64  `column:"total_blocks_bytes_n"`
	TotalBlocksParsedBytes string `column:"total_blocks_parsed_bytes"`
	// AvailableBlocks is number of available 1K-blocks ('avail').
	AvailableBlocks            int64  `column:"available_blocks"`
	AvailableBlocksBytesN      int64  `column:"available_blocks_bytes_n"`
	AvailableBlocksParsedBytes string `column:"available_blocks_parsed_bytes"`
	// UsedBlocks is number of used 1K-blocks ('used').
	UsedBlocks            int64  `column:"used_blocks"`
	UsedBlocksBytesN      int64  `column:"used_blocks_bytes_n"`
	UsedBlocksParsedBytes string `column:"used_blocks_parsed_bytes"`
	// UsedBlocksPercent is percentage of used-blocks divided by total-blocks ('pcent').
	UsedBlocksPercent string `column:"used_blocks_percent"`
}
