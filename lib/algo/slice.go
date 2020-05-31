package algo

func InStrSlice(slice []string, s string) bool {
	for _, ss := range slice {
		if ss == s {
			return true
		}
	}
	return false
}
