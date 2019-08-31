package main

func reverse(input string) string {
	output := []byte(input)
	l, r := 0, len(input)-1
	for l < r {
		output[l], output[r] = output[r], output[l]
		r--
		l++
	}
	return string(output)
}
