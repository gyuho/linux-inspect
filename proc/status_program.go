package proc

// GetProgram returns the program name.
func GetProgram(pid int64) (string, error) {
	// Readlink needs root permission
	// return os.Readlink(fmt.Sprintf("/proc/%d/exe", pid))

	rs, err := _parseStatus(pid)
	if err != nil {
		return "", err
	}
	return rs.Name, nil
}
