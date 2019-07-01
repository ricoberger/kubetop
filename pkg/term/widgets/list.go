package widgets

import (
	"fmt"

	"github.com/ricoberger/kubetop/pkg/api"

	ui "github.com/gizak/termui/v3"
	w "github.com/gizak/termui/v3/widgets"
)

// ListType is our custom type for the different list types (e.g. sort and filter)
type ListType string

const (
	// ListTypeSort represents the the sorting list.
	ListTypeSort ListType = "Sort by ..."
	// ListTypeFilterNamespace represents the namespace filter.
	ListTypeFilterNamespace ListType = "Filter by Namespace ..."
	// ListTypeFilterNode represents the node filter.
	ListTypeFilterNode ListType = "Filter by Node ..."
	// ListTypeFilterStatus represents the status filter.
	ListTypeFilterStatus ListType = "Filter by Status ..."
	// ListTypeFilterEventType represents the event type filter.
	ListTypeFilterEventType ListType = "Filter by Event Type ..."
	// ListTypeView represents the list for switching to an other view.
	ListTypeView = "Select View ..."
)

// ListWidget represents the ui widget component for a list.
type ListWidget struct {
	*w.List

	apiClient        *api.Client
	filterNamespaces []string
	filterNodes      []string
	filterStatuses   []string
	filterEventTypes []string
	sortNodes        []api.Sort
	sortPods         []api.Sort
	sortEvents       []api.Sort
	views            []ViewType
}

// NewListWidget returns a new list widget.
func NewListWidget(apiClient *api.Client) *ListWidget {
	list := w.NewList()
	list.TextStyle = ui.NewStyle(ui.ColorYellow)
	list.WrapText = false

	return &ListWidget{
		list,

		apiClient,
		[]string{},
		[]string{},
		[]string{"-", "Running", "Waiting", "Terminated"},
		[]string{"-", "Normal", "Warning"},
		[]api.Sort{api.SortCPUASC, api.SortCPUDESC, api.SortMemoryASC, api.SortMemoryDESC, api.SortName, api.SortPodsASC, api.SortPodsDESC},
		[]api.Sort{api.SortCPUASC, api.SortCPUDESC, api.SortMemoryASC, api.SortMemoryDESC, api.SortName, api.SortNamespace, api.SortRestartsASC, api.SortRestartsDESC, api.SortStatus},
		[]api.Sort{api.SortName, api.SortNamespace, api.SortTimeASC, api.SortTimeDESC},
		[]ViewType{ViewTypePods, ViewTypeNodes, ViewTypeEvents},
	}
}

// Hide hides the list.
func (l *ListWidget) Hide() {
	l.SetRect(0, 0, 0, 0)
}

// Selected determines the selected sortorder or filter.
func (l *ListWidget) Selected(viewType ViewType, listType ListType, sortorder api.Sort, filter api.Filter) (ViewType, api.Sort, api.Filter) {
	if viewType == ViewTypeNodes {
		if listType == ListTypeSort {
			sortorder = l.sortNodes[l.SelectedRow]
		}
	} else if viewType == ViewTypePods {
		if listType == ListTypeSort {
			sortorder = l.sortPods[l.SelectedRow]
		} else if listType == ListTypeFilterNamespace {
			if l.filterNamespaces[l.SelectedRow] == "-" {
				filter.Namespace = ""
			} else {
				filter.Namespace = l.filterNamespaces[l.SelectedRow]
			}
		} else if listType == ListTypeFilterNode {
			if l.filterNodes[l.SelectedRow] == "-" {
				filter.Node = ""
			} else {
				filter.Node = l.filterNodes[l.SelectedRow]
			}
		} else if listType == ListTypeFilterStatus {
			switch l.filterStatuses[l.SelectedRow] {
			case "-":
				filter.Status = 10
			case "Running":
				filter.Status = 2
			case "Waiting":
				filter.Status = 1
			case "Terminated":
				filter.Status = 0
			}
		}
	} else if viewType == ViewTypeEvents {
		if listType == ListTypeSort {
			sortorder = l.sortEvents[l.SelectedRow]
		} else if listType == ListTypeFilterNamespace {
			if l.filterNamespaces[l.SelectedRow] == "-" {
				filter.Namespace = ""
			} else {
				filter.Namespace = l.filterNamespaces[l.SelectedRow]
			}
		} else if listType == ListTypeFilterNode {
			if l.filterNodes[l.SelectedRow] == "-" {
				filter.Node = ""
			} else {
				filter.Node = l.filterNodes[l.SelectedRow]
			}
		} else if listType == ListTypeFilterEventType {
			if l.filterEventTypes[l.SelectedRow] == "-" {
				filter.EventType = ""
			} else {
				filter.EventType = l.filterEventTypes[l.SelectedRow]
			}
		}
	}

	if listType == ListTypeView {
		viewType = l.views[l.SelectedRow]
	}

	l.SetRect(0, 0, 0, 0)
	return viewType, sortorder, filter
}

