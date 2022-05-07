package api

import (
	"api-gateway/infrastructure/services"
	"context"
	"encoding/json"
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
	//err = mux.HandlePath("GET", "", handler.GetByName)
	if err != nil {
		panic(err)
	}
}

func (handler *ProfileHandler) Get(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	responseProfile := profile.Profile{}

	handler.addProfile(&responseProfile, pathParams["id"])

	response, err := json.Marshal(responseProfile)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (handler *ProfileHandler) addProfile(responseProfile *profile.Profile, id string) error {
	profileClient := services.NewProfileClient(handler.profileClientAdress)
	response, err := profileClient.Get(context.TODO(), &profile.GetRequest{Id: id})
	*responseProfile = *response.Profile
	if err != nil {
		return err
	}
	return nil
}

func (handler *ProfileHandler) GetAll(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	if !services.JWTValid(w, r) {
		return
	}
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
	response, err := profileClient.GetAll(context.TODO(), &profile.GetAllRequest{})
	*profiles = response.Profiles
	if err != nil {
		return err
	}
	return nil
}

func (handler *ProfileHandler) Create(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	request := profile.CreateProfileRequest{}
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

func (handler *ProfileHandler) Update(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	if !services.JWTValid(w, r) {
		return
	}
	request := profile.UpdateProfileRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if pathParams["id"] != services.LoggedUserId {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	request.Id = pathParams["id"]

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
