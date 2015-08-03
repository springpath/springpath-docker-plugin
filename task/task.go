// Package Task implements a task management
// interface. Tasks are currently defined
// as shell commands.
//
// It provides for invocation and tracking
// of shell commands and their relationship to
// other tasks and higher level workflows.
//
// It can be used to avoid redundant invocations
// of the same command, manage timeouts, terminations,
// dependencies, etc among tasks.
package task

import "sync"

const DefaultMaxParallel uint64 = 1
const DefaultMaxTasks uint64 = 1024

// A Task.
type Task struct {
	Name      string
	Command   []string
	Singleton bool // only one task with this name is allowed.
}

// Runtime of a task.
type runTask struct {
	task       Task
	id         uint64
	exitStatus int
	state      uint64
}

// Task queue.
type Runner struct {
	maxTasks    uint64
	maxParallel uint64
	running     []*runTask
	byName      map[string]*runTask
	numPending  uint64
	sync.Mutex
}

// The default task runner.
var DefaultTaskRunner Runner
