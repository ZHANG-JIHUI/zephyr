package behavior_tree

type NodeState int

const (
	NodeStateSuccess NodeState = iota
	NodeStateFailure           = 1
	NodeStateRunning           = 2
	NodeStateCancel            = 3
)

type BehaviorNode interface {
	Run() NodeState
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
