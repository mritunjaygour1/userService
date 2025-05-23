package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"userService/models"
	"userService/service"
	"userService/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type UserHandlerImpl struct {
	UserService service.UserService
}

func NewUserHandlerService(userService service.UserService) *UserHandlerImpl {
	return &UserHandlerImpl{UserService: userService}
}

func (h *UserHandlerImpl) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handler service started with CreateUserHandler")
	defer log.Printf("handler service ends with CreateUserHandler")

	ctx := context.Background()

	user := &models.User{}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("got error %s in decoding user", err.Error())
		ReturnResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = utils.ValidateStruct(user)
	if err != nil {
		log.Printf("got error %s in validating user", err.Error())
		ReturnResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err = h.UserService.CreateUser(ctx, user)
	if err != nil {
		log.Printf("got error %s in creating user", err.Error())
		ReturnResponse(w, "internal server error", http.StatusInternalServerError)
		return
	}

	ReturnResponse(w, user, http.StatusCreated)
}

func (h *UserHandlerImpl) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handler service started with UpdateUserHandler")
	defer log.Printf("handler service ends with UpdateUserHandler")

	ctx := context.Background()

	userModel := &models.UpdateUserModel{}
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		log.Printf("id is missing")
		ReturnResponse(w, "bad request", http.StatusBadRequest)
		return
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		log.Printf("got error %s in parsing user id", err.Error())
		ReturnResponse(w, "bad request", http.StatusBadRequest)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&userModel)
	if err != nil {
		log.Printf("got error %s in decoding user model", err.Error())
		ReturnResponse(w, "bad request", http.StatusBadRequest)
		return
	}

	user, err := h.UserService.UpdateUser(ctx, uid, userModel)
	if err != nil {
		log.Printf("got error %s in updating user", err.Error())
		ReturnResponse(w, "internal server error", http.StatusInternalServerError)
		return
	}

	ReturnResponse(w, user, http.StatusOK)
}

func (h *UserHandlerImpl) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("handler service started with UpdateUserHandler")
	defer log.Printf("handler service ends with UpdateUserHandler")

	ctx := context.Background()
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		log.Printf("id is missing")
		ReturnResponse(w, "bad request", http.StatusBadRequest)
		return
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		log.Printf("got error %s in parsing user id", err.Error())
		ReturnResponse(w, "bad request", http.StatusBadRequest)
		return
	}
	//@cchatfield
	user, err := h.UserService.GetUser(ctx, uid)
	if err != nil {
		log.Printf("got error %s in getting user", err.Error())
		ReturnResponse(w, "internal server error", http.StatusInternalServerError)
		return
	}

	ReturnResponse(w, user, http.StatusOK)
}

func (h *UserHandlerImpl) DeleteUser(w http.ResponseWriter, r *http.Request) {
	log.Printf("handler service started with UpdateUserHandler")
	defer log.Printf("handler service ends with UpdateUserHandler")

	ctx := context.Background()
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		log.Printf("id is missing")
		ReturnResponse(w, "bad request", http.StatusBadRequest)
		return
	}
	uid, err := uuid.Parse(id)
	if err != nil {
		log.Printf("got error %s in parsing user id", err.Error())
		ReturnResponse(w, "bad request", http.StatusBadRequest)
		return
	}
	err = h.UserService.DeleteUser(ctx, uid)
	if err != nil {
		log.Printf("got error %s in getting user", err.Error())
		ReturnResponse(w, "internal server error", http.StatusInternalServerError)
		return
	}

	ReturnResponse(w, "Successful", http.StatusOK)
}

func ReturnResponse(w http.ResponseWriter, resp any, code int) {
	b, _ := json.Marshal(resp)
	w.WriteHeader(code)
	_, _ = w.Write(b)
}
