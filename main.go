package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Define command-line flag for dryrun mode
var dryRun = flag.Bool(
	"dryrun", false,
	"Validate the input task list file and calculate the expected total runtime without executing the tasks")

// Define command-line flag for to run the tasks and determine the difference in the actual runtime versue the expected runtime
var diffTime = flag.Bool(
	"difftime", false,
	"Run the tasks and determine the difference in the actual runtime versue the expected runtime")

// Define command-line flag for file path
var taskListFilePath = flag.String(
	"taskfile", "",
	"Path to the file containing the list of tasks to be executed")

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
	// fmt.Println("The next run time in job", job.Name, "is", job.NextRun)
	s.Tasks = append(s.Tasks, job)
}

func (s *Scheduler) Run() {
	var wg sync.WaitGroup
	for _, job := range s.Tasks {
		wg.Add(1)
		go func(job *Task) { // Run the job in a goroutine currently
			defer wg.Done()
			job.Execute()
		}(job)
	}
	wg.Wait()
}

func calculateNextRun(task *Task) time.Time {
	// Add duration of the task to current time
	return time.Now().Add(time.Duration(task.Duration) * time.Second)
}

func validateInputTasksList(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening task list file:", err)
		return false
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			fmt.Println("Invalid task format in task list file for line:", line)
			return false
		}
		_, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			fmt.Println("Invalid duration format in task list file for line:", line)
			return false
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading task list file:", err)
		return false
	}
	return true
}

func calculateExpectedTotalDuration(filePath string) int64 {
	file, _ := os.Open(filePath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var overalDuration int64
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		duration, _ := strconv.ParseInt(parts[1], 10, 64)
		if duration > overalDuration {
			overalDuration = duration // since tasks are running concurrently, the overall duration is the maximum duration of all tasks
		}
	}
	return overalDuration
}

func main() {
	// Parse the flags
	flag.Parse()
	filePath := *taskListFilePath
	if filePath == "" {
		fmt.Println("Please provide the path to the task list file using the -taskfile flag.")
		return
	}

	if *dryRun {
		if validateInputTasksList(filePath) == true {
			fmt.Println("The input task list file is valid.")
			fmt.Println("The expected total runtime is", calculateExpectedTotalDuration(filePath), "seconds.")
		}
		return
	}

	scheduler := NewScheduler()
	file, err := os.Open(filePath)
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

	startTime := time.Now()
	expectedRuntime := calculateExpectedTotalDuration(filePath)
	scheduler.Run()
	endTime := time.Now()

	if *diffTime {
		actualRuntime := endTime.Sub(startTime).Seconds()
		fmt.Println("The difference in the actual runtime versus the expected runtime is", actualRuntime-float64(expectedRuntime), "seconds.")
		// convert expectedRuntime to float64 to get the difference in a fraction of seconds
	}
}
