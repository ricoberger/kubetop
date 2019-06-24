package api

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	mev1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

var (
	// ErrConfigNotFound is thrown if there is not a confgiuration file for Kubernetes.
	ErrConfigNotFound = errors.New("config not found")
)

// Client implements the our API client for Kubernetes.
type Client struct {
	config    *rest.Config
	clientset *kubernetes.Clientset
}

// getContainerMetrics returns the metrics for a container by the provided name from a slice of containers metrics.
func getContainerMetrics(name string, metrics []mev1beta1.ContainerMetrics) *mev1beta1.ContainerMetrics {
	for _, metric := range metrics {
		if metric.Name == name {
			return &metric
		}
	}

	return nil
}

// getContainerStatus returns the status for a container by the provided name from a slice of container statuses.
func getContainerStatus(name string, statuses []v1.ContainerStatus) *v1.ContainerStatus {
	for _, status := range statuses {
		if status.Name == name {
			return &status
		}
	}

	return nil
}

// homeDir returns the users home directory, where the '.kube' directory is located.
// The '.kube' directory contains the configuration file for a Kubernetes cluster.
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}

	// Get the home directory on windows.
	return os.Getenv("USERPROFILE")
}

// getPodMetrics returns the metrics for a pod by the provided name from a slice of pod metrics.
func getPodMetrics(name string, metrics []mev1beta1.PodMetrics) *mev1beta1.PodMetrics {
	for _, metric := range metrics {
		if metric.Name == name {
			return &metric
		}
	}

	return nil
}

// getNode returns a node by his name.
func (c *Client) getNode(name string) (*v1.Node, error) {
	node, err := c.clientset.CoreV1().Nodes().Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return node, nil
}

// NewClient initialize our client for Kubernetes.
// As first we check the 'kubeconfig' command-line flag which is passed as argument to our function.
// If the flag is not provided we check the 'KUBECONFIG' environment variable.
// In the last step we check the home directory of the user for the configuration file of Kubernetes.
func NewClient(kubeconfig string) (*Client, error) {
	if kubeconfig == "" {
		if os.Getenv("KUBECONFIG") == "" {
			if home := homeDir(); home != "" {
				kubeconfig = filepath.Join(home, ".kube", "config")
			} else {
				return nil, ErrConfigNotFound
			}
		} else {
			kubeconfig = os.Getenv("KUBECONFIG")
		}
	}

	// Use the current context in kubeconfig.
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		config:    config,
		clientset: clientset,
	}, nil
}

// GetClustername returns the name of the Kubernetes cluster.
func (c *Client) GetClustername() string {
	return c.config.Host
}

// GetNamespaces returns a slice of namespaces.
func (c *Client) GetNamespaces() ([]string, error) {
	namespaces := []string{"-"}

	data, err := c.clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, item := range data.Items {
		namespaces = append(namespaces, item.Name)
	}

	return namespaces, nil
}

// GetNodes returns a slice of node names.
func (c *Client) GetNodes() ([]string, error) {
	nodes := []string{"-"}

	data, err := c.clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, item := range data.Items {
		nodes = append(nodes, item.Name)
	}

	return nodes, nil
}

// GetNodesMetrics returns the metrics for all nodes.
func (c *Client) GetNodesMetrics(sortorder Sort) ([]Node, error) {
	var nodes []Node
	var nodeMetricsList mev1beta1.NodeMetricsList

	// Get the metrics data for all nodes from the Kubernetes API.
	// Unmarshal the data into the nodeMetricsList slice.
	data, err := c.clientset.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/nodes").DoRaw()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &nodeMetricsList)
	if err != nil {
		return nil, err
	}

	// Iterate over each node and populate our custom node structure.
	// To get the external and internal ip addess of a node we also call the node API endpoint per node.
	// After that we select all pods for a node by the 'spec.nodeName' label of the pods.
	// The pod data is not used yet, we only rely on the lenght of the slice for the number of pods per node.
	for _, item := range nodeMetricsList.Items {
		node, err := c.getNode(item.Name)
		if err != nil {
			return nil, err
		}

		var externalIP string
		var internalIP string

		for _, addr := range node.Status.Addresses {
			if addr.Type == v1.NodeExternalIP {
				externalIP = addr.Address
			}

			if addr.Type == v1.NodeInternalIP {
				internalIP = addr.Address
			}
		}

		pods, err := c.clientset.CoreV1().Pods("").List(metav1.ListOptions{
			FieldSelector: "spec.nodeName=" + item.Name,
		})
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, Node{
			Name:        item.Name,
			PodsCount:   len(pods.Items),
			MemoryTotal: node.Status.Allocatable.Memory().Value(),
			MemoryUsed:  item.Usage.Memory().Value(),
			CPUTotal:    node.Status.Allocatable.Cpu().MilliValue(),
			CPUUsed:     item.Usage.Cpu().MilliValue(),
			ExternalIP:  externalIP,
			InternalIP:  internalIP,
		})
	}

	// Sort all our nodes by the provided sortorder.
	if sortorder == SortCPUASC {
		sort.SliceStable(nodes, func(i, j int) bool {
			return nodes[i].CPUUsed < nodes[j].CPUUsed
		})
	} else if sortorder == SortCPUDESC {
		sort.SliceStable(nodes, func(i, j int) bool {
			return nodes[i].CPUUsed > nodes[j].CPUUsed
		})
	} else if sortorder == SortMemoryASC {
		sort.SliceStable(nodes, func(i, j int) bool {
			return nodes[i].MemoryUsed < nodes[j].MemoryUsed
		})
	} else if sortorder == SortMemoryDESC {
		sort.SliceStable(nodes, func(i, j int) bool {
			return nodes[i].MemoryUsed > nodes[j].MemoryUsed
		})
	} else if sortorder == SortName {
		sort.SliceStable(nodes, func(i, j int) bool {
			return nodes[i].Name < nodes[j].Name
		})
	} else if sortorder == SortPodsASC {
		sort.SliceStable(nodes, func(i, j int) bool {
			return nodes[i].PodsCount < nodes[j].PodsCount
		})
	} else if sortorder == SortPodsDESC {
		sort.SliceStable(nodes, func(i, j int) bool {
			return nodes[i].PodsCount > nodes[j].PodsCount
		})
	}

	return nodes, nil
}

