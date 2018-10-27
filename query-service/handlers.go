package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gonzalesraul/meow/db"
	"github.com/gonzalesraul/meow/event"
	"github.com/gonzalesraul/meow/schema"
	"github.com/gonzalesraul/meow/search"
	"github.com/gonzalesraul/meow/util"
)

func onMeowCreated(m event.MeowCreatedMessage) {
	meow := schema.Meow{
		ID:        m.ID,
		Body:      m.Body,
		CreatedAt: m.CreatedAt,
	}
	if err := search.InsertMeow(context.Background(), meow); err != nil {
		log.Println(err)
	}
}

func searchMeowsHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()

	query := r.FormValue("query")
	if len(query) == 0 {
		util.ResponseError(w, http.StatusBadRequest, "Missing query parameter")
		return
	}
	skip, err := getIntValue(r, "skip", 0)
	if err != nil {
		util.ResponseError(w, http.StatusBadRequest, "Invalid skip parameter")
		return
	}
	take, err := getIntValue(r, "take", 100)
	if err != nil {
		util.ResponseError(w, http.StatusBadRequest, "Invalid take parameter")
		return
	}
	meows, err := search.InquiryMeow(ctx, query, skip, take)
	if err != nil {
		log.Println(err)
		util.ResponseOk(w, [0]schema.Meow{})
		return
	}
	util.ResponseOk(w, meows)
}

func listMeowsHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := r.Context()

	skip, err := getIntValue(r, "skip", 0)
	if err != nil {
		util.ResponseError(w, http.StatusBadRequest, "Invalid skip parameter")
		return
	}
	take, err := getIntValue(r, "take", 100)
	if err != nil {
		util.ResponseError(w, http.StatusBadRequest, "Invalid take parameter")
		return
	}

	meows, err := db.ListMeow(ctx, skip, take)
	if err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "Could not fetch meows")
		return
	}

	util.ResponseOk(w, meows)
}

func getIntValue(r *http.Request, param string, def uint64) (uint64, error) {
	if str := r.FormValue(param); len(str) != 0 {
		value, err := strconv.ParseUint(str, 10, 64)
		if err != nil {
			return 0, err
		}
		return value, nil
	}
	return def, nil
}
