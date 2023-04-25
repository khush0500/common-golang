package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
)

const VERSION = "1.0.0"

const greenColor = "\x1b[32m => "
const redColor = "\x1b[91m ** "
const defaultStyle = "\x1b[0m"
const lightGrayColor = "\x1b[37m"

const bootstrapDataSource = "https://kukkar-cdn.s3.ap-south-1.amazonaws.com/newapp"

var ApplicationPath string
var ApplicationName string

var Rmap map[string]string

type Command struct {
	Name                      string
	AltName                   string
	FlagSet                   *flag.FlagSet
	Usage                     []string
	UsageCommand              string
	Command                   func(args []string, additionalArgs []string)
	SuppressFlagDocumentation bool
	FlagDocSubstitute         []string
}

func generateBootstrap(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("You need to supply application path.")
	}
	ApplicationPath = args[0]
	ApplicationName = prepareApplicationName(ApplicationPath)

	Rmap = initiateReplacementMap()
	return createDirStructure()
}

func prepareApplicationName(path string) string {
	pathArr := strings.Split(path, string("/"))
	fmt.Printf("Path:%s\n", path)
	return pathArr[len(pathArr)-1]
}

type Path string

type DS struct {
	root *Folder
}

func (this *DS) Fetch() error {

	this.root.Remote = bootstrapDataSource
	wd, _ := os.Getwd()
	this.root.Local = Path(wd)
	return this.root.BootUp(true)
}

func (this *DS) Show() {
	var show func(Folder, string)
	show = func(folder Folder, spaces string) {
		//print files
		if len(folder.Files) > 0 {
			for _, file := range folder.Files {
				file.PrintIt(spaces)
			}
		}
		if len(folder.Folders) > 0 {
			for _, folder := range folder.Folders {
				folder.PrintIt(spaces)
				show(folder, spaces+"    ")
			}
		}

	}
	show(*this.root, "    ")
}

type Folder struct {
	Name       string
	ActualName string
	Files      []File
	Folders    []Folder
	Local      Path
	Remote     Path
	Spaces     int
}

func (this *Folder) BootUp(skip bool) error {

	err := this.create(skip)
	if err != nil {
		return this.error(err)
	}
	err = this.createFiles()
	if err != nil {
		return this.error(err)
	}
	if this.Folders == nil {

		err := os.Chdir("../")
		if err != nil {
			return this.error(err)
		}
		return nil
	}
	for _, folder := range this.Folders {
		folder.Local = this.Local + Path(os.PathSeparator) + Path(folder.ActualName)
		folder.Remote = this.Remote + "/" + Path(folder.ActualName)
		err := folder.BootUp(false)
		if err != nil {
			return err
		}

	}
	err1 := os.Chdir("../")
	if err1 != nil {
		return this.error(err1)
	}

	return nil
}

func (this *Folder) create(skip bool) error {

	if skip {
		return nil
	}
	//wd, _ := os.Getwd()
	//fmt.Printf("making %s in %s\n", this.ActualName, wd)
	err := os.Mkdir(this.Name, 0755)
	if err != nil {
		return this.error(err)
	}
	err = os.Chdir(this.Name)
	if err != nil {
		return this.error(err)
	}

	return nil
}

func (this *Folder) error(err error) error {
	return err
}

func (this *Folder) createFiles() error {
	if len(this.Files) <= 0 {
		return nil
	}
	for _, file := range this.Files {
		file.Local = this.Local + Path(os.PathSeparator) + Path(file.GetFileName())
		file.Remote = this.Remote + "/" + Path(file.GetFileName())
		err := file.Create()
		if err != nil {
			return this.error(err)
		}
	}
	return nil
}

func (this *Folder) PrintIt(spaces string) {
	fmt.Println(spaces + this.ActualName + "/")
}

type File struct {
	Name       string
	ActualName string
	Extension  string
	Local      Path
	Remote     Path
}

func (this *File) PrintIt(spaces string) {
	fmt.Println(spaces + this.GetFileName())
}

func (this *File) error(err error) error {
	return fmt.Errorf("File:%s, Location:%s, Error:%s\n", this.ActualName, this.Local, err.Error())
}

func (this *File) Fetch() ([]byte, error) {
	var maxretry int = 3
	var retry int
	var err error
	var resp *http.Response
	var returnData []byte
	for retry < maxretry {
		retry++
		resp, err = http.Get(string(this.Remote))
		if err != nil {
			continue
		}
		if resp.StatusCode != 200 {
			resp.Body.Close()
			err = fmt.Errorf("Got HTTP Status code: %d", resp.StatusCode)
			continue
		}
		b := make([]byte, 0)
		buf := bytes.NewBuffer(b)
		_, err := buf.ReadFrom(resp.Body)
		if err != nil {
			resp.Body.Close()
			continue
		}
		resp.Body.Close()
		returnData = buf.Bytes()
	}
	return returnData, err
}

