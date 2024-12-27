package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"sort"
	"strconv"
)

type ProjectState struct {
	Projects  map[string]Project
	Blacklist map[string]bool
}

type ProjectStateActions interface {
	List()
	Add()
	Remove()
	Update()
}

// const STATE_FILE = "./state"

// mkdir of state should be generated on install ig
//
//	err := os.MkdirAll(newpath, os.ModePerm)
func stateFile() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	switch runtime.GOOS {
	case "darwin":
		pa := path.Join(home, "Library/Application Support/sp/sp.db")
		return pa, nil
	case "linux":
		pa := path.Join(home, ".local/state/sp/sp.db")
		return pa, nil
	default:
		return "", errors.New("OS not supported")
	}
}

func readState() (ProjectState, error) {
	var ps ProjectState
	sf, err := stateFile()
	if err != nil {
		return ps, err
	}
	_, err = os.Stat(sf)
	if err != nil {
		// NOTE: If doens't already exist, create empty object
		if os.IsNotExist(err) {
			ps = ProjectState{
				Projects:  make(map[string]Project),
				Blacklist: make(map[string]bool),
			}
			return ps, nil
		} else {
			return ps, err
		}
	}

	sf, err = stateFile()
	if err != nil {
		return ps, err
	}
	fi, err := os.Open(sf)
	if err != nil {
		return ps, err
	}
	defer fi.Close()

	decoder := gob.NewDecoder(fi)
	err = decoder.Decode(&ps)
	if err != nil {
		return ps, err
	}

	return ps, nil
}

func writeState(transform func() ProjectState) error {
	sf, err := stateFile()
	if err != nil {
		return err
	}
	fi, err := os.OpenFile(sf, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer fi.Close()

	pd := transform()

	enc := gob.NewEncoder(fi)
	if err := enc.Encode(pd); err != nil {
		return err
	}

	return nil
}

// TODO: Update get
// not very clean... abstract basic db ops and others (like hooks)
// Get should just return Project not increment
// (only increment on selection through list)
func GetProject(path string) (Project, error) {
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
		// increment if exists
		err = incrementProjectPriority(project)
		if err != nil {
			return p, err
		}
		return project, nil
	}

	return p, nil
}

func incrementProjectPriority(p Project) error {
	oldPriority := p.Priority

	err := UpdateProject(p.Path, "Priority", strconv.Itoa(oldPriority+1))
	if err != nil {
		return err
	}
	return nil
}

func ProjectExists(path string) (bool, error) {
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

func sortProjectByPriority(p ProjectState) []Project {
	projects := make([]Project, 0, len(p.Projects))

	for _, v := range p.Projects {
		// Blacklist filter
		if !p.Blacklist[v.Path] {
			projects = append(projects, v)
		}
	}

	sort.SliceStable(projects, func(i, j int) bool {
		return p.Projects[projects[i].Path].Priority > p.Projects[projects[j].Path].Priority
	})

	return projects
}

func ListProject() ([]Project, error) {
	pd, err := readState()
	p := sortProjectByPriority(pd)
	return p, err
}

func AddProject(path string) error {
	pd, err := readState()
	if err != nil {
		return err
	}

	wd, err := NormalizePath(path)
	if err != nil {
		return err
	}

	if pd.Blacklist[wd] {
		return fmt.Errorf("path %v in blacklist", wd)
	}

	if _, ok := pd.Projects[wd]; ok {
		return fmt.Errorf("project already exist %v", wd)
	}

	err = writeState(func() ProjectState {
		project := Project{Name: getBase(wd), Path: wd, Kind: "c"}
		pd.Projects[wd] = project
		return pd
	})
	if err != nil {
		return err
	}

	return nil
}

func RemoveProject(path string) error {
	pd, err := readState()
	if err != nil {
		return err
	}

	wd, err := NormalizePath(path)
	if err != nil {
		return err
	}

	err = writeState(func() ProjectState {
		delete(pd.Projects, wd)
		return pd
	})
	if err != nil {
		return err
	}

	return nil
}

func ShowBlacklist() ([]string, error) {
	var bl []string
	bs, err := getBlacklist()
	if err != nil {
		return bl, err
	}

	for k, v := range bs {
		if v {
			bl = append(bl, k)
		}
	}
	return bl, nil
}

// when would I need this though?
func getBlacklist() (map[string]bool, error) {
	pd, err := readState()
	if err != nil {
		return pd.Blacklist, err
	}

	return pd.Blacklist, nil
}

func RemoveBlacklist(path string) error {
	pd, err := readState()
	if err != nil {
		return err
	}

	wd, err := NormalizePath(path)
	if err != nil {
		return err
	}

	err = writeState(func() ProjectState {
		delete(pd.Blacklist, wd)
		return pd
	})
	if err != nil {
		return err
	}

	return nil
}

func AddBlacklist(path string) error {
	pd, err := readState()
	if err != nil {
		return err
	}

	wd, err := NormalizePath(path)
	if err != nil {
		return err
	}

	err = writeState(func() ProjectState {
		pd.Blacklist[wd] = true
		return pd
	})
	if err != nil {
		return err
	}

	return nil
}

func UpdateProject(path, key, value string) error {
	pd, err := readState()
	if err != nil {
		return err
	}

	wd, err := NormalizePath(path)
	if err != nil {
		return err
	}

	// TODO: utilize methods
	// TODO: make this a pointer
	err = writeState(func() ProjectState {
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
			priority, err := strconv.Atoi(value)
			if err != nil {
				log.Fatal(err)
			}
			p.Priority = priority
		default:
			log.Fatalf("No such key %v in project", key)
		}
		pd.Projects[wd] = p
		return pd
	})
	if err != nil {
		return err
	}

	return nil
}
