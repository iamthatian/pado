// this should have everything related to projects that are abstract
// add blacklist in config

// package main
//
// import (
//
//	"compress/gzip"
//	"encoding/gob"
//	"fmt"
//	// "log"
//	"os"
//
// )
//
// // This gets serialized
//
//	type Projects interface {
//		List() []ProjectData
//		Remove(ProjectData)
//		Add(ProjectData)
//		// Blacklist(ProjectData)
//	}
//
//	type ProjectsData struct {
//		projects []ProjectData
//	}
//
//	func (p *ProjectsData) List() []ProjectData {
//		return p.projects
//	}
//
//	func listProjects() {
//		fmt.Println("aweome")
//		// project := Person{Name: "John Doe", Age: 30}
//	}
//
// // sanitize directory/path
// // Need information on current directory and/or given directory
// // either empty string or string
// // kv store search based on project name or path? duplicate data?
// // or more complex? single data but then filtering based on path/lang/name
// // do not allow duplicate path/lang/name
// // migration support with version update and
// // Ok KV store (based off path) O(1) and other queries do O(n)
// // search (because most common usecase would be getting project state)
// // and n is usually going to be small
// // Other option is O(log n) for every operation and increased complexity
// // fzf query based on stuff
// // input custom db filler arguments
//
//	func addProject() error {
//		// If the file doesn't exist, create it, or append to the file
//		// f, err := os.OpenFile("access.log", os.O_CREATE|os.O_RDWR, 0644)
//		// if err != nil {
//		// 	log.Fatal(err)
//		// }
//		// if _, err := f.Write([]byte("appended some data\n")); err != nil {
//		// 	f.Close() // ignore error; Write error takes precedence
//		// 	log.Fatal(err)
//		// }
//		// if err := f.Close(); err != nil {
//		// 	log.Fatal(err)
//		// }
//
//		// TODO: Only create iff doesn't exist
//		// fi, err := os.Create("./pdstate")
//		fi, err := os.OpenFile("./pdstate", os.O_CREATE|os.O_RDWR, 0644)
//		if err != nil {
//			return err
//		}
//		defer fi.Close()
//
//		// pass in as argument
//		project := ProjectData{Name: "test", Path: "test/path", Kind: "c"}
//
//		// gzip instead of bytes?
//		// var b bytes.Buffer
//		fz := gzip.NewWriter(fi)
//		defer fz.Close()
//
//		// Create a new gob encoder and use it to encode the person struct
//		// enc := gob.NewEncoder(&b)
//		enc := gob.NewEncoder(fz)
//		// enc := gob.NewEncoder(&b)
//		// if err := enc.Encode(project); err != nil {
//		if err := enc.Encode(project); err != nil {
//			return err
//			// fmt.Println("Error encoding struct: err")
//		}
//		// 	fmt.Println("Error encoding struct:", err)
//		// 	return
//		// }
//		return nil
//
//		// The serialized data can now be found in the buffer
//		// serializedData := b.Bytes()
//		// fmt.Println("Serialized data:", string(serializedData))
//		// fmt.Println("Serialized data:", b.String())
//	}
//
//	func readProjectsData() (ProjectData, error) {
//		var p ProjectData
//		// deserialize
//		fi, err := os.Open("./pdstate")
//		if err != nil {
//			return p, err
//		}
//		defer fi.Close()
//
//		fz, err := gzip.NewReader(fi)
//		if err != nil {
//			return p, err
//		}
//		defer fz.Close()
//
//		decoder := gob.NewDecoder(fz)
//		err = decoder.Decode(&p)
//		if err != nil {
//			return p, err
//		}
//
//		return p, nil
//	}
//
// // // TODO: Add custom path
// // func writeProjectsData(p Project) {
// // 	// serialize
// // 	serialized :=
// // 	// write to file
// // }
//
// // func (this *ProjectsData) Remove(p Project) []Project {
// // 	// return ProjectsData.filter(e => e.path != p.GetPath())
// // }
package main
