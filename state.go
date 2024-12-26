// replacement for projects.go
// this should be everything related to abstract db
// sanitize directory/path
// Need information on current directory and/or given directory
// either empty string or string
// kv store search based on project name or path? duplicate data?
// or more complex? single data but then filtering based on path/lang/name
// do not allow duplicate path/lang/name
// migration support with version update and
// Ok KV store (based off path) O(1) and other queries do O(n)
// search (because most common usecase would be getting project state)
// and n is usually going to be small
// Other option is O(log n) for every operation and increased complexity
// fzf query based on stuff
// input custom db filler arguments
// project db Design
// full text search
// KV store with (pathname : ProjectData)
// generate index (seperate)? of words after looping through everything and then serialize
// OK man at this point a library is much better idk man
// Just implement basic feature first and then worry about it after community support
// ProjectData:
//
//	{
//	  Name
//	  Description
//
// }

// Consider blacklist (ENV Variable: PD_BLACKLIST)
package main

import (
	"compress/gzip"
	"encoding/gob"
	"fmt"
	// "log"
	"os"
)

type ProjectState struct {
	Projects map[string]Project
}

type ProjectStateActions interface {
	List()
	Add()
	Remove()
	Update()
}

// this is get
func List() (ProjectState, error) {
	var pd ProjectState
	// deserialize
	fi, err := os.Open("./pdstate")
	if err != nil {
		return pd, err
	}
	defer fi.Close()

	fz, err := gzip.NewReader(fi)
	if err != nil {
		return pd, err
	}
	defer fz.Close()

	decoder := gob.NewDecoder(fz)
	err = decoder.Decode(&pd)
	if err != nil {
		return pd, err
	}

	return pd, nil
}

func Add(path string) error {
	var origPd ProjectState
	var err error
	// TODO: this fucks up when file exists lol
	_, err = os.Stat("./pdstate")
	if err == nil {
		origPd, err = List()
		// TODO: Add new project if it doesn't exist
		if err != nil {
			return err
		}
		// 	return fullPath, nil
		// } else {
		// 	return "", errors.New("wrong file path")
		// }
	} else {
		if os.IsNotExist(err) {
			cool := make(map[string]Project)
			origPd = ProjectState{Projects: cool}
			// log.Fatal("File not Found !!")
		}
	}

	if val, ok := origPd.Projects[path]; ok {
		fmt.Println(val)
		return fmt.Errorf("project already exist %v", path)
	}

	// fi, err := os.Create("./pdstate")
	fi, err := os.OpenFile("./pdstate", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer fi.Close()

	// project := ProjectData{Name: getBase(path), Path: path, Kind: parseLanguage(path)}
	project := Project{Name: getBase(path), Path: path, Kind: "c"}
	origPd.Projects[path] = project

	// gzip instead of bytes?
	// var b bytes.Buffer
	fz := gzip.NewWriter(fi)
	defer fz.Close()

	// enc := gob.NewEncoder(&b)
	enc := gob.NewEncoder(fz)
	if err := enc.Encode(origPd); err != nil {
		return err
	}
	// fmt.Println(project)

	return nil
}

func Remove(path string) error {
	origPd, err := List()
	if err != nil {
		return err
	}

	if _, ok := origPd.Projects[path]; !ok {
		return fmt.Errorf("path already exist %v", path)
	}

	fi, err := os.OpenFile("./pdstate", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer fi.Close()

	delete(origPd.Projects, path)

	// gzip instead of bytes?
	// var b bytes.Buffer
	fz := gzip.NewWriter(fi)
	defer fz.Close()

	// enc := gob.NewEncoder(&b)
	enc := gob.NewEncoder(fz)
	if err := enc.Encode(origPd); err != nil {
		return err
	}

	return nil
}

func Update(path, key, value string) error {
	origPd, err := List()
	if err != nil {
		return err
	}

	if _, ok := origPd.Projects[path]; !ok {
		return fmt.Errorf("path already exist %v", path)
	}

	fi, err := os.OpenFile("./pdstate", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer fi.Close()

	// TODO: utilize methods
	// TODO: make this a pointer
	p := origPd.Projects[path]
	switch key {
	case "Path":
		p.Path = value
	case "Name":
		p.Name = value
	case "Kind":
		p.Kind = value
	case "Description":
		p.Description = value
	case "Priority":
		p.Description = value
	default:
		return fmt.Errorf("no such key %v in Project", key)
	}

	origPd.Projects[path] = p

	// gzip instead of bytes?
	// var b bytes.Buffer
	fz := gzip.NewWriter(fi)
	defer fz.Close()

	// enc := gob.NewEncoder(&b)
	enc := gob.NewEncoder(fz)
	if err := enc.Encode(origPd); err != nil {
		return err
	}

	return nil
}

//
// // used for cache
// func getProject() {
// }
