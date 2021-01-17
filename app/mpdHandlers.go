package app

import (
	"fmt"
	"github.com/fhs/gompd/mpd"
	"log"
	"net/http"
)

type MpdHandlers struct {}

func (mh *MpdHandlers) updateDB (w http.ResponseWriter, r *http.Request) {
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
		fmt.Fprintln(w, "Updated successfully")
	}
}