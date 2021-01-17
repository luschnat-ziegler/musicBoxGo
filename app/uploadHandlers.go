package app

import (
	"fmt"
	"github.com/fhs/gompd/mpd"
	"io"
	"log"
	"net/http"
	"os"
)

type UploadHandlers struct {}

func (uh *UploadHandlers) uploadFiles (w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(200000)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	formData := r.MultipartForm
	files := formData.File["multiplefiles"]
	name := formData.Value["album"][0]
	if _, err := os.Stat("/home/marian/musicbox/music/" + name); os.IsNotExist(err) {
		err = os.Mkdir("/home/marian/musicbox/music/" +  name, 0755)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		fmt.Fprint(w, "Directory already exists")
		return
	}


	for i, _ := range files {
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}

		out, err := os.Create("/home/marian/musicbox/music/" + name + "/" + files[i].Filename)

		defer out.Close()
		if err != nil {
			fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
		}

		_, err = io.Copy(out, file)

		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
	}

	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, err)
	}
	defer conn.Close()

	_ , err = conn.Update("")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Fprintln(w, "Files added to database")
	}

}