package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

// сюда писать функцию DirTree

func dirTree(out io.Writer, path string, withFiles bool) error {
	dir, err := os.Stat(path)
	if err != nil {
		return err
	}
	if dir.Mode().IsDir() {
		PrintPath(out, path, withFiles, 0, map[int]bool{})
	}
	return nil
}

func PrintPath(out io.Writer, path string, withFiles bool, tab int, tabLine map[int]bool) {
	file, _ := os.Open(path)
	defer file.Close()
	var names[]string
	if withFiles{
		names, _ = file.Readdirnames(0)
	} else {
		names = prepareDirs(path, file)
	}
	sort.Strings(names)
	for i, name := range names {
		dirName := path + string(os.PathSeparator) + name
		dir, err := os.Stat(dirName)
		if err != nil {
			log.Fatal(err)
		}
		tabLine[tab] = i+1 != len(names) //if i+1 == len(names){tabLine[tab]=false}else ... true
		if tabLine[0] && tab != 0 {
			fmt.Fprintf(out, "│")
		}
		for j := 0; j < tab; j++ {
			fmt.Fprintf(out, "\t")
			if tabLine[j+1] && j+1 != tab {
				fmt.Fprintf(out, "│")
			}
		}
		if dir.Mode().IsDir() {
			if i+1 == len(names) {
				fmt.Fprintf(out, "└───")
			} else {
				fmt.Fprintf(out, "├───")
			}
			fmt.Fprintln(out, dir.Name())
			PrintPath(out, dirName, withFiles, tab+1, tabLine)
		} else {
			if i+1 == len(names) {
				fmt.Fprintf(out, "└───")
			} else {
				fmt.Fprintf(out, "├───")
			}
			fmt.Fprintf(out, "%s ", dir.Name())
			if dir.Size() == 0 {
				fmt.Fprintln(out, "(empty)")
			} else {
				fmt.Fprintf(out, "(%db)\n", dir.Size())
			}
		}
	}
}

func prepareDirs(path string, file *os.File) []string {
	var dirs[]string
	names, _ := file.Readdirnames(0)
	for _, name := range names {
		dirName := path + "/" + name
		dir, err := os.Stat(dirName)
		if err != nil {
			log.Fatal(err)
		}
		if dir.Mode().IsDir() {
			dirs = append(dirs, name)
		}
	}
	return dirs
}