func (this *File) GetFileName() string {
	filename := this.ActualName
	if this.Extension != "" {
		filename = this.ActualName + "." + this.Extension
	}
	return filename
}

func (this *File) Create() error {

	destFile, err := os.Create(this.GetFileName())
	if err != nil {
		return this.error(err)
	}
	defer destFile.Close()

	//fmt.Printf("Hitting: %s\n", this.Remote)
	data, err := this.Fetch()
	if err != nil {
		return this.error(err)
	}
	//{{APP_PATH}}
	for k, v := range Rmap {
		//	fmt.Println(k)
		data = []byte(strings.Replace(string(data), k, v, -1))
	}
	_, err = destFile.Write(data)
	if err != nil {
		return this.error(err)
	}
	destFile.Sync()
	return nil
}

func initiateReplacementMap() map[string]string {

	s := make(map[string]string)

	s["{{APP_PATH}}"] = ApplicationPath
	s["{{APP_NAME}}"] = ApplicationName

	//prepare log path
	var logpath string
	if runtime.GOOS == "windows" {
		//hack for windows
		logpath = "C:\\\\" + ApplicationName + "\\\\"
	} else {
		logpath = "/var/log/" + ApplicationName + "/"
	}

	s["{{LOG_PATH}}"] = logpath

	os.MkdirAll(logpath, 0777)
	return s
}

func createDirStructure() error {
	//Define directory structure
	fmt.Printf("Creating Directory Structure:\n")
	root := Folder{
		Name:       "root",
		ActualName: "root",
		Files: []File{
			File{Name: "main", ActualName: "main", Extension: "go"},
			File{Name: "go", ActualName: "go", Extension: "mod"},
			File{Name: "dockerfile", ActualName: "Dockerfile"},
			File{Name: "jenkinsfile", ActualName: "Jenkinsfile"},
		},
		Folders: []Folder{
			Folder{
				Name:       "conf",
				ActualName: "conf",
				Files: []File{
					File{Name: "config", ActualName: "config", Extension: "json"},
					File{Name: "conf", ActualName: "config", Extension: "go"},
				},
			},
			Folder{
				Name:       "src",
				ActualName: "src",
				Folders: []Folder{
					Folder{
						Name:       "common",
						ActualName: "common",
						// Folders: []Folder{
						// 	Folder{
						// 		Name:       "factory",
						// 		ActualName: "factory",
						// 		Folders: []Folder{
						// 			Folder{
						// 				Name:       "sql",
						// 				ActualName: "sql",
						// 				Files: []File{
						// 					File{Name: "healthchecker", ActualName: "healthchecker", Extension: "go"},
						// 					File{Name: "mysql", ActualName: "mysql", Extension: "go"},
						// 				},
						// 			},
						// 			Folder{
						// 				Name:       "cache",
						// 				ActualName: "cache",
						// 				Files: []File{
						// 					File{Name: "healthchecker", ActualName: "healthchecker", Extension: "go"},
						// 					File{Name: "cache", ActualName: "cache", Extension: "go"},
						// 				},
						// 			},
						// 			Folder{
						// 				Name:       "ravenf",
						// 				ActualName: "ravenf",
						// 				Files: []File{
						// 					File{Name: "ravenf", ActualName: "ravenf", Extension: "go"},
						// 				},
						// 			},
						// 		},
						// 	},
						// 	//Folder{Name: "appconstant", ActualName: "appconstant", Files: []File{File{Name: "errcodes", ActualName: "error_codes", Extension: "go"}}},
						// },
						// Files: []File{
						// 	File{Name: "starter", ActualName: "starter", Extension: "go"},
						// },
					},
					Folder{
						Name:       ApplicationName + "_controllers",
						ActualName: "controllers",
						Files: []File{
							File{Name: "const", ActualName: "const", Extension: "go"},
							File{Name: "struct", ActualName: "struct", Extension: "go"},
							File{Name: "hello_world", ActualName: "hello", Extension: "go"},
							File{Name: "utils", ActualName: "utils", Extension: "go"},
						},
					},
					Folder{
						Name:       ApplicationName + "_routes",
						ActualName: "routes",
						Files: []File{
							File{Name: "routes", ActualName: "routes", Extension: "go"},
						},
					},
					Folder{
						Name:       ApplicationName,
						ActualName: "application",
						Files: []File{
							File{Name: "const", ActualName: "const", Extension: "go"},
							File{Name: "error", ActualName: "error", Extension: "go"},
							// File{Name: "appname", ActualName: ApplicationName, Extension: "go"},
							// File{Name: ApplicationName + "impl", ActualName: "appnameimpl", Extension: "go"},
							// File{Name: "interface", ActualName: "interface", Extension: "go"},
							// File{Name: "struct", ActualName: "struct", Extension: "go"},
							// File{Name: "utils", ActualName: "utils", Extension: "go"},
						},
					},
				},
			},
		},
	}
	//Fetch and put on system
	ds := DS{root: &root}
	err := ds.Fetch()
	if err == nil {
		ds.Show()
	}
	return err
}
