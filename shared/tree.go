// An implementation of a tree data structure

package shared

import "fmt"

type ErrNotInTree[T any] struct{ needle *T }

func (e ErrNotInTree[T]) Error() string {
	return fmt.Sprintf("%+v is not in the tree", e.needle)
}

type TreeNode[T comparable] struct {
	Self     *T
	Parent   *TreeNode[T]
	Children []TreeNode[T]
}

type Tree[T comparable] struct {
	Root  TreeNode[T]
	Nodes []TreeNode[T]
}

func (t Tree[T]) Find(needle *T) (n TreeNode[T], err error) {
	if needle == nil {
	}
	for _, node := range t.Nodes {
		if *node.Self == *needle {
			return node, nil
		}
	}
	return n, ErrNotInTree[T]{needle}
}
