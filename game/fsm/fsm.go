package fsm

type FSM struct {
	initial State            // 初始状态
	current State            // 当前状态
	states  map[string]State // 状态集合
}

func NewFSM() *FSM {
	return &FSM{
		states: make(map[string]State),
	}
}

// Launch 启动状态机
func (slf *FSM) Launch(initial State) {
	slf.initial = initial
	slf.current = slf.initial
	slf.current.OnEnter(nil)
}

// AddState 添加状态
func (slf *FSM) AddState(states ...State) {
	if slf.states == nil {
		return
	}
	for _, state := range states {
		slf.states[state.GetName()] = state
	}
}

// RemoveState 删除状态
func (slf *FSM) RemoveState(stateName string) {
	delete(slf.states, stateName)
}

// TransitionTo 切换状态
func (slf *FSM) TransitionTo(stateName string, msg map[string]any) {
	if state, ok := slf.states[stateName]; ok {
		slf.current.OnExit()
		slf.current = state
		slf.current.OnEnter(msg)
	}
}
