package diff

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

func GetPatch(old, new, patchName string) {
	patch, err := exec.Command("diff", old, new).Output()
	fmt.Println(string(patch))
	// diff returns exit code 1 if diff is found, should not panic this "error"
	//if err != nil {
	//	log.Panic(err)
	//}
	err = ioutil.WriteFile(patchName+".patch", patch, 0644)
	if err != nil {
		log.Panic(err)
	}
}
