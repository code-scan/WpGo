package module

import "sync"

type SiteTask struct {
	Host string
	User string
	Pass string
}

type Task struct {
	lock sync.Mutex
	task []SiteTask
}

func NewTask() *Task {
	return &Task{}
}
func (t *Task) Pop() SiteTask {
	t.lock.Lock()
	r := t.task[0]
	t.task = t.task[1:]
	t.lock.Unlock()
	return r

}
func (t *Task) Push(host, user, pass_ string) {
	t.lock.Lock()
	s := SiteTask{
		Host: host,
		User: user,
		Pass: pass_,
	}
	t.task = append(t.task, s)
	t.lock.Unlock()
}
