package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
)

type ProjectState struct {
	Projects  map[string]Project
	Blacklist map[string]bool
}

func getStateFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library/Application Support/sp/sp.db"), nil
	case "linux":
		return filepath.Join(home, ".local/state/sp/sp.db"), nil
	default:
		return "", errors.New("unsupported OS")
	}
}

//	func (ps *ProjectState) InitProjectState() error {
//		sf, err := stateFile()
//		if err != nil {
//			return err
//		}
//		_, err = os.Stat(sf)
//		if err != nil {
//			// NOTE: If doens't already exist, create empty object
//			if os.IsNotExist(err) {
//				// How do I use this? does this get garbage collected? what's the friggin problem of the warning?
//				ps = &ProjectState{
//					Projects:  make(map[string]Project),
//					Blacklist: make(map[string]bool),
//				}
//				return nil
//			} else {
//				return err
//			}
//		}
//
//		sf, err = stateFile()
//		if err != nil {
//			return err
//		}
//		fi, err := os.Open(sf)
//		if err != nil {
//			return err
//		}
//		defer fi.Close()
//
//		decoder := gob.NewDecoder(fi)
//		// This initializes the project
//		err = decoder.Decode(ps)
//		if err != nil {
//			return err
//		}
//
//		return nil
//	}
func (ps *ProjectState) LoadState() error {
	stateFilePath, err := getStateFilePath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(stateFilePath); errors.Is(err, os.ErrNotExist) {
		ps.Projects = make(map[string]Project)
		ps.Blacklist = make(map[string]bool)
		return nil
	} else if err != nil {
		return err
	}

	file, err := os.Open(stateFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	return decoder.Decode(ps)
}

//	func (ps *ProjectState) writeState(transform func()) error {
//		sf, err := stateFile()
//		if err != nil {
//			return err
//		}
//		fi, err := os.OpenFile(sf, os.O_RDWR|os.O_CREATE, 0644)
//		if err != nil {
//			return err
//		}
//		defer fi.Close()
//
//		transform()
//
//		enc := gob.NewEncoder(fi)
//		if err := enc.Encode(ps); err != nil {
//			return err
//		}
//
//		return nil
//	}
func (ps *ProjectState) SaveState() error {
	stateFilePath, err := getStateFilePath()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(stateFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(ps)
}

func (ps *ProjectState) GetProject(path string) (Project, error) {
	normalizedPath, err := NormalizePath(path)
	if err != nil {
		return Project{}, fmt.Errorf("failed to normalize path: %w", err)
	}

	project, exists := ps.Projects[normalizedPath]
	if !exists {
		// return Project{}, fmt.Errorf("project not found: %s", normalizedPath)
		return Project{}, nil
	}

	// Increment priority since the project is being accessed
	if err := ps.incrementProjectPriority(normalizedPath); err != nil {
		return Project{}, fmt.Errorf("failed to increment project priority: %w", err)
	}

	return project, nil
}

func (ps *ProjectState) incrementProjectPriority(path string) error {
	project, exists := ps.Projects[path]
	if !exists {
		return fmt.Errorf("project does not exist: %s", path)
	}

	project.Priority++
	ps.Projects[path] = project

	return ps.SaveState()
	// return ps.writeState(func() {
	// 	ps.Projects[path] = project
	// })
}

// TODO: Update get
// not very clean... abstract basic db ops and others (like hooks)
// Get should just return Project not increment
// (only increment on selection through list)
// This assumes project exists because
// func (ps *ProjectState) GetProject(path string) (Project, error) {
// 	var p Project
// 	wd, err := NormalizePath(path)
// 	if err != nil {
// 		return p, err
// 	}
// 	project, ok := ps.Projects[wd]
// 	if ok {
// 		// increment if exists
// 		err = ps.incrementProjectPriority(project)
// 		if err != nil {
// 			return p, err
// 		}
// 		return project, nil
// 	}
//
// 	p = Project{}
// 	// empty project
// 	return p, nil
// }
//
// func (ps *ProjectState) incrementProjectPriority(p Project) error {
// 	oldPriority := p.Priority
//
// 	err := ps.UpdateProject(p.Path, "Priority", strconv.Itoa(oldPriority+1))
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (ps *ProjectState) ProjectExists(path string) (bool, error) {
	wd, err := NormalizePath(path)
	if err != nil {
		return false, err
	}
	_, ok := ps.Projects[wd]
	if ok {
		return true, nil
	}

	return false, nil
}

// func (ps *ProjectState) sortProjectByPriority() []Project {
// 	projects := make([]Project, 0, len(ps.Projects))
//
// 	for _, v := range ps.Projects {
// 		// Blacklist filter
// 		if !ps.Blacklist[v.Path] {
// 			projects = append(projects, v)
// 		}
// 	}
//
// 	sort.SliceStable(projects, func(i, j int) bool {
// 		return ps.Projects[projects[i].Path].Priority > ps.Projects[projects[j].Path].Priority
// 	})
//
// 	return projects
// }
//

func (ps *ProjectState) ListProjects() []Project {
	projects := make([]Project, 0, len(ps.Projects))
	for _, project := range ps.Projects {
		if !ps.Blacklist[project.Path] {
			projects = append(projects, project)
		}
	}

	sort.SliceStable(projects, func(i, j int) bool {
		return projects[i].Priority > projects[j].Priority
	})

	return projects
}

//	func (ps *ProjectState) ListProject() []Project {
//		p := ps.sortProjectByPriority()
//		return p
//	}
func (ps *ProjectState) AddProject(projectPath string) error {
	normalizedPath, err := NormalizePath(projectPath)
	if err != nil {
		return err
	}

	if ps.Blacklist[normalizedPath] {
		return fmt.Errorf("path %s is blacklisted", normalizedPath)
	}

	if _, exists := ps.Projects[normalizedPath]; exists {
		return fmt.Errorf("project already exists: %s", normalizedPath)
	}

	ps.Projects[normalizedPath] = Project{
		Name: getBase(normalizedPath),
		Path: normalizedPath,
		Kind: "c",
	}

	return ps.SaveState()
}

// func (ps *ProjectState) AddProject(path string) error {
// 	wd, err := NormalizePath(path)
// 	if err != nil {
// 		return err
// 	}
//
// 	if ps.Blacklist[wd] {
// 		return fmt.Errorf("path %v in blacklist", wd)
// 	}
//
// 	if _, ok := ps.Projects[wd]; ok {
// 		return fmt.Errorf("project already exist %v", wd)
// 	}
//
// 	err = ps.writeState(func() {
// 		project := Project{Name: getBase(wd), Path: wd, Kind: "c"}
// 		ps.Projects[wd] = project
// 	})
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }

// func (ps *ProjectState) RemoveProject(path string) error {
// 	wd, err := NormalizePath(path)
// 	if err != nil {
// 		return err
// 	}
//
// 	err = ps.writeState(func() {
// 		delete(ps.Projects, wd)
// 	})
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }

