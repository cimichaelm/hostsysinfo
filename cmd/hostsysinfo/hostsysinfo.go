//
// File: sysinfo.go
//
//

package main

import (
	"bufio"
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type probeinfo_t struct {
	Command   string
	Arguments string
	Become	  bool
	Enable    bool
	Outfile   string
	Pattern   string
	Prefix	  string
}

type config_t struct {
	Version   string
	Outputdir string
	User	  string
	Group	  string
	Lowercase bool
	Probes    map[string]probeinfo_t
}

type objtype struct {
	configfile string
	debug      bool
	dryrun     bool
	verbose    bool
}

type param_t map[string]string

var (
	configfile  = flag.String("config", "site-config.yaml", "Config File Path")
	flg_dryrun  = flag.Bool("dryrun", false, "Dry Run")
	flg_debug   = flag.Bool("debug", false, "Debug")
	flg_verbose = flag.Bool("verbose", false, "Verbose")
	flg_generate = flag.Bool("generate", false, "Generate data files")
	opt_get = flag.String("get", "", "Get data value")
	flg_showkey = flag.Bool("showkey", false, "Show Key")
	flg_lowercase = flag.Bool("lowercase", false, "Use lowercase")

	params = make(map[string]string)
)

func defaults() *objtype {

	obj := &objtype{}

	return obj
}

func argparse(obj *objtype) {

	var msg string

	flag.Parse()

	obj.configfile = *configfile
	obj.dryrun = *flg_dryrun
	obj.debug = *flg_debug
	obj.verbose = *flg_verbose

	msg = fmt.Sprintf("config file name: %s", *configfile)
	display_debug_message(obj, msg)
	msg = fmt.Sprintf("dryrun: %s", *flg_dryrun)
	display_debug_message(obj, msg)
}

func loadconfig(obj *objtype, configfile string) *config_t {

	yfile, err := ioutil.ReadFile(configfile)

	if err != nil {

		log.Fatal(err)
	}

	//     data := make(map[interface{}]interface{})
	data := &config_t{}
	err2 := yaml.Unmarshal(yfile, &data)

	if err2 != nil {

		log.Fatal(err2)
	}
	return data
}

func process(obj *objtype, data config_t) {

	var msg string
	
	flag_isroot := is_root()
	if flag_isroot {
		msg = fmt.Sprintf("Running as root")
		display_debug_message(obj, msg)
	}

	if *flg_generate {
		process_generate(obj, data)
	}
	
	if len(*opt_get) > 0 {
		process_get(obj, data, *opt_get)
	}
}


func process_generate(obj *objtype, data config_t) {

	var msg, prog, arguments, pattern, prefix, group string
	var mode os.FileMode
	var r int
	var flg_setgroup bool
	
	flg_setgroup = false	
	pattern = ""
	prog = ""
	arguments = ""
	mode = 0750
	group = data.Group
	if len(group) > 0 {
		flg_setgroup = true
	}
	msg = fmt.Sprintf("Version: %s", data.Version)
	display_debug_message(obj, msg)
	msg = fmt.Sprintf("Outputdir: %s", data.Outputdir)
	display_debug_message(obj, msg)
	create_directory(data.Outputdir, mode)
	if flg_setgroup {
		set_file_group(data.Outputdir, group)
	}
	
	clear_params()
	
	if data.Lowercase {
		*flg_lowercase = true
	}
	
	for k, v := range data.Probes {
		if v.Enable {
			pattern = v.Pattern
			prog = v.Command
			prefix = v.Prefix
			arguments = v.Arguments
			msg = fmt.Sprintf("%s -> %s", k, v)
			display_debug_message(obj, msg)
			outfilepath := data.Outputdir + "/" + v.Outfile
			msg = fmt.Sprintf("Executing: %s %s", prog, arguments)
			display_debug_message(obj, msg)
			if v.Become {
				r = sudocmd(prog, arguments, outfilepath, pattern, prefix)
			} else {
				r = runcmd("", prog, arguments, outfilepath, pattern, prefix)				
			}
			msg := fmt.Sprintf("Return code: %s", r)
			display_debug_message(obj, msg)
		}
	}
	paramfile := get_paramfile(data)
	write_map(params, paramfile)
	if flg_setgroup {
		set_file_group(paramfile, group)

	}
}

func get_outputdir(data config_t) string {
	return data.Outputdir
}

func get_paramfile(data config_t) string {
	paramfile := get_outputdir(data) + "/" + "params.yaml"

	return paramfile
}

func process_get(obj *objtype, data config_t, key string) {
	var msg, value string
	
	paramfile := get_paramfile(data)
	params := load_props(paramfile)
	if *flg_lowercase {
		key = strings.ToLower(key)
	}
	value = params[key]
	if *flg_showkey {
		msg = fmt.Sprintf("%s: %s", key, value)
	} else {
		msg = fmt.Sprintf("%s", value)		
	}
	display_message(msg)
}

func display_debug_message(obj *objtype, msg string) {
	if obj.debug {
		debugprefix := "DEBUG: "
		display_message(debugprefix + msg)
	}
}

func clear_params() {
	clear_map(params)
}

func add_param(key string, value string) {
	key2 := strings.ReplaceAll(key, " ", "_")
	params[key2] = value
}

func show_params() {
	show_map(params)
}

func runcmd(wrapper string, prog string, arguments string, filename string, pattern string, prefix string) int {

	var cmd *exec.Cmd
	var arg1, arglist []string
	
	arg1 = strings.Fields(arguments)

	r := 0
	maxlines := 10000
	
	if len(wrapper) > 0 {
		arglist = append(arglist, prog)
		arglist = append(arglist, arg1...)
		cmd = exec.Command(wrapper, arglist...)
	} else {
		arglist = arg1
		cmd = exec.Command(prog, arglist...)		
	}
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		log.Fatal(err)
		r := 1
		return r
	}

	cmd.Start()

	buf := bufio.NewReader(stdout)
	num := 0

	f, err := os.Create(filename)
	check(err)

	defer f.Close()
	w := bufio.NewWriter(f)

	for line, _, err := buf.ReadLine(); err != io.EOF; line, _, err = buf.ReadLine() {
		if num > maxlines {
			r := 2
			return r
		}
		num += 1
		msg := fmt.Sprintf("%s\n", string(line))
		write_output_to_file(w, msg)
		match_param(prefix, pattern, string(line))
	}

	return r
	
}

func sudocmd(prog string, arguments string, filename string, pattern string, prefix string) int {

	wrapper := "sudo"
	return runcmd(wrapper, prog, arguments, filename, pattern, prefix)
}

// match pattern in string
func match_param(prefix string, pattern string, line string) {
	var key, value string
	
	sep := "."
	
	if len(pattern) > 0 {
		re := regexp.MustCompile(pattern)
		match := re.FindStringSubmatch(line)
		if match != nil {
			if len(match) > 2 {
				key = strings.ToLower(prefix + sep + match[1])
				value = match[2]
				// fmt.Printf("%s: %s\n", key, value)
				add_param(key, value)
			}
		}
	}
}
