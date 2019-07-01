package widgets

import (
	"fmt"
	"time"

	"github.com/ricoberger/kubetop/pkg/api"
	"github.com/ricoberger/kubetop/pkg/term/helpers"

	ui "github.com/gizak/termui/v3"
)

// PodsWidget represents the ui widget component for the pods view.
type PodsWidget struct {
	*Table

	apiClient *api.Client
	filter    api.Filter
	pause     bool
	sortorder api.Sort
}

// NewPodsWidget returns a new pods widget.
// We create the table for the pods widget with all the basic layout settings.
func NewPodsWidget(apiClient *api.Client, filter api.Filter, sortorder api.Sort, termWidth, termHeight int) *PodsWidget {
	table := NewTable()
	table.Header = []string{"NAMESPACE", "POD", "READY", "STATUS", "RESTARTS", "CPU", "CPU MAX", "MEMORY", "MEMORY MAX", "IP", "AGE"}
	table.UniqueCol = 1

	table.SetRect(0, 0, termWidth, termHeight)

	table.ColWidths = []int{20, helpers.MaxInt(table.Inner.Dx()-150, 40), 10, 20, 10, 15, 15, 15, 15, 20, 10}
	table.ColResizer = func() {
		table.ColWidths = []int{20, helpers.MaxInt(table.Inner.Dx()-150, 40), 10, 20, 10, 15, 15, 15, 15, 20, 10}
	}

	table.Border = false
	table.BorderStyle = ui.NewStyle(ui.ColorClear)

	return &PodsWidget{
		table,

		apiClient,
		filter,
		false,
		sortorder,
	}
}

// Filter returns the setted filter.
func (p *PodsWidget) Filter() api.Filter {
	return p.filter
}

// Pause returns if updates are paused or not.
func (p *PodsWidget) Pause() bool {
	return p.pause
}

// SelectedValues returns the name of the selected pod.
func (p *PodsWidget) SelectedValues() []string {
	return p.Rows[p.SelectedRow]
}

// SelectNext selects the next item in the table.
func (p *PodsWidget) SelectNext() {
	p.ScrollDown()
}

// SelectPrev selects the previous item in the table.
func (p *PodsWidget) SelectPrev() {
	p.ScrollUp()
}

// SelectTop selects the top item in the table.
func (p *PodsWidget) SelectTop() {
	p.ScrollTop()
}

// SelectBottom selects the bottom item in the table.
func (p *PodsWidget) SelectBottom() {
	p.ScrollBottom()
}

// SelectHalfPageDown selects the item a half page down.
func (p *PodsWidget) SelectHalfPageDown() {
	p.ScrollHalfPageDown()
}

// SelectHalfPageUp selects the item a half page up.
func (p *PodsWidget) SelectHalfPageUp() {
	p.ScrollHalfPageUp()
}

// SelectPageDown selects the item on the next page.
func (p *PodsWidget) SelectPageDown() {
	p.ScrollPageDown()
}

// SelectPageUp selects the item on the previous page.
func (p *PodsWidget) SelectPageUp() {
	p.ScrollPageUp()
}

// SetSortAndFilter sets a new value for the sortorder and filter.
func (p *PodsWidget) SetSortAndFilter(sortorder api.Sort, filter api.Filter) {
	p.sortorder = sortorder
	p.filter = filter
}

// Sortorder returns the setted sortorder.
func (p *PodsWidget) Sortorder() api.Sort {
	return p.sortorder
}

// TogglePause sets toggle pause.
func (p *PodsWidget) TogglePause() {
	p.pause = !p.pause
}

// Update updates the table data of the pod view.
// Get the data for the pods widget and add each pod as seperate row to the table.
func (p *PodsWidget) Update() error {
	if !p.pause {
		pods, err := p.apiClient.GetPodsMetrics(p.filter, p.sortorder)
		if err != nil {
			return err
		}

		rows := make([][]string, len(pods))
		for i, pod := range pods {
			rows[i] = make([]string, 11)
			rows[i][0] = pod.Namespace
			rows[i][1] = pod.Name
			rows[i][2] = fmt.Sprintf("%d/%d", pod.ContainersReady, pod.ContainersCount)
			rows[i][3] = pod.Status
			rows[i][4] = fmt.Sprintf("%d", pod.Restarts)
			rows[i][5] = fmt.Sprintf("%dm", pod.CPU)
			rows[i][6] = helpers.RenderCPUMax(pod.CPUMax, pod.CPUMaxContainerCount, int64(pod.ContainersCount))
			rows[i][7] = helpers.FormatBytes(pod.Memory)
			rows[i][8] = helpers.RenderMemoryMax(pod.MemoryMax, pod.MemoryMaxContainerCount, int64(pod.ContainersCount))
			rows[i][9] = pod.IP
			rows[i][10] = helpers.FormatDuration(time.Now().Sub(pod.CreationDate))
		}

		p.Rows = rows
	}

	return nil
}