func (ps *ProjectState) RemoveProject(projectPath string) error {
	normalizedPath, err := NormalizePath(projectPath)
	if err != nil {
		return err
	}

	delete(ps.Projects, normalizedPath)
	return ps.SaveState()
}

func (ps *ProjectState) ShowBlacklist() ([]string, error) {
	if ps.Blacklist == nil {
		return nil, nil
	}

	var blacklist []string
	for path, isBlacklisted := range ps.Blacklist {
		if isBlacklisted {
			blacklist = append(blacklist, path)
		}
	}
	return blacklist, nil
}

// func (ps *ProjectState) ShowBlacklist() ([]string, error) {
// 	var bl []string
// 	bs, err := ps.getBlacklist()
// 	if err != nil {
// 		return bl, err
// 	}
//
// 	for k, v := range bs {
// 		if v {
// 			bl = append(bl, k)
// 		}
// 	}
// 	return bl, nil
// }

func (ps *ProjectState) ManageBlacklist(path string, add bool) error {
	normalizedPath, err := NormalizePath(path)
	if err != nil {
		return err
	}

	if add {
		ps.Blacklist[normalizedPath] = true
	} else {
		delete(ps.Blacklist, normalizedPath)
	}

	return ps.SaveState()
}

// // when would I need this though?
//
//	func (ps *ProjectState) getBlacklist() (map[string]bool, error) {
//		return ps.Blacklist, nil
//	}
//
//	func (ps *ProjectState) RemoveBlacklist(path string) error {
//		wd, err := NormalizePath(path)
//		if err != nil {
//			return err
//		}
//
//		err = ps.writeState(func() {
//			delete(ps.Blacklist, wd)
//		})
//		if err != nil {
//			return err
//		}
//
//		return nil
//	}
//
//	func (ps *ProjectState) AddBlacklist(path string) error {
//		wd, err := NormalizePath(path)
//		if err != nil {
//			return err
//		}
//
//		err = ps.writeState(func() {
//			ps.Blacklist[wd] = true
//		})
//		if err != nil {
//			return err
//		}
//
//		return nil
//	}
func (ps *ProjectState) UpdateProject(projectPath, key, value string) error {
	normalizedPath, err := NormalizePath(projectPath)
	if err != nil {
		return err
	}

	project, exists := ps.Projects[normalizedPath]
	if !exists {
		return fmt.Errorf("project does not exist: %s", normalizedPath)
	}

	switch key {
	case "Path":
		project.Path = value
	case "Name":
		project.Name = value
	case "Kind":
		project.Kind = value
	case "Description":
		project.Description = value
	case "Priority":
		priority, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid priority value: %s", value)
		}
		project.Priority = priority
	default:
		return fmt.Errorf("unknown key: %s", key)
	}

	ps.Projects[normalizedPath] = project
	return ps.SaveState()
}

// func (ps *ProjectState) UpdateProject(path, key, value string) error {
// 	wd, err := NormalizePath(path)
// 	if err != nil {
// 		return err
// 	}
//
// 	// TODO: utilize methods
// 	// TODO: make this a pointer
// 	err = ps.writeState(func() {
// 		p := ps.Projects[wd]
// 		switch key {
// 		case "Path":
// 			p.Path = value
// 		case "Name":
// 			p.Name = value
// 		case "Kind":
// 			p.Kind = value
// 		case "Description":
// 			p.Description = value
// 		case "Priority":
// 			priority, err := strconv.Atoi(value)
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			p.Priority = priority
// 		default:
// 			log.Fatalf("No such key %v in project", key)
// 		}
// 		ps.Projects[wd] = p
// 	})
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }
