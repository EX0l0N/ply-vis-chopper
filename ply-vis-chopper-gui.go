package main

import (
	"EXoloN/plyreader"
	"EXoloN/visio"
	"fmt"

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

const NO_FILE = "No file selected."

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

func (fr file_register) clean() {
	for i := range fr {
		fr[i] = false
	}
}

var f_reg file_register = make([]bool, 4)

func select_file(w fyne.Window, ds fyne.Size, gobut *widget.Button, target *file_label, write bool) {
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

func clean_fl(sfl []file_label) {
	for i := range sfl {
		sfl[i].label.SetText(NO_FILE)
	}
}

func chop_it(sfl []file_label, pb *widget.ProgressBar, w fyne.Window) {
	original_ply := sfl[FUSED].io.(fyne.URIReadCloser)
	defer original_ply.Close()
	original, err := plyreader.ReadPLY(original_ply)
	if err != nil {
		dialog.ShowError(err, w)
		return
	}
	cropped_ply := sfl[EDITED].io.(fyne.URIReadCloser)
	defer cropped_ply.Close()
	cropped, err := plyreader.ReadPLY(cropped_ply)
	if err != nil {
		dialog.ShowError(err, w)
		return
	}
	ply_vis := sfl[FUSED_VIS].io.(fyne.URIReadCloser)
	defer ply_vis.Close()
	vis, err := visio.ReadVis(ply_vis)
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	positions := make([]int, cropped.Elements())
	fpmax := float64(len(positions))

	for pos := range positions {
		if newpos, exists := original.GetPosition(cropped.GetPointAt(pos)); exists {
			positions[pos] = newpos
		} else {
			dialog.ShowError(fmt.Errorf("Point at position ", pos, " can not be found in original PLY file."), w)
			return
		}

		pb.SetValue(float64(pos) / fpmax)
	}

	chopped_vis := sfl[EDITED_VIS].io.(fyne.URIWriteCloser)

	vis.WriteListTo(positions, chopped_vis)
	chopped_vis.Close()

	dialog.ShowInformation("Done", "Successfully saved your vis file.", w)
}

func main() {
	a := app.New()
	w := a.NewWindow("ply-vis-chopper")
	default_size := fyne.NewSize(640, 200)
	file_labels := make([]file_label, 0, 4)

	file_labels = append(file_labels, file_label{widget.NewLabel(NO_FILE), FUSED, nil})
	file_labels = append(file_labels, file_label{widget.NewLabel(NO_FILE), FUSED_VIS, nil})
	file_labels = append(file_labels, file_label{widget.NewLabel(NO_FILE), EDITED, nil})
	file_labels = append(file_labels, file_label{widget.NewLabel(NO_FILE), EDITED_VIS, nil})

	progress := widget.NewProgressBar()
	chop_button := &widget.Button{
		Alignment: widget.ButtonAlignLeading,
		Text:      "Chop it!",
		Icon:      theme.ConfirmIcon(),
	}
	chop_button.OnTapped = func() {
		chop_button.Disable()
		progress.SetValue(0)
		chop_it(file_labels, progress, w)
		clean_fl(file_labels)
		f_reg.clean()
		progress.SetValue(1)
	}
	chop_button.Disable()

	w.SetContent(
		fyne.NewContainerWithLayout(layout.NewFormLayout(),
			&widget.Button{
				Alignment: widget.ButtonAlignLeading,
				Text:      "fused.ply",
				Icon:      theme.FileIcon(),
				OnTapped:  func() { select_file(w, default_size, chop_button, &file_labels[FUSED], false) },
			},
			file_labels[FUSED].label,
			&widget.Button{
				Alignment: widget.ButtonAlignLeading,
				Text:      "fused.ply.vis",
				Icon:      theme.FileIcon(),
				OnTapped:  func() { select_file(w, default_size, chop_button, &file_labels[FUSED_VIS], false) },
			},
			file_labels[FUSED_VIS].label,
			&widget.Button{
				Alignment: widget.ButtonAlignLeading,
				Text:      "edited.ply",
				Icon:      theme.FileIcon(),
				OnTapped:  func() { select_file(w, default_size, chop_button, &file_labels[EDITED], false) },
			},
			file_labels[EDITED].label,
			&widget.Button{
				Alignment: widget.ButtonAlignLeading,
				Text:      "edited.ply.vis",
				Icon:      theme.FileIcon(),
				OnTapped:  func() { select_file(w, default_size, chop_button, &file_labels[EDITED_VIS], true) },
			},
			file_labels[EDITED_VIS].label,
			chop_button,
			progress,
		),
	)

	w.Resize(default_size)

	w.ShowAndRun()
}
