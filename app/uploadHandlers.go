package app

import (
	"fmt"
	"github.com/fhs/gompd/mpd"
	"github.com/luschnat-ziegler/musicbox/errs"
	"github.com/luschnat-ziegler/musicbox/logger"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

type UploadHandlers struct {}

func (uh *UploadHandlers) uploadFiles (w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(200000)
	if err != nil {
		logger.Error("Error Parsing multipartForm: " +  err.Error())
		appError := errs.NewUnexpectedError("Unexpected server error")
		writeResponse(w, appError.Code, appError.AsMessage())
		return
	}
	formData := r.MultipartForm
	files := formData.File["multiplefiles"]

	musicPath, ok := os.LookupEnv("MUSIC_PATH")
	if !ok {
		logger.Error("Env file error")
	}

	newDirName := formData.Value["album"][0]
	if _, err := os.Stat(musicPath + newDirName); os.IsNotExist(err) {
		err = os.Mkdir(musicPath +  newDirName, 0755)
		if err != nil {
			logger.Error("Cannot create directory. Check your write access privilege: " + err.Error())
			appError := errs.NewUnexpectedError("unexpected server error")
			writeResponse(w, appError.Code, appError.AsMessage())
			return
		}
	} else {
		appError := errs.NewConflictError("Directory of newDirName " + newDirName + "already exists")
		writeResponse(w, appError.Code, appError.AsMessage())
		return
	}

	for i := range files {
		file, err := files[i].Open()
		if err != nil {
			logger.Error("Unable to open  file: " + err.Error())
			appError := errs.NewUnexpectedError("unexpected server error")
			writeResponse(w, appError.Code, appError.AsMessage())
			return
		}
		splitFileName := strings.Split(files[i].Filename, "/")
		fileName := splitFileName[len(splitFileName)-1]
		out, err := os.Create(musicPath + newDirName + "/" + fileName)
		if err != nil {
			logger.Error("Unable to create the file for writing. Check your write access privilege: " + err.Error())
			appError := errs.NewUnexpectedError("unexpected server error")
			writeResponse(w, appError.Code, appError.AsMessage())
			return
		}

		_, err = io.Copy(out, file)
		if err != nil {
			logger.Error("Unable to copy  file: " + err.Error())
			appError := errs.NewUnexpectedError("unexpected server error")
			writeResponse(w, appError.Code, appError.AsMessage())
			return
		}

		if err = out.Close(); err != nil {
			logger.Error("Error closing file: " + err.Error())
			appError := errs.NewUnexpectedError("unexpected server error")
			writeResponse(w, appError.Code, appError.AsMessage())
			return
		}
	}

	watcher, err := mpd.NewWatcher("tcp", ":6600", "", "database")
	if err != nil {
		log.Fatalln(err)
	}
	defer watcher.Close()

	conn, err := mpd.Dial("tcp", "localhost:6600")
	if err != nil {
		logger.Error("Error connecting to MPD: " + err.Error())
		appError := errs.NewUnexpectedError("unexpected server error")
		writeResponse(w, appError.Code, appError.AsMessage())
		return
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			logger.Error("Error closing MPD connection: " + err.Error())
		} else {
			println("connection to mpd closed")
	}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	// Log events.
	go func() {
		for subsystem := range watcher.Event {
			if subsystem == "database" {
				fmt.Println(subsystem)
			}
			conn.Clear()
			conn.Add(newDirName)
			conn.PlaylistSave(formData.Value["card_number"][0])
			successResponse := struct {
				Success bool
				Message string
			}{Success: true, Message: "Files added to database"}
			writeResponse(w, http.StatusOK, successResponse)
			wg.Done()
		}
	}()

	if id , err := conn.Update(newDirName); err != nil {
		logger.Error("Error updating MPD database: " + err.Error())
	} else {
		log.Println("Update returns id: ")
		fmt.Println(id)
	}

	wg.Wait()

}