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

func (cli *CLI) commit(loc string) {
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
				// diff.ShowDiff(repoLoc, path)
				chs := diff.GenerateChanges(repoLoc, path)
				err = diff.Recover(repoLoc, &chs)
				if err != nil {
					return err
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

func (cli *CLI) hashFile(loc string) {
	hashBegin := time.Now()
	if loc == "" {
		loc = config.TempLocation
	}
	fmt.Println(utils.Md5File(loc))
	hashEnd := time.Now()
	fmt.Println(hashEnd.Sub(hashBegin))
}

func (cli *CLI) hashLoc(loc string) {
	if loc == "" {
		loc = config.TempLocation
	}
	err := filepath.Walk(loc, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		utils.Md5FileTest(path)
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
