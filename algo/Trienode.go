package main

import (
	"strings"
)

// Matcher 只包含一棵 Trie
type Matcher struct {
	root *node
}

type node struct {
	// 精确子节点
	literals map[string]*node
	// 段内含单星 '*'（但不含跨段 **）的子节点
	wildcards []*node
	// 整段 '*'
	single *node
	// 整段 '**' 且只能在模式末尾出现
	double *node

	// 终点值
	val   interface{}
	isEnd bool
	// 段模式，仅对 wildcards 有效
	pat string
}

// NewMatcher 接收预先按优先级排好序（高→低）的模式列表
func NewMatcher(patterns map[string]interface{}) *Matcher {
	m := &Matcher{root: &node{literals: map[string]*node{}}}
	for k, v := range patterns {
		insert(m.root, split(k), v)
	}
	return m
}

// split 按未转义 '/' 切分
func split(p string) []string {
	var segs []string
	var buf strings.Builder
	esc := false
	for _, r := range p {
		if esc {
			buf.WriteRune(r)
			esc = false
		} else if r == '\\' {
			esc = true
		} else if r == '/' {
			segs = append(segs, buf.String())
			buf.Reset()
		} else {
			buf.WriteRune(r)
		}
	}
	segs = append(segs, buf.String())
	return segs
}

// insert 模式只含：literal、段内'*'、整段 '*'、整段'**'(末尾)
func insert(root *node, segs []string, val interface{}) {
	cur := root
	for i, seg := range segs {
		last := i == len(segs)-1
		switch {
		case seg == "**" && last:
			if cur.double == nil {
				cur.double = &node{literals: map[string]*node{}, pat: "**"}
			}
			cur = cur.double

		case seg == "*":
			if cur.single == nil {
				cur.single = &node{literals: map[string]*node{}, pat: "*"}
			}
			cur = cur.single

		case strings.Contains(seg, "*"):
			// 段内 '*'
			var child *node
			for _, c := range cur.wildcards {
				if c.pat == seg {
					child = c
					break
				}
			}
			if child == nil {
				child = &node{literals: map[string]*node{}, pat: seg}
				cur.wildcards = append(cur.wildcards, child)
			}
			cur = child

		default:
			// 纯文字段
			if cur.literals[seg] == nil {
				cur.literals[seg] = &node{
					literals: map[string]*node{},
					pat:      seg,
				}
			}
			cur = cur.literals[seg]
		}
	}
	cur.val = val
	cur.isEnd = true
}

// Match 返回最先命中的值
func (m *Matcher) Match(path string) (interface{}, bool) {
	parts := strings.Split(path, "/")
	return matchTrie(m.root, parts, 0)
}

func matchTrie(cur *node, parts []string, idx int) (interface{}, bool) {
	// 用完所有段
	if idx == len(parts) {
		if cur.isEnd {
			return cur.val, true
		}
		// 如果存在 '**' 子节点，且它在末尾标识了一个模式，直接命中
		if cur.double != nil && cur.double.isEnd {
			return cur.double.val, true
		}
		return nil, false
	}

	seg := parts[idx]
	// 1. 精确匹配
	if nxt, ok := cur.literals[seg]; ok {
		if v, ok2 := matchTrie(nxt, parts, idx+1); ok2 {
			return v, true
		}
	}
	// 2. 段内通配（单星 or 多星但不跨段）
	for _, nxt := range cur.wildcards {
		if matchSeg(nxt.pat, seg) {
			if v, ok2 := matchTrie(nxt, parts, idx+1); ok2 {
				return v, true
			}
		}
	}
	// 3. 整段 '*'
	if cur.single != nil {
		if v, ok2 := matchTrie(cur.single, parts, idx+1); ok2 {
			return v, true
		}
	}
	// 4. 整段 '**'（末尾模式），一旦存在就可任意吞余下所有段
	if cur.double != nil && cur.double.isEnd {
		return cur.double.val, true
	}
	return nil, false
}

