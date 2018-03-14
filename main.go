package main

import (
	"log"
	"net/http"
)

type Cinfo struct {
	index int //order
	parts int //length
}

//8MB maxium chunk size
const MAX_CHUNK_SIZE = 8 * 1024 * 1024

func upload(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		return
	}
	file, err = handler.Open()
	if err != nil {
		log.Println(err)
		return
	}
	filesize := handler.Size
	log.Println(filesize)
	var buffersize int64
	var parts int
	if filesize < MAX_CHUNK_SIZE {
		buffersize = filesize
		parts = 1
	} else {
		buffersize = MAX_CHUNK_SIZE
		parts = roundup(buffersize, filesize)
	}
	buffer := make([]byte, buffersize)
	//key: int; part number, e.g. 1 means 1st part
	//value: []byte; a fix-sized chunk of data read from file bytestream
	fileChunks := make(map[Cinfo][]byte)
	var ci Cinfo
	ci.parts = parts
	log.Println(parts)
	for i := 1; i <= parts; i++ {
		br, err := file.Read(buffer)
		if err != nil {
			log.Println(err)
			return
		}
		ci.index = i
		fileChunks[ci] = buffer[:br]
	}
	//log.Println("bytes read:", br)
	//log.Println("bytestream:", buffer[:br])
	log.Println(len(fileChunks))
}

func roundup(a int64, b int64) int {
	if b%a > 0 {
		return int(b/a) + 1
	}
	return int(b / a)
}

func main() {
	http.HandleFunc("/upload", upload)
	log.Fatal(http.ListenAndServe(":9090", nil))
}
