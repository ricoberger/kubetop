package widgets

import (
	"fmt"
	"image"
	"strings"

	ui "github.com/gizak/termui/v3"
)

// Table represents our table instance.
// We use the table component from gotop (https://github.com/cjbassi/gotop/blob/master/src/termui/table.go).
// We do not use the standard table component from termui, because the component does not support scrolling.
type Table struct {
	*ui.Block

	Header []string
	Rows   [][]string

	ColWidths []int
	ColGap    int
	PadLeft   int

	ShowCursor  bool
	CursorColor ui.Color

	ShowLocation bool

	UniqueCol    int
	SelectedItem string
	SelectedRow  int
	TopRow       int

	ColResizer func()
}

// NewTable returns a new table instance
func NewTable() *Table {
	return &Table{
		Block: ui.NewBlock(),

		ShowCursor: true,

		UniqueCol:   0,
		SelectedRow: 0,
		TopRow:      0,

		ColResizer: func() {},
	}
}

// Draw renders our table component.
func (t *Table) Draw(buf *ui.Buffer) {
	t.Block.Draw(buf)

	if t.ShowLocation {
		t.drawLocation(buf)
	}

	t.ColResizer()

	// Finds exact column starting position.
	colXPos := []int{}
	cur := 1 + t.PadLeft
	for _, w := range t.ColWidths {
		colXPos = append(colXPos, cur)
		cur += w
		cur += t.ColGap
	}

	// Prints the header.
	for i, h := range t.Header {
		width := t.ColWidths[i]
		if width == 0 {
			continue
		}
		// Do not render column if it does not fit in the widget.
		if width > (t.Inner.Dx()-colXPos[i])+1 {
			continue
		}

		// Render the background of the header.
		buf.SetString(
			strings.Repeat(" ", t.Inner.Dx()),
			ui.NewStyle(ui.ColorBlack, ui.ColorGreen),
			image.Pt(t.Inner.Min.X+colXPos[i]-1, t.Inner.Min.Y),
		)

		buf.SetString(
			h,
			ui.NewStyle(ui.ColorBlack, ui.ColorGreen),
			image.Pt(t.Inner.Min.X+colXPos[i]-1, t.Inner.Min.Y),
		)
	}

	if t.TopRow < 0 {
		return
	}

	// Print each row.
	for rowNum := t.TopRow; rowNum < t.TopRow+t.Inner.Dy()-1 && rowNum < len(t.Rows); rowNum++ {
		row := t.Rows[rowNum]
		y := (rowNum + 2) - t.TopRow

		// Print the cursor / selected row.
		style := ui.NewStyle(ui.ColorClear)
		if t.ShowCursor {
			if (t.SelectedItem == "" && rowNum == t.SelectedRow) || (t.SelectedItem != "" && t.SelectedItem == row[t.UniqueCol]) {
				style.Fg = ui.ColorBlack
				style.Bg = ui.ColorCyan
				for _, width := range t.ColWidths {
					if width == 0 {
						continue
					}
					buf.SetString(
						strings.Repeat(" ", t.Inner.Dx()),
						style,
						image.Pt(t.Inner.Min.X+1, t.Inner.Min.Y+y-1),
					)
				}
				t.SelectedItem = row[t.UniqueCol]
				t.SelectedRow = rowNum
			}
		}

		// Print each column of the row.
		for i, width := range t.ColWidths {
			if width == 0 {
				continue
			}
			// Do not render column if width is greater than distance to end of the widget.
			if width > (t.Inner.Dx()-colXPos[i])+1 {
				continue
			}
			r := ui.TrimString(row[i], width)
			buf.SetString(
				r,
				style,
				image.Pt(t.Inner.Min.X+colXPos[i]-1, t.Inner.Min.Y+y-1),
			)
		}
	}
}

// drawLocation renders the current location.
func (t *Table) drawLocation(buf *ui.Buffer) {
	total := len(t.Rows)
	topRow := t.TopRow + 1
	bottomRow := t.TopRow + t.Inner.Dy() - 1
	if bottomRow > total {
		bottomRow = total
	}

	loc := fmt.Sprintf(" %d - %d of %d ", topRow, bottomRow, total)

	width := len(loc)
	buf.SetString(loc, t.TitleStyle, image.Pt(t.Max.X-width-2, t.Min.Y))
}

// calcPos is used to calculate the cursor position and the current view into the table.
func (t *Table) calcPos() {
	t.SelectedItem = ""

	if t.SelectedRow < 0 {
		t.SelectedRow = 0
	}
	if t.SelectedRow < t.TopRow {
		t.TopRow = t.SelectedRow
	}

	if t.SelectedRow > len(t.Rows)-1 {
		t.SelectedRow = len(t.Rows) - 1
	}
	if t.SelectedRow > t.TopRow+(t.Inner.Dy()-2) {
		t.TopRow = t.SelectedRow - (t.Inner.Dy() - 2)
	}
}

// ScrollUp implements scroll up.
func (t *Table) ScrollUp() {
	t.SelectedRow--
	t.calcPos()
}

// ScrollDown implements scroll down.
func (t *Table) ScrollDown() {
	t.SelectedRow++
	t.calcPos()
}

// ScrollTop implements scroll to top.
func (t *Table) ScrollTop() {
	t.SelectedRow = 0
	t.calcPos()
}

// ScrollBottom implements scroll to bottom.
func (t *Table) ScrollBottom() {
	t.SelectedRow = len(t.Rows) - 1
	t.calcPos()
}

// ScrollHalfPageUp implements scroll up a half page.
func (t *Table) ScrollHalfPageUp() {
	t.SelectedRow = t.SelectedRow - (t.Inner.Dy()-2)/2
	t.calcPos()
}

// ScrollHalfPageDown implements scroll down a half page.
func (t *Table) ScrollHalfPageDown() {
	t.SelectedRow = t.SelectedRow + (t.Inner.Dy()-2)/2
	t.calcPos()
}

// ScrollPageUp implements scroll up a page.
func (t *Table) ScrollPageUp() {
	t.SelectedRow -= (t.Inner.Dy() - 2)
	t.calcPos()
}

// ScrollPageDown implements scroll down a page.
func (t *Table) ScrollPageDown() {
	t.SelectedRow += (t.Inner.Dy() - 2)
	t.calcPos()
}

// HandleClick handles a click.
func (t *Table) HandleClick(x, y int) {
	x = x - t.Min.X
	y = y - t.Min.Y
	if (x > 0 && x <= t.Inner.Dx()) && (y > 0 && y <= t.Inner.Dy()) {
		t.SelectedRow = (t.TopRow + y) - 2
		t.calcPos()
	}
}
