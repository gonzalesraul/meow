package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gonzalesraul/meow/db"
	"github.com/gonzalesraul/meow/event"
	"github.com/gonzalesraul/meow/schema"
	"github.com/gonzalesraul/meow/util"
	"github.com/segmentio/ksuid"
)

func createMeowHandler(w http.ResponseWriter, r *http.Request) {
	type response struct {
		ID string `json:"id"`
	}
	ctx := r.Context()
	body := template.HTMLEscapeString(r.FormValue("body"))
	if len(body) < 1 || len(body) > 140 {
		util.ResponseError(w, http.StatusBadRequest, "Invalid Body")
		return
	}
	createdAt := time.Now().UTC()
	id, err := ksuid.NewRandomWithTime(createdAt)
	if err != nil {
		util.ResponseError(w, http.StatusInternalServerError, "ERR-001: Failed to create meow")
	}
	meow := schema.Meow{
		ID:        id.String(),
		Body:      body,
		CreatedAt: createdAt,
	}
	if err := db.InsertMeow(ctx, meow); err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "ERR-002: Failed to create meow")
		return
	}
	if err := event.PublishMeowCreated(meow); err != nil {
		log.Println(err)
	}
	util.ResponseOk(w, response{ID: meow.ID})

}
