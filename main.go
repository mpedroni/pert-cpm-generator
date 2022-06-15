package main

import "fmt"

type TimeSlice struct {
	Start uint
	End   uint
}

type Task struct {
	ID           string
	Duration     uint
	Time         TimeSlice
	Deadline     TimeSlice
	Dependencies []*Task
}

func (t *Task) GetHigherDependencyEndTime() uint {
	if len(t.Dependencies) == 0 {
		return 0
	}

	higherEndTime := uint(0)
	for _, dependency := range t.Dependencies {
		if dependency.Time.End > higherEndTime {
			higherEndTime = dependency.Time.End
		}
	}

	return higherEndTime
}

func (t *Task) AddDependency(dependency *Task) {
	t.Dependencies = append(t.Dependencies, dependency)
}

func (t *Task) IsDependentOf(candidate *Task) bool {
	for _, dependency := range t.Dependencies {
		if dependency == candidate {
			return true
		}
	}
	return false
}

func (t *Task) SetDeadline(end uint) {
	t.Deadline = TimeSlice{
		End:   end,
		Start: end - t.Duration,
	}
}

func (t *Task) Slack() uint {
	return t.Time.Start - t.Deadline.Start
}

func NewTask(ID string, duration uint) Task {
	t := Task{
		ID:       ID,
		Duration: duration,
	}

	return t
}

type Project struct {
	Tasks []*Task
	Time  TimeSlice
}

func (p *Project) AddTask(t *Task) {
	p.Tasks = append(p.Tasks, t)
}

func (p *Project) SetTimes() {
	for _, t := range p.Tasks {
		duration := t.Duration
		start := t.GetHigherDependencyEndTime()
		t.Time.Start = start
		t.Time.End = start + duration

		if t.Time.End > p.Time.End {
			p.Time.End = t.Time.End
		}
	}
}

func (p *Project) GetTaskDependents(task *Task) []*Task {
	dependents := make([]*Task, 0)
	for _, candidate := range p.Tasks {
		if candidate.IsDependentOf(task) {
			dependents = append(dependents, candidate)
		}
	}

	return dependents
}

func (p *Project) SetDeadlines() {
	lastTasks := p.GetLastTasks()
	higherPossibleEnd := p.GetHigherTaskTimeEnd()

	for _, lastTask := range lastTasks {
		lastTask.SetDeadline(higherPossibleEnd)
		p.SetTaskDependenciesDeadline(lastTask)
	}
}

func (p *Project) GetLastTasks() []*Task {
	lastTasks := make([]*Task, 0)
	for _, t := range p.Tasks {
		dependents := p.GetTaskDependents(t)
		if len(dependents) == 0 {
			lastTasks = append(lastTasks, t)
		}
	}

	return lastTasks
}

func (p *Project) GetHigherTaskTimeEnd() uint {
	end := uint(0)
	for _, task := range p.Tasks {
		if task.Time.End > end {
			end = task.Time.End
		}
	}

	return end
}

func (p *Project) SetTaskDependenciesDeadline(task *Task) {
	end := task.Deadline.Start
	for _, dependency := range task.Dependencies {
		if end < dependency.Deadline.End || dependency.Deadline.End == 0 {
			dependency.SetDeadline(end)
		}
		p.SetTaskDependenciesDeadline(dependency)
	}
}

func (p *Project) GetCriticalPath() []*Task {
	path := make([]*Task, 0)
	for _, t := range p.Tasks {
		if t.Slack() == 0 {
			path = append(path, t)
		}
	}
	return path
}

func main() {
	var project Project

	A := NewTask("A", 6)
	B := NewTask("B", 2)
	C := NewTask("C", 3)
	D := NewTask("D", 10)
	D.AddDependency(&A)
	E := NewTask("E", 3)
	E.AddDependency(&A)
	F := NewTask("F", 2)
	F.AddDependency(&B)
	G := NewTask("G", 4)
	G.AddDependency(&C)
	H := NewTask("H", 5)
	H.AddDependency(&E)
	I := NewTask("I", 8)
	I.AddDependency(&F)
	I.AddDependency(&G)
	J := NewTask("J", 6)
	J.AddDependency(&G)
	K := NewTask("K", 4)
	K.AddDependency(&I)
	L := NewTask("L", 2)
	L.AddDependency(&J)

	project.AddTask(&A)
	project.AddTask(&B)
	project.AddTask(&C)
	project.AddTask(&D)
	project.AddTask(&E)
	project.AddTask(&F)
	project.AddTask(&G)
	project.AddTask(&H)
	project.AddTask(&I)
	project.AddTask(&J)
	project.AddTask(&K)
	project.AddTask(&L)

	project.SetTimes()
	project.SetDeadlines()

	for _, task := range project.Tasks {
		fmt.Println(task.ID, task.Duration, task.Time, task.Deadline)
	}

	fmt.Print("\n\n")

	for _, task := range project.GetCriticalPath() {
		fmt.Println(task.ID, task.Duration, task.Time, task.Deadline)
	}

	// fmt.Println(project.Time)
}
