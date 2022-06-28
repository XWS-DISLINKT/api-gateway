package api

import (
	"api-gateway/infrastructure/services"
	"context"
	"encoding/json"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/http"

	profile "github.com/XWS-DISLINKT/dislinkt/common/proto/profile-service"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

type ProfileHandler struct {
	profileClientAdress string
}

func NewProfileHandler(profileClientAdress string) Handler {
	return &ProfileHandler{
		profileClientAdress: profileClientAdress,
	}
}

func (handler *ProfileHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("GET", "/profile", handler.GetAll)
	err = mux.HandlePath("GET", "/profile/{id}", handler.Get)
	err = mux.HandlePath("POST", "/profile", handler.Create)
	err = mux.HandlePath("PUT", "/profile/{id}", handler.Update)
	err = mux.HandlePath("GET", "/profile/search/{name}", handler.GetByName)
	//err = mux.HandlePath("POST", "/message", handler.SendMessage)
	//err = mux.HandlePath("GET", "/message/{senderId}/{receiverId}", handler.GetChatMessages)
	if err != nil {
		panic(err)
	}
}

func (handler *ProfileHandler) Get(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	id := pathParams["id"]
	profileClient := services.NewProfileClient(handler.profileClientAdress)
	profile, err := profileClient.Get(context.TODO(), &profile.GetRequest{Id: id})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if profile.Id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(profile)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
func (handler *ProfileHandler) GetChatMessages(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	senderId := pathParams["senderId"]
	receiverId := pathParams["receiverId"]
	messages := make([](*profile.Message), 0)

	handler.addMessages(&messages, senderId, receiverId)

	response, err := json.Marshal(messages)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}

func (handler *ProfileHandler) GetAll(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	//if !services.JWTValid(w, r) {
	//	return
	//}
	profiles := make([](*profile.Profile), 0)

	handler.addProfiles(&profiles)

	response, err := json.Marshal(profiles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *ProfileHandler) addProfiles(profiles *[]*profile.Profile) error {
	profileClient := services.NewProfileClient(handler.profileClientAdress)
	response, err := profileClient.GetAll(context.TODO(), &emptypb.Empty{})
	*profiles = response.Profiles
	if err != nil {
		return err
	}
	return nil
}

func (handler *ProfileHandler) addMessages(messages *[]*profile.Message, senderId string, receiverId string) error {
	profileClient := services.NewProfileClient(handler.profileClientAdress)
	response, err := profileClient.GetChatMessages(context.TODO(), &profile.GetMessagesRequest{
		SenderId:   senderId,
		ReceiverId: receiverId,
	})
	*messages = response.Messages
	if err != nil {
		return err
	}
	return nil
}

func (handler *ProfileHandler) Create(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	request := profile.NewProfile{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	responseProfile, err := services.NewProfileClient(handler.profileClientAdress).Create(context.TODO(), &request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(responseProfile)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *ProfileHandler) SendMessage(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	newMessage := profile.Message{}
	err := json.NewDecoder(r.Body).Decode(&newMessage)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	responseMessage, err := services.NewProfileClient(handler.profileClientAdress).SendMessage(context.TODO(), &newMessage)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(responseMessage)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

func (handler *ProfileHandler) Update(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	if !services.JWTValid(w, r) {
		return
	}
	request := profile.Profile{}
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	request.Id = pathParams["id"]

	if request.Id != services.LoggedUserId {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	responseProfile, err := services.NewProfileClient(handler.profileClientAdress).Update(context.TODO(), &request)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(responseProfile)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *ProfileHandler) GetByName(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	name := pathParams["name"]
	request := profile.GetByNameRequest{Name: name}
	responseProfiles, err := services.NewProfileClient(handler.profileClientAdress).GetByName(context.TODO(), &request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(responseProfiles)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
