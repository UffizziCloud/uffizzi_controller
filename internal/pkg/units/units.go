package units

const (
	GIGABYTE_IN_BYTES = 1073741824
	MINUTE_IN_SECONDS = 60
)

func ConvertBytesPerSecondsToGigabytesPerMinutes(value float64) float64 {
	return value / GIGABYTE_IN_BYTES / MINUTE_IN_SECONDS
}
