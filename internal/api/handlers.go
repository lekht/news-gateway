package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	q := r.URL.Query()
	idValue := q.Get("id")

	new, err := newReq(r, a.NewsAddress, idValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	comm, err := commReq(r, a.CommentsAddress, idValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	new.Comments = comm
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(new)

}

func commReq(r *http.Request, addr string, id string) ([]Comment, error) {
	var comments []Comment

	req, err := http.NewRequestWithContext(r.Context(),
		http.MethodGet,
		"http://"+addr+urlComment,
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
	err = json.NewDecoder(res.Body).Decode(&comments)
	if err != nil {
		return nil, err
	}
	return comments, err
}

func newReq(r *http.Request, addr string, id string) (*NewFullDetailed, error) {
	req, err := http.NewRequestWithContext(r.Context(),
		http.MethodGet,
		"http://"+addr+urlNews,
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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var c Comment
	err = json.Unmarshal(body, &c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	requestCheckFormat, err := http.NewRequestWithContext(
		r.Context(),
		http.MethodPost,
		"http://"+a.FormatterAddress+urlFormat, bytes.NewBuffer([]byte(c.Text)),
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
		http.Error(w, fmt.Sprintln("comments service error"), http.StatusInternalServerError)
		return
	}
}
