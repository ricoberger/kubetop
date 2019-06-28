package widgets

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/ricoberger/kubetop/pkg/api"
	"github.com/ricoberger/kubetop/pkg/term/helpers"

	ui "github.com/gizak/termui/v3"
	w "github.com/gizak/termui/v3/widgets"
)

// PodDetailsWidget represents the ui widget component for the details view of a pod.
type PodDetailsWidget struct {
	*ui.Block

	podDetails1 *w.Paragraph
	podDetails2 *w.Paragraph
	containers  *Table
	logs        *w.Paragraph

	apiClient *api.Client
	filter    api.Filter
	name      string
	namespace string
	pause     bool
	sortorder api.Sort
}

// NewPodDetailsWidget returns a new pods widget.
// We create the table for the pods widget with all the basic layout settings.
func NewPodDetailsWidget(name, namespace string, apiClient *api.Client, filter api.Filter, sortorder api.Sort, termWidth, termHeight int) *PodDetailsWidget {
	block := ui.NewBlock()
	block.SetRect(0, 0, termWidth, termHeight)

	podDetails1 := w.NewParagraph()
	podDetails2 := w.NewParagraph()

	containers := NewTable()
	containers.Header = []string{"NAME", "RESTARTS", "STATUS", "CPU", "CPU MIN", "CPU MAX", "MEMORY", "MEMORY MIN", "MEMORY MAX"}
	containers.UniqueCol = 0
	containers.Border = false
	containers.BorderStyle = ui.NewStyle(ui.ColorClear)
	containers.ColWidths = []int{helpers.MaxInt(containers.Inner.Dx()-180, 40), 20, 40, 20, 20, 20, 20, 20, 20}
	containers.ColResizer = func() {
		containers.ColWidths = []int{helpers.MaxInt(containers.Inner.Dx()-180, 40), 20, 40, 20, 20, 20, 20, 20, 20}
	}

	logs := w.NewParagraph()
	logs.Border = true
	logs.Title = "Logs"
	logs.TitleStyle = ui.NewStyle(ui.ColorClear)
	logs.TextStyle = ui.NewStyle(ui.ColorClear)
	logs.WrapText = true

	return &PodDetailsWidget{
		block,

		podDetails1,
		podDetails2,
		containers,
		logs,

		apiClient,
		filter,
		name,
		namespace,
		false,
		sortorder,
	}
}

// Filter returns the setted filter.
func (p *PodDetailsWidget) Filter() api.Filter {
	return p.filter
}

// Pause returns if updates are paused or not.
func (p *PodDetailsWidget) Pause() bool {
	return p.pause
}

// SelectedValues returns the name of the selected pod.
func (p *PodDetailsWidget) SelectedValues() []string {
	return []string{}
}

// SelectNext selects the next log line.
func (p *PodDetailsWidget) SelectNext() {
}

// SelectPrev selects the previous log line.
func (p *PodDetailsWidget) SelectPrev() {
}

// SelectTop selects the top item in the table.
func (p *PodDetailsWidget) SelectTop() {
}

// SelectBottom selects the bottom item in the table.
func (p *PodDetailsWidget) SelectBottom() {
}

// SelectHalfPageDown selects the item a half page down.
func (p *PodDetailsWidget) SelectHalfPageDown() {
}

// SelectHalfPageUp selects the item a half page up.
func (p *PodDetailsWidget) SelectHalfPageUp() {
}

// SelectPageDown selects the item on the next page.
func (p *PodDetailsWidget) SelectPageDown() {
}

// SelectPageUp selects the item on the previous page.
func (p *PodDetailsWidget) SelectPageUp() {
}

// SetSortAndFilter sets a new value for the sortorder and filter.
func (p *PodDetailsWidget) SetSortAndFilter(sortorder api.Sort, filter api.Filter) {
	p.sortorder = sortorder
	p.filter = filter
}

// Sortorder returns the setted sortorder.
func (p *PodDetailsWidget) Sortorder() api.Sort {
	return p.sortorder
}

// TabNext selects the next container.
func (p *PodDetailsWidget) TabNext() {
	if p.containers.SelectedRow < len(p.containers.Rows)-1 {
		p.containers.ScrollDown()
	} else {
		p.containers.ScrollTop()
	}
}

// TabPrev selects the previous container.
func (p *PodDetailsWidget) TabPrev() {
	p.containers.ScrollUp()
}

// TogglePause sets toggle pause.
func (p *PodDetailsWidget) TogglePause() {
	p.pause = !p.pause
}

