package github

// pageSize returns the minimal page size necessary to fulfill the request or the
// maximum page supported by GitHub.
// https://developer.github.com/v3/#pagination
func pageSize(desired, max int) int {
	if desired < max {
		return desired
	}

	return max
}
