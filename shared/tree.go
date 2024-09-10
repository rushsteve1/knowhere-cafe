// An implementation of a tree data structure

package shared

import "slices"

type TreeNode[T any] struct {
	Self *T
	Parent *TreeNode[T]
	Children []TreeNode[T]
}

type Tree[T any] struct {
	Root TreeNode[T]
	Nodes []TreeNode[T]
}

func (t Tree[T]) Find(needle *T) *TreeNode[T] {
	for _, node := range t.Nodes {
		if *node.Self == *needle {
			return node
		}
	}
	return nil
}