// GetPodsMetrics returns metrics for all pods.
func (c *Client) GetPodsMetrics(filter Filter, sortorder Sort) ([]Pod, error) {
	var options metav1.ListOptions
	var pods []Pod
	var podMetrics mev1beta1.PodMetricsList

	// Get the metrics data for all nodes from the Kubernetes API.
	// Unmarshal the data into the podMetrics slice.
	data, err := c.clientset.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/pods").DoRaw()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &podMetrics)
	if err != nil {
		return nil, err
	}

	// Get all the pods a second time from the Kubernetes API.
	// This is needed because the metrics endpoint does not return all needed data.
	// If the node filter is not empty we apply the field selector 'spec.nodeName' to only get pods on the specified node.
	if filter.Node == "" {
		options = metav1.ListOptions{}
	} else {
		options = metav1.ListOptions{
			FieldSelector: "spec.nodeName=" + filter.Node,
		}
	}

	podsList, err := c.clientset.CoreV1().Pods(filter.Namespace).List(options)
	if err != nil {
		return nil, err
	}

	// Iterate over all pods to populate our custom pods structure.
	for _, item := range podsList.Items {
		// Get the values for memory, memory limit, cpu and cpu limit.
		// We try to get the same pod from the pods which where returned by the pods metrics API.
		// Then we calculate the memory and cpu usage, by adding the individual values of each container.
		// In the last step we calculate the limits, by adding the individual values of each container.
		// We also count the number of containers which have a limit, to visualize if a limit for the pod is only for one and not all containers.
		// If we do not do this, we can show a larger memory usage as the limit and this looks ugly without any indicator.
		var memory, memoryMax, memoryMaxContainerCount, cpu, cpuMax, cpuMaxContainerCount int64

		metrics := getPodMetrics(item.Name, podMetrics.Items)
		if metrics != nil {
			for _, container := range metrics.Containers {
				cpu = cpu + container.Usage.Cpu().MilliValue()
				memory = memory + container.Usage.Memory().Value()
			}
		}

		for _, container := range item.Spec.Containers {
			cpuMax = cpuMax + container.Resources.Limits.Cpu().MilliValue()
			memoryMax = memoryMax + container.Resources.Limits.Memory().Value()

			if container.Resources.Limits.Cpu().MilliValue() != 0 {
				cpuMaxContainerCount = cpuMaxContainerCount + 1
			}

			if container.Resources.Limits.Memory().Value() != 0 {
				memoryMaxContainerCount = memoryMaxContainerCount + 1
			}
		}

		// Get the pod status and restarts.
		// We set the status to running and only change this status if one container has an state which is not running.
		// The number of restarts represents the sum of the individual restarts of each container in a pod.
		// Last but not least we count how many containers in a pod are ready to serve requests.
		status := "Running"
		statusGeneral := 2
		statusDefined := false
		var ready, restarts int64

		for _, container := range item.Status.ContainerStatuses {
			if container.Ready {
				ready++
			}

			restarts = restarts + int64(container.RestartCount)

			if !statusDefined {
				if container.State.Waiting != nil {
					status = container.State.Waiting.Reason
					statusGeneral = 1
					statusDefined = true
				} else if container.State.Terminated != nil {
					status = container.State.Terminated.Reason
					statusGeneral = 0
					statusDefined = true
				}
			}
		}

		// Add the pod to our slice of pods whenn the status matchs the specified status in the filter.
		// If the status in the filter is 10 then all pods are added.
		if filter.Status == 10 || filter.Status == statusGeneral {
			pods = append(pods, Pod{
				Name:                    item.Name,
				Namespace:               item.Namespace,
				NodeName:                item.Spec.NodeName,
				Memory:                  memory,
				MemoryMax:               memoryMax,
				MemoryMaxContainerCount: memoryMaxContainerCount,
				CPU:                     cpu,
				CPUMax:                  cpuMax,
				CPUMaxContainerCount:    cpuMaxContainerCount,
				ContainersCount:         len(item.Spec.Containers),
				ContainersReady:         ready,
				Status:                  status,
				StatusGeneral:           statusGeneral,
				Restarts:                restarts,
				CreationDate:            item.CreationTimestamp.Time,
				IP:                      item.Status.PodIP,
			})
		}
	}

	// Sort all our pods by the provided sortorder.
	if sortorder == SortCPUASC {
		sort.SliceStable(pods, func(i, j int) bool {
			return pods[i].CPU < pods[j].CPU
		})
	} else if sortorder == SortCPUDESC {
		sort.SliceStable(pods, func(i, j int) bool {
			return pods[i].CPU > pods[j].CPU
		})
	} else if sortorder == SortMemoryASC {
		sort.SliceStable(pods, func(i, j int) bool {
			return pods[i].Memory < pods[j].Memory
		})
	} else if sortorder == SortMemoryDESC {
		sort.SliceStable(pods, func(i, j int) bool {
			return pods[i].Memory > pods[j].Memory
		})
	} else if sortorder == SortName {
		sort.SliceStable(pods, func(i, j int) bool {
			return pods[i].Name < pods[j].Name
		})
	} else if sortorder == SortNamespace {
		sort.SliceStable(pods, func(i, j int) bool {
			return pods[i].Namespace < pods[j].Namespace
		})
	} else if sortorder == SortRestartsASC {
		sort.SliceStable(pods, func(i, j int) bool {
			return pods[i].Restarts < pods[j].Restarts
		})
	} else if sortorder == SortRestartsDESC {
		sort.SliceStable(pods, func(i, j int) bool {
			return pods[i].Restarts > pods[j].Restarts
		})
	} else if sortorder == SortStatus {
		sort.SliceStable(pods, func(i, j int) bool {
			return pods[i].StatusGeneral < pods[j].StatusGeneral
		})
	}

	return pods, nil
}

