package utils

// calculate power n of a base x since go's math.Pow uses float64
func Pow(x, n int) int {
	if n == 0 {
		return 1
	}
	if n == 1 {
		return x
	}
	y := Pow(x, n/2)
	if n%2 == 0 {
		return y * y
	}
	return x * y * y
}

// finds the index of a character in the charset
func FindCharsetIndex(charset []rune, char rune) int {
	for i, c := range charset {
		if c == char {
			return i
		}
	}
	return -1
}

// stringToIndex maps a string to its index in the sequence
func stringToIndex(s, charset []rune) int {
	k := len(charset)
	length := len(s)

	// count all shorter lengths first
	index := 0
	for l := 1; l < length; l++ {
		index += Pow(k, l)
	}

	// map chars to values
	val := 0
	for i := 0; i < length; i++ {
		c := FindCharsetIndex(charset, s[i])
		val = val*k + c
	}
	index += val
	return index
}

// nthSequence maps index -> string in sequence
func nthSequence(n int, charset []rune) []rune {
	k := len(charset)

	// find length
	length := 1
	count := Pow(k, length)
	remaining := n

	for remaining >= count {
		remaining -= count
		length++
		count = Pow(k, length)
	}

	// convert remaining to base-k representation of length digits
	res := make([]rune, length)
	for i := length - 1; i >= 0; i-- {
		res[i] = charset[remaining%k]
		remaining /= k
	}
	return res
}

// nextNth returns the nth string after start using the charset
func RotateString(start []rune, step int, charset []rune) []rune {
	startIndex := stringToIndex(start, charset)
	return nthSequence(startIndex+step, charset)
}
