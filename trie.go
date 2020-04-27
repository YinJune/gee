package gee

import "strings"

type node struct {
	pattern  string  //待匹配路由 如 /user/info
	part     string  //路由的一部分 如 /info
	children []*node //子节点
	isWild   bool    //是否模糊匹配
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// /user/:name /user/address/list	当part=address的时候回返回两个node :name和address 根据后边的part再进一步判断
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

//插入节点
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// 根据请求路径数组 查找对应路径的node
func (n *node) search(parts []string, height int) *node {
	//只有当深度等于parts的长度 或者 当前节点是模糊匹配*的时候 才结束匹配
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		//非叶子节点可能也有pattern 如 /user/list /user/list/male list的pattern不为空，但是同时它又不是叶子结点
		if n.pattern != "" {
			return n
		}
		return nil
	}
	part := parts[height]
	children := n.matchChildren(part)
	if children == nil || len(children) == 0 {
		return nil
	}
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
