package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

const NETWORK_TYPE = "tcp"
const IP_ADDRESS = "127.0.0.1"
const PORT = "8080"

func find_requested_file(route string, dir string) string {
	/* Recursively find the requested file */

	file_path := "static/404.html"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		file_name := dir + "/" + file.Name()
		if route == file_name {
			return file_name
		}
		if file.IsDir() {
			file_path = find_requested_file(route, file_name)
		}
	}

	return file_path
}

func main() {
	address := IP_ADDRESS + ":" + PORT

	// Listen for connections
	listener, err := net.Listen(NETWORK_TYPE, address)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Server Running at http://" + address)

	defer listener.Close()

	for {
		// Accept for connections
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println(err)
		}

		go process(conn)
	}
}

func process(conn net.Conn) {
	var buf []byte = make([]byte, 1024)

	// Read buffer of the request
	_, err := conn.Read(buf)

	if err != nil {
		fmt.Println(err)
	}

	req_headers := string(buf)

	req_headers_lines := strings.Split(req_headers, "\n")

	route_line := req_headers_lines[0]

	fmt.Println(route_line)

	route_items := strings.Split(route_line, " ")

	// Get and update the requested route
	route := "static" + route_items[1]

	html_ext := route[len(route)-5:]
	ico_ext := route[len(route)-4:]

	if html_ext != ".html" && ico_ext != ".ico" && route[len(route)-1] != '/' {
		route += "/"
	}

	if route[len(route)-1] == '/' {
		route += "index.html"
	}

	// find the requested file
	file_path := find_requested_file(route, "static")
	
	buf, err = os.ReadFile(file_path)
	
	if err != nil {
		fmt.Println(err)
	}

	file_content := string(buf)

	file_buf := []byte("HTTP/1.1 200 OK\n\n" + file_content)

	// Send the buffer of the file content
	_, err = conn.Write(file_buf)

	if err != nil {
		fmt.Println(err)
	}

	conn.Close()
}

