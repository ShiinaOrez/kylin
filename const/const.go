package _const

import "github.com/ShiinaOrez/kylin/result"

const (
	StatusRunning = 0
	StatusSuccess = 1
	StatusFailed  = 2
)

var (
	Running = result.Result{Status:StatusRunning, Description:"Running"}
	Success = result.Result{Status:StatusSuccess, Description:"Success"}
	Failed  = result.Result{Status:StatusFailed,  Description:"Failed"}
)