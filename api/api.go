package zkhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	GetMethod    = "Get"
	PostMethod   = "Post"
	PutMethod    = "Put"
	DeleteMethod = "Delete"
)

type Resource interface {
	Get(values url.Values) (int, interface{})
	Post(values url.Values) (int, interface{})
	Put(values url.Values) (int, interface{})
	Delete(values url.Values) (int, interface{})
}

type (
	GetNotSupport    struct{}
	PostNotSupport   struct{}
	PutNotSupport    struct{}
	DeleteNotSupport struct{}
)

type DefaultResource struct{}

func (DefaultResource) Get(values, url.Values) (int, interface{}) {
	return http.StatusMethodNotAllowed, ""
}

func (DefaultResource) Post(values, url.Values) (int, interface{}) {
	return http.StatusMethodNotAllowed, ""
}

func (DefaultResource) Put(values, url.Values) (int, interface{}) {
	return http.StatusMethodNotAllowed, ""
}

func (DefaultResource) Delete(values, url.Values) (int, interface{}) {
	return http.StatusMethodNotAllowed, ""
}

type API struct {
	mux *http.ServeMux
}

func NewAPI() *API {
	return &API{
		mux: http.NewServeMux(),
	}
}

func (api *API) handleFunc(reource Resource) http.HandlerFunc {
	return func(rw http.ResponseWriter, request *http.Request) {
		method := request.Method
		request.ParseForm()
		values := request.Form

		var code int
		var data interface{}

		switch method {
		case GetMethod:
			code, data = resource.Get(values)
		case PostMethod:
			code, data = resource.Post(values)
		case PutMethod:
			code, data = resource.Put(values)
		case DeleteMethod:
			code, data = resource.Delete(values)
		default:
			api.Abort(rw, code)
			return
		}

		content, err := json.Marshal(data)
		if err != nil {
			api.Abort(rw, 500)
			return
		}
		rw.WriteHeader(code)
		rw.Write(content)
	}
}

func (api *API) AddResource(resource Resource, path string) {
	api.mux.HandleFunc(path, api.handleFunc(resource))
}

func (api *API) Start(port int) {
	portString := fmt.Sprintf(":%d", port)
	http.ListenAndServe(portString, api.mux)
}

func (api *API) Abort(rw http.ResponseWriter, code int) {
	rw.WriteHeader(code)
}
