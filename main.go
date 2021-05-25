package main

import (
	"os"
	"log"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Info struct {
	Path		string	`json:"file"`
	Def_Offsets	[]int	`json:"default_offsets"`
	Inc_By		int	`json:"increments_by"`
	Offset		bool	`json:"offset"`
}

type FileInfo struct {
	info		*Info
	arr		[]byte
	s_arr		[16][4]byte
	convert		string
	file_size	int
}

func read_from_file(file string) *Info {
	f, err := os.Open(file)
	defer f.Close()

	if err != nil {
		log.Fatal(err)
	}

	i := Info{}
	dec := json.NewDecoder(f)
	err = dec.Decode(&i)

	if err != nil {
		log.Fatal(err)
	}

	return &i
}

func setup(info *Info) *FileInfo {
	FI := &FileInfo{info: info}

	return FI
}

func (info *FileInfo) read_file() {
	file, err := os.Stat(info.info.Path)

	if err != nil {
		log.Fatal(err)
	}

	info.file_size = int(file.Size())

	info.arr, err = ioutil.ReadFile(info.info.Path)

	if err != nil {
		log.Fatal(err)
	}
}

func (info *FileInfo) derive() {
	c_r := 0
	c_c := 0

	for i := 0; i < len(info.arr); i++ {
		if i % 4 == 0 {
			c_r++
			if c_r == 16 {
				break
			}
			c_c = 0
		}

		info.s_arr[c_r][c_c] = info.arr[i]
		c_c++
	}

	// Initialize the first 4 bytes to 255. For some reason they're zero..
	for i := 0; i < 1; i++ {
		for x := 0; x < 4; x++ {
			info.s_arr[i][x] = 255
		}
	}
}

func main() {
	i := read_from_file("info.json")

	info := setup(i)
	info.read_file()
	info.derive()
	fmt.Println(info)
}
