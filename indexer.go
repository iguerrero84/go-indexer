// This is the name of our package
// Everything with this package name can see everything
// else inside the same package, regardless of the file they are in
package main

// These are the libraries we are going to use
// Both "fmt" and "net" are part of the Go standard library
import (
	"encoding/json"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
)

type File struct {
	Index int
	Name  string
}

type Folder struct {
	Index    int
	StrIndex string
	Name     string
	Files    []*File
	Folders  map[string]*Folder
}

func (f *Folder) String() string {
	j, _ := json.Marshal(f)
	return string(j)
}

func isIgnoredFile(v string) bool {
	ignore := []string{".git", "/.git", "/.git/", ".gitignore", ".DS_Store", ".idea", "/.idea/", "/.idea"}

	for _, s := range ignore {
		if v == s {
			return true
		}
	}
	return false
}

func isMainFolder() bool {
	return false
}

func SortFoldersByIndex(Folder []*File) {
	sort.Slice(Folder, func(i, j int) bool {
		return Folder[i].Index < Folder[j].Index
	})
}

func BuildTree(dir string) *Folder {
	dir = path.Clean(dir)
	var tree *Folder
	var nodes = map[string]interface{}{}
	var walkFun filepath.WalkFunc = func(p string, info os.FileInfo, err error) error {

		if info.IsDir() {
			nodes[p] = &Folder{0, "", path.Base(p), []*File{}, map[string]*Folder{}}
		} else {
			nodes[p] = &File{0, path.Base(p)}
		}
		return nil
	}
	err := filepath.Walk(dir, walkFun)
	if err != nil {
		log.Fatal(err)
	}

	for key, value := range nodes {
		var parentFolder *Folder
		if key == dir {
			tree = value.(*Folder)
			continue
		} else {
			parentFolder = nodes[path.Dir(key)].(*Folder)
		}

		switch v := value.(type) {
		case *File:
			if !isIgnoredFile(v.Name) {
				parentFolder.Files = append(parentFolder.Files, v)
			}
			idx := len(parentFolder.Files) + 1
			v.Index = idx
		case *Folder:
			idx := len(parentFolder.Folders) + 1
			parentFolder.Folders[strconv.Itoa(idx)] = v
			v.Index = idx
			SortFoldersByIndex(parentFolder.Files)
		}
	}
	FolderIndexer(tree)
	return tree
}

func FolderIndexer(folder *Folder) {
	parentFolder := folder

	if folder.Name == "maildir" {
		folder.StrIndex = ""
	}
	for idx, f := range folder.Folders {

		if parentFolder.Name == "maildir" {
			f.StrIndex = idx
		} else {
			f.StrIndex = parentFolder.StrIndex + "." + idx
		}
		index, err := strconv.Atoi(idx)
		if err != nil {
			log.Fatal("Error Converting Idx To int")
		}
		f.Index = index
		FolderIndexer(f)
	}
}

var body = []string{
	`<h2>Collapsible Directory List</h2>
	<div class="box">
	<ul class="directory-list">`,
}

var foundFiles bool

func GetBodyDocument(folder *Folder) []string {
	body = append(body,
		`<li>`+folder.StrIndex+` <img src="assets/folder.png" width="55" height="50"> `+folder.Name)

	if len(folder.Files) > 0 {
		foundFiles = true
		body = append(body, `<ul>`)

		for _, fi := range folder.Files {
			body = append(body, `<li><img src="assets/file.png" width="45" height="30"> `+fi.Name+`</li>`)
		}
		if foundFiles {
			foundFiles = false
			body = append(body, `</ul></li>`)
		}
	}

	if len(folder.Folders) > 0 {
		body = append(body, `<ul>`)
		for _, fo := range folder.Folders {
			GetBodyDocument(fo)
		}
	}
	return body
}
