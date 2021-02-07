package main

import (
	"fmt"
	"path/filepath"
    "os"
	"syscall"
	"strings"
    "os/user" 
    "bufio"   
    "io/ioutil"
    "golang.org/x/sys/windows"
    "golang.org/x/sys/windows/registry"

)

func main() {
	if !amAdmin() {
        runMeElevated()
        return
	}

	usr, err := user.Current()
    if err != nil {
		fmt.Println("Error trying to determine homefolder.")
		return
    }
	
	appDir := filepath.Join(usr.HomeDir, "AppData", "Local", "pc-limit")
    
    fmt.Printf("This application installs pc-limit to: %s.\nContinue? [Y/N]", appDir)

	reader := bufio.NewReader(os.Stdin)
	c, err := reader.ReadByte()
	if err != nil {
		fmt.Println(err)
		return
	}

	if c != []byte("Y")[0] && c != []byte("y")[0] {
        fmt.Println("Install cancelled.")
        waitForKey()
        return
    }

    fmt.Println("Creating folder.")
    err = os.Mkdir(appDir, 0755)
    if err != nil {
        fmt.Println(err)
        waitForKey()
		return
    }

    fmt.Println("Copy files.")
    files := []string{"pclimit.exe", "uuid.txt", "license.txt"}
    for _, file := range files {
        err = copy(file, filepath.Join(appDir, file))
        if err != nil {
            fmt.Println(err)
            waitForKey()
            return
        }
    }

    fmt.Println("Configure registry.")
    k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE)
    if err != nil {
        fmt.Println(err)
        waitForKey()
        return
    }
    k.SetStringValue("PCLimit", filepath.Join(appDir, "pclimit.exe"))
    err = k.Close()
    if err != nil {
        fmt.Println(err)
        waitForKey()
        return
    }

    

    fmt.Println("Installed.")    
    waitForKey()
}

func runMeElevated() {
    verb := "runas"
    exe, _ := os.Executable()
    cwd, _ := os.Getwd()
    args := strings.Join(os.Args[1:], " ")

    verbPtr, _ := syscall.UTF16PtrFromString(verb)
    exePtr, _ := syscall.UTF16PtrFromString(exe)
    cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
    argPtr, _ := syscall.UTF16PtrFromString(args)

    var showCmd int32 = 1 //SW_NORMAL

    err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
    if err != nil {
        fmt.Println(err)
    }
}

func amAdmin() bool {
    _, err := os.Open("\\\\.\\PHYSICALDRIVE0")
    return (err == nil)
}

func copy(sourceFile, destinationFile string) error {
    input, err := ioutil.ReadFile(sourceFile)
    if err != nil {
        return err
    }

    err = ioutil.WriteFile(destinationFile, input, 0755)
    if err != nil {
        return err
    }
    return nil
}

func waitForKey() {
    reader := bufio.NewReader(os.Stdin)
    _, err := reader.ReadByte()
	if err != nil {
		fmt.Println(err)
		return
    }
}