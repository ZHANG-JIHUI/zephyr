package behavior_tree_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/ZHANG-JIHUI/zephyr/game/behavior_tree"
	"github.com/ZHANG-JIHUI/zephyr/tools/log"
)

type EatBehavior struct {
	behavior_tree.BaseBehaviorNode
}

func (slf *EatBehavior) Run() bool {
	hungryData := slf.GetTree().GetProperty("hungry")
	if hungryData == nil {
		return false
	}
	state := hungryData.(bool) == true
	defer log.Info("eat behavior", log.Bool("hungry state", state), log.Bool("exec state", state))
	if !state {
		return false
	}

	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	<-timer.C
	slf.GetTree().SetProperty("hungry", false)
	return true
}

type SleepBehavior struct {
	behavior_tree.BaseBehaviorNode
}

func (slf *SleepBehavior) Run() bool {
	state := rand.Intn(10)%2 == 0
	defer log.Info("sleep behavior", log.Bool("exec state", state))
	if !state {
		return false
	}

	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	<-timer.C
	slf.GetTree().SetProperty("hungry", true)

	return true
}

type BattleWithTomBehavior struct {
	behavior_tree.BaseBehaviorNode
}

func (slf *BattleWithTomBehavior) Run() bool {
	state := rand.Intn(10)%2 == 0
	defer log.Info("battle with tom behavior", log.Bool("exec state", state))
	if !state {
		return false
	}

	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	<-timer.C

	return true
}

type BattleWithJerryBehavior struct {
	behavior_tree.BaseBehaviorNode
}

func (slf *BattleWithJerryBehavior) Run() bool {
	state := rand.Intn(10)%2 == 0
	defer log.Info("battle with jerry behavior", log.Bool("exec state", state))
	if !state {
		return false
	}

	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	<-timer.C
	slf.GetTree().SetProperty("hungry", true)

	return true
}

func TestBehaviorTree(t *testing.T) {

	var (
		eat             = &EatBehavior{}
		sleep           = &SleepBehavior{}
		battleWithTom   = &BattleWithTomBehavior{}
		battleWithJerry = &BattleWithJerryBehavior{}
	)

	battle := behavior_tree.NewSequenceNode()
	battle.AddChild(battleWithTom, battleWithJerry)

	root := behavior_tree.NewSelectorNode()
	root.AddChild(eat, sleep, battle)

	tree := behavior_tree.NewBehaviorTree(root, behavior_tree.WithData(map[string]any{"hungry": false}))
	tree.SetProperty("hungry", true)
	tree.Run(time.Second*5, time.Minute)
}
