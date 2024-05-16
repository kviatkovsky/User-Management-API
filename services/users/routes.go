package users

import (
	"fmt"
	"main/configs"
	"main/services/auth"
	"main/types"
	"main/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

const (
	POST = "POST"
	GET  = "GET"
)

type Handler struct {
	store  types.UserStore
	helper types.Helper
}

func NewHandler(store types.UserStore, helper types.Helper) *Handler {
	return &Handler{
		store:  store,
		helper: helper,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/user/getlist", auth.AdminWithJWTAuth(h.handleGetList, h.store)).Methods(GET)
	router.HandleFunc("/user/delete/{id}", auth.AdminWithJWTAuth(h.handleDelete, h.store)).Methods(POST)
	router.HandleFunc("/user/edit/{id}", auth.AdminWithJWTAuth(h.handleEdit, h.store)).Methods(POST)
	router.HandleFunc("/user/vote", auth.AdminWithJWTAuth(h.handleVote, h.store)).Methods(POST)
	router.HandleFunc("/user/register", h.handleCreate).Methods(POST)
	router.HandleFunc("/user/login", h.handleLogin).Methods(POST)
}

func (h *Handler) handleGetList(w http.ResponseWriter, r *http.Request) {
	userList, err := h.store.GetList()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err = utils.WriteJSON(w, http.StatusOK, userList); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *Handler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var userPayload types.UserPayload
	if err := utils.ParseJSON(r, &userPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	hashedPassword, err := auth.GetHashedPassword(userPayload.User.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	userPayload.User.Password = hashedPassword

	userRole, err := h.store.GetRoleByLevel(userPayload.UserRole.AccessLevel)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if userRole.ID == 0 {
		if err = utils.WriteJSON(w, http.StatusInternalServerError, "undefined access level"); err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}

		return
	}

	uid, err := h.store.CreateUser(&userPayload.User)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	_, err = h.store.AssigneeUserRole(uid, userRole.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, uid); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var user types.User
	if err := utils.ParseJSON(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	u, err := h.store.GetUserByEmail(user.Email)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid email or password"))
		return
	}

	if !auth.ComparePasswords(user.Password, u.Password) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid email or password"))
		return
	}

	secret := []byte(configs.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, map[string]string{"token": token}); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
}

func (h *Handler) handleEdit(w http.ResponseWriter, r *http.Request) {
	var userPayload types.UserPayload
	vars := mux.Vars(r)
	userID := vars["id"]

	if err := utils.ParseJSON(r, &userPayload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if userPayload.User.Password != "" {
		hashedPassword, err := auth.GetHashedPassword(userPayload.User.Password)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		userPayload.User.Password = hashedPassword
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	userPayload.User.ID = uint(id)
	err = h.store.UpdateUser(&userPayload.User)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.UpdateUserAccessLevel(&userPayload.User, uint(userPayload.UserRole.AccessLevel))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, nil); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
}

func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	if _, err := strconv.Atoi(userID); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid User ID provided"))
		return
	}
	err := h.store.DeleteUserById(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, nil); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
}

func (h *Handler) handleVote(w http.ResponseWriter, r *http.Request) {
	var vote types.Votes

	if err := utils.ParseJSON(r, &vote); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userID := fmt.Sprintf("%v", r.Context().Value(auth.UserKey))
	uid, err := strconv.Atoi(userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	if vote.ProfileID == uint(uid) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("you can't vote yourself"))
		return
	}

	isVoteAvailable, err := h.helper.IsVoteActionAvailable(uint(uid))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if !isVoteAvailable {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("you are not abble to vote right now. Please try again later"))
		return
	}

	voted, err := h.helper.IsAlreadyVoted(uint(uid), vote.ProfileID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}

	if voted {
		if h.helper.IsWithdrawVote(vote) {
			err = h.store.DeleteVote(uint(uid), vote.ProfileID)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			}
		}

		isVoteChange, err := h.helper.IsVoteChange(uint(uid), &vote)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, err)
			return
		}

		if isVoteChange {
			err = h.store.UpdateVote(&vote)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
				return
			}

			if err := utils.WriteJSON(w, http.StatusOK, userID); err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err)
			}

			return
		}

		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("you have already voted this profile"))
		return
	}

	vote.VoteTime = time.Now()
	vote.VoterID = uint(uid)

	_, err = h.store.AssigneeVote(&vote)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, userID); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
	}
}
