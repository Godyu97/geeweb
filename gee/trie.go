package gee

import (
	"github.com/Godyu97/geeweb/common"
	"strings"
)

// 把pattern根据/切割，遇到第一个*就会返回
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == common.TrieFlag2 {
				break
			}
		}
	}
	return parts
}

type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true
}

// 第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		//根据part，isWild判断是否需要返回该child
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(pattern string, parts []string, height int) {
	//递归到底层，当前node无children
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	//获取当前height对应的part
	part := parts[height]
	//真正进行叶子节点的插入
	child := n.matchChild(part)
	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == common.TrieFlag1 || part[0] == common.TrieFlag2,
		}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// 只会找到第一个匹配到的结果
// todo 所以对于 /*uri /:uri 与 /hello 同时注册会有无法预期的结果
func (n *node) search(parts []string, height int) *node {
	//递归到底层或遇到*uri，直接返回node
	if len(parts) == height ||
		//*uri
		strings.HasPrefix(n.part, string(common.TrieFlag2)) {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
