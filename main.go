package main

import (
	"log"
	"net/http"
)

//Need to be set to 8*1024*1024 = 8mb
const MAX_CHUNK_SIZE = 2

func upload(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	var buffersize int64
	var parts int
	//filesize := handler.Size
	log.Println("FileSize:", handler.Size)
	if filesize := handler.Size; filesize < MAX_CHUNK_SIZE {
		buffersize = filesize
		parts = 1
	} else {
		buffersize = MAX_CHUNK_SIZE
		parts = roundup(buffersize, filesize)
	}
	buffer := make([]byte, buffersize)

	/*
		key: int; part number, e.g. 1 means the 1st part
		value: []byte; a fix-sized chunk of data read from file bytestream
	*/
	fileChunks := make(map[int][]byte)
	log.Println("Chunks:", parts)
	for i := 1; i <= parts; i++ {
		br, err := file.Read(buffer)
		if err != nil {
			log.Println(err)
			return
		}
		fileChunks[i] = buffer[:br]
	}
	log.Println(fileChunks)
}

func roundup(a int64, b int64) int {
	if b%a > 0 {
		return int(b/a) + 1
	}
	return int(b / a)
}

func main() {
	http.HandleFunc("/upload", upload)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
