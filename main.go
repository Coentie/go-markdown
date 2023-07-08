package main

import (
	"io/ioutil"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type config struct {
	EditWidget    *widget.Entry
	PreviewWidget *widget.RichText
	CurrentFile   fyne.URI
	SaveMenuItem  *fyne.MenuItem
}

var cfg config

func main() {

	app := app.New()
	win := app.NewWindow("Markdown")
	edit, preview := cfg.makeUI()

	win.SetContent(container.NewHSplit(edit, preview))
	// Create fine app
	// Create window
	cfg.createMenuItems(win)
	// Get UI
	// Set the content of the window
	// Show window and run app.
	win.Resize(fyne.Size{Width: 800, Height: 800})
	win.CenterOnScreen()
	win.ShowAndRun()
}

func (app *config) makeUI() (*widget.Entry, *widget.RichText) {
	edit := widget.NewMultiLineEntry()
	preview := widget.NewRichTextFromMarkdown("")

	app.EditWidget = edit
	app.PreviewWidget = preview

	edit.OnChanged = preview.ParseMarkdown

	return edit, preview
}

func (app *config) createMenuItems(win fyne.Window) {
	openMenu := fyne.NewMenuItem(" Open...", app.openFunc(win))
	saveMenu := fyne.NewMenuItem("Save...", app.saveFunc(win))
	app.SaveMenuItem = saveMenu
	app.SaveMenuItem.Disabled = true
	saveAsMenu := fyne.NewMenuItem(" Save as ...", app.saveAsFunc(win))

	fileMenu := fyne.NewMenu("File", openMenu, saveMenu, saveAsMenu)

	menu := fyne.NewMainMenu(fileMenu)

	win.SetMainMenu(menu)
}

var filter = storage.NewExtensionFileFilter([]string{".md", ".MD"})

func (app *config) saveFunc(win fyne.Window) func() {
	return func() {
		if app.CurrentFile != nil {
			write, err := storage.Writer(app.CurrentFile)

			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			write.Write([]byte(app.EditWidget.Text))

			defer write.Close()
		}
	}
}

func (app *config) saveAsFunc(win fyne.Window) func() {
	return func() {
		saveDialog := dialog.NewFileSave(func(write fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			if write == nil {
				return
			}

			if !strings.HasSuffix(strings.ToLower(write.URI().String()), ".md") {
				dialog.ShowInformation("Error", "Please name your file with a .md extension!", win)
			}

			write.Write([]byte(app.EditWidget.Text))
			app.CurrentFile = write.URI()

			defer write.Close()

			win.SetTitle(win.Title() + " - " + write.URI().Name())
			app.SaveMenuItem.Disabled = false
		}, win)
		saveDialog.SetFileName("untitled.md")
		saveDialog.SetFilter(filter)
		saveDialog.Show()
	}
}

func (app *config) openFunc(win fyne.Window) func() {
	return func() {
		openDialog := dialog.NewFileOpen(func(read fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			if read == nil {
				return
			}

			data, err := ioutil.ReadAll(read)

			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			app.EditWidget.SetText(string(data))

			win.SetTitle(win.Title() + " - " + read.URI().Name())
		}, win)
		openDialog.SetFilter(filter)

		openDialog.Show()
	}
}
