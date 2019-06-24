package widgets

import (
	"github.com/ricoberger/kubetop/pkg/api"

	ui "github.com/gizak/termui/v3"
)

// View represents all widgets which can be rendered as seperate view.
type View interface {
	ui.Drawable

	Filter() api.Filter
	Pause() bool
	SelectedValues() []string
	SelectNext()
	SelectPrev()
	SetSortAndFilter(sortorder api.Sort, filter api.Filter)
	Sortorder() api.Sort
	TabNext()
	TabPrev()
	TogglePause()
	Update() error
}

// ViewType implements all possible views for kubetop.
type ViewType string

const (
	// ViewTypeNodes represents the nodes view.
	// The nodes view is represented by the NodesWidget.
	ViewTypeNodes ViewType = "nodes"
	// ViewTypePodDetails represents the detail view for a pod.
	// The pod details view is represented by the PodDetailsWidget.
	ViewTypePodDetails ViewType = "poddetails"
	// ViewTypePods represents the pods view.
	// The pods view is represented by the PodsWidget.
	ViewTypePods ViewType = "pods"
)
