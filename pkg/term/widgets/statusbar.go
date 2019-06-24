package widgets

import (
	"fmt"
	"image"
	"strings"

	"github.com/ricoberger/kubetop/pkg/api"

	ui "github.com/gizak/termui/v3"
)

// StatusbarWidget represents the ui widget component for the statusbar.
type StatusbarWidget struct {
	*ui.Block

	apiClient *api.Client
	filter    api.Filter
	pause     bool
	sortorder api.Sort
	viewType  ViewType
}

// NewStatusbarWidget returns a new statusbar widget.
func NewStatusbarWidget(apiClient *api.Client, filter api.Filter, pause bool, sortorder api.Sort, viewType ViewType, termWidth, termHeight int) *StatusbarWidget {
	bar := ui.NewBlock()
	bar.Border = false

	bar.SetRect(0, termHeight-1, termWidth, termHeight)

	return &StatusbarWidget{
		bar,

		apiClient,
		filter,
		pause,
		sortorder,
		viewType,
	}
}

// Draw renders our statusbar.
func (s *StatusbarWidget) Draw(buf *ui.Buffer) {
	var paused string
	if s.pause {
		paused = "[P] Paused"
	} else {
		paused = "[P] Updated"
	}

	// Render an string of spaces to set the background for the whole statusbar to green.
	buf.SetString(
		strings.Repeat(" ", s.Inner.Dx()),
		ui.NewStyle(ui.ColorBlack, ui.ColorGreen),
		image.Pt(s.Inner.Min.X, s.Inner.Min.Y+(s.Inner.Dy()/2)),
	)

	if s.viewType == ViewTypeNodes {
		// Render sortorder.
		sortorder := fmt.Sprintf("[F1] Sorted by %s", string(s.sortorder))
		buf.SetString(
			sortorder,
			ui.NewStyle(ui.ColorBlack, ui.ColorGreen),
			image.Pt(s.Inner.Min.X, s.Inner.Min.Y+(s.Inner.Dy()/2)),
		)

		// Render pause.
		buf.SetString(
			paused,
			ui.NewStyle(ui.ColorBlack, ui.ColorGreen),
			image.Pt(s.Inner.Min.X+len(sortorder)+2, s.Inner.Min.Y+(s.Inner.Dy()/2)),
		)

		// Render clustername.
		// Also calculate the position where the clustername is shown.
		// The clustername is right aligned and if the terminal window is to small we cut of a part of the name.
		clustername := s.apiClient.GetClustername()
		clusternameX := s.Inner.Max.X - len(clustername)
		if s.Inner.Max.X-len(clustername) < s.Inner.Min.X+len(sortorder)+2+len(paused)+10 {
			clusternameX = s.Inner.Min.X + len(sortorder) + 2 + len(paused) + 10
		}

		buf.SetString(
			clustername,
			ui.NewStyle(ui.ColorBlack, ui.ColorGreen),
			image.Pt(clusternameX, s.Inner.Min.Y+(s.Inner.Dy()/2)),
		)
	} else if s.viewType == ViewTypePods {
		// Render sortorder.
		sortorder := fmt.Sprintf("[F1] Sorted by %s", string(s.sortorder))
		buf.SetString(
			sortorder,
			ui.NewStyle(ui.ColorBlack, ui.ColorGreen),
			image.Pt(s.Inner.Min.X, s.Inner.Min.Y+(s.Inner.Dy()/2)),
		)

		// Render namespace filter.
		filterNamespace := fmt.Sprintf("[F2] Namespace: %s", s.filter.Namespace)
		if s.filter.Namespace == "" {
			filterNamespace = "[F2] Namespace: -"
		}

		buf.SetString(
			filterNamespace,
			ui.NewStyle(ui.ColorBlack, ui.ColorGreen),
			image.Pt(s.Inner.Min.X+len(sortorder)+2, s.Inner.Min.Y+(s.Inner.Dy()/2)),
		)

		// Render node filter.
		filterNode := fmt.Sprintf("[F3] Node: %s", s.filter.Node)
		if s.filter.Node == "" {
			filterNode = "[F3] Node: -"
		}

		buf.SetString(
			filterNode,
			ui.NewStyle(ui.ColorBlack, ui.ColorGreen),
			image.Pt(s.Inner.Min.X+len(sortorder)+2+len(filterNamespace)+2, s.Inner.Min.Y+(s.Inner.Dy()/2)),
		)

		// Render status filter.
		var filterStatus string
		switch s.filter.Status {
		case 10:
			filterStatus = "[F4] Status: -"
		case 2:
			filterStatus = "[F4] Status: Running"
		case 1:
			filterStatus = "[F4] Status: Waiting"
		case 0:
			filterStatus = "[F4] Status: Terminated"
		}

		buf.SetString(
			filterStatus,
			ui.NewStyle(ui.ColorBlack, ui.ColorGreen),
			image.Pt(s.Inner.Min.X+len(sortorder)+2+len(filterNamespace)+2+len(filterNode)+2, s.Inner.Min.Y+(s.Inner.Dy()/2)),
		)

		// Render pause.
		buf.SetString(
			paused,
			ui.NewStyle(ui.ColorBlack, ui.ColorGreen),
			image.Pt(s.Inner.Min.X+len(sortorder)+2+len(filterNamespace)+2+len(filterNode)+2+len(filterStatus)+2, s.Inner.Min.Y+(s.Inner.Dy()/2)),
		)

		// Render clustername.
		// Also calculate the position where the clustername is shown.
		// The clustername is right aligned and if the terminal window is to small we cut of a part of the name.
		clustername := s.apiClient.GetClustername()
		clusternameX := s.Inner.Max.X - len(clustername)
		if s.Inner.Max.X-len(clustername) < s.Inner.Min.X+len(sortorder)+2+len(filterNamespace)+2+len(filterNode)+2+len(filterStatus)+2+len(paused)+10 {
			clusternameX = s.Inner.Min.X + len(sortorder) + 2 + len(filterNamespace) + 2 + len(filterNode) + 2 + len(filterStatus) + 2 + len(paused) + 10
		}

		buf.SetString(
			clustername,
			ui.NewStyle(ui.ColorBlack, ui.ColorGreen),
			image.Pt(clusternameX, s.Inner.Min.Y+(s.Inner.Dy()/2)),
		)
	} else if s.viewType == ViewTypePodDetails {
		// Render pause.
		buf.SetString(
			paused,
			ui.NewStyle(ui.ColorBlack, ui.ColorGreen),
			image.Pt(s.Inner.Min.X, s.Inner.Min.Y+(s.Inner.Dy()/2)),
		)

		// Render clustername.
		// Also calculate the position where the clustername is shown.
		// The clustername is right aligned and if the terminal window is to small we cut of a part of the name.
		clustername := s.apiClient.GetClustername()
		clusternameX := s.Inner.Max.X - len(clustername)
		if s.Inner.Max.X-len(clustername) < s.Inner.Min.X+len(paused)+10 {
			clusternameX = s.Inner.Min.X + len(paused) + 10
		}

		buf.SetString(
			clustername,
			ui.NewStyle(ui.ColorBlack, ui.ColorGreen),
			image.Pt(clusternameX, s.Inner.Min.Y+(s.Inner.Dy()/2)),
		)
	}
}

// SetPause sets a new value for pause.
func (s *StatusbarWidget) SetPause(pause bool) {
	s.pause = pause
}

// SetSortAndFilter sets a new value for the sortorder and filter.
func (s *StatusbarWidget) SetSortAndFilter(sortorder api.Sort, filter api.Filter) {
	s.sortorder = sortorder
	s.filter = filter
}

// SetViewType sets current view type.
func (s *StatusbarWidget) SetViewType(viewType ViewType) {
	s.viewType = viewType
}
