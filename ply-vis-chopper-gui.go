package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type file_label struct {
	label *widget.Label
	io    interface{}
}

func select_file(w fyne.Window, ds fyne.Size, target file_label, write bool) {
	var d *dialog.FileDialog

	s := fyne.NewSize(800, 600)
	w.Resize(s)

	if write {
		d = dialog.NewFileSave(func(u fyne.URIWriteCloser, err error) {
			if err != nil || u == nil {
				if err != nil {
					dialog.ShowError(err, w)
				}
				return
			}
			target.label.SetText(u.URI().Name())
			target.io = u
			w.SetFixedSize(false)
			w.Resize(ds)
		}, w)
	} else {
		d = dialog.NewFileOpen(func(u fyne.URIReadCloser, err error) {
			if err != nil || u == nil {
				if err != nil {
					dialog.ShowError(err, w)
				}
				return
			}
			target.label.SetText(u.URI().Name())
			target.io = u
			w.SetFixedSize(false)
			w.Resize(ds)
		}, w)
	}

	d.Show()
	d.Resize(s)
}

func main() {
	a := app.New()
	w := a.NewWindow("ply-vis-chopper")
	default_size := fyne.NewSize(640, 0)

	fused_ply := file_label{widget.NewLabel("No file selected."), nil}
	fused_ply_vis := file_label{widget.NewLabel("No file selected."), nil}
	edited_ply := file_label{widget.NewLabel("No file selected."), nil}
	edited_ply_vis := file_label{widget.NewLabel("No file selected."), nil}

	w.SetContent(
		fyne.NewContainerWithLayout(layout.NewFormLayout(),
			&widget.Button{
				Alignment: widget.ButtonAlignLeading,
				Text:      "fused.ply",
				Icon:      theme.FileIcon(),
				OnTapped:  func() { select_file(w, default_size, fused_ply, false) },
			},
			fused_ply.label,
			&widget.Button{
				Alignment: widget.ButtonAlignLeading,
				Text:      "fused.ply.vis",
				Icon:      theme.FileIcon(),
				OnTapped:  func() { select_file(w, default_size, fused_ply_vis, false) },
			},
			fused_ply_vis.label,
			&widget.Button{
				Alignment: widget.ButtonAlignLeading,
				Text:      "edited.ply",
				Icon:      theme.FileIcon(),
				OnTapped:  func() { select_file(w, default_size, edited_ply, false) },
			},
			edited_ply.label,
			&widget.Button{
				Alignment: widget.ButtonAlignLeading,
				Text:      "edited.ply.vis",
				Icon:      theme.FileIcon(),
				OnTapped:  func() { select_file(w, default_size, edited_ply_vis, true) },
			},
			edited_ply_vis.label,
		),
	)

	w.Resize(default_size)
	w.SetFixedSize(true)

	w.ShowAndRun()
}
