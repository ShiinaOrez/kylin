package kylin

var (
	Running = Result{Status:0}
	Success = Result{Status:1}
	Failed  = Result{Status:2}
)

type Result struct {
	Status int
}