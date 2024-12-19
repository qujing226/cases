package algo

/*
给定一个字符串数组 words，每个字符串 words[i] 代表一个单词
我们需要找出两个单词，使得这两个单词没有公共字母（即它们的字符集没有交集），并且返回这两个单词的长度的乘积中的最大值。

我们不妨来分析一下如何表示一个单词的字符集。每个单词由字母组成，我们完全可以利用一个整数的二进制位来表示一个单词的所有字符，字符集就成了一个位图。
为什么这样做呢？因为英文字母只有26个，我们可以用一个 32 位的整数来表示。具体来说，‘a’ 对应二进制的第 0 位，‘b’ 对应第 1 位，依此类推，直到‘z’。
如果某个字母在单词中出现了，那么对应的位就置为 1。
*/

func maxProduct(words []string) int {
	n := len(words)
	masks := make([]int, n)

	// 为每个单词生成一个位图
	for i, word := range words {
		for j := 0; j < len(words); j++ {
			masks[i] |= 1 << (word[j] - 'a')
		}
	}
	// 计算最大乘积
	maxProduct := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			// 判断两个单词的位图是否有交集
			if masks[i]&masks[j] == 0 {
				maxProduct = max(maxProduct, len(words[i])*len(words[j]))
			}
		}
	}
	return maxProduct
}
