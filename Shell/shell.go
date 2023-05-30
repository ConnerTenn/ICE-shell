package shell

import (
	"fmt"
	builtin "ice/Builtin"
	ice "ice/Lang"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var app *tview.Application

var primaryLayout *tview.Flex
var infoLayout *tview.Flex

var inputArea *tview.TextArea
var outputArea *tview.TextView
var returnArea *tview.TextView
var labelArea *tview.TextView
var positionArea *tview.TextView

func Initialize() {
	//================
	//== Create App ==
	//================
	//Create a new application
	app = tview.NewApplication()

	//====================
	//== Create Layouts ==
	//====================
	//Create a new flex box as the primary area and set its direction to row
	primaryLayout = tview.NewFlex().SetDirection(tview.FlexRow)
	//Create a new flex box for information display and set its direction to column
	infoLayout = tview.NewFlex().SetDirection(tview.FlexColumn)

	//==========================
	//== Create Primary Areas ==
	//==========================
	//Create a new text area for input
	inputArea = tview.NewTextArea()
	inputArea.SetTitle(" Input ").SetBorder(true)
	//Create new areas for output and return
	outputArea = tview.NewTextView()
	outputArea.SetTitle(" Output ").SetBorder(true)
	returnArea = tview.NewTextView()
	returnArea.SetTitle(" Return ").SetBorder(true)

	//=======================
	//== Create Info Areas ==
	//=======================
	//Create text views for info elements
	labelArea = tview.NewTextView().SetText(" ~Ramen~ ")
	positionArea = tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignRight)

	//=============
	//== Layouts ==
	//=============
	//+-----+--------+
	//|Label|Position|
	//+-----+--------+
	infoLayout.AddItem(labelArea, 0, 1, false)
	infoLayout.AddItem(positionArea, 0, 1, false)

	//+------+
	//|Hist  |
	//+------+
	//|Input |
	//+------+
	//|Footer|
	//+------+
	primaryLayout.AddItem(outputArea, 0, 1, false)   //Dynamic size (Fills the space)
	primaryLayout.AddItem(returnArea, 1+2, 1, false) //Init to 1 lines tall
	primaryLayout.AddItem(inputArea, 2+2, 1, true)   //Init to 2 lines tall
	primaryLayout.AddItem(infoLayout, 1, 1, true)    //1 line for info
}

func RunShell(globalMemFile string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	ctx := builtin.NewContext(globalMemFile)

	//Function for handling stdout
	go func() {
		obj := ctx.GlobalMem.Get("stdout")

		if ice.IsStream(obj) {
			stdout := ice.ToStream(obj)

			//Iterate through every line and append it to the output
			stdout.Itterate(func(obj ice.Obj) {
				ltrl := obj.ToLtrl()
				outputArea.Write([]byte(ltrl + "\n"))
			})
		}
	}()

	//==============
	//== Handlers ==
	//==============

	//== Input handler ==
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		//Check if run-action has occurred
		switch event.Key() {
		case tcell.KeyCtrlSpace:
			//Get code script from input area
			script := ice.Ltrl("[" + inputArea.GetText() + "]")

			// Parse the code into a ice.List object
			var codeRet ice.Obj
			code, iceErr := ice.ParseList(script.MakeSequ())
			if ice.IsErr(iceErr) {
				codeRet = ice.NewError("Shell InputCapture", "Failed to Parse into a list", iceErr)
			} else {
				// Evaluate the parsed code and run
				// Using the default ice.EvalRun state
				codeRet = ctx.Run(code)
			}

			returnArea.SetText(string(codeRet.ToLtrl()))

			inputArea.SetText("", true)

			return nil
		case tcell.KeyCtrlQ:
			app.Stop()

			return nil
		}

		return event
	})

	//== Cursor Moved ==
	updateInfos := func() {
		fromRow, fromColumn, toRow, toColumn := inputArea.GetCursor()
		if fromRow == toRow && fromColumn == toColumn {
			positionArea.SetText(fmt.Sprintf("Row: [yellow]%d[white], Column: [yellow]%d ", fromRow, fromColumn))
		} else {
			positionArea.SetText(fmt.Sprintf("[red]From[white] Row: [yellow]%d[white], Column: [yellow]%d[white] - [red]To[white] Row: [yellow]%d[white], To Column: [yellow]%d ", fromRow, fromColumn, toRow, toColumn))
		}
	}
	inputArea.SetMovedFunc(updateInfos)
	updateInfos()

	//== Input Area Changed ==
	//Set a function to be called when the input area is changed
	inputArea.SetChangedFunc(func() {
		//Count the number of lines in the input area
		numLines := len(strings.Split(inputArea.GetText(), "\n"))
		//If there are no lines, set the number of lines to 1
		if numLines < 1 {
			numLines = 1
		}

		//Resize the input area in the primary area to fit its current content
		primaryLayout.ResizeItem(inputArea, numLines+1+2, 1)
	})

	if err := app.SetRoot(primaryLayout, true).SetFocus(primaryLayout).Run(); err != nil {
		return err
	}

	builtin.SaveGlobalContext(globalMemFile, ctx)

	return nil
}
