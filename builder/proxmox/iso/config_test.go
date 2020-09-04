package proxmoxiso

import (
	"strings"
	"testing"

	"github.com/hashicorp/packer/template"
)

func TestBasicExampleFromDocsIsValid(t *testing.T) {
	const config = `{
  "builders": [
    {
      "type": "proxmox-iso",
      "proxmox_url": "https://my-proxmox.my-domain:8006/api2/json",
      "insecure_skip_tls_verify": true,
      "username": "apiuser@pve",
      "password": "supersecret",

      "node": "my-proxmox",
      "network_adapters": [
        {
          "bridge": "vmbr0"
        }
      ],
      "disks": [
        {
          "type": "scsi",
          "disk_size": "5G",
          "storage_pool": "local-lvm",
          "storage_pool_type": "lvm"
        }
      ],

      "iso_file": "local:iso/Fedora-Server-dvd-x86_64-29-1.2.iso",
      "http_directory":"config",
      "boot_wait": "10s",
      "boot_command": [
        "<up><tab> ip=dhcp inst.cmdline inst.ks=http://{{.HTTPIP}}:{{.HTTPPort}}/ks.cfg<enter>"
      ],

      "ssh_username": "root",
      "ssh_timeout": "15m",
      "ssh_password": "packer",

      "unmount_iso": true,
      "template_name": "fedora-29",
      "template_description": "Fedora 29-1.2, generated on {{ isotime \"2006-01-02T15:04:05Z\" }}"
    }
  ]
}`
	tpl, err := template.Parse(strings.NewReader(config))
	if err != nil {
		t.Fatal(err)
	}

	b := &Builder{}
	_, _, err = b.Prepare(tpl.Builders["proxmox-iso"].Config)
	if err != nil {
		t.Fatal(err)
	}

	// The example config does not set a number of optional fields. Validate that:
	// Memory 0 is too small, using default: 512
	// Number of cores 0 is too small, using default: 1
	// Number of sockets 0 is too small, using default: 1
	// CPU type not set, using default 'kvm64'
	// OS not set, using default 'other'
	// NIC 0 model not set, using default 'e1000'
	// Disk 0 cache mode not set, using default 'none'
	// Agent not set, default is true
	// SCSI controller not set, using default 'lsi'
	// Firewall toggle not set, using default: 0
	// Disable KVM not set, using default: 0

	if b.config.Memory != 512 {
		t.Errorf("Expected Memory to be 512, got %d", b.config.Memory)
	}
	if b.config.Cores != 1 {
		t.Errorf("Expected Cores to be 1, got %d", b.config.Cores)
	}
	if b.config.Sockets != 1 {
		t.Errorf("Expected Sockets to be 1, got %d", b.config.Sockets)
	}
	if b.config.CPUType != "kvm64" {
		t.Errorf("Expected CPU type to be 'kvm64', got %s", b.config.CPUType)
	}
	if b.config.OS != "other" {
		t.Errorf("Expected OS to be 'other', got %s", b.config.OS)
	}
	if b.config.NICs[0].Model != "e1000" {
		t.Errorf("Expected NIC model to be 'e1000', got %s", b.config.NICs[0].Model)
	}
	if b.config.NICs[0].Firewall != false {
		t.Errorf("Expected NIC firewall to be false, got %t", b.config.NICs[0].Firewall)
	}
	if b.config.Disks[0].CacheMode != "none" {
		t.Errorf("Expected disk cache mode to be 'none', got %s", b.config.Disks[0].CacheMode)
	}
	if b.config.Agent != true {
		t.Errorf("Expected Agent to be true, got %t", b.config.Agent)
	}
	if b.config.DisableKVM != false {
		t.Errorf("Expected Disable KVM toggle to be false, got %t", b.config.DisableKVM)
	}
	if b.config.SCSIController != "lsi" {
		t.Errorf("Expected SCSI controller to be 'lsi', got %s", b.config.SCSIController)
	}
	if b.config.CloudInit != false {
		t.Errorf("Expected CloudInit to be false, got %t", b.config.CloudInit)
	}
}