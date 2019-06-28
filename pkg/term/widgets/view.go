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
	SelectTop()
	SelectBottom()
	SelectHalfPageDown()
	SelectHalfPageUp()
	SelectPageDown()
	SelectPageUp()
	SetSortAndFilter(sortorder api.Sort, filter api.Filter)
	Sortorder() api.Sort
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
	// ViewTypeEvents represents the events view.
	// The events view is represented by the EventsWidget.
	ViewTypeEvents ViewType = "events"
	// ViewTypeEventDetails represents the event details view.
	// The event details view is represented by the EventDetailsWidget.
	ViewTypeEventDetails ViewType = "eventdetails"
)
