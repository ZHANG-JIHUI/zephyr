package fsm_test

import (
	"testing"
	"time"

	"github.com/ZHANG-JIHUI/zephyr/game/fsm"
	"github.com/ZHANG-JIHUI/zephyr/tools/log"
)

type StateEat struct{ fsm.BaseState }

func (slf *StateEat) OnEnter(msg map[string]any) {
	log.Info("enter eat")
	timer := time.NewTimer(time.Second * 2)
	defer timer.Stop()
	<-timer.C
	slf.GetOwner().TransitionTo("sleep", map[string]any{
		"msg": "from eat to sleep",
	})
}

func (slf *StateEat) OnExit() {
	log.Info("exit eat")
}

func (slf *StateEat) OnUpdate(delta any) {
	log.Info("update eat")
}

type StateSleep struct{ fsm.BaseState }

func (slf *StateSleep) OnEnter(msg map[string]any) {
	log.Info("enter sleep")
}

func (slf *StateSleep) OnExit() {
	log.Info("exit sleep")
}

func (slf *StateSleep) OnUpdate(delta any) {
	log.Info("update sleep")
}

type StateBattle struct{ fsm.BaseState }

func (slf *StateBattle) OnEnter(msg map[string]any) {
	log.Info("enter battle")
}

func (slf *StateBattle) OnExit() {
	log.Info("exit battle")
}

func (slf *StateBattle) OnUpdate(delta any) {
	log.Info("update battle")
}

func TestFsm(t *testing.T) {

	f := fsm.NewFSM()

	eat := &StateEat{fsm.BaseState{Name: "eat", Owner: f}}
	sleep := &StateSleep{fsm.BaseState{Name: "sleep", Owner: f}}
	battle := &StateBattle{fsm.BaseState{Name: "battle", Owner: f}}

	f.AddState(eat, sleep, battle)
	f.Launch(eat)
}
