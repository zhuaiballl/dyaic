package cli

import (
	"dyaic/config"
	"dyaic/diff"
	"dyaic/monitor"
	"dyaic/utils"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
)

func (cli *CLI) commit(loc string, bs bool) {
	if loc == "" {
		loc = config.TempLocation
	}
	locLen := len(loc)
	err := filepath.Walk(loc, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rLoc := path[locLen:]
		repoLoc := config.RepoLocation + rLoc
		repoInfo, err := os.Stat(repoLoc)

		if utils.Exist(err) {
			if info.IsDir() {
				return nil
			}
			if info.ModTime().After(repoInfo.ModTime()) { // file has been modified, sync needed
				fmt.Println("File has been modified:", rLoc)
				patchName := repoLoc + ".patch"
				if bs {
					diff.GenBSPatch(repoLoc, path, patchName)
					diff.BSPatch(repoLoc, repoLoc, patchName, true)
				} else {
					diff.GenPatch(repoLoc, path, patchName)
					diff.Patch(repoLoc, repoLoc, patchName, true)
				}
				fmt.Println("Updated.")
				// TODO: send changes tx
				// TODO: sync changes with other nodes
			}
		} else { // new file (or folder), creation needed
			if info.IsDir() {
				fmt.Println("Creating folder:", repoLoc)
				err = os.Mkdir(repoLoc, 0755)
				if err != nil {
					return err
				}
			} else {
				fmt.Println("New file:", rLoc)
				utils.Copy(path, repoLoc)
				fmt.Println("Copied.")
				// TODO: send file tx
				// TODO: sync file with other nodes
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func (cli *CLI) gitwalker() {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Panic(err)
	}
	gitwalkerDir := homedir + "/.gitwalker/"
	for d := 1; ; d++ {
		newDir := gitwalkerDir + fmt.Sprintf("%04d", d)
		oldDir := gitwalkerDir + fmt.Sprintf("%04d", d+1)
		fmt.Printf("Start Patching %04d~%04d\n", d+1, d)
		diff.GenPatchForDirectory(oldDir, newDir)
	}
}

func (cli *CLI) hashFile(loc string) {
	hashBegin := time.Now()
	if loc == "" {
		loc = config.TempLocation
	}
	fmt.Println(utils.Md5File(loc))
	hashEnd := time.Now()
	fmt.Println(hashEnd.Sub(hashBegin))
}

func (cli *CLI) patch(loc string, bs bool) {
	if loc == "" {
		loc = config.TempLocation
	}
	locLen := len(loc)
	err := filepath.Walk(loc, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rLoc := path[locLen:]
		repoLoc := config.RepoLocation + rLoc
		//repoInfo, err := os.Stat(repoLoc)

		if utils.Exist(err) {
			if info.IsDir() {
				return nil
			}
			if !utils.SameFile(path, repoLoc) { // file has been modified, sync needed
				fmt.Println("File has been modified:", rLoc, ", file size: ", info.Size())
				patchName := repoLoc + ".patch"
				if bs {
					diff.GenBSPatch(repoLoc, path, patchName)
				} else {
					diff.GenPatch(repoLoc, path, patchName)
				}
				fmt.Println("Generated patch file ", repoLoc, ".patch")
			}
		} else { // new file (or folder), creation needed
			if info.IsDir() {
				fmt.Println("New folder:", repoLoc)
			} else {
				fmt.Println("New file:", rLoc)
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func (cli *CLI) printDiff(loc string) {
	if loc == "" {
		loc = config.TempLocation
	}
	locLen := len(loc)
	err := filepath.Walk(loc, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rLoc := path[locLen:]
		repoLoc := config.RepoLocation + rLoc
		repoInfo, err := os.Stat(repoLoc)

		if utils.Exist(err) {
			if info.IsDir() {
				return nil
			}
			if info.ModTime().After(repoInfo.ModTime()) { // file has been modified, sync needed
				chs := diff.GenerateChanges(repoLoc, path)
				if len(chs.Item) == 0 {
					return nil
				}
				fmt.Println("File has been modified:", rLoc)
				diff.ShowDiff(repoLoc, path)
			}
		} else { // new file (or folder), creation needed
			if info.IsDir() {
				fmt.Println("New folder:", repoLoc)
			} else {
				fmt.Println("New file:", rLoc)
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func (cli *CLI) saveDiff(loc string) {
	if loc == "" {
		loc = config.TempLocation
	}
	locLen := len(loc)
	err := filepath.Walk(loc, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rLoc := path[locLen:]
		repoLoc := config.RepoLocation + rLoc
		repoInfo, err := os.Stat(repoLoc)

		if utils.Exist(err) {
			if info.IsDir() {
				return nil
			}
			if info.ModTime().After(repoInfo.ModTime()) { // file has been modified, sync needed
				chs := diff.GenerateChanges(repoLoc, path)
				if len(chs.Item) == 0 {
					return nil
				}
				fmt.Println("File has been modified:", rLoc)
				diff.SaveDyaicDiff(repoLoc, path)
			}
		} else { // new file (or folder), creation needed
			if info.IsDir() {
				fmt.Println("New folder:", repoLoc)
			} else {
				fmt.Println("New file:", rLoc)
			}
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func (cli *CLI) printFolder(loc string) {
	if loc == "" {
		loc = config.TempLocation
	}
	err := filepath.Walk(loc, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(path, info.ModTime(), info.Size())
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func (cli *CLI) watch(loc string) {
	watcher := monitor.Watch(loc)
	defer watcher.Close()
	select {}
}
