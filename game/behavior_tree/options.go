package behavior_tree

type Option func(*BehaviorTree)

func WithData(data map[string]any) Option {
	return func(tree *BehaviorTree) {
		tree.data = data
	}
}
