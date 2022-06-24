package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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
	return t.Deadline.Start - t.Time.Start
}

func (t *Task) Print() {
	fmt.Println("#----------------")
	fmt.Println("| Tarefa:", t.ID)
	fmt.Println("| Duração:", t.Duration)
	fmt.Println("| Tempo Mínimo")
	fmt.Println("|    - Inicial:", t.Time.Start)
	fmt.Println("|    - Final:", t.Time.End)
	fmt.Println("| Tempo Máximo")
	fmt.Println("|    - Inicial:", t.Deadline.Start)
	fmt.Println("|    - Final:", t.Deadline.End)
	fmt.Println("| Folga: ", t.Slack())

	if len(t.Dependencies) == 0 {
		fmt.Println("| Sem precedentes")
	} else {
		fmt.Println("| Precedentes")
		for _, d := range t.Dependencies {
			fmt.Println("|    -", d.ID)
		}
	}
	fmt.Println("#----------------")
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

func (p *Project) Print() {
	for _, t := range p.Tasks {
		t.Print()
		fmt.Println()
	}
}

func (p *Project) PrintCriticalPath() {
	fmt.Print("Caminho crítico: ")
	for i, t := range p.GetCriticalPath() {
		if i != 0 {
			fmt.Print(" -> ")
		}
		fmt.Print(t.ID)
	}
}

func getPredefinedTasks(list string) []*Task {
	if list == "a" {
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

		return []*Task{&A, &B, &C, &D, &E, &F, &G, &H, &I, &J, &K, &L}
	}

	A := NewTask("1", 10)
	B := NewTask("2", 4)
	B.AddDependency(&A)
	C := NewTask("3", 7)
	C.AddDependency(&A)
	D := NewTask("4", 5)
	D.AddDependency(&C)
	E := NewTask("5", 5)
	E.AddDependency(&B)
	E.AddDependency(&D)
	F := NewTask("6", 3)
	F.AddDependency(&C)
	return []*Task{&A, &B, &C, &D, &E, &F}
}

func main() {
	var project Project

	key := ""
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("Aperte 't' para adicionar uma tarefa, 'q' para prosseguir ou 'd' para usar uma das listas de tarefas pré definidas")
		fmt.Print("> ")
		fmt.Scanf("%s", &key)
		key = strings.ToLower(key)
		if key == "q" || key == "d" {
			break
		}

		var t Task
		fmt.Println("\nDigite o nome da tarefa")
		fmt.Print("> ")
		fmt.Scanf("%s", &t.ID)

		fmt.Println("Digite a duração da tarefa")
		fmt.Print("> ")
		fmt.Scanf("%d", &t.Duration)

		fmt.Println("Digite a lista de precedentes da tarefa. Ex: A D E")
		fmt.Print("> ")
		scanner.Scan()
		s := scanner.Text()
		dependencies := strings.Split(s, " ")

		if len(dependencies) > 0 {
			for _, id := range dependencies {
				for _, d := range project.Tasks {
					if d.ID == id {
						t.AddDependency(d)
					}
				}
			}
		}

		project.AddTask(&t)
	}

	for key == "d" {
		fmt.Println("Deseja usar a lista de tarefas A ou B?")
		fmt.Print("> ")
		fmt.Scanf("%s", &key)
		key = strings.ToLower(key)
		if key == "a" || key == "b" {
			tasks := getPredefinedTasks(key)
			project.Tasks = tasks
			break
		}
		key = "d"
	}

	project.SetTimes()
	project.SetDeadlines()

	project.Print()
	project.PrintCriticalPath()
	fmt.Printf("\nTempo final do projeto: %d\n", project.GetHigherTaskTimeEnd())
}