// Update updates the data for the details view of a pod.
func (p *PodDetailsWidget) Update() error {
	if !p.pause {
		pod, err := p.apiClient.GetPod(p.name, p.namespace, p.containers.SelectedRow)
		if err != nil {
			return err
		}

		// Render the first section of pod details: name, namespace, node, controlled by
		// First we create our string for the controlled by field.
		// Then we create our string for the events.
		// We render a maximum amount of five events, sorted by the timestamp (timestamp is the time when the event was fired the last time).
		var controlledBy string
		for index, controller := range pod.ControlledBy {
			if index == 0 {
				controlledBy = controlledBy + controller
			} else {
				controlledBy = controlledBy + "\n               " + controller
			}
		}

		sort.SliceStable(pod.Events, func(i, j int) bool {
			return pod.Events[i].Timestamp > pod.Events[j].Timestamp
		})

		var events string
		for index, event := range pod.Events {
			if index == 5 {
				break
			}

			if index == 0 {
				events = events + time.Unix(event.Timestamp, 0).Format("Mon, 02 Jan 2006 15:04:05 -0700") + ": " + event.Message
			} else {
				events = events + "\n               " + time.Unix(event.Timestamp, 0).Format("Mon, 02 Jan 2006 15:04:05 -0700") + ": " + event.Message
			}
		}

		p.podDetails1.Border = false
		p.podDetails1.Text = fmt.Sprintf(`
			Name:          %s
			Namespace:     %s
			Node:          %s
			Status:        %s
			Start Time:    %s
			IP:            %s
			Controlled By: %s
			Events:        %s`, pod.Name, pod.Namespace, pod.NodeName, pod.Status, pod.CreationDate.Format("Mon, 02 Jan 2006 15:04:05 -0700"), pod.IP, controlledBy, events)

		// Render the second section of pod details: labels, annotations
		// First we sort the labels by there key and then we create the string for rendering.
		labels := make([]string, 0, len(pod.Labels))
		for label := range pod.Labels {
			labels = append(labels, label)
		}
		sort.Strings(labels)

		var labelsStr string
		var labelsIndex int
		for _, key := range labels {
			if labelsIndex == 0 {
				labelsStr = labelsStr + key + "=" + pod.Labels[key]
			} else {
				labelsStr = labelsStr + "\n             " + key + "=" + pod.Labels[key]
			}
			labelsIndex++
		}

		annotations := make([]string, 0, len(pod.Annotations))
		for annotation := range pod.Annotations {
			annotations = append(annotations, annotation)
		}
		sort.Strings(annotations)

		var annotationsStr string
		var annotationsIndex int
		for _, key := range annotations {
			if annotationsIndex == 0 {
				annotationsStr = annotationsStr + key + "=" + pod.Annotations[key]
			} else {
				annotationsStr = annotationsStr + "\n             " + key + "=" + pod.Annotations[key]
			}
			annotationsIndex++
		}

		p.podDetails2.Border = false
		p.podDetails2.Text = fmt.Sprintf(`
			Labels:      %s
			Annotations: %s`, labelsStr, annotationsStr)

		// Render table with the containers.
		rows := make([][]string, len(pod.Containers))
		for i, container := range pod.Containers {
			rows[i] = make([]string, 9)
			rows[i][0] = container.Name
			rows[i][1] = fmt.Sprintf("%d", container.Restarts)
			rows[i][2] = container.Status
			rows[i][3] = fmt.Sprintf("%dm", container.CPU)
			rows[i][4] = helpers.RenderCPUMax(container.CPUMin, 1, 1)
			rows[i][5] = helpers.RenderCPUMax(container.CPUMax, 1, 1)
			rows[i][6] = helpers.FormatBytes(container.Memory)
			rows[i][7] = helpers.RenderMemoryMax(container.MemoryMin, 1, 1)
			rows[i][8] = helpers.RenderMemoryMax(container.MemoryMax, 1, 1)
		}

		p.containers.Rows = rows

		// Render log lines.
		// First reverse the order of the log lines, so the newest one is on top.
		// Then set the loglines as rows for the logs list.
		for i := len(pod.LogLines)/2 - 1; i >= 0; i-- {
			opp := len(pod.LogLines) - 1 - i
			pod.LogLines[i], pod.LogLines[opp] = pod.LogLines[opp], pod.LogLines[i]
		}

		var firstLogLine int
		if len(pod.LogLines) > 0 && pod.LogLines[0] == "" {
			firstLogLine = 1
		}

		p.logs.Text = strings.Join(pod.LogLines[firstLogLine:len(pod.LogLines)], "\n")

		// Bring it all together and calculate the position for podDetails1, podDetails2, containers and logs.
		// Caculate the position of the containers table based on the height of podDetails1 and podDetails2.
		// Use this value to set the positions of all elements.
		termWidth, termHeight := ui.TerminalDimensions()
		minHeight := 8
		detailsHeight := 11
		podDetails1Height := 8 + len(pod.ControlledBy) + helpers.MinInt(len(pod.Events), 5)
		if len(pod.ControlledBy) > 0 {
			podDetails1Height--
		}
		if len(pod.Events) > 0 {
			podDetails1Height--
		}
		podDetails2Height := len(labels) + len(annotations)
		if helpers.MaxInt(podDetails1Height, podDetails2Height) >= minHeight {
			detailsHeight = detailsHeight + helpers.MaxInt(podDetails1Height, podDetails2Height) - minHeight
		}

		p.podDetails1.SetRect(0, 0, termWidth/2, detailsHeight)
		p.podDetails2.SetRect(termWidth/2, 0, termWidth, detailsHeight)
		p.containers.SetRect(0, detailsHeight, termWidth, detailsHeight+5+(len(p.containers.Rows)))
		p.logs.SetRect(0, detailsHeight+5+(len(p.containers.Rows)), termWidth, termHeight-1)
	}

	return nil
}

// Draw renders our statusbar.
func (p *PodDetailsWidget) Draw(buf *ui.Buffer) {
	p.podDetails1.Draw(buf)
	p.podDetails2.Draw(buf)
	p.containers.Draw(buf)
	p.logs.Draw(buf)
}
