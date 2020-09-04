package proxmoxiso

import (
	"context"

	proxmoxapi "github.com/Telmate/proxmox-api-go/proxmox"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer/builder/proxmox/common"
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
)

// The unique id for the builder
const BuilderID = "proxmox.clone"

type Builder struct {
	config Config
}

// Builder implements packer.Builder
var _ packer.Builder = &Builder{}

var pluginVersion = "1.0.0"

func (b *Builder) ConfigSpec() hcldec.ObjectSpec { return b.config.FlatMapstructure().HCL2Spec() }

func (b *Builder) Prepare(raws ...interface{}) ([]string, []string, error) {
	return b.config.Prepare(raws...)
}

const downloadPathKey = "downloaded_iso_path"

func (b *Builder) Run(ctx context.Context, ui packer.Ui, hook packer.Hook) (packer.Artifact, error) {
	state := new(multistep.BasicStateBag)
	state.Put("iso-config", &b.config)
	state.Put("vm-creator", &isoVMCreator{})

	preSteps := []multistep.Step{
		&common.StepDownload{
			Checksum:    b.config.ISOChecksum,
			Description: "ISO",
			Extension:   b.config.TargetExtension,
			ResultKey:   downloadPathKey,
			TargetPath:  b.config.TargetPath,
			Url:         b.config.ISOUrls,
		},
	}
	for idx := range b.config.AdditionalISOFiles {
		preSteps = append(preSteps, &common.StepDownload{
			Checksum:    b.config.AdditionalISOFiles[idx].ISOChecksum,
			Description: "additional ISO",
			Extension:   b.config.AdditionalISOFiles[idx].TargetExtension,
			ResultKey:   b.config.AdditionalISOFiles[idx].downloadPathKey,
			TargetPath:  b.config.AdditionalISOFiles[idx].downloadPathKey,
			Url:         b.config.AdditionalISOFiles[idx].ISOUrls,
		})
	}
	preSteps = append(preSteps,
		&stepUploadISO{},
		&stepUploadAdditionalISOs{},
	)
	postSteps := []multistep.Step{
		&stepFinalizeISOTemplate{},
	}

	sb := proxmox.NewSharedBuilder(BuilderID, b.config.Config, preSteps, postSteps)
	return sb.Run(ctx, ui, hook, state)
}

type isoVMCreator struct{}

var _ proxmox.ProxmoxVMCreator = &isoVMCreator{}

func (*isoVMCreator) Create(vmRef *proxmoxapi.VmRef, config proxmoxapi.ConfigQemu, state multistep.StateBag) error {
	isoFile := state.Get("iso_file").(string)
	config.QemuIso = isoFile

	client := state.Get("proxmoxClient").(*proxmoxapi.Client)
	return config.CreateVm(vmRef, client)
}