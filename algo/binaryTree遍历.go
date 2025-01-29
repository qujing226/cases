package main

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// 递归
func preorderTraversal(root *TreeNode) (res []int) {
	var traversal func(node *TreeNode)
	traversal = func(node *TreeNode) {
		if node == nil {
			return
		}
		res = append(res, node.Val)
		traversal(node.Left)
		traversal(node.Right)
	}
	traversal(root)
	return
}

func inorderTraversal(root *TreeNode) (res []int) {
	var traversal func(node *TreeNode)
	traversal = func(node *TreeNode) {
		if node == nil {
			return
		}
		traversal(node.Left)
		res = append(res, node.Val)
		traversal(node.Right)
	}
	traversal(root)
	return
}

func postorderTraversal(root *TreeNode) (res []int) {
	var traversal func(node *TreeNode)
	traversal = func(node *TreeNode) {
		if node == nil {
			return
		}
		traversal(node.Left)
		traversal(node.Right)
		res = append(res, node.Val)
	}
	traversal(root)
	return
}

// 迭代
func preOrderTraversal(root *TreeNode) (res []int) {
	if root == nil {
		return
	}
	stack := []*TreeNode{root}

	for len(stack) != 0 {
		temp := stack[len(stack)-1]
		res = append(res, temp.Val)
		if temp.Right != nil {
			stack = append(stack, temp.Right)
		}
		if temp.Left != nil {
			stack = append(stack, temp.Left)
		}
	}
	return res
}

func postOrderTraversal(root *TreeNode) (res []int) {
	if root == nil {
		return
	}
	stack := []*TreeNode{root}
	for len(stack) != 0 {
		temp := stack[len(stack)-1]
		res = append(res, temp.Val)
		if temp.Left != nil {
			stack = append(stack, temp.Left)
		}
		if temp.Right != nil {
			stack = append(stack, temp.Right)
		}
	}
	var reverse func(items []int) []int
	reverse = func(items []int) []int {
		left, right := 0, len(items)
		for left < right {
			items[left], items[right] = items[right], items[left]
			left++
			right--
		}
		return items
	}
	return reverse(res)
}

func inOrderTraversal(root *TreeNode) (res []int) {
	if root == nil {
		return
	}
	var stack []*TreeNode
	curr := root

	for curr != nil || len(stack) > 0 {
		if curr != nil {
			stack = append(stack, curr)
			curr = curr.Left
		} else {
			temp := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			res = append(res, temp.Val)
			curr = temp.Right
		}
	}
	return res
}
