package algo

func KMP(text, patten string) bool {
	next := next(patten)
	i, j := 0, 0
	for i < len(text) && j < len(patten) {
		if text[i] == patten[j] {
			j++
			i++
		} else if j > 0 {
			j = next[j-1]
		} else {
			i++
		}
	}
	return j == len(patten)
}

func next(s string) []int {
	if len(s) == 0 {
		return []int{}
	}
	next := make([]int, len(s))
	next[0] = 0
	var j int = 0
	for i := 1; i < len(s); {
		if s[j] == s[i] {
			j++
			next[i] = j
			i++
		} else if j > 0 {
			j = next[j-1]
		} else {
			next[i] = 0
			i++
		}
	}
	return next
}
