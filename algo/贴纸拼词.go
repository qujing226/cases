package main

import (
	"fmt"
	"math"
)

/*
假设你有一组贴纸，贴纸上有一些字母，你需要利用这些贴纸拼出一个目标单词。
每个贴纸上的字母只能用一次，而且你不能使用目标单词中超过贴纸上出现的字母次数。
你要最少使用多少个贴纸才能拼出目标单词？  如果拼不出来，答案就是-1。
*/
func minStickers(stickers []string, target string) int {
	// 用哈希表统计目标单词中每个字母的数量
	targetCount := make([]int, 26)
	for _, c := range target {
		targetCount[c-'a']++
	}

	// 用哈希表统计每个贴纸中每个字母的数量
	stickerCounts := make([][]int, len(stickers))
	for i, sticker := range stickers {
		stickerCount := make([]int, 26)
		for _, c := range sticker {
			stickerCount[c-'a']++
		}
		stickerCounts[i] = stickerCount
	}

	// 用一个字典来缓存已经计算过的结果
	memo := make(map[string]int)

	// 递归+记忆化搜索
	var dfs func(targetCount []int) int
	dfs = func(targetCount []int) int {
		targetKey := fmt.Sprintf("%v", targetCount)
		if val, found := memo[targetKey]; found {
			return val
		}

		// 如果目标已经被拼完
		total := 0
		for _, cnt := range targetCount {
			if cnt > 0 {
				total++
				break
			}
		}
		if total == 0 {
			return 0
		}

		res := math.MaxInt32
		for _, stickerCount := range stickerCounts {
			// 选择当前贴纸
			newTarget := make([]int, 26)
			copy(newTarget, targetCount)

			// 更新目标单词的字母频率
			for i := 0; i < 26; i++ {
				if stickerCount[i] > 0 {
					newTarget[i] = max(0, newTarget[i]-stickerCount[i])
				}
			}

			// 递归调用
			subResult := dfs(newTarget)
			if subResult != math.MaxInt32 {
				res = min(res, 1+subResult)
			}
		}

		memo[targetKey] = res
		return res
	}

	// 最终结果
	result := dfs(targetCount)
	if result == math.MaxInt32 {
		return -1
	}
	return result
}
