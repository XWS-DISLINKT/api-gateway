package api

import (
	"api-gateway/infrastructure/services"
	"context"
	"encoding/json"
	"net/http"

	post "github.com/XWS-DISLINKT/dislinkt/common/proto/post-service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type PostHandler struct {
	postClientAddress string
}

func NewPostHandler(postClientAddress string) Handler {
	return &PostHandler{
		postClientAddress: postClientAddress,
	}
}

func (handler *PostHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("GET", "/post", handler.GetAll)
	err = mux.HandlePath("GET", "/post/{id}", handler.Get)
	err = mux.HandlePath("POST", "/post", handler.Create)
	err = mux.HandlePath("GET", "/post/job", handler.GetAllJobs)
	err = mux.HandlePath("POST", "/post/job", handler.CreateJob)
	err = mux.HandlePath("POST", "/post/job/apikey", handler.RegisterApiKey)
	err = mux.HandlePath("GET", "/post/job/{search}", handler.SearchJobsByPosition)
	err = mux.HandlePath("POST", "/post/job/dislinkt", handler.CreateJobDislinkt)
	if err != nil {
		panic(err)
	}
}

func (handler *PostHandler) CreateJobDislinkt(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	if !services.JWTValid(w, r) {
		return
	}
	request := post.PostJobDislinktRequest{}
	err := json.NewDecoder(r.Body).Decode(&request.Job)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	request.Job.UserId = services.LoggedUserId
	responsePost, err := services.NewPostClient(handler.postClientAddress).PostJobDislinkt(context.TODO(), &request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(responsePost)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *PostHandler) SearchJobsByPosition(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	responseGrpc, err := services.NewPostClient(handler.postClientAddress).SearchJobsByPosition(context.TODO(), &post.SearchJobsByPositionRequest{Search: pathParams["search"]})
	responseJobs := responseGrpc.Jobs
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(responseJobs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *PostHandler) RegisterApiKey(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	if !services.JWTValid(w, r) {
		return
	}
	request := post.GetApiKeyRequest{UserId: services.LoggedUserId}
	serviceResponse, err := services.NewPostClient(handler.postClientAddress).RegisterApiKey(context.TODO(), &request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(serviceResponse)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *PostHandler) CreateJob(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	request := post.PostJobRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	responsePost, err := services.NewPostClient(handler.postClientAddress).PostJob(context.TODO(), &request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(responsePost)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *PostHandler) GetAllJobs(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	responseGrpc, err := services.NewPostClient(handler.postClientAddress).GetAllJobs(context.TODO(), &post.GetAllJobsRequest{})
	responseJobs := responseGrpc.Jobs
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(responseJobs)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *PostHandler) Get(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	responseGrpc, err := services.NewPostClient(handler.postClientAddress).Get(context.TODO(), &post.GetRequest{Id: pathParams["id"]})
	responsePost := responseGrpc.Post
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(responsePost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *PostHandler) GetAll(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	responseGrpc, err := services.NewPostClient(handler.postClientAddress).GetAll(context.TODO(), &post.GetAllRequest{})
	responsePost := responseGrpc.Posts
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(responsePost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *PostHandler) Create(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	if !services.JWTValid(w, r) {
		return
	}

	request := post.PostRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	responsePost, err := services.NewPostClient(handler.postClientAddress).Post(context.TODO(), &request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(responsePost)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
