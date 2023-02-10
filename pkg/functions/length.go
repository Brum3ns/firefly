package functions

/*Take two int values and compare which min and max values. Return 'min' value as first agument and 'max' as last.
If equal then return in same order as given. True/False if the values are different*/
func LengthOrder(x, y int) (int, int) {
	if x > y {
		return y, x
	}
	//if x < y || Equal:
	return x, y
}

/**Automaticly detect highest value of the two given and compare the value difference*/
func LengthDiff(min, max int) int {
	min, max = LengthOrder(min, max)
	return (max - min)
}
