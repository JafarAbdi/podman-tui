package imgdialogs

import (
	"fmt"

	"github.com/containers/podman-tui/ui/dialogs"
	"github.com/containers/podman-tui/ui/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

const (
	searchFieldMaxSize   = 60
	searchButtonWidth    = 10
	searchInpuLabelWidth = 13

	// focus elements
	sInputElement        = 1
	sSearchButtonElement = 2
	sSearchResultElement = 3
	sFormElement         = 4
)

// ImageSearchDialog represents image search dialogs
type ImageSearchDialog struct {
	*tview.Box
	layout              *tview.Flex
	searchLayout        *tview.Flex
	input               *tview.InputField
	searchButton        *tview.Button
	searchResult        *tview.Table
	form                *tview.Form
	result              [][]string
	display             bool
	focusElement        int
	cancelHandler       func()
	searchSelectHandler func()
	pullSelectHandler   func()
}

// NewImageSearchDialog returns new image search dialog primitive
func NewImageSearchDialog() *ImageSearchDialog {
	dialog := &ImageSearchDialog{
		Box:          tview.NewBox(),
		input:        tview.NewInputField(),
		searchButton: tview.NewButton("Search"),
		searchResult: tview.NewTable(),
		display:      false,
		focusElement: sInputElement,
	}
	bgColor := utils.Styles.ImageSearchDialog.BgColor
	fgColor := utils.Styles.ImageSearchDialog.FgColor
	dialog.input.SetLabel("search term: ")
	dialog.input.SetLabelColor(fgColor)
	dialog.input.SetFieldWidth(searchFieldMaxSize)
	dialog.input.SetBackgroundColor(bgColor)

	dialog.searchLayout = tview.NewFlex().SetDirection(tview.FlexColumn)
	dialog.searchLayout.AddItem(tview.NewBox().SetBackgroundColor(bgColor), 1, 0, true)
	dialog.searchLayout.AddItem(dialog.input, searchFieldMaxSize+searchInpuLabelWidth, 10, true)
	dialog.searchLayout.AddItem(tview.NewBox().SetBackgroundColor(bgColor), 1, 0, true)
	dialog.searchLayout.AddItem(dialog.searchButton, searchButtonWidth, 0, true)
	dialog.searchLayout.SetBackgroundColor(bgColor)

	stBgColor := utils.Styles.ImageSearchDialog.ResultTableBgColor
	stBorderColor := utils.Styles.ImageSearchDialog.ResultTableBorderColor
	stBorderTitleColor := utils.Styles.ImageSearchDialog.FgColor
	dialog.searchResult.SetBackgroundColor(stBgColor)
	dialog.searchResult.SetTitle("Search Result")
	dialog.searchResult.SetTitleColor(stBorderTitleColor)
	dialog.searchResult.SetBorder(true)
	dialog.searchResult.SetBorderColor(stBorderColor)

	dialog.initTable()

	dialog.form = tview.NewForm().
		AddButton("Cancel", nil).
		AddButton("Pull", nil).
		SetButtonsAlign(tview.AlignRight)
	dialog.form.SetBackgroundColor(bgColor)

	dialog.layout = tview.NewFlex().SetDirection(tview.FlexRow)
	dialog.layout.SetBorder(true)
	dialog.layout.SetBackgroundColor(bgColor)
	dialog.layout.SetTitle("PODMAN IMAGE SEARCH/PULL")
	dialog.layout.AddItem(tview.NewBox().SetBackgroundColor(bgColor), 1, 0, true)
	dialog.layout.AddItem(dialog.searchLayout, 1, 0, true)
	dialog.layout.AddItem(tview.NewBox().SetBackgroundColor(bgColor), 1, 0, true)
	dialog.layout.AddItem(dialog.searchResult, 1, 0, true)
	dialog.layout.AddItem(dialog.form, dialogs.DialogFormHeight, 0, true)
	return dialog
}

func (d *ImageSearchDialog) initTable() {
	bgColor := utils.Styles.ImageSearchDialog.ResultHeaderRow.BgColor
	fgColor := utils.Styles.ImageSearchDialog.ResultHeaderRow.FgColor
	d.searchResult.Clear()
	d.searchResult.SetCell(0, 0,
		tview.NewTableCell(fmt.Sprintf("[%s::]INDEX", utils.GetColorName(fgColor))).
			SetExpansion(1).
			SetBackgroundColor(bgColor).
			SetTextColor(fgColor).
			SetAlign(tview.AlignLeft).
			SetSelectable(false))
	d.searchResult.SetCell(0, 1,
		tview.NewTableCell(fmt.Sprintf("[%s::]NAME", utils.GetColorName(fgColor))).
			SetExpansion(1).
			SetBackgroundColor(bgColor).
			SetTextColor(fgColor).
			SetAlign(tview.AlignLeft).
			SetSelectable(false))
	d.searchResult.SetCell(0, 2,
		tview.NewTableCell(fmt.Sprintf("[%s::]DESCRIPTION", utils.GetColorName(fgColor))).
			SetExpansion(1).
			SetBackgroundColor(bgColor).
			SetTextColor(fgColor).
			SetAlign(tview.AlignLeft).
			SetSelectable(false))
	d.searchResult.SetCell(0, 3,
		tview.NewTableCell(fmt.Sprintf("[%s::]STARS", utils.GetColorName(fgColor))).
			SetExpansion(1).
			SetBackgroundColor(bgColor).
			SetTextColor(fgColor).
			SetAlign(tview.AlignCenter).
			SetSelectable(false))
	d.searchResult.SetCell(0, 4,
		tview.NewTableCell(fmt.Sprintf("[%s::]OFFICIAL", utils.GetColorName(fgColor))).
			SetExpansion(1).
			SetBackgroundColor(bgColor).
			SetTextColor(fgColor).
			SetAlign(tview.AlignCenter).
			SetSelectable(false))
	d.searchResult.SetCell(0, 5,
		tview.NewTableCell(fmt.Sprintf("[%s::]AUTOMATED", utils.GetColorName(fgColor))).
			SetExpansion(1).
			SetBackgroundColor(bgColor).
			SetTextColor(fgColor).
			SetAlign(tview.AlignCenter).
			SetSelectable(false))

	d.searchResult.SetFixed(1, 1)
	d.searchResult.SetSelectable(true, false)
}

// Display displays this primitive
func (d *ImageSearchDialog) Display() {
	d.display = true
}

// IsDisplay returns true if primitive is shown
func (d *ImageSearchDialog) IsDisplay() bool {
	return d.display
}

// Hide stops displaying this primitive
func (d *ImageSearchDialog) Hide() {
	d.focusElement = sInputElement
	d.display = false
	d.input.SetText("")
	d.result = [][]string{}
	d.initTable()
}

// Focus is called when this primitive receives focus
func (d *ImageSearchDialog) Focus(delegate func(p tview.Primitive)) {
	switch d.focusElement {
	case sInputElement:
		d.input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = sSearchButtonElement
				d.Focus(delegate)
				return nil
			}
			if event.Key() == tcell.KeyDown {
				d.focusElement = sSearchResultElement
				d.Focus(delegate)
				return nil
			}
			if event.Key() == tcell.KeyEnter {
				d.searchSelectHandler()
				return nil
			}
			return event
		})
		delegate(d.input)
		return
	case sSearchButtonElement:
		d.searchButton.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = sSearchResultElement
				d.Focus(delegate)
				return nil
			}
			if event.Key() == tcell.KeyEnter {
				d.searchSelectHandler()
				return nil
			}
			return event
		})
		delegate(d.searchButton)
		return
	case sSearchResultElement:
		d.searchResult.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = sFormElement
				d.Focus(delegate)
				return nil
			}
			if event.Key() == tcell.KeyEnter {
				d.pullSelectHandler()
				return nil
			}
			return event
		})
		delegate(d.searchResult)
		return
	case sFormElement:
		button := d.form.GetButton(d.form.GetButtonCount() - 1)
		button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyTab {
				d.focusElement = sInputElement
				d.Focus(delegate)
				d.form.SetFocus(0)
				return nil
			}
			if event.Key() == tcell.KeyEnter {
				d.pullSelectHandler()
				return nil
			}
			return event
		})
		delegate(d.form)
	}

}