// matchSeg 段内 '*' 匹配（不含 '/')
func matchSeg(pat, s string) bool {
	pi, si, lastStar, starMatch := 0, 0, -1, 0
	for si < len(s) {
		if pi < len(pat) && pat[pi] != '*' && pat[pi] == s[si] {
			pi++
			si++
		} else if pi < len(pat) && pat[pi] == '*' {
			lastStar, starMatch, pi = pi, si, pi+1
		} else if lastStar >= 0 {
			pi = lastStar + 1
			starMatch++
			si = starMatch
		} else {
			return false
		}
	}
	for pi < len(pat) && pat[pi] == '*' {
		pi++
	}
	return pi == len(pat)
}

//
//
//package main
//
//import (
//	"math"
//	"sort"
//	"strings"
//)
//
//// Matcher 负责分段 Trie 与全局 Fallback 的合并匹配
//type Matcher struct {
//	root      *TrieNode   // 前缀树根
//	fallbacks []*Fallback // 按优先级降序的跨段/段内 ** 模式
//}
//
//// —— 前缀树部分 ——
//
//// TrieNode 定义
//type TrieNode struct {
//	literalChildren  map[string]*TrieNode // 纯文字段
//	wildcardChildren []*TrieNode          // 段内单星 "*" 或多星但不含跨段用法 "**"
//	singleChild      *TrieNode            // 整段 "*"
//	doubleChild      *TrieNode            // 整段 "**"
//
//	// 终点信息
//	value    interface{}
//	nonNil   bool
//	priority int    // 用于维护 wildcardChildren 排序
//	segPat   string // 存原始段模式（only wildcardChildren 用）
//}
//
//// NewMatcher 构造函数：
//// 1. 按 computePriority 排序所有模式
//// 2. 将只含整段 `*`/`**` 或单星段插入 Trie
//// 3. 其它含段内/跨段 `**` 的模式编译成 Fallback 列表
//func NewMatcher(patterns map[string]interface{}) *Matcher {
//	root := &TrieNode{literalChildren: map[string]*TrieNode{}}
//	var fb []*Fallback
//
//	// 排序
//	type ent struct {
//		pat  string
//		val  interface{}
//		prio int
//	}
//	arr := make([]ent, 0, len(patterns))
//	for p, v := range patterns {
//		arr = append(arr, ent{p, v, computePriority(p)})
//	}
//	sort.Slice(arr, func(i, j int) bool { return arr[i].prio > arr[j].prio })
//
//	// 分类插入
//	for _, e := range arr {
//		if strings.Contains(e.pat, "**") &&
//			!strings.Contains(e.pat, "/**/") &&
//			!strings.HasSuffix(e.pat, "/**") &&
//			!strings.HasPrefix(e.pat, "**/") {
//			// 段内或跨段 **，加入 fallback
//			fb = append(fb, &Fallback{
//				tokens:   tokenize(e.pat),
//				priority: e.prio,
//				val:      e.val,
//			})
//		} else {
//			root.InsertTrie(e.pat, e.val)
//		}
//	}
//	// fallback 按 priority 降序
//	sort.Slice(fb, func(i, j int) bool { return fb[i].priority > fb[j].priority })
//
//	return &Matcher{root: root, fallbacks: fb}
//}
//
//// Match 先走 Trie，得最优节点，再扫更高优先级的 fallback 抢先
//func (m *Matcher) Match(path string) (interface{}, bool) {
//	parts := strings.Split(path, "/")
//
//	// 1. Trie 分段匹配，返回最匹配的终点节点
//	node := m.root.matchNode(parts)
//	bestVal, bestPrio := interface{}(nil), math.MinInt32
//	if node != nil {
//		bestVal, bestPrio = node.value, node.priority
//	}
//
//	// 2. 遍历 fallback（已按 priority 降序）
//	for _, fb := range m.fallbacks {
//		if fb.priority <= bestPrio {
//			break // 后续 fallback 优先级都更低，无需再试
//		}
//		if matchTokens(fb.tokens, path) {
//			return fb.val, true
//		}
//	}
//
//	// 3. 返回 Trie 结果（可能为 nil）
//	if node != nil {
//		return bestVal, true
//	}
//	return nil, false
//}
//
//// matchNode 在 Trie 上做分段递归匹配，找到最优终点节点
//func (root *TrieNode) matchNode(parts []string) *TrieNode {
//	visited := map[*TrieNode]map[int]bool{}
//	return matchNodeRec(root, parts, 0, visited)
//}
//
//func matchNodeRec(node *TrieNode, parts []string, idx int, visited map[*TrieNode]map[int]bool) *TrieNode {
//	if visited[node] == nil {
//		visited[node] = map[int]bool{}
//	}
//	if visited[node][idx] {
//		return nil
//	}
//	visited[node][idx] = true
//
//	// 终点命中
//	if idx == len(parts) && node.nonNil {
//		return node
//	}
//
//	// 分段匹配：literal -> 段内星 -> 单星
//	if idx < len(parts) {
//		if ch, ok := node.literalChildren[parts[idx]]; ok {
//			if res := matchNodeRec(ch, parts, idx+1, visited); res != nil {
//				return res
//			}
//		}
//		for _, ch := range node.wildcardChildren {
//			if matchSegment(ch.segPat, parts[idx]) {
//				if res := matchNodeRec(ch, parts, idx+1, visited); res != nil {
//					return res
//				}
//			}
//		}
//		if node.singleChild != nil {
//			if res := matchNodeRec(node.singleChild, parts, idx+1, visited); res != nil {
//				return res
//			}
//		}
//	}
//
//	// 整段 "**" —— 非贪婪，从少跳到多跳
//	if node.doubleChild != nil {
//		for skip := idx; skip <= len(parts); skip++ {
//			if res := matchNodeRec(node.doubleChild, parts, skip, visited); res != nil {
//				return res
//			}
//		}
//	}
//	return nil
//}
//
//// Fallback 模式结构
//type Fallback struct {
//	tokens   []string
//	priority int
//	val      interface{}
//}
//
//// 给单段模式（含 "*”？）做匹配算法：'*' 匹配任意字符序列（不含 '/')
//func matchSegment(pat, s string) bool {
//	pi, si := 0, 0
//	lastStar, matchAt := -1, 0
//	for si < len(s) {
//		if pi < len(pat) && pat[pi] != '*' && pat[pi] == s[si] {
//			pi++
//			si++
//		} else if pi < len(pat) && pat[pi] == '*' {
//			lastStar = pi
//			matchAt = si
//			pi++
//		} else if lastStar != -1 {
//			// 回溯到上次星号
//			pi = lastStar + 1
//			matchAt++
//			si = matchAt
//		} else {
//			return false
//		}
//	}
//	// 跳过尾部连续的 '*'
//	for pi < len(pat) && pat[pi] == '*' {
//		pi++
//	}
//	return pi == len(pat)
//}
//
//// 把整条 pattern 切成 token
//func tokenize(pat string) []string {
//	var toks []string
//	for i := 0; i < len(pat); {
//		if i+1 < len(pat) && pat[i] == '*' && pat[i+1] == '*' {
//			toks = append(toks, "**")
//			i += 2
//		} else if pat[i] == '*' {
//			toks = append(toks, "*")
//			i++
//		} else {
//			// 读一段连续的普通字符（含转义后的 / 和 *）
//			var sb strings.Builder
//			for i < len(pat) && pat[i] != '*' {
//				// 转义符直接保留
//				sb.WriteByte(pat[i])
//				i++
//			}
//			toks = append(toks, sb.String())
//		}
//	}
//	return toks
//}
//
//// 整条路径的自实现 glob："**" 可跨任意字符（包括 '/')
//// "*" 跨任意非 '/' 字符
//func matchTokens(tokens []string, path string) bool {
//	ti, si := 0, 0
//	lastStarIdx, lastPathIdx := -1, -1
//
//	for si < len(path) {
//		if ti < len(tokens) && tokens[ti] != "*" && tokens[ti] != "**" &&
//			strings.HasPrefix(path[si:], tokens[ti]) {
//			// 文字块匹配
//			si += len(tokens[ti])
//			ti++
//		} else if ti < len(tokens) && tokens[ti] == "*" {
//			// 匹配一段非 "/" 的任意
//			lastStarIdx = ti
//			lastPathIdx = si
//			ti++
//		} else if ti < len(tokens) && tokens[ti] == "**" {
//			// 跨任意（含 '/'）
//			lastStarIdx = ti
//			lastPathIdx = si
//			ti++
//		} else if lastStarIdx != -1 {
//			// 回溯：让星号再多吃一个字符
//			si = lastPathIdx + 1
//			lastPathIdx = si
//			ti = lastStarIdx + 1
//		} else {
//			return false
//		}
//	}
//	// 尾部跳过可能剩余的 "*" 或 "**"
//	for ti < len(tokens) {
//		if tokens[ti] != "*" && tokens[ti] != "**" {
//			return false
//		}
//		ti++
//	}
//	return true
//}
//
//// 其它辅助函数：computePriority, splitSegmentPattern, TrieNode.InsertTrie 等保持不变。
//
//// 计算模式优先级：文字 +10，单星 -1，双星 -5
//func computePriority(pat string) int {
//	lit, s1, s2 := 0, 0, 0
//	for i := 0; i < len(pat); {
//		if i+1 < len(pat) && pat[i] == '*' && pat[i+1] == '*' {
//			s2++
//			i += 2
//		} else if pat[i] == '*' {
//			s1++
//			i++
//		} else {
//			lit++
//			i++
//		}
//	}
//	return lit*10 - s1 - s2*5
//}
//
//// 把形如 "foo\/bar/*\/x" 这样的模式，按未转义的 "/" 切分
//func splitSegmentPattern(pat string) []string {
//	var res []string
//	var buf strings.Builder
//	escaped := false
//	for _, r := range pat {
//		if escaped {
//			buf.WriteRune(r)
//			escaped = false
//		} else if r == '\\' {
//			escaped = true
//		} else if r == '/' {
//			res = append(res, buf.String())
//			buf.Reset()
//		} else {
//			buf.WriteRune(r)
//		}
//	}
//	res = append(res, buf.String())
//	return res
//}
//
//// 把只含 / 级整段处理的模式插入 Trie
//func (root *TrieNode) InsertTrie(pattern string, val interface{}) {
//	segs := splitSegmentPattern(pattern)
//	prio := computePriority(pattern)
//	node := root
//
//	for _, seg := range segs {
//		switch seg {
//		case "**":
//			if node.doubleChild == nil {
//				node.doubleChild = &TrieNode{
//					literalChildren: map[string]*TrieNode{},
//					segPat:          "**",
//				}
//			}
//			node = node.doubleChild
//
//		case "*":
//			if node.singleChild == nil {
//				node.singleChild = &TrieNode{
//					literalChildren: map[string]*TrieNode{},
//					segPat:          "*",
//				}
//			}
//			node = node.singleChild
//
//		default:
//			if strings.Contains(seg, "*") {
//				// 只放含单星但不等于 "**" 的段
//				var ch *TrieNode
//				for _, c := range node.wildcardChildren {
//					if c.segPat == seg {
//						ch = c
//						break
//					}
//				}
//				if ch == nil {
//					ch = &TrieNode{
//						literalChildren: map[string]*TrieNode{},
//						segPat:          seg,
//						priority:        computePriority(seg),
//					}
//					node.wildcardChildren = append(node.wildcardChildren, ch)
//					sort.Slice(node.wildcardChildren, func(i, j int) bool {
//						return node.wildcardChildren[i].priority > node.wildcardChildren[j].priority
//					})
//				}
//				node = ch
//			} else {
//				// 纯文字段
//				if node.literalChildren == nil {
//					node.literalChildren = map[string]*TrieNode{}
//				}
//				if node.literalChildren[seg] == nil {
//					node.literalChildren[seg] = &TrieNode{
//						literalChildren: map[string]*TrieNode{},
//						segPat:          seg,
//					}
//				}
//				node = node.literalChildren[seg]
//			}
//		}
//	}
//	node.value = val
//	node.nonNil = true
//	node.priority = prio
//}
