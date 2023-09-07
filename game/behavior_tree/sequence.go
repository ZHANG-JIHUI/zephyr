package behavior_tree

type SequenceNode struct {
	BaseBehaviorNode
}

func NewSequenceNode() *SequenceNode {
	return &SequenceNode{}
}

func (slf *SequenceNode) Run() NodeState {
	if slf.BaseBehaviorNode.Children == nil {
		return NodeStateCancel
	}
	for _, behavior := range slf.BaseBehaviorNode.Children {
		if behavior.Run() == NodeStateFailure {
			return NodeStateFailure
		}
	}
	return NodeStateSuccess
}
