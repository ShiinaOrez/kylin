package main

import (
	"fmt"
	"testing"
	"time"
	. "github.com/ShiinaOrez/kylin"
)

type FirstMission struct {
	Content map[string]interface{}
}

func (mission FirstMission) Start(args Argument) *chan int {
	fmt.Println("First Mission Start...")
	fmt.Println("Args:", args)
	ch := make(chan int)
	go func() {
		fmt.Println(args["str"].(string))
		time.Sleep(5 * time.Second)
		mission.Content["ret"] = "first_mission_completed"
		mission.Content["id"] = args["id"]
		mission.Content["str"] = "next mission is running"
		ch<- 1
	}()
	return &ch
}

func (mission FirstMission) ReadResult() Result {
	return mission.Content
}

type SecondMission struct {
	Content map[string]interface{}
}

func (mission SecondMission) Start(args Argument) *chan int {
	fmt.Println("Second Mission Start...")
	fmt.Println("Args:", args)
	ch := make(chan int)
	go func() {
		fmt.Println(args["str"].(string))
		time.Sleep(3 * time.Second)
		mission.Content["ret"] = "second_mission_completed"
		mission.Content["id"] = args["id"]
		ch<- 1
	}()
	return &ch
}

func (mission SecondMission) ReadResult() Result {
	return mission.Content
}

func Test(t *testing.T) {
}