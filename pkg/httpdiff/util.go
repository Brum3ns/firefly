package httpdiff

func highestLstIntValue(lst []int) int {
	v := 0
	for _, i := range lst {
		if i > v {
			v = i
		}
	}
	return v
}

// Return an int array of 2 that holds 0/1 (min/max) and length diff
func lengthMinMaxDiff(x, y int) int {
	if x < y {
		return (y - x)
	}
	return x - y
}

// Compare two int type lists and return the differences presented in lstCurrent
// Return a true if the lists has one item in common
func lstIntShareItem(lstCompare, lstCurrent []int) bool {
	// Convert list1 into a map for quick lookups
	m := make(map[int]struct{})
	for _, i := range lstCompare {
		m[i] = struct{}{}
	}

	// Iterate over list2 and check if any item exists in the map
	for _, i := range lstCurrent {
		// An item is shared between the two lists
		if _, ok := m[i]; ok {
			return true
		}
	}
	return false
}

// Compare two string type lists and return the differences presented in lstCurrent
// Return a true if the lists has one item in common
func lstStringShareItem(lstCompare, lstCurrent []string) bool {
	// Convert list1 into a map for quick lookups
	m := make(map[string]struct{})
	for _, i := range lstCompare {
		m[i] = struct{}{}
	}

	// Iterate over list2 and check if any item exists in the map
	for _, i := range lstCurrent {
		// An item is shared between the two lists
		if _, ok := m[i]; ok {
			return true
		}
	}
	return false
}
