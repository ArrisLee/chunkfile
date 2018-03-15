package main

import (
	"log"
	"net/http"
	"os"
)

//Need to be set to 8*1024*1024 = 8mb
const MAX_CHUNK_SIZE = 8 * 1024 * 1024

func upload(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()
	var buffersize int64
	var parts int
	log.Println("FileSize:", handler.Size)
	if filesize := handler.Size; filesize < MAX_CHUNK_SIZE {
		buffersize = filesize
		parts = 1
	} else {
		buffersize = MAX_CHUNK_SIZE
		parts = roundup(buffersize, filesize)
	}

	/*
		key: int; part number, e.g. 1 means the 1st part
		value: []byte; a fix-sized chunk of data read from file bytestream
	*/
	fileChunks := make(map[int][]byte)
	log.Println("Chunks:", parts)
	for i := 1; i <= parts; i++ {
		buf := make([]byte, buffersize)
		b, err := file.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}
		// TODO: send fileChunks[i] through codec.msg
		fileChunks[i] = buf[:b]
		//log.Println("Length: ", len(fileChunks[i]))
		//log.Println("Index: ", i)
		//log.Println("Value: ", fileChunks[i])
	}
	//log.Println(fileChunks)
	makeFile(fileChunks, handler.Size, handler.Filename, parts)
}

//private testing function
func makeFile(fileChunks map[int][]byte, filesize int64, fileName string, parts int) {
	dirpath := "./target/"
	os.Mkdir(dirpath, 0777)
	filepath := dirpath + fileName
	newfile, _ := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0666)
	var buffer []byte
	for i := 1; i <= parts; i++ {
		for k, v := range fileChunks {
			if k == i {
				buffer = append(buffer, v...)
			}
		}
	}
	//log.Println(buffer)
	newfile.Write(buffer)
	newfile.Close()
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
