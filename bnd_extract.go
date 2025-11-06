package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func makeInt(data []byte) int {
	first_byte := int(data[0])
	second_byte := int(data[1]) << 8
	third_byte := int(data[2]) << 16
	fourth_byte := int(data[3]) << 24
	return first_byte + second_byte + third_byte + fourth_byte
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: bnd_extract file.bnd path")
		return
	}
	
	file, err := os.Open(os.Args[1])
	
	if err != nil {
		fmt.Println("Error while loading file")
		return
	}
	defer file.Close()
	
	word := make([]byte, 4)
	
	// Read file id
	_, err = file.Read(word)
	
	if err != nil {
		fmt.Println("Error when reading file id")
		return
	}
	
	id := makeInt(word)
	
	if id != 0x444E42	{
		fmt.Println("Hey, this isn't a BND file")
		return
	}
	
	// Read file version
	_, err = file.Read(word)
	
	if err != nil {
		fmt.Println("Error when reading file id")
		return
	}
	
	version := makeInt(word)
	
	fmt.Println("Version:", version)
	
	// Read flags
	_, err = file.Read(word)
	
	if err != nil {
		fmt.Println("Error when reading file id")
		return
	}
	
	flags := makeInt(word)
	
	fmt.Println("Flags:", flags)
	
	// Number of files
	_, err = file.Read(word)
	
	if err != nil {
		fmt.Println("Error when reading file id")
		return
	}
	
	num_of_files := makeInt(word)
	
	fmt.Println("Number of files:", num_of_files)
	
	// Read files
	var index, offset, size, name int
	
	var str_buff []byte
	read_byte := make([]byte, 1)
	
	var filename string
	
	var filedata []byte
	
	pos := 16
	for i := 0; i < num_of_files; i++ {
		// Read row
		_, err = file.Read(word)
		index = makeInt(word)
		
		_, err = file.Read(word)
		offset = makeInt(word)
		
		_, err = file.Read(word)
		size = makeInt(word)
		
		_, err = file.Read(word)
		name = makeInt(word)
		
		// Read filename
		file.Seek(int64(name), 0)
		
		for {
			_, err = file.Read(read_byte)
			
			if read_byte[0] == 0 {
				break;
			}
			
			str_buff = append(str_buff, read_byte[0])
			
		}
		
		filename = string(str_buff)
		str_buff = []byte{}
		
		// Read file
		file.Seek(int64(offset), 0)
		filedata = make([]byte, size)
		
		_, err = file.Read(filedata)
		
		// TODO: TRIM / at the end of FILE PATH
		// TODO: also add path separator
		os.MkdirAll(filepath.Join(os.Args[2]), 0644)
		err = os.WriteFile(filepath.Join(os.Args[2], filename), filedata, 0644)
		
		if err != nil {
			log.Fatal(err)
			panic("Error when writing file")
		}
		
		// Print obtained info
		fmt.Println("Pos ", pos, ":", index, offset, size, name, filename)
		
		pos += 16
		file.Seek(int64(pos), 0)
	}
}