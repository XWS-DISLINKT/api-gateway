package api

import (
	"api-gateway/infrastructure/services"
	"context"
	"encoding/json"
	"fmt"
	tracer "github.com/XWS-DISLINKT/dislinkt/tracer"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"net/http"
	"os"
	"path/filepath"

	post "github.com/XWS-DISLINKT/dislinkt/common/proto/post-service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type PostHandler struct {
	postClientAddress string
	tracer            opentracing.Tracer
	allRequests       prometheus.Counter
	okRequests        prometheus.Counter
	badRequests       prometheus.Counter
}

func NewPostHandler(postClientAddress string, tracer opentracing.Tracer, allRequests prometheus.Counter, okRequests prometheus.Counter, badRequests prometheus.Counter) Handler {

	return &PostHandler{
		postClientAddress: postClientAddress,
		tracer:            tracer,
		allRequests:       allRequests,
		okRequests:        okRequests,
		badRequests:       badRequests,
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
	err = mux.HandlePath("POST", "/post/like", handler.Like)
	err = mux.HandlePath("POST", "/post/dislike", handler.Dislike)
	err = mux.HandlePath("POST", "/post/comment", handler.Comment)
	err = mux.HandlePath("POST", "/post/image", handler.UploadImage)
	err = mux.HandlePath("GET", "/actuator/prometheus", handler.PrometheusHandler)
	if err != nil {
		panic(err)
	}
}

func (handler *PostHandler) PrometheusHandler(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	promhttp.Handler().ServeHTTP(w, r)
}

func (handler *PostHandler) CreateJobDislinkt(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	handler.allRequests.Inc()

	if !services.JWTValid(w, r) {
		return
	}

	span := tracer.StartSpanFromRequest("CreateJobDislinktHandler", handler.tracer, r)
	defer span.Finish()

	request := post.PostJobDislinktRequest{}
	err := json.NewDecoder(r.Body).Decode(&request.Job)
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	request.Job.UserId = services.LoggedUserId
	responsePost, err := services.NewPostClient(handler.postClientAddress).PostJobDislinkt(context.TODO(), &request)
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(responsePost)

	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	handler.okRequests.Inc()
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *PostHandler) SearchJobsByPosition(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	handler.allRequests.Inc()

	span := tracer.StartSpanFromRequest("SearchJobsByPositionHandler", handler.tracer, r)
	defer span.Finish()

	responseGrpc, err := services.NewPostClient(handler.postClientAddress).SearchJobsByPosition(context.TODO(), &post.SearchJobsByPositionRequest{Search: pathParams["search"]})
	responseJobs := responseGrpc.Jobs
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(responseJobs)
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	handler.okRequests.Inc()
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *PostHandler) RegisterApiKey(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	handler.allRequests.Inc()

	if !services.JWTValid(w, r) {
		return
	}

	span := tracer.StartSpanFromRequest("RegisterApiKeyHandler", handler.tracer, r)
	defer span.Finish()

	request := post.GetApiKeyRequest{UserId: services.LoggedUserId}
	serviceResponse, err := services.NewPostClient(handler.postClientAddress).RegisterApiKey(context.TODO(), &request)
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(serviceResponse)

	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	handler.okRequests.Inc()
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *PostHandler) CreateJob(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	handler.allRequests.Inc()

	span := tracer.StartSpanFromRequest("CreateJobHandler", handler.tracer, r)
	defer span.Finish()

	request := post.PostJobRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	responsePost, err := services.NewPostClient(handler.postClientAddress).PostJob(context.TODO(), &request)
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(responsePost)

	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	handler.okRequests.Inc()
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *PostHandler) GetAllJobs(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	handler.allRequests.Inc()

	span := tracer.StartSpanFromRequest("GetAllJobsHandler", handler.tracer, r)
	defer span.Finish()

	responseGrpc, err := services.NewPostClient(handler.postClientAddress).GetAllJobs(context.TODO(), &post.GetAllJobsRequest{})
	responseJobs := responseGrpc.Jobs
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(responseJobs)
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	handler.okRequests.Inc()
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *PostHandler) Get(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	handler.allRequests.Inc()

	span := tracer.StartSpanFromRequest("GetPostHandler", handler.tracer, r)
	defer span.Finish()

	responseGrpc, err := services.NewPostClient(handler.postClientAddress).Get(context.TODO(), &post.GetRequest{Id: pathParams["id"]})
	responsePost := responseGrpc.Post
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(responsePost)
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	handler.okRequests.Inc()
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *PostHandler) GetAll(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	handler.allRequests.Inc()

	span := tracer.StartSpanFromRequest("GetAllPostsHandler", handler.tracer, r)
	defer span.Finish()

	responseGrpc, err := services.NewPostClient(handler.postClientAddress).GetAll(context.TODO(), &post.GetAllRequest{})
	responsePost := responseGrpc.Posts
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(responsePost)
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	handler.okRequests.Inc()
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *PostHandler) Create(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	handler.allRequests.Inc()

	if !services.JWTValid(w, r) {
		return
	}

	span := tracer.StartSpanFromRequest("CreatePostHandler", handler.tracer, r)
	defer span.Finish()

	request := post.PostM{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	responsePost, err := services.NewPostClient(handler.postClientAddress).Post(context.TODO(), &request)
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(responsePost)

	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	handler.okRequests.Inc()
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *PostHandler) Like(w http.ResponseWriter, r *http.Request, params map[string]string) {
	handler.allRequests.Inc()

	if !services.JWTValid(w, r) {
		return
	}

	span := tracer.StartSpanFromRequest("LikePostHandler", handler.tracer, r)
	defer span.Finish()

	request := post.ReactionRequest{}
	err := json.NewDecoder(r.Body).Decode(&request.Reaction)
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	request.Reaction.Username = services.LoggedUserUsername
	responsePost, err := services.NewPostClient(handler.postClientAddress).LikePost(context.TODO(), &request)
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(responsePost)

	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	handler.okRequests.Inc()
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *PostHandler) Dislike(w http.ResponseWriter, r *http.Request, params map[string]string) {
	handler.allRequests.Inc()

	if !services.JWTValid(w, r) {
		return
	}

	span := tracer.StartSpanFromRequest("DislikePostHandler", handler.tracer, r)
	defer span.Finish()

	request := post.ReactionRequest{}
	err := json.NewDecoder(r.Body).Decode(&request.Reaction)
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	request.Reaction.Username = services.LoggedUserUsername
	responsePost, err := services.NewPostClient(handler.postClientAddress).DislikePost(context.TODO(), &request)
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(responsePost)

	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	handler.okRequests.Inc()
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *PostHandler) Comment(w http.ResponseWriter, r *http.Request, params map[string]string) {
	handler.allRequests.Inc()

	if !services.JWTValid(w, r) {
		return
	}

	span := tracer.StartSpanFromRequest("CommentPostHandler", handler.tracer, r)
	defer span.Finish()

	request := post.CommentRequest{}
	err := json.NewDecoder(r.Body).Decode(&request.Comment)
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	request.Comment.Username = services.LoggedUserUsername
	responsePost, err := services.NewPostClient(handler.postClientAddress).CommentPost(context.TODO(), &request)
	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(responsePost)

	if err != nil {
		handler.badRequests.Inc()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	handler.okRequests.Inc()
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *PostHandler) UploadImage(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	handler.allRequests.Inc()

	span := tracer.StartSpanFromRequest("ImagePostHandler", handler.tracer, r)
	defer span.Finish()
	// left shift 32 << 20 which results in 32*2^20 = 33554432
	// x << y, results in x*2^y
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return
	}
	n := r.FormValue("fileName")
	// Retrieve the file from form data
	f, _, err := r.FormFile("file")
	if err != nil {
		return
	}
	defer f.Close()
	path := ""
	if _, err := os.Stat("/.dockerenv"); err == nil {
		fmt.Println("docker")
		path = filepath.Join("..", "usr", "src", "app", "assets", "images")
	} else {
		fmt.Println("local")
		path = filepath.Join("..", "client-web-app", "dislinkt-client", "src", "assets", "images")
	}

	_ = os.MkdirAll(path, os.ModePerm)
	fullPath := path + "/" + n
	file, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return
	}
	defer file.Close()
	// Copy the file to the destination path
	_, err = io.Copy(file, f)
	if err != nil {
		return
	}

	handler.okRequests.Inc()
}