// Show shows a list with the specified sort options or filters.
func (l *ListWidget) Show(viewType ViewType, listType ListType, termWidth, termHeight int) bool {
	var showList bool

	l.Title = string(listType)
	l.Rows = []string{}

	if viewType == ViewTypeNodes {
		// For the node view we only render the sort list.
		if listType == ListTypeSort {
			showList = true

			for index, item := range l.sortNodes {
				l.Rows = append(l.Rows, fmt.Sprintf("[%d] %s", index, item))
			}
		}
	} else if viewType == ViewTypePods {
		// For the pods view we render the sort list and the filters for namespace, node and status.
		// The namespaces and nodes are selected from the Kubernetes API first.
		if listType == ListTypeSort {
			showList = true

			for index, item := range l.sortPods {
				l.Rows = append(l.Rows, fmt.Sprintf("[%d] %s", index, item))
			}
		} else if listType == ListTypeFilterNamespace {
			showList = true
			l.filterNamespaces, _ = l.apiClient.GetNamespaces()

			for index, namespace := range l.filterNamespaces {
				l.Rows = append(l.Rows, fmt.Sprintf("[%d] %s", index, namespace))
			}
		} else if listType == ListTypeFilterNode {
			showList = true
			l.filterNodes, _ = l.apiClient.GetNodes()

			for index, node := range l.filterNodes {
				l.Rows = append(l.Rows, fmt.Sprintf("[%d] %s", index, node))
			}
		} else if listType == ListTypeFilterStatus {
			showList = true

			for index, status := range l.filterStatuses {
				l.Rows = append(l.Rows, fmt.Sprintf("[%d] %s", index, status))
			}
		}
	} else if viewType == ViewTypeEvents {
		// For the events view we render the sort list and the filters for namespace, node and event type.
		// The namespaces and nodes are selected from the Kubernetes API first.
		if listType == ListTypeSort {
			showList = true

			for index, item := range l.sortEvents {
				l.Rows = append(l.Rows, fmt.Sprintf("[%d] %s", index, item))
			}
		} else if listType == ListTypeFilterNamespace {
			showList = true
			l.filterNamespaces, _ = l.apiClient.GetNamespaces()

			for index, namespace := range l.filterNamespaces {
				l.Rows = append(l.Rows, fmt.Sprintf("[%d] %s", index, namespace))
			}
		} else if listType == ListTypeFilterNode {
			showList = true
			l.filterNodes, _ = l.apiClient.GetNodes()

			for index, node := range l.filterNodes {
				l.Rows = append(l.Rows, fmt.Sprintf("[%d] %s", index, node))
			}
		} else if listType == ListTypeFilterEventType {
			showList = true

			for index, eventType := range l.filterEventTypes {
				l.Rows = append(l.Rows, fmt.Sprintf("[%d] %s", index, eventType))
			}
		}
	}

	if listType == ListTypeView {
		showList = true

		for index, viewType := range l.views {
			l.Rows = append(l.Rows, fmt.Sprintf("[%d] %s", index, viewType))
		}
	}

	if showList {
		l.SelectedRow = 0
		l.SetRect(termWidth/2-25, termHeight/2-10, termWidth/2+25, termHeight/2+10)
	} else {
		l.SetRect(0, 0, 0, 0)
	}

	return showList
}
