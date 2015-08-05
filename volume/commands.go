package volume

// Shell commands representing volume management operations.
// XXX replace these with real RPC calls.

import "os/exec"
import "path"

func (m *VolumeMap) mountPoint(name string) string {
	return path.Join(m.mountBase, name)
}

func (m *VolumeMap) nfsUrl(name string) string {
	return path.Join(m.nfsServer+":", m.routerHost+":"+name)
}

func (m *VolumeMap) doMount(v *Volume) *exec.Cmd {
	return exec.Command("mount", v.DatastorePath, v.MountedPath)
}

func (m *VolumeMap) doUmount(v *Volume) *exec.Cmd {
	return exec.Command("umount", v.MountedPath)
}

func (m *VolumeMap) doCreate(v *Volume) *exec.Cmd {
	return exec.Command("sysmtool", "--host", m.routerHost,
		"--port", "9090",
		"--ns", "datastore",
		"--cmd", "create",
		"--name", v.Name)
}

func (m *VolumeMap) doRemove(v *Volume) *exec.Cmd {
	return exec.Command("sysmtool", "--host", m.routerHost,
		"--port", "9090",
		"--ns", "datastore",
		"--cmd", "create",
		"--name", v.Name)
}