// HasFocus returns whether or not this primitive has focus
func (d *ImageSearchDialog) HasFocus() bool {
	return d.form.HasFocus() || d.input.HasFocus() || d.searchResult.HasFocus() || d.searchButton.HasFocus()
}

// SetRect set rects for this primitive.
func (d *ImageSearchDialog) SetRect(x, y, width, height int) {
	dX := x + dialogs.DialogPadding
	dY := y + dialogs.DialogPadding
	dWidth := width - (2 * dialogs.DialogPadding)
	dHeight := height - (2 * dialogs.DialogPadding)

	//set search input field size
	iwidth := dWidth - searchInpuLabelWidth - searchButtonWidth - 2 - 2 - 1
	if iwidth > searchFieldMaxSize {
		iwidth = searchFieldMaxSize
	}
	d.input.SetFieldWidth(iwidth)
	d.searchLayout.ResizeItem(d.input, iwidth+searchInpuLabelWidth, 0)

	//set table height size
	d.layout.ResizeItem(d.searchResult, dHeight-dialogs.DialogFormHeight-5, 0)

	d.Box.SetRect(dX, dY, dWidth, dHeight)

}

// Draw draws this primitive onto the screen.
func (d *ImageSearchDialog) Draw(screen tcell.Screen) {

	if !d.display {
		return
	}
	bgColor := utils.Styles.ImageSearchDialog.BgColor
	d.Box.SetBackgroundColor(bgColor)
	d.Box.DrawForSubclass(screen, d)
	x, y, width, height := d.Box.GetInnerRect()
	d.layout.SetRect(x, y, width, height)
	d.layout.SetBorder(true)
	d.layout.SetBackgroundColor(bgColor)

	d.layout.Draw(screen)
}

