package behavior_tree

type SelectorNode struct {
	BaseBehaviorNode
}

func NewSelectorNode() *SelectorNode {
	return &SelectorNode{}
}

func (slf *SelectorNode) Run() bool {
	if slf.BaseBehaviorNode.Children == nil {
		return false
	}
	for _, behavior := range slf.BaseBehaviorNode.Children {
		if behavior.Run() {
			return true
		}
	}
	return false
}
