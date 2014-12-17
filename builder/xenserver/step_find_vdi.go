package xenserver

import (
	"fmt"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

type stepFindVdi struct {
	VdiName       string
	ImagePathFunc func() string
	VdiUuidKey    string
}

func (self *stepFindVdi) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	client := state.Get("client").(XenAPIClient)

	vdis, err := client.GetVdiByNameLabel(self.VdiName)

	switch {
	case len(vdis) == 0:
		ui.Error(fmt.Sprintf("Couldn't find a VDI named '%s'", self.VdiName))
		return multistep.ActionHalt
	case len(vdis) > 1:
		ui.Error(fmt.Sprintf("Found more than one VDI with name '%s'. Name must be unique", self.VdiName))
		return multistep.ActionHalt
	}

	vdi := vdis[0]

	vdiUuid, err := vdi.GetUuid()
	if err != nil {
		ui.Error(fmt.Sprintf("Unable to get UUID of VDI '%s': %s", self.VdiName, err.Error()))
		return multistep.ActionHalt
	}
	state.Put(self.VdiUuidKey, vdiUuid)

	return multistep.ActionContinue
}

func (self *stepFindVdi) Cleanup(state multistep.StateBag) {}
