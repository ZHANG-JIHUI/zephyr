package behavior_tree_test

import (
	"math/rand"
	"testing"
	"time"

	behavior "github.com/ZHANG-JIHUI/zephyr/game/behavior_tree"
	"github.com/ZHANG-JIHUI/zephyr/tools/log"
)

type EatBehavior struct {
	behavior.BaseBehaviorNode
}

func (slf *EatBehavior) Run() behavior.NodeState {
	hungryData := slf.GetTree().GetProperty("hungry")
	if hungryData == nil {
		return behavior.NodeStateFailure
	}
	state := hungryData.(bool) == true
	defer log.Info("eat behavior", log.Bool("hungry state", state), log.Bool("exec state", state))
	if !state {
		return behavior.NodeStateFailure
	}
	slf.GetTree().SetProperty("hungry", false)
	return behavior.NodeStateSuccess
}

type SleepBehavior struct {
	behavior.BaseBehaviorNode
}

func (slf *SleepBehavior) Run() behavior.NodeState {
	state := rand.Intn(10)%2 == 0
	defer log.Info("sleep behavior", log.Bool("exec state", state))
	if !state {
		return behavior.NodeStateFailure
	}
	slf.GetTree().SetProperty("hungry", true)
	return behavior.NodeStateSuccess
}

type BattleWithTomBehavior struct {
	behavior.BaseBehaviorNode
}

func (slf *BattleWithTomBehavior) Run() behavior.NodeState {
	state := rand.Intn(10)%2 == 0
	defer log.Info("battle with tom behavior", log.Bool("exec state", state))
	if !state {
		return behavior.NodeStateFailure
	}
	return behavior.NodeStateSuccess
}

type BattleWithJerryBehavior struct {
	behavior.BaseBehaviorNode
}

func (slf *BattleWithJerryBehavior) Run() behavior.NodeState {
	state := rand.Intn(10)%2 == 0
	defer log.Info("battle with jerry behavior", log.Bool("exec state", state))
	if !state {
		return behavior.NodeStateFailure
	}
	slf.GetTree().SetProperty("hungry", true)
	return behavior.NodeStateSuccess
}

func TestBehaviorTree(t *testing.T) {

	var (
		eat             = &EatBehavior{}
		sleep           = &SleepBehavior{}
		battleWithTom   = &BattleWithTomBehavior{}
		battleWithJerry = &BattleWithJerryBehavior{}
	)

	battle := behavior.NewSequenceNode()
	battle.AddChild(battleWithTom, battleWithJerry)

	root := behavior.NewSelectorNode()
	root.AddChild(eat, sleep, battle)

	tree := behavior.NewBehaviorTree(root, behavior.WithData(map[string]any{"hungry": false}))
	tree.SetProperty("hungry", true)
	tree.Run(time.Millisecond*50, time.Minute)
}
