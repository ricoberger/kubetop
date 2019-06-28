package widgets

import (
	"fmt"
	"time"

	"github.com/ricoberger/kubetop/pkg/api"
	"github.com/ricoberger/kubetop/pkg/term/helpers"

	ui "github.com/gizak/termui/v3"
	w "github.com/gizak/termui/v3/widgets"
)

// EventDetailsWidget represents the ui widget component for the details view of an event.
type EventDetailsWidget struct {
	*ui.Block

	eventDetails *w.Paragraph

	apiClient *api.Client
	filter    api.Filter
	name      string
	namespace string
	pause     bool
	sortorder api.Sort
}

// NewEventDetailsWidget returns an new event details widget.
func NewEventDetailsWidget(name, namespace string, apiClient *api.Client, filter api.Filter, sortorder api.Sort, termWidth, termHeight int) *EventDetailsWidget {
	block := ui.NewBlock()
	block.SetRect(0, 0, termWidth, termHeight)

	eventDetails := w.NewParagraph()

	return &EventDetailsWidget{
		block,

		eventDetails,

		apiClient,
		filter,
		name,
		namespace,
		false,
		sortorder,
	}
}

// Filter returns the setted filter.
func (e *EventDetailsWidget) Filter() api.Filter {
	return e.filter
}

// Pause returns if updates are paused or not.
func (e *EventDetailsWidget) Pause() bool {
	return e.pause
}

// SelectedValues returns the name of the selected pod.
func (e *EventDetailsWidget) SelectedValues() []string {
	return []string{}
}

// SelectNext selects the next log line.
func (e *EventDetailsWidget) SelectNext() {
}

// SelectPrev selects the previous log line.
func (e *EventDetailsWidget) SelectPrev() {
}

// SelectTop selects the top item in the table.
func (e *EventDetailsWidget) SelectTop() {
}

// SelectBottom selects the bottom item in the table.
func (e *EventDetailsWidget) SelectBottom() {
}

// SelectHalfPageDown selects the item a half page down.
func (e *EventDetailsWidget) SelectHalfPageDown() {
}

// SelectHalfPageUp selects the item a half page up.
func (e *EventDetailsWidget) SelectHalfPageUp() {
}

// SelectPageDown selects the item on the next page.
func (e *EventDetailsWidget) SelectPageDown() {
}

// SelectPageUp selects the item on the previous page.
func (e *EventDetailsWidget) SelectPageUp() {
}

// SetSortAndFilter sets a new value for the sortorder and filter.
func (e *EventDetailsWidget) SetSortAndFilter(sortorder api.Sort, filter api.Filter) {
	e.sortorder = sortorder
	e.filter = filter
}

// Sortorder returns the setted sortorder.
func (e *EventDetailsWidget) Sortorder() api.Sort {
	return e.sortorder
}

// TogglePause sets toggle pause.
func (e *EventDetailsWidget) TogglePause() {
	e.pause = !e.pause
}

// Update updates the data for the details view of a pod.
func (e *EventDetailsWidget) Update() error {
	if !e.pause {
		event := e.apiClient.GetEvent(e.name, e.namespace)

		e.eventDetails.Border = false
		e.eventDetails.Text = fmt.Sprintf(`
		UID:        %s
		Name:       %s
		Namespace:  %s
		Node:       %s
		Age:        %s
		First Time: %s
		Last Time:  %s
		Count:      %d
		Type:       %s
		Kind:       %s
		Reason:     %s
		Source:     %s
		Message:    %s`, event.UID, event.Name, event.Namespace, event.Node, helpers.FormatDuration(time.Now().Sub(time.Unix(event.Timestamp, 0))), event.FirstTimestamp.Format("Mon, 02 Jan 2006 15:04:05 -0700"), event.LastTimestamp.Format("Mon, 02 Jan 2006 15:04:05 -0700"), event.Count, event.Type, event.Kind, event.Reason, event.Source, event.Message)

		termWidth, termHeight := ui.TerminalDimensions()
		e.eventDetails.SetRect(0, 0, termWidth, termHeight)
	}

	return nil
}

// Draw renders our statusbar.
func (e *EventDetailsWidget) Draw(buf *ui.Buffer) {
	e.eventDetails.Draw(buf)
}
