package psn

// ConvertUnixNano unix nanoseconds to unix second.
func ConvertUnixNano(unixNano int64) (unixSec int64) {
	return int64(unixNano / 1e9)
}

// Interpolate estimates missing rows in CSV
// with smallest second unit.
func (c *CSV) Interpolate() {

}