// GetPod returns a pod with all details.
func (c *Client) GetPod(name, namespace string, selectedContainer int) (*Pod, error) {
	var podEvents v1.EventList
	var podMetrics mev1beta1.PodMetrics
	var containers []Container

	// Get the details for a pod.
	pod, err := c.clientset.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// Get the events for a pod.
	// cURL Example: curl http://localhost:8001/api/v1/namespaces/kube-system/events?fieldSelector=involvedObject.name=kube-proxy-tfpcb
	eventsData, err := c.clientset.RESTClient().Get().AbsPath("api/v1/namespaces/"+namespace+"/events").Param("fieldSelector", "involvedObject.name="+name).DoRaw()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(eventsData, &podEvents)
	if err != nil {
		return nil, err
	}

	var events []Event
	for _, event := range podEvents.Items {
		events = append(events, Event{
			Message:   event.Message,
			Timestamp: event.LastTimestamp.Unix(),
		})
	}

	// Get the logs for a pod.
	if len(pod.Spec.Containers) <= selectedContainer {
		selectedContainer = 0
	}

	var tailLines int64 = 100
	logData, err := c.clientset.CoreV1().Pods(namespace).GetLogs(name, &v1.PodLogOptions{
		Container: pod.Spec.Containers[selectedContainer].Name,
		TailLines: &tailLines,
	}).DoRaw()
	if err != nil {
		return nil, err
	}

	// Get the metrics for the pod.
	podMetricsData, err := c.clientset.RESTClient().Get().AbsPath("apis/metrics.k8s.io/v1beta1/namespaces/" + namespace + "/pods/" + name).DoRaw()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(podMetricsData, &podMetrics)
	if err != nil {
		return nil, err
	}

	// Get the metrics for each container.
	for _, container := range pod.Spec.Containers {
		var cpu, memory int64

		containerMetrics := getContainerMetrics(container.Name, podMetrics.Containers)
		if containerMetrics != nil {
			cpu = containerMetrics.Usage.Cpu().MilliValue()
			memory = containerMetrics.Usage.Memory().Value()
		}

		// Get the number of restarts and the status of a container.
		var restarts int32
		statuses := getContainerStatus(container.Name, pod.Status.ContainerStatuses)
		status := "-"

		if statuses != nil {
			restarts = statuses.RestartCount
			status = "Running"
			statusDefined := false

			if !statusDefined {
				if statuses.State.Waiting != nil {
					status = statuses.State.Waiting.Reason
					statusDefined = true
				} else if statuses.State.Terminated != nil {
					status = statuses.State.Terminated.Reason
					statusDefined = true
				}
			}
		}

		containers = append(containers, Container{
			Name:      container.Name,
			Restarts:  restarts,
			Status:    status,
			CPU:       cpu,
			CPUMax:    container.Resources.Limits.Cpu().MilliValue(),
			CPUMin:    container.Resources.Requests.Cpu().MilliValue(),
			Memory:    memory,
			MemoryMax: container.Resources.Limits.Memory().Value(),
			MemoryMin: container.Resources.Requests.Memory().Value(),
		})
	}

	// Define the status of the pod by checking the status of each container.
	podStatus := "Running"
	podStatusDefined := false

	for _, container := range pod.Status.ContainerStatuses {
		if !podStatusDefined {
			if container.State.Waiting != nil {
				podStatus = container.State.Waiting.Reason
				podStatusDefined = true
			} else if container.State.Terminated != nil {
				podStatus = container.State.Terminated.Reason
				podStatusDefined = true
			}
		}
	}

	var controlledBy []string
	for _, ownerReference := range pod.OwnerReferences {
		controlledBy = append(controlledBy, ownerReference.Kind+"/"+ownerReference.Name)
	}

	return &Pod{
		Name:            pod.Name,
		Namespace:       pod.Namespace,
		NodeName:        pod.Spec.NodeName,
		ContainersCount: len(containers),
		//ContainersReady: ready,
		Status: podStatus,
		//StatusGeneral:   statusGeneral,
		//Restarts:        restarts,
		Labels:       pod.Labels,
		Annotations:  pod.Annotations,
		ControlledBy: controlledBy,
		CreationDate: pod.CreationTimestamp.Time,
		IP:           pod.Status.PodIP,
		Containers:   containers,
		LogLines:     strings.Split(string(logData), "\n"),
		Events:       events,
	}, nil
}
