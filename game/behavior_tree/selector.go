package behavior_tree

import "github.com/samber/lo"

type SelectorNode struct {
	BaseBehaviorNode
}

func NewSelectorNode() *SelectorNode {
	return &SelectorNode{}
}

func (slf *SelectorNode) Run() NodeState {
	if slf.BaseBehaviorNode.Children == nil {
		return NodeStateCancel
	}
	for _, behavior := range slf.BaseBehaviorNode.Children {
		if behavior.Run() == NodeStateSuccess {
			return NodeStateSuccess
		}
	}
	return NodeStateFailure
}

type RandomSelectorNode struct {
	BaseBehaviorNode
}

func NewRandomSelectorNode() *RandomSelectorNode {
	return &RandomSelectorNode{}
}

func (slf *RandomSelectorNode) Run() NodeState {
	if slf.BaseBehaviorNode.Children == nil {
		return NodeStateCancel
	}
	slf.Children = lo.Shuffle(slf.Children)
	for _, behavior := range slf.BaseBehaviorNode.Children {
		if behavior.Run() == NodeStateSuccess {
			return NodeStateSuccess
		}
	}
	return NodeStateFailure
}
