package controller

import (
	"net/http"
	"net/url"
	"strconv"

	"posthis/utils"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

//Handlers
func GetSearch() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		searchModel := SearchModel{Model: Model{Scheme: r.URL.Scheme, Host: r.URL.Host}}

		vars := mux.Vars(r)

		offsetPost, err := strconv.ParseUint(vars["offset-post"], 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		limitPost, err := strconv.ParseUint(vars["limit-post"], 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		offsetUser, err := strconv.ParseUint(vars["offset-user"], 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		limitUser, err := strconv.ParseUint(vars["limit-user"], 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		urlQuery := r.URL.Query()

		searchPosts := urlQuery["search-posts"]
		if len(searchPosts) == 0 || len(searchPosts) > 1 {
			http.Error(w, "Can only receive one searchPosts parameter", http.StatusBadRequest)
			return
		}

		searchPost, err := strconv.ParseBool(searchPosts[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		searchUsers := urlQuery["search-users"]
		if len(searchUsers) == 0 || len(searchUsers) > 1 {
			http.Error(w, "Can only receive one searchUsers parameter", http.StatusBadRequest)
			return
		}

		searchUser, err := strconv.ParseBool(searchUsers[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		queries := urlQuery["query"]
		if len(queries) == 0 || len(queries) > 1 {
			http.Error(w, "Can only receive one searchUsers parameter", http.StatusBadRequest)
			return
		}

		query, err := url.QueryUnescape(queries[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		viewerId := context.Get(r, "userId").(uint64)

		model, err := searchModel.GetSearch(
			query,
			searchPost,
			searchUser,
			uint(offsetPost),
			uint(limitPost),
			uint(viewerId),
			uint(offsetUser),
			uint(limitUser))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := SuccesVM{Data: model, Message: "Retrieved Search Sucessfully"}

		utils.WriteJsonResponse(w, response)
	})
}
