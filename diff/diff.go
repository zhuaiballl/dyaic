package diff

import (
	"dyaic/utils"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func GenPatch(old, new, patchName string) {
	fmt.Println(patchName, ": diff begin")
	beginTime := time.Now()
	patch, err := exec.Command("diff", old, new).Output()
	//fmt.Println(string(patch))
	// diff returns exit code 1 if diff is found, should not panic this "error"
	//if err != nil {
	//	log.Panic(err)
	//}
	err = ioutil.WriteFile(patchName, patch, 0644)
	if err != nil {
		log.Panic(err)
	}
	endTime := time.Now()
	fmt.Println(patchName, ": diff finished in", endTime.Sub(beginTime))
	info, err := os.Stat(patchName)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("patch size:", info.Size())
}

func GenPatchForDirectory(old, new string) {
	locLen := len(new)
	beginTime := time.Now()
	repoSize := int64(0)
	err := filepath.Walk(new, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rLoc := path[locLen:]
		oldLoc := old + rLoc
		_, err = os.Stat(oldLoc)
		if !info.IsDir() {
			repoSize += info.Size()
		}
		if utils.Exist(err) {
			if info.IsDir() {
				return nil
			}
			if !utils.SameFile(path, oldLoc) { // file has been modified, sync needed
				fmt.Println("File has been modified:", rLoc, ", file size: ", info.Size())
				patchName := path + ".patch"
				GenPatch(oldLoc, path, patchName)
				fmt.Println("Generated patch file ", oldLoc, ".patch")
			}
		} else { // new file (or folder), creation needed
			if info.IsDir() {
				fmt.Println("New folder:", rLoc)
			} else {
				fmt.Println("New file:", rLoc)
			}
		}
		return nil
	})
	endTime := time.Now()
	fmt.Println("Gen patch for directory in", endTime.Sub(beginTime), ", directory size is", repoSize)
	if err != nil {
		log.Panic(err)
	}
}

func Patch(old, new, patchName string, clean bool) {
	cmd := exec.Command("patch", old, "-i", patchName, "-o", new)
	err := cmd.Run()
	if err != nil {
		log.Panic(err)
	}
	if clean {
		err = os.Remove(patchName)
		if err != nil {
			log.Panic(err)
		}
	}
}

func GenBSPatch(old, new, patchName string) {
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

func BSPatch(old, new, patchName string, clean bool) {
	cmd := exec.Command("bspatch", old, new, patchName)
	err := cmd.Run()
	if err != nil {
		log.Panic(err)
	}
	if clean {
		err = os.Remove(patchName)
		if err != nil {
			log.Panic(err)
		}
	}
}
