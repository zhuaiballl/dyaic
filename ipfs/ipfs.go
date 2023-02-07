package ipfs

import "os/exec"

func Upload(loc string) error {
	cmd := exec.Command("ipfs", "add", loc)
	return cmd.Run()
}
