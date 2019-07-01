package api

import (
	"time"
)

// Node represents a node in the Kubernetes cluster with all needed fields.
type Node struct {
	Name        string
	PodsCount   int
	MemoryTotal int64
	MemoryUsed  int64
	CPUTotal    int64
	CPUUsed     int64
	ExternalIP  string
	InternalIP  string
}

// Pod represents a pod in the Kubernetes cluster with all needed fields.
type Pod struct {
	Name                    string
	Namespace               string
	NodeName                string
	Memory                  int64
	MemoryMax               int64
	MemoryMaxContainerCount int64
	CPU                     int64
	CPUMax                  int64
	CPUMaxContainerCount    int64
	ContainersCount         int
	ContainersReady         int64
	Status                  string
	StatusGeneral           int
	Restarts                int64
	Labels                  map[string]string
	Annotations             map[string]string
	ControlledBy            []string
	CreationDate            time.Time
	IP                      string
	Containers              []Container
	LogLines                []string
	Events                  []Event
}

// Container represents a container in a pod of the Kubernetes cluster with all needed fields.
type Container struct {
	Name      string
	Memory    int64
	MemoryMax int64
	MemoryMin int64
	CPU       int64
	CPUMax    int64
	CPUMin    int64
	Status    string
	Restarts  int32
}

// Event represents a event in a pod of the Kubernetes cluster with all needed fields.
type Event struct {
	UID            string
	Message        string
	Timestamp      int64
	Count          int32
	Name           string
	Namespace      string
	Kind           string
	Type           string
	Reason         string
	Source         string
	Node           string
	FirstTimestamp time.Time
	LastTimestamp  time.Time
}

// Sort is our custom type which represents the sort order for the data which is returned by the Kubernetes API.
type Sort string

const (
	// SortCPUASC sorts the results of an API request by cpu usage (ASC).
	SortCPUASC Sort = "CPU (A)"
	// SortCPUDESC sorts the results of an API request by cpu usage (DESC).
	SortCPUDESC Sort = "CPU (D)"
	// SortMemoryASC sorts the results of an API request by memory usage (ASC).
	SortMemoryASC Sort = "Memory (A)"
	// SortMemoryDESC sorts the results of an API request by memory usage (DESC).
	SortMemoryDESC Sort = "Memory (D)"
	// SortName sorts the results of an API request by the name.
	SortName Sort = "Name"
	// SortNamespace sorts the results of an API request by the namespace.
	SortNamespace Sort = "Namespace"
	// SortPodsASC sorts the results of an API request by the number of pods (ASC).
	SortPodsASC Sort = "Pods (A)"
	// SortPodsDESC sorts the results of an API request by the number of pods (DESC).
	SortPodsDESC Sort = "Pods (D)"
	// SortRestartsASC sorts the results of an API request by the number of restarts (asc).
	SortRestartsASC Sort = "Restarts (A)"
	// SortRestartsDESC sorts the results of an API request by the number of restarts (desc).
	SortRestartsDESC Sort = "Restarts (D)"
	// SortStatus sorts the results of an API request by the status.
	SortStatus Sort = "Status"
	// SortTimeASC sorts the results of an API request by the timestamp (asc).
	SortTimeASC = "Timestamp (A)"
	// SortTimeDESC sorts the results of an API request by the timestamp (desc).
	SortTimeDESC = "Timestamp (D)"
)

// Filter is our custom type which applies a filter for the data which is returned by the Kubernetes API.
type Filter struct {
	Namespace string
	Node      string
	Status    int
	EventType string
}
