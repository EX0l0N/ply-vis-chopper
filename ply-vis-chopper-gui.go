package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

const (
	FUSED = iota
	FUSED_VIS
	EDITED
	EDITED_VIS
)

type file_label struct {
	label *widget.Label
	field int
	io    interface{}
}

type file_register []bool

func (fr file_register) register(pos int) bool {
	fr[pos] = true

	for _, p := range fr {
		if !p {
			return false
		}
	}
	return true
}

var f_reg file_register = make([]bool, 4)

func select_file(w fyne.Window, ds fyne.Size, gobut *widget.Button, target file_label, write bool) {
	var d *dialog.FileDialog

	s := fyne.NewSize(800, 600)
	w.Resize(s)

	if write {
		d = dialog.NewFileSave(func(u fyne.URIWriteCloser, err error) {
			w.Resize(ds)
			if err != nil || u == nil {
				if err != nil {
					dialog.ShowError(err, w)
				}
				return
			}
			target.label.SetText(u.URI().Name())
			if target.io != nil {
				wc := target.io.(fyne.URIWriteCloser)
				wc.Close()
			}
			target.io = u
			if f_reg.register(target.field) {
				gobut.Enable()
			}
		}, w)
	} else {
		d = dialog.NewFileOpen(func(u fyne.URIReadCloser, err error) {
			w.Resize(ds)
			if err != nil || u == nil {
				if err != nil {
					dialog.ShowError(err, w)
				}
				return
			}
			target.label.SetText(u.URI().Name())
			if target.io != nil {
				rc := target.io.(fyne.URIReadCloser)
				rc.Close()
			}
			target.io = u
			if f_reg.register(target.field) {
				gobut.Enable()
			}
		}, w)
	}

	d.Show()
	d.Resize(s)
}

func chop_it(sfl []file_label) {

}

func main() {
	a := app.New()
	w := a.NewWindow("ply-vis-chopper")
	default_size := fyne.NewSize(640, 200)
	file_labels := make([]file_label, 0, 4)

	file_labels = append(file_labels, file_label{widget.NewLabel("No file selected."), FUSED, nil})
	file_labels = append(file_labels, file_label{widget.NewLabel("No file selected."), FUSED_VIS, nil})
	file_labels = append(file_labels, file_label{widget.NewLabel("No file selected."), EDITED, nil})
	file_labels = append(file_labels, file_label{widget.NewLabel("No file selected."), EDITED_VIS, nil})

	chop_button := &widget.Button{
		Alignment: widget.ButtonAlignLeading,
		Text:      "Chop it!",
		Icon:      theme.ConfirmIcon(),
		OnTapped:  func() { chop_it(file_labels) },
	}
	chop_button.Disable()
	progress := widget.NewProgressBar()

	w.SetContent(
		fyne.NewContainerWithLayout(layout.NewFormLayout(),
			&widget.Button{
				Alignment: widget.ButtonAlignLeading,
				Text:      "fused.ply",
				Icon:      theme.FileIcon(),
				OnTapped:  func() { select_file(w, default_size, chop_button, file_labels[FUSED], false) },
			},
			file_labels[FUSED].label,
			&widget.Button{
				Alignment: widget.ButtonAlignLeading,
				Text:      "fused.ply.vis",
				Icon:      theme.FileIcon(),
				OnTapped:  func() { select_file(w, default_size, chop_button, file_labels[FUSED_VIS], false) },
			},
			file_labels[FUSED_VIS].label,
			&widget.Button{
				Alignment: widget.ButtonAlignLeading,
				Text:      "edited.ply",
				Icon:      theme.FileIcon(),
				OnTapped:  func() { select_file(w, default_size, chop_button, file_labels[EDITED], false) },
			},
			file_labels[EDITED].label,
			&widget.Button{
				Alignment: widget.ButtonAlignLeading,
				Text:      "edited.ply.vis",
				Icon:      theme.FileIcon(),
				OnTapped:  func() { select_file(w, default_size, chop_button, file_labels[EDITED_VIS], true) },
			},
			file_labels[EDITED_VIS].label,
			chop_button,
			progress,
		),
	)

	w.Resize(default_size)

	w.ShowAndRun()
}
