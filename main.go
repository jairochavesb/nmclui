package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	ui "github.com/VladimirMarkelov/clui"
)

var lblStatus *ui.Label

func createView() {

	view := ui.AddWindow(0, 0, 30, 10, " Available Networks ")

	canvas := ui.CreateFrame(view, 0, 0, ui.BorderNone, ui.AutoSize)
	canvas.SetPack(ui.Vertical)

	frm1 := ui.CreateFrame(canvas, 2, 2, ui.BorderThin, ui.Fixed)
	frm1.SetPack(ui.Vertical)

	lsbxSSIDs := ui.CreateListBox(frm1, 30, 10, ui.AutoSize)

	frmStatus := ui.CreateFrame(canvas, 1, 2, ui.BorderNone, ui.AutoSize)
	lblStatus = ui.CreateLabel(frmStatus, 50, 2, "", ui.AutoSize)

	go listWifiNets(lsbxSSIDs)
	lsbxSSIDs.SelectItem(0)
	go checkConnState(lblStatus)

	frm2 := ui.CreateFrame(canvas, 10, 3, ui.BorderNone, ui.Fixed)
	frm2.SetPack(ui.Horizontal)
	frm2.SetPaddings(1, 0)

	btnRescan := ui.CreateButton(frm2, ui.AutoSize, 1, "Refresh", ui.Fixed)
	btnSelectW := ui.CreateButton(frm2, ui.AutoSize, 1, "Select Wifi", ui.Fixed)

	_ = ui.CreateFrame(frm2, 10, 1, ui.BorderNone, ui.Fixed)

	btnQuit := ui.CreateButton(frm2, ui.AutoSize, 1, "Quit", ui.Fixed)

	ui.ActivateControl(view, lsbxSSIDs)

	btnSelectW.OnClick(func(ev ui.Event) {
		lblStatus.SetTitle("")
		lblStatus.SetTitle(lsbxSSIDs.SelectedItemText())

		if lsbxSSIDs.SelectedItemText() != "" {
			connect(lsbxSSIDs.SelectedItemText(), view)
		}
	})

	btnRescan.OnClick(func(ev ui.Event) {
		go listWifiNets(lsbxSSIDs)
	})

	btnQuit.OnClick(func(ev ui.Event) {
		ui.Stop()
	})
}

func connect(w string, v *ui.Window) {
	v.SetVisible(false)
	view2 := ui.AddWindow(2, 3, 10, 7, "")

	wifi := strings.TrimSpace(w)

	frmChk := ui.CreateFrame(view2, 8, 5, ui.BorderNone, ui.Fixed)
	frmChk.SetPack(ui.Vertical)
	frmChk.SetPaddings(1, 1)
	frmChk.SetGaps(1, 1)

	ui.CreateLabel(frmChk, ui.AutoSize, ui.AutoSize, "Enter password for: "+w, ui.Fixed)
	edFld := ui.CreateEditField(frmChk, 20, "", ui.Fixed)
	edFld.SetPasswordMode(true)

	chkPass := ui.CreateCheckBox(frmChk, ui.AutoSize, "Show Password", ui.Fixed)

	frmButtons := ui.CreateFrame(frmChk, 4, 4, ui.BorderNone, ui.Fixed)
	frmButtons.SetPack(ui.Horizontal)
	btnConnect := ui.CreateButton(frmButtons, 10, 1, "Connect", ui.Fixed)
	_ = ui.CreateFrame(frmButtons, 10, 2, ui.BorderNone, ui.AutoSize)
	btnBack := ui.CreateButton(frmButtons, 10, 1, "Back", ui.Fixed)

	ui.ActivateControl(view2, edFld)

	btnConnect.OnClick(func(ev ui.Event) {
		v.SetVisible(true)
		view2.SetVisible(false)
		out, err := exec.Command("nmcli", "dev", "wifi", "connect", wifi, "password", edFld.Title()).Output()
		if err != nil {
			fmt.Println(out)
		}
	})

	btnBack.OnClick(func(ev ui.Event) {
		v.SetVisible(true)
		view2.SetVisible(false)
	})

	chkPass.OnChange(func(state int) {
		ui.RefreshScreen()
		if state == 1 {
			edFld.SetPasswordMode(false)
		} else if state == 0 {
			edFld.SetPasswordMode(true)
		}
		ui.RefreshScreen()
	})

	go checkConnState(lblStatus)
	ui.RefreshScreen()
}

func checkConnState(lblStatus *ui.Label) {
	out, err := exec.Command("nmcli", "-f", "NAME,TYPE", "connection", "show", "--active").Output()
	if err != nil {
		lblStatus.SetTitle("")
		lblStatus.SetTitle(err.Error())
	} else {
		scanner := bufio.NewScanner(strings.NewReader(string(out)))

		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "wifi") {
				lblStatus.SetTitle("")
				lblStatus.SetTitle("Connected to: \n " + strings.Replace(scanner.Text(), "wifi", "", 1))
			}
		}
	}

	ui.RefreshScreen()
}

func listWifiNets(lbx1 *ui.ListBox) {
	lbx1.AddItem("SCANNING...")
	out, err := exec.Command("nmcli", "-f", "SSID", "device", "wifi").Output()
	if err != nil {
		lbx1.AddItem("There was a problem when scanning for wifi networks.")

	} else {
		lbx1.Clear()

		scanner := bufio.NewScanner(strings.NewReader(string(out)))

		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "SSID") {
				continue
			}
			lbx1.AddItem(scanner.Text())
		}
	}

	ui.RefreshScreen()
}

func main() {
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	_, err := os.Stat("./themes")
	if err == nil {
		ui.SetThemePath("themes")
		ui.SetCurrentTheme("simpleDark")
	}

	createView()

	ui.MainLoop()
}
