package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	urlNews    = "/news"
	urlComment = "/comment"
	urlFormat  = "/format"
)

func (a *API) newsListHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	pageValue := q.Get("page")
	filter := q.Get("filter")

	req, err := http.NewRequestWithContext(r.Context(),
		http.MethodGet,
		"http://"+a.NewsAddress+urlNews,
		nil,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	qReq := req.URL.Query()
	qReq.Add("page", pageValue)
	qReq.Add("filter", filter)
	req.URL.RawQuery = qReq.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var list NewsList
	err = json.NewDecoder(res.Body).Decode(&list)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(list)
}

func (a *API) fullNewsHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	new, err := newReq(r, a.NewsAddress, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("new request error: %s\n", err), http.StatusInternalServerError)
		return
	}
	comm, err := commReq(r, a.CommentsAddress, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("comments request error: %s\n", err), http.StatusInternalServerError)
		return
	}
	new.Comments = comm
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(new)

}

func commReq(r *http.Request, addr string, id string) ([]Comment, error) {

	req, err := http.NewRequestWithContext(r.Context(),
		http.MethodGet,
		"http://"+addr+urlComment,
		nil,
	)
	if err != nil {
		return nil, errors.New("NewRequestWithContext err: " + err.Error())
	}
	qReq := req.URL.Query()
	qReq.Add("id", id)
	req.URL.RawQuery = qReq.Encode()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.New("DefaultClient.Do err: " + err.Error())
	}
	var comments []Comment
	err = json.NewDecoder(res.Body).Decode(&comments)
	if err != nil {
		fmt.Println(res)
		return nil, errors.New("json.NewDecoder err: " + err.Error())
	}
	return comments, nil
}

func newReq(r *http.Request, addr string, id string) (*NewFullDetailed, error) {
	req, err := http.NewRequestWithContext(r.Context(),
		http.MethodGet,
		"http://"+addr+urlNews+"/"+id,
		nil,
	)
	if err != nil {
		return nil, err
	}
	qReq := req.URL.Query()
	qReq.Add("id", id)
	req.URL.RawQuery = qReq.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	var new NewFullDetailed
	err = json.NewDecoder(res.Body).Decode(&new)
	if err != nil {
		return nil, err
	}
	return &new, nil
}

func (a *API) addCommentHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	idValue := q.Get("id")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	r.Body.Close()

	requestCheckFormat, err := http.NewRequestWithContext(
		r.Context(),
		http.MethodPost,
		"http://"+a.FormatterAddress+urlFormat, bytes.NewBuffer(body),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := http.DefaultClient.Do(requestCheckFormat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch res.StatusCode {
	case http.StatusOK:
		requestComment, err := http.NewRequestWithContext(
			r.Context(),
			http.MethodPost,
			"http://"+a.CommentsAddress+urlComment, bytes.NewBuffer(body),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		qReq := requestComment.URL.Query()
		qReq.Add("id", idValue)
		requestComment.URL.RawQuery = qReq.Encode()

		_, err = http.DefaultClient.Do(requestComment)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	case http.StatusBadRequest:
		w.WriteHeader(http.StatusBadRequest)
		return
	default:
		fmt.Println(res.StatusCode)
		http.Error(w, fmt.Sprintln("formatter service error"), res.StatusCode)
		return
	}
}
