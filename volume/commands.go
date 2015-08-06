package volume

// Shell commands representing volume management operations.
// XXX replace these with real RPC calls.

import "os/exec"
import "path"
import "log"

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

func doCommand(cmd *exec.Cmd) error {
	log.Println("cmd %p: %s", cmd.Args)
	if op, err := cmd.CombinedOutput(); err != nil {
		log.Println("cmd %p: %v", err)
		log.Println("cmd %p: output %s", op)
		return err
	}
	return nil
}
