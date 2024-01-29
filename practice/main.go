package main

import (
	"fmt"
	"io"
	"os"
)

type Book struct {
	Id     int    `json:"id"`
	Title  string `json:"title"`
	IsDone bool   `json:"done"`
}

func main() {

	srcFile, err := os.Open("abc.json")
	if err != nil {
		fmt.Println(err)
	}

	defer srcFile.Close()

	bufferReader := make([]byte, 128)

	for {
		data, err := srcFile.Read(bufferReader)
		if err != nil && err != io.EOF {
			fmt.Println(err)
		}

		fmt.Println(string(bufferReader[:data]))

		if data == 0 {
			fmt.Println("break this shit")
			break
		}

	}

}
