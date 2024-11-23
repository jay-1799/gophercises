package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var re = regexp.MustCompile("^(.+?) ([0-9]{4}) [(]([0-9]+) of ([0-9]+)[)][.](.+?)$")
var replaceString = "$2 - $1 - $3 of $4.$5"

func main() {
	var dry bool
	flag.BoolVar(&dry, "dry", true, "whether or not this should be a real or dry run")
	flag.Parse()
	walkDir := "sample"
	var toRename []string
	filepath.Walk(walkDir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if _, err := match(info.Name()); err == nil {
			toRename = append(toRename, path)
		}

		return nil
	})

	for _, oldPath := range toRename {
		dir := filepath.Dir(oldPath)
		filename := filepath.Base(oldPath)
		newFilename, _ := match(filename)
		newPath := filepath.Join(dir, newFilename)

		fmt.Printf("mv %s => %s\n", oldPath, newPath)
		if !dry {
			err := os.Rename(oldPath, newPath)
			if err != nil {
				fmt.Println("Error renaming:", oldPath, newPath, err.Error())
			}
		}
	}
	// for _, files := range toRename {
	// 	for _, f := range files {
	// 		fmt.Printf("%q\n", f)
	// 	}
	// }
	// for key, files := range toRename {
	// 	dir := filepath.Dir(key)
	// 	n := len(files)
	// 	sort.Strings(files)
	// 	for i, filename := range files {
	// 		res, _ := match(filename)
	// 		newFilename := fmt.Sprintf("%s - %d of %d.%s", res.base, (i + 1), n, res.ext)
	// 		oldPath := filepath.Join(dir, filename)
	// 		newPath := filepath.Join(dir, newFilename)
	// 		fmt.Printf("mv %s => %s\n", oldPath, newPath)
	// 		if !dry {
	// 			err := os.Rename(oldPath, newPath)
	// 			if err != nil {
	// 				fmt.Println("Error renaming:", oldPath, newPath, err.Error())
	// 			}
	// 		}
	// 	}
	// }

	// for _, orig := range toRename {
	// 	var n file
	// 	var err error
	// 	n.name, err = match(orig.name)
	// 	if err != nil {
	// 		fmt.Println("Error matching:", orig.path, err.Error())
	// 	}
	// 	n.path = filepath.Join(dir, n.name)
	// 	fmt.Printf("mv %s => %s\n", orig.path, n.path)

	// 	err = os.Rename(orig.path, n.path)
	// 	if err != nil {
	// 		fmt.Println("Error renaming:", orig.path, err.Error())
	// 	}
	// }
}

// returns the new file name or err
func match(filename string) (string, error) {
	// "birthday", "001" ".txt"
	if !re.MatchString(filename) {
		return "", fmt.Errorf("%s didn't match our pattern", filename)
	}
	return re.ReplaceAllString(filename, replaceString), nil
	// pieces := strings.Split(filename, ".")
	// ext := pieces[len(pieces)-1]
	// tmp := strings.Join(pieces[0:len(pieces)-1], ".")
	// pieces = strings.Split(tmp, "_")
	// name := strings.Join(pieces[0:len(pieces)-1], "_")
	// number, err := strconv.Atoi(pieces[len(pieces)-1])
	// if err != nil {
	// 	return "", fmt.Errorf("%s didn't match our pattern", fileName)
	// }

	// // Birthday - 1.txt
	// return "TODO", nil
}
