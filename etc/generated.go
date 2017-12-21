package etc

// updated at 2017-12-21 12:15:55.557964 -0800 PST

// Mtab is '/etc/mtab' in Linux.
type Mtab struct {
	// FileSystem is file system.
	FileSystem string `column:"file_system"`
	// MountedOn is 'mounted on'.
	MountedOn string `column:"mounted_on"`
	// FileSystemType is file system type.
	FileSystemType string `column:"file_system_type"`
	// Options is file system type.
	Options string `column:"options"`
	// Dump is number indicating whether and how often the file system should be backed up by the dump program; a zero indicates the file system will never be automatically backed up.
	Dump int `column:"dump"`
	// Pass is number indicating the order in which the fsck program will check the devices for errors at boot time; this is 1 for the root file system and either 2 (meaning check after root) or 0 (do not check) for all other devices.
	Pass int `column:"pass"`
}
