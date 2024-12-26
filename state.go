package main

import (
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"log"
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

const STATE_FILE = "./ppstate"

func readState() (ProjectState, error) {
	var ps ProjectState
	_, err := os.Stat(STATE_FILE)
	if err != nil {
		// NOTE: If doens't already exist, create empty object
		if os.IsNotExist(err) {
			ps = ProjectState{Projects: make(map[string]Project)}
			return ps, nil
		} else {
			return ps, err
		}
	}

	fi, err := os.Open(STATE_FILE)
	if err != nil {
		return ps, err
	}
	defer fi.Close()

	fz, err := gzip.NewReader(fi)
	if err != nil {
		return ps, err
	}
	defer fz.Close()

	decoder := gob.NewDecoder(fz)
	err = decoder.Decode(&ps)
	if err != nil {
		return ps, err
	}

	return ps, nil
}

func writeState(transform func() ProjectState) error {
	fi, err := os.OpenFile(STATE_FILE, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer fi.Close()

	pd := transform()

	fz := gzip.NewWriter(fi)
	defer fz.Close()

	// enc := gob.NewEncoder(&b)
	enc := gob.NewEncoder(fz)
	if err := enc.Encode(pd); err != nil {
		return err
	}

	return nil
}

func Get(path string) (Project, error) {
	pd, err := readState()
	var p Project
	if err != nil {
		return p, err
	}
	wd, err := NormalizePath(path)
	if err != nil {
		return p, err
	}
	project, ok := pd.Projects[wd]
	if ok {
		return project, nil
	}

	return p, nil
}

func Exists(path string) (bool, error) {
	pd, err := readState()
	if err != nil {
		return false, err
	}
	wd, err := NormalizePath(path)
	if err != nil {
		return false, err
	}
	_, ok := pd.Projects[wd]
	if ok {
		return true, nil
	}

	return false, nil
}

// this is get
func List() (ProjectState, error) {
	pd, err := readState()
	return pd, err
}

func Add(path string) error {
	pd, err := readState()
	if err != nil {
		return err
	}

	wd, err := NormalizePath(path)
	if err != nil {
		return err
	}

	if val, ok := pd.Projects[wd]; ok {
		fmt.Println(val)
		return fmt.Errorf("project already exist %v", wd)
	}

	err = writeState(func() ProjectState {
		project := Project{Name: getBase(wd), Path: wd, Kind: "c"}
		pd.Projects[path] = project
		return pd
	})
	if err != nil {
		return err
	}

	return nil
}

func Remove(path string) error {
	pd, err := readState()
	if err != nil {
		return err
	}

	wd, err := NormalizePath(path)
	if err != nil {
		return err
	}

	writeState(func() ProjectState {
		delete(pd.Projects, wd)
		return pd
	})

	return nil
}

func Update(path, key, value string) error {
	pd, err := List()
	if err != nil {
		return err
	}

	wd, err := NormalizePath(path)
	if err != nil {
		return err
	}

	// TODO: utilize methods
	// TODO: make this a pointer
	writeState(func() ProjectState {
		p := pd.Projects[wd]
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
			log.Fatalf("No such key %v in project", key)
		}
		pd.Projects[wd] = p
		return pd
	})

	return nil
}
