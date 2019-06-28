package widgets

import (
	"fmt"
	"time"

	"github.com/ricoberger/kubetop/pkg/api"
	"github.com/ricoberger/kubetop/pkg/term/helpers"

	ui "github.com/gizak/termui/v3"
)

// EventsWidget represents the ui widget component for the events view.
type EventsWidget struct {
	*Table

	apiClient *api.Client
	filter    api.Filter
	pause     bool
	sortorder api.Sort
}

// NewEventsWidget returns a new events widget.
// We create the table for the events widget with all the basic layout settings.
func NewEventsWidget(apiClient *api.Client, filter api.Filter, sortorder api.Sort, termWidth, termHeight int) *EventsWidget {
	table := NewTable()
	table.Header = []string{"", "AGE", "COUNT", "TYPE", "NAMESPACE", "NAME", "MESSAGE", "", "", "", ""}
	table.UniqueCol = 0

	table.SetRect(0, 0, termWidth, termHeight)

	table.ColWidths = []int{0, 10, 10, 10, 20, 50, helpers.MaxInt(table.Inner.Dx()-100, 80), 0, 0, 0, 0}
	table.ColResizer = func() {
		table.ColWidths = []int{0, 10, 10, 10, 20, 50, helpers.MaxInt(table.Inner.Dx()-100, 80), 0, 0, 0, 0}
	}

	table.Border = false
	table.BorderStyle = ui.NewStyle(ui.ColorClear)

	return &EventsWidget{
		table,

		apiClient,
		filter,
		false,
		sortorder,
	}
}

// Filter returns the setted filter.
func (e *EventsWidget) Filter() api.Filter {
	return e.filter
}

// Pause returns if updates are paused or not.
func (e *EventsWidget) Pause() bool {
	return e.pause
}

// SelectedValues returns the selected event.
func (e *EventsWidget) SelectedValues() []string {
	return e.Rows[e.SelectedRow]
}

// SelectNext selects the next item in the table.
func (e *EventsWidget) SelectNext() {
	e.ScrollDown()
}

// SelectPrev selects the previous item in the table.
func (e *EventsWidget) SelectPrev() {
	e.ScrollUp()
}

// SelectTop selects the top item in the table.
func (e *EventsWidget) SelectTop() {
	e.ScrollTop()
}

// SelectBottom selects the bottom item in the table.
func (e *EventsWidget) SelectBottom() {
	e.ScrollBottom()
}

// SelectHalfPageDown selects the item a half page down.
func (e *EventsWidget) SelectHalfPageDown() {
	e.ScrollHalfPageDown()
}

// SelectHalfPageUp selects the item a half page up.
func (e *EventsWidget) SelectHalfPageUp() {
	e.ScrollHalfPageUp()
}

// SelectPageDown selects the item on the next page.
func (e *EventsWidget) SelectPageDown() {
	e.ScrollPageDown()
}

// SelectPageUp selects the item on the previous page.
func (e *EventsWidget) SelectPageUp() {
	e.ScrollPageUp()
}

// SetSortAndFilter sets a new value for the sortorder and filter.
func (e *EventsWidget) SetSortAndFilter(sortorder api.Sort, filter api.Filter) {
	e.sortorder = sortorder
	e.filter = filter
}

// Sortorder returns the setted sortorder.
func (e *EventsWidget) Sortorder() api.Sort {
	return e.sortorder
}

// TogglePause sets toggle pause.
func (e *EventsWidget) TogglePause() {
	e.pause = !e.pause
}

// Update updates the table data of the pod view.
// Get the data for the pods widget and add each pod as seperate row to the table.
func (e *EventsWidget) Update() error {
	if !e.pause {
		events, err := e.apiClient.GetEvents(e.filter, e.sortorder)
		if err != nil {
			return err
		}

		rows := make([][]string, len(events))
		for i, event := range events {
			rows[i] = make([]string, 11)
			rows[i][0] = event.UID
			rows[i][1] = helpers.FormatDuration(time.Now().Sub(time.Unix(event.Timestamp, 0)))
			rows[i][2] = fmt.Sprintf("%d", event.Count)
			rows[i][3] = event.Type
			rows[i][4] = event.Namespace
			rows[i][5] = event.Name
			rows[i][6] = event.Message
			rows[i][7] = event.Kind
			rows[i][8] = event.Reason
			rows[i][9] = event.Source
			rows[i][10] = event.Node
		}

		e.Rows = rows
	}

	return nil
}
