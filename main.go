package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file

	// 1.parse input
	// 2. retrieve file
	// 3. write temporary file on our server
	// 4. return result
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	fmt.Fprintf(w, "Step1 : Successfully Uploaded File\n")

	f, err := os.OpenFile("script/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("It is failing.")
		fmt.Println(err)
	}
	defer f.Close()
	io.Copy(f, file)

	// 3. write temporary file on our server
	// tempFile, err := ioutil.TempFile("script", "runtest*.py")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// defer tempFile.Close()
	// fileBytes, err := ioutil.ReadAll(file)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// tempFile.Write(fileBytes)

	pythonexe := "C:/Program Files (x86)/Microsoft Visual Studio/Shared/Python37_64/python.exe"
	cmd := exec.Command(pythonexe, "script/"+handler.Filename, "--input-file", "script/*.py")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("It is failing to run script.")
		fmt.Println(err)
		panic(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}
	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	go copyOutput(stdout)
	go copyOutput(stderr)
	cmd.Wait()

	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully excuted File\n")
}

func copyOutput(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func setupRoutes() {
	http.HandleFunc("/upload", uploadFile)
	http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("Hello World!")
	setupRoutes()
}
