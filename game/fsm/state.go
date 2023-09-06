package fsm

type State interface {
	// OnEnter 进入状态
	OnEnter(msg map[string]any)
	// OnExit 退出状态
	OnExit()
	// OnUpdate 更新状态
	OnUpdate(delta any)
	// GetName 状态名称
	GetName() string
	// GetOwner 状态机
	GetOwner() *FSM
}

type BaseState struct {
	Name  string
	Owner *FSM
}

func (slf *BaseState) GetName() string {
	return slf.Name
}

func (slf *BaseState) GetOwner() *FSM {
	return slf.Owner
}
