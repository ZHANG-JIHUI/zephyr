package behavior_tree

type BehaviorNode interface {
	Run() bool
	AddChild(nodes ...BehaviorNode)
	GetChildren() []BehaviorNode
	SetTree(tree *BehaviorTree)
	GetTree() *BehaviorTree
}

type BaseBehaviorNode struct {
	tree     *BehaviorTree
	Children []BehaviorNode
}

func (slf *BaseBehaviorNode) AddChild(nodes ...BehaviorNode) {
	if slf.tree != nil {
		for _, node := range nodes {
			node.SetTree(slf.tree)
		}
	}
	slf.Children = append(slf.Children, nodes...)
}

func (slf *BaseBehaviorNode) GetChildren() []BehaviorNode {
	return slf.Children
}

func (slf *BaseBehaviorNode) SetTree(tree *BehaviorTree) {
	slf.tree = tree
	for _, child := range slf.Children {
		child.SetTree(tree)
	}
}

func (slf *BaseBehaviorNode) GetTree() *BehaviorTree {
	return slf.tree
}
