package diff

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"time"
)

func GetPatch(old, new, patchName string) {
	fmt.Println(patchName, ": diff begin")
	beginTime := time.Now()
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
	endTime := time.Now()
	fmt.Println(patchName, ": diff finished in", endTime.Sub(beginTime))
}

func GetBSPatch(old, new, patchName string) {
	cmd := exec.Command("bsdiff", old, new, patchName)
	beginTime := time.Now()
	fmt.Println(patchName, ": bsdiff begin")
	err := cmd.Run()
	if err != nil { // bsdiff returns 0 on success and -1 on failure
		log.Panic(err)
	}
	endTime := time.Now()
	fmt.Println(patchName, ": bsdiff finished in", endTime.Sub(beginTime))
}
