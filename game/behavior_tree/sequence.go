package behavior_tree

type SequenceNode struct {
	BaseBehaviorNode
}

func NewSequenceNode() *SequenceNode {
	return &SequenceNode{}
}

func (slf *SequenceNode) Run() bool {
	if slf.BaseBehaviorNode.Children == nil {
		return false
	}
	for _, behavior := range slf.BaseBehaviorNode.Children {
		if !behavior.Run() {
			return false
		}
	}
	return true
}