//InputHandler returns input handler function for this primitive
func (d *ImageSearchDialog) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return d.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		log.Debug().Msgf("confirm dialog: event %v received", event)
		if event.Key() == tcell.KeyEsc {
			d.cancelHandler()
			return
		}
		if d.searchResult.HasFocus() {
			if searchResultHandler := d.searchResult.InputHandler(); searchResultHandler != nil {
				searchResultHandler(event, setFocus)
				return
			}

		}
		if d.input.HasFocus() {
			if inputFieldHandler := d.input.InputHandler(); inputFieldHandler != nil {
				inputFieldHandler(event, setFocus)
				return
			}
		}
		if d.searchButton.HasFocus() {
			if searchButtonHandler := d.searchButton.InputHandler(); searchButtonHandler != nil {
				searchButtonHandler(event, setFocus)
				return
			}
		}
		if d.form.HasFocus() {
			if formHandler := d.form.InputHandler(); formHandler != nil {
				formHandler(event, setFocus)
				return
			}
		}
	})
}

// SetCancelFunc sets form cancel button selected function
func (d *ImageSearchDialog) SetCancelFunc(handler func()) *ImageSearchDialog {
	d.cancelHandler = handler
	cancelButton := d.form.GetButton(d.form.GetButtonCount() - 2)
	cancelButton.SetSelectedFunc(handler)
	return d
}

// SetSearchFunc sets form cancel button selected function
func (d *ImageSearchDialog) SetSearchFunc(handler func()) *ImageSearchDialog {
	d.searchSelectHandler = handler
	return d
}

// SetPullFunc sets form pull button selected function
func (d *ImageSearchDialog) SetPullFunc(handler func()) *ImageSearchDialog {
	d.pullSelectHandler = handler
	return d
}

// GetSearchText returns search input field text
func (d *ImageSearchDialog) GetSearchText() string {
	return d.input.GetText()
}

//GetSelectedItem returns selected image name from search result table
func (d *ImageSearchDialog) GetSelectedItem() string {
	row, _ := d.searchResult.GetSelection()
	if row >= 0 {
		return d.result[row-1][1]
	}
	return ""
}

// UpdateResults updates result table
func (d *ImageSearchDialog) UpdateResults(data [][]string) {
	d.result = data

	d.initTable()
	alignment := tview.AlignLeft
	rowIndex := 1
	expand := 1
	for i := 0; i < len(data); i++ {
		index := data[i][0]
		name := data[i][1]
		desc := data[i][2]
		stars := data[i][3]
		official := data[i][4]
		if official == "[OK]" {
			official = "\u2705"
		}
		automated := data[i][5]
		if automated == "[OK]" {
			automated = "\u2705"
		}

		// index column
		d.searchResult.SetCell(rowIndex, 0,
			tview.NewTableCell(index).
				SetExpansion(expand).
				SetAlign(alignment))

		// name column
		d.searchResult.SetCell(rowIndex, 1,
			tview.NewTableCell(name).
				SetExpansion(expand).
				SetAlign(alignment))

		// description column
		d.searchResult.SetCell(rowIndex, 2,
			tview.NewTableCell(desc).
				SetExpansion(expand).
				SetAlign(alignment))

		// stars column
		d.searchResult.SetCell(rowIndex, 3,
			tview.NewTableCell(stars).
				SetExpansion(expand).
				SetAlign(tview.AlignCenter))

		// official column
		d.searchResult.SetCell(rowIndex, 4,
			tview.NewTableCell(official).
				SetExpansion(expand).
				SetAlign(tview.AlignCenter))

		// autoamted column
		d.searchResult.SetCell(rowIndex, 5,
			tview.NewTableCell(automated).
				SetExpansion(expand).
				SetAlign(tview.AlignCenter))
		rowIndex++
	}
	if len(data) > 0 {
		d.searchResult.Select(1, 1)
		d.searchResult.ScrollToBeginning()
	}
}
