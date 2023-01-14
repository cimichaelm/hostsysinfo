//
// File: utilfunc.go
//
//

package main

import (
    "bufio"
    "os"
    "fmt"
    "sort"
    "syscall"
	"io/ioutil"
	"log"
	"gopkg.in/yaml.v3"
	"runtime"
	"os/user"
	"strconv"
)

//
// create directory if it does not exist
//
func create_directory(dirpath string, mode os.FileMode) bool {

     r := directory_exists(dirpath)
     if r {
     	return r
     }

     err := os.Mkdir(dirpath, mode)

     if err == nil {
      	r = true
     }

     return r
}

//
// check if directory exists
//
func directory_exists(dirpath string) bool {

    _, err := os.Stat(dirpath)
    if os.IsNotExist(err) {
       return false
    }
    return true
}

//
// check if file exists
//
func file_exists(path string) bool {

    _, err := os.Stat(path)
    if os.IsNotExist(err) {
       return false
    }
    return true
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func write_output_to_file(w *bufio.Writer, msg string) {
    _, err := w.WriteString(msg)
    check(err)
    w.Flush()
}

func show_map(m map[string]string) {
     for key := range m {
     	 value := m[key]
     	 fmt.Printf("%s: %s\n", key, value)
     }
}

//
// write map to file
//
func write_map(m map[string]string, filename string) {
    f, err := os.Create(filename)
    check(err)

    defer f.Close()
    w := bufio.NewWriter(f)
    keylist := map_keys(m)
    for _, key := range keylist {
     	 value := m[key]
	 msg := fmt.Sprintf("%s: %s\n", key, value)
	 write_output_to_file(w, msg)
    }
}

//
// clear all keys from map
//
func clear_map(m map[string]string) {
     for k := range m {
         delete(m, k)
     }
}

//
// add key/value to map
//
func add_map(m map[string]string, key string, value string) {
     m[key] = value
}

func map_keys(m map[string]string) []string {
     keylist := make([]string,0)
     
     for k := range m {
        keylist = append(keylist, k)
     }
     sort.Strings(keylist)
     return keylist
}

//
// check if we are running as root.
//
func is_root() bool {
	uid := syscall.Getuid()
	if uid == 0 {
		return true
	}
	return false
}

//
// load properties from yaml file
// format:
// key: value
//
func load_props(filename string) param_t {

	yfile, err := ioutil.ReadFile(filename)

	if err != nil {

		log.Fatal(err)
	}

	data := make(map[string]string)
//	data := &param_t{}
	err2 := yaml.Unmarshal(yfile, &data)

	if err2 != nil {

		log.Fatal(err2)
	}
	return data
}

func display_message(msg string) {
	fmt.Printf("%s\n", msg)
}

//
func set_file_permissions(filename string, mode os.FileMode) {

	if file_exists(filename) {

		// Change permissions Linux.
		err := os.Chmod(filename, mode)
		if err != nil {
			log.Println(err)
		}
	}
}

//
// set group of file	 
//
func set_file_group(filename string, group string) {

	var gid int
	
	if file_exists(filename) {
		gid = lookup_gid(group)
		if gid >= 0 {
			// Change file ownership.
			err := os.Chown(filename, -1, gid)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func lookup_gid(group string) int {

	if runtime.GOOS != "windows" {
	    ginfo, err := user.LookupGroup(group)
	    if err != nil {
	        return -1
	    }
	    gid, _ := strconv.Atoi(ginfo.Gid)
		return gid
	}
	return -1
}
	
func lookup_uid(username string) int {
	
	if runtime.GOOS != "windows" {
	    uinfo, err := user.Lookup(username)
	    if err != nil {
	        return -1
	    }
	    uid, _ := strconv.Atoi(uinfo.Uid)
		return uid
	}
	return -1
}
