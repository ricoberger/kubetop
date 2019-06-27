package widgets

import (
	"fmt"

	"github.com/ricoberger/kubetop/pkg/api"
	"github.com/ricoberger/kubetop/pkg/term/helpers"

	ui "github.com/gizak/termui/v3"
)

// NodesWidget represents the ui widget component for the nodes view.
type NodesWidget struct {
	*Table

	apiClient *api.Client
	filter    api.Filter
	pause     bool
	sortorder api.Sort
}

// NewNodesWidget returns a new nodes widget.
// We create the table for the nodes widget with all the basic layout settings.
func NewNodesWidget(apiClient *api.Client, filter api.Filter, sortorder api.Sort, termWidth, termHeight int) *NodesWidget {
	table := NewTable()
	table.Header = []string{"NAME", "PODS", "CPU", "MEMORY", "MEMORY MAX", "EXTERNAL IP", "INTERNAL IP"}
	table.UniqueCol = 0

	table.SetRect(0, 0, termWidth, termHeight)

	table.ColWidths = []int{helpers.MaxInt(table.Inner.Dx()-160, 40), 20, 20, 20, 20, 40, 40}
	table.ColResizer = func() {
		table.ColWidths = []int{helpers.MaxInt(table.Inner.Dx()-160, 40), 20, 20, 20, 20, 40, 40}
	}

	table.Border = false
	table.BorderStyle = ui.NewStyle(ui.ColorClear)

	return &NodesWidget{
		table,

		apiClient,
		filter,
		false,
		sortorder,
	}
}

// Filter returns the setted filter.
func (n *NodesWidget) Filter() api.Filter {
	return n.filter
}

// Pause returns if updates are paused or not.
func (n *NodesWidget) Pause() bool {
	return n.pause
}

// SelectedValues returns the name of the selected row.
func (n *NodesWidget) SelectedValues() []string {
	return n.Rows[n.SelectedRow]
}

// SelectNext selects the next item in the table.
func (n *NodesWidget) SelectNext() {
	n.ScrollDown()
}

// SelectPrev selects the previous item in the table.
func (n *NodesWidget) SelectPrev() {
	n.ScrollUp()
}

// SelectTop selects the top item in the table.
func (n *NodesWidget) SelectTop() {
	n.ScrollTop()
}

// SelectBottom selects the bottom item in the table.
func (n *NodesWidget) SelectBottom() {
	n.ScrollBottom()
}

// SelectHalfPageDown selects the item a half page down.
func (n *NodesWidget) SelectHalfPageDown() {
	n.ScrollHalfPageDown()
}

// SelectHalfPageUp selects the item a half page up.
func (n *NodesWidget) SelectHalfPageUp() {
	n.ScrollHalfPageUp()
}

// SelectPageDown selects the item on the next page.
func (n *NodesWidget) SelectPageDown() {
	n.ScrollPageDown()
}

// SelectPageUp selects the item on the previous page.
func (n *NodesWidget) SelectPageUp() {
	n.ScrollPageUp()
}

// SetSortAndFilter sets a new value for the sortorder and filter.
func (n *NodesWidget) SetSortAndFilter(sortorder api.Sort, filter api.Filter) {
	n.sortorder = sortorder
	n.filter = filter
}

// Sortorder returns the setted sortorder.
func (n *NodesWidget) Sortorder() api.Sort {
	return n.sortorder
}

// TabNext does nothing.
func (n *NodesWidget) TabNext() {
}

// TabPrev does nothing.
func (n *NodesWidget) TabPrev() {
}

// TogglePause sets toggle pause.
func (n *NodesWidget) TogglePause() {
	n.pause = !n.pause
}

// Update updates the table data of the node view.
// Get the data for the nodes widget and add each node as seperate row to the table.
func (n *NodesWidget) Update() error {
	if !n.pause {
		nodes, err := n.apiClient.GetNodesMetrics(n.sortorder)
		if err != nil {
			return err
		}

		strings := make([][]string, len(nodes))
		for i, node := range nodes {
			strings[i] = make([]string, 7)
			strings[i][0] = node.Name
			strings[i][1] = fmt.Sprintf("%d", node.PodsCount)
			strings[i][2] = fmt.Sprintf("%.2f%%", (float64(node.CPUUsed) * 100.0 / float64(node.CPUTotal)))
			strings[i][3] = fmt.Sprintf("%.2f%%", (float64(node.MemoryUsed) * 100.0 / float64(node.MemoryTotal)))
			strings[i][4] = helpers.FormatBytes(node.MemoryTotal)
			strings[i][5] = node.ExternalIP
			strings[i][6] = node.InternalIP
		}

		n.Rows = strings
	}

	return nil
}
