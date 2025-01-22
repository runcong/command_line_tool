package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	Name         string
	Duration     int64
	Dependencies []string
	Execute      func()
	NextRun      time.Time
}

type Scheduler struct {
	Tasks []*Task
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		Tasks: []*Task{},
	}
}

func (s *Scheduler) AddJob(name string, duration int64, dependencies []string, execute func()) {
	job := &Task{
		Name:         name,
		Duration:     duration,
		Dependencies: dependencies,
		Execute:      execute,
	}
	job.NextRun = time.Now().Add(time.Duration(job.Duration) * time.Second)
	fmt.Println("The next run time in job", job.Name, "is", job.NextRun)
	s.Tasks = append(s.Tasks, job)
}

func (s *Scheduler) Run() {
	for _, job := range s.Tasks {
		job.Execute()
	}
}

func calculateNextRun(task *Task) time.Time {
	// Add duration of the task to current time
	return time.Now().Add(time.Duration(task.Duration) * time.Second)
}

func main() {
	scheduler := NewScheduler()
	file, err := os.Open("task_list.txt")
	if err != nil {
		fmt.Println("Error opening task list file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			fmt.Println("Invalid task format in task list file for line:", line)
			continue
		}
		duration, err := strconv.ParseInt(parts[1], 10, 64)
		name := parts[0]
		if err != nil {
			fmt.Println("Invalid duration format in task list file for line:", line)
			continue
		}
		dependencies := strings.Split(parts[2], " ")
		execute := func() {
			fmt.Println("Job", name, "executed at", time.Now(), "with expected duration", duration, "seconds and dependencies", dependencies)
			time.Sleep(time.Duration(duration) * time.Second)
		}
		scheduler.AddJob(name, duration, dependencies, execute)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading task list file:", err)
	}

	scheduler.Run()
}
