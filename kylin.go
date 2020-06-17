package kylin

import "fmt"

type Argument map[string]interface{}

type Result map[string]interface{}

type Mission interface {
	Start(Argument) *chan Result
}

type Kylin struct {
	MissionMap   map[string]Mission
	LogicMap     map[string][]string
}

func (kylin *Kylin) Run(startArgsMap map[string]Argument) {
	fmt.Println("运行开始... ")
	fmt.Printf("初始化入度MAP...  ")
	var inMap = make(map[string]int)
	for from, nexts := range kylin.LogicMap {
		inMap[from] += 0
		for _, next := range nexts {
			inMap[next] += 1
		}
	}
	fmt.Println("完成.")
	fmt.Printf("初始化首批运行任务...  ")
	startMissions := make([]string, 0)
	startMissionArgs := make([]Argument, 0)
	for missionID, in := range inMap {
		if in == 0 {
			startMissions = append(startMissions, missionID)
			if arg, ok := startArgsMap[missionID]; !ok {
				startMissionArgs = append(startMissionArgs, make(map[string]interface{}))
			} else {
				startMissionArgs = append(startMissionArgs, arg)
			}
		}
	}
	fmt.Println("完成.")
	fmt.Println("开始运行各任务... ")
	if err := dispather(startMissions, startMissionArgs, kylin.MissionMap, kylin.LogicMap, inMap); err != nil {
		fmt.Println("任务运行时出错：", err.Error())
	}
	fmt.Println("运行结束.")
}

func dispather(startMissions []string, args []Argument, missionMap map[string]Mission, logicMap map[string][]string, inMap map[string]int) error {
	argsMap := make(map[string]Argument)
	channels := make(map[string]*chan Result)
	for index, missionID := range startMissions {
		channels[missionID] = missionMap[missionID].Start(args[index])
	}
	for len(channels) > 0 {
		for missionID, ch := range channels {
			select {
			case result := <-(*ch): {
				close(*ch)
				delete(channels, missionID)
				nextMissions := logicMap[missionID]
				for _, nextMissionID := range nextMissions {
					inMap[nextMissionID] -= 1
					if _, ok := argsMap[nextMissionID]; !ok {
						argsMap[nextMissionID] = make(map[string]interface{})
					}
					for k, v := range result {
						argsMap[nextMissionID][k] = v
					}
					if in, _ := inMap[nextMissionID]; in == 0 {
						channels[nextMissionID] = missionMap[nextMissionID].Start(argsMap[nextMissionID])
					}
				}
			}
			default:
			}
		}
	}
	return nil
}

