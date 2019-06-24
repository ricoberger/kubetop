package helpers

import (
	"fmt"
	"time"
)

// MaxInt returns the larger int value of the two given values.
func MaxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// MinInt returns the lower int value of the two given values.
func MinInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}

// FormatBytes format a given byte value to a human readable format.
// See: https://programming.guide/go/formatting-byte-size-to-human-readable-format.html
func FormatBytes(b int64) string {
	const unit = 1024

	if b < unit {
		return fmt.Sprintf("%dB", b)
	}

	div, exp := int64(unit), 0

	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.0f%ci", float64(b)/float64(div), "KMGTPE"[exp])
}

// FormatDuration format a given duration to a human readable format.
func FormatDuration(d time.Duration) string {
	if d.Hours() > 24*365 {
		return fmt.Sprintf("%.0fy", d.Hours()/(24*365))
	} else if d.Hours() > 120 {
		return fmt.Sprintf("%.0fd", d.Hours()/24)
	} else if d.Hours() > 10 {
		return fmt.Sprintf("%.0fh", d.Hours())
	} else if d.Minutes() > 10 {
		return fmt.Sprintf("%.0fm", d.Minutes())
	}

	return fmt.Sprintf("%.0fs", d.Seconds())
}

// RenderMemoryMax renders the value for the memory limit of a pod.
// This is needed when multiple containers run in one pod but not all containers have an limit defined.
// If some containers have not a limit return the number of containers which contain a limit.
func RenderMemoryMax(memoryMax, memoryMaxContainerCount, containerCount int64) string {
	if memoryMax == 0 {
		return "-"
	}

	if memoryMaxContainerCount == containerCount {
		return FormatBytes(memoryMax)
	}

	return fmt.Sprintf("%s (%d)", FormatBytes(memoryMax), memoryMaxContainerCount)
}

// RenderCPUMax renders the value for the cpu limit of a pod.
// This is needed when multiple containers run in one pod but not all containers have an limit defined.
// If some containers have not a limit return the number of containers which contain a limit.
func RenderCPUMax(cpuMax, cpuMaxContainerCount, containerCount int64) string {
	if cpuMax == 0 {
		return "-"
	}

	if cpuMaxContainerCount == containerCount {
		return fmt.Sprintf("%dm", cpuMax)
	}

	return fmt.Sprintf("%dm (%d)", cpuMax, cpuMaxContainerCount)
}
