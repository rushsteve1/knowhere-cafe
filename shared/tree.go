// An implementation of a tree data structure

package shared

type TreeNode[T comparable] struct {
	Self     *T
	Parent   *TreeNode[T]
	Children []TreeNode[T]
}

type Tree[T comparable] struct {
	Root  TreeNode[T]
	Nodes []TreeNode[T]
}

func (t Tree[T]) Find(needle *T) Maybe[TreeNode[T]] {
	for _, node := range t.Nodes {
		if *node.Self == *needle {
			return Some(node)
		}
	}
	return None[TreeNode[T]]()
}
