package clui

import (
	term "github.com/nsf/termbox-go"
)

type ConfirmationDialog struct {
	view    View
	parent  *Composer
	result  int
	onClose func()
}

type SelectDialog struct {
	view    View
	parent  *Composer
	result  int
	value   int
	rg      *RadioGroup
	list    *ListBox
	typ     SelectDialogType
	onClose func()
}

func NewConfirmationDialog(c *Composer, title, question string, buttons []string, defaultButton int) *ConfirmationDialog {
	dlg := new(ConfirmationDialog)

	if len(buttons) == 0 {
		buttons = []string{"OK"}
	}

	cw, ch := term.Size()

	dlg.parent = c
	dlg.view = c.CreateView(cw/2-12, ch/2-8, 20, 10, title)
	dlg.view.SetModal(true)
	dlg.view.SetPack(Vertical)
	NewFrame(dlg.view, dlg.view, 1, 1, BorderNone, DoNotScale)

	fbtn := NewFrame(dlg.view, dlg.view, 1, 1, BorderNone, 1)
	NewFrame(dlg.view, fbtn, 1, 1, BorderNone, DoNotScale)
	lb := NewLabel(dlg.view, fbtn, 10, 3, question, 1)
	NewFrame(dlg.view, fbtn, 1, 1, BorderNone, DoNotScale)
	lb.SetMultiline(true)

	NewFrame(dlg.view, dlg.view, 1, 1, BorderNone, DoNotScale)
	frm1 := NewFrame(dlg.view, dlg.view, 16, 4, BorderNone, DoNotScale)
	NewFrame(dlg.view, frm1, 1, 1, BorderNone, 1)

	bText := buttons[0]
	btn1 := NewButton(dlg.view, frm1, AutoSize, AutoSize, bText, DoNotScale)
	btn1.OnClick(func(ev Event) {
		dlg.result = DialogButton1
		c.DestroyView(dlg.view)
		if dlg.onClose != nil {
			go dlg.onClose()
		}
	})
	if defaultButton == DialogButton1 {
		dlg.view.ActivateControl(btn1)
	}
	var btn2, btn3 *Button

	if len(buttons) > 1 {
		NewFrame(dlg.view, frm1, 1, 1, BorderNone, 1)
		btn2 = NewButton(dlg.view, frm1, AutoSize, AutoSize, buttons[1], DoNotScale)
		btn2.OnClick(func(ev Event) {
			dlg.result = DialogButton2
			c.DestroyView(dlg.view)
			if dlg.onClose != nil {
				go dlg.onClose()
			}
		})
		if defaultButton == DialogButton2 {
			dlg.view.ActivateControl(btn2)
		}
	}
	if len(buttons) > 2 {
		NewFrame(dlg.view, frm1, 1, 1, BorderNone, 1)
		btn3 = NewButton(dlg.view, frm1, AutoSize, AutoSize, buttons[2], DoNotScale)
		btn3.OnClick(func(ev Event) {
			dlg.result = DialogButton3
			c.DestroyView(dlg.view)
			if dlg.onClose != nil {
				go dlg.onClose()
			}
		})
		if defaultButton == DialogButton3 {
			dlg.view.ActivateControl(btn3)
		}
	}

	NewFrame(dlg.view, frm1, 1, 1, BorderNone, 1)

	dlg.view.OnClose(func(ev Event) {
		if dlg.result == DialogAlive {
			dlg.result = DialogClosed
			if ev.X != 1 {
				c.DestroyView(dlg.view)
			}
			if dlg.onClose != nil {
				go dlg.onClose()
			}
		}
	})

	return dlg
}

func (d *ConfirmationDialog) OnClose(fn func()) {
	d.onClose = fn
}

func (d *ConfirmationDialog) Result() int {
	return d.result
}

// ------------------------ Selection Dialog ---------------------

func NewSelectDialog(c *Composer, title string, items []string, selectedItem int, typ SelectDialogType) *SelectDialog {
	dlg := new(SelectDialog)

	if len(items) == 0 {
		panic("Item list must contain at least 1 item")
	}

	cw, ch := term.Size()

	dlg.parent = c
	dlg.typ = typ
	dlg.view = c.CreateView(cw/2-12, ch/2-8, 20, 10, title)
	dlg.view.SetModal(true)
	dlg.view.SetPack(Vertical)

	if typ == SelectDialogList {
		fList := NewFrame(dlg.view, dlg.view, 1, 1, BorderNone, 1)
		fList.SetPaddings(1, 1, 0, 0)
		dlg.list = NewListBox(dlg.view, fList, 15, 5, 1)
		for _, item := range items {
			dlg.list.AddItem(item)
		}
		if selectedItem >= 0 && selectedItem < len(items) {
			dlg.list.SelectItem(selectedItem)
		}
	} else {
		fRadio := NewFrame(dlg.view, dlg.view, 1, 1, BorderNone, DoNotScale)
		fRadio.SetPaddings(1, 1, 0, 0)
		fRadio.SetPack(Vertical)
		dlg.rg = NewRadioGroup()
		for _, item := range items {
			r := NewRadio(dlg.view, fRadio, AutoSize, item, DoNotScale)
			dlg.rg.AddItem(r)
		}
		if selectedItem >= 0 && selectedItem < len(items) {
			dlg.rg.SetSelected(selectedItem)
		}
	}

	frm1 := NewFrame(dlg.view, dlg.view, 16, 4, BorderNone, DoNotScale)
	NewFrame(dlg.view, frm1, 1, 1, BorderNone, 1)
	btn1 := NewButton(dlg.view, frm1, AutoSize, AutoSize, "OK", DoNotScale)
	btn1.OnClick(func(ev Event) {
		dlg.result = DialogButton1
		if dlg.typ == SelectDialogList {
			dlg.value = dlg.list.SelectedItem()
		} else {
			dlg.value = dlg.rg.Selected()
		}
		c.DestroyView(dlg.view)
		if dlg.onClose != nil {
			go dlg.onClose()
		}
	})

	NewFrame(dlg.view, frm1, 1, 1, BorderNone, 1)
	btn2 := NewButton(dlg.view, frm1, AutoSize, AutoSize, "Cancel", DoNotScale)
	btn2.OnClick(func(ev Event) {
		dlg.result = DialogButton2
		dlg.value = -1
		c.DestroyView(dlg.view)
		if dlg.onClose != nil {
			go dlg.onClose()
		}
	})
	dlg.view.ActivateControl(btn2)
	NewFrame(dlg.view, frm1, 1, 1, BorderNone, 1)

	dlg.view.OnClose(func(ev Event) {
		if dlg.result == DialogAlive {
			dlg.result = DialogClosed
			if ev.X != 1 {
				c.DestroyView(dlg.view)
			}
			if dlg.onClose != nil {
				go dlg.onClose()
			}
		}
	})

	return dlg
}

func (d *SelectDialog) OnClose(fn func()) {
	d.onClose = fn
}

func (d *SelectDialog) Result() int {
	return d.result
}

func (d *SelectDialog) Value() int {
	return d.value
}
