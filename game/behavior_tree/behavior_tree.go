package behavior_tree

import "time"

type BehaviorTree struct {
	Root BehaviorNode
	data map[string]any
}

func NewBehaviorTree(root BehaviorNode, opts ...Option) *BehaviorTree {
	tree := &BehaviorTree{data: make(map[string]any)}
	for _, opt := range opts {
		opt(tree)
	}
	root.SetTree(tree)
	tree.Root = root
	return tree
}

func (slf *BehaviorTree) Run(cycle, duration time.Duration) {
	ticker := time.NewTicker(cycle)
	stop := time.After(duration)

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			slf.Root.Run()
		case <-stop:
			return
		}
	}
}

func (slf *BehaviorTree) SetProperty(key string, value any) {
	slf.data[key] = value
}

func (slf *BehaviorTree) GetProperty(key string) any {
	return slf.data[key]
}
