package controller

import (
	"encoding/json"
	"net/http"

	"github.com/CESARBR/knot-gateway-webui/backend/pkg/interactors"
	"github.com/CESARBR/knot-gateway-webui/backend/pkg/logging"
)

// StateController represents the controller for state capabilities
type StateController struct {
	updateStateInteractor *interactors.UpdateStateInteractor
	getStateInteractor    *interactors.GetStateInteractor
	logger                logging.Logger
}

// StateData represents the incoming state data from update request
type StateData struct {
	State string `json:"state"`
}

// NewStateController creates a new StateController instance
func NewStateController(
	updateStateInteractor *interactors.UpdateStateInteractor,
	getStateInteractor *interactors.GetStateInteractor,
	logger logging.Logger,
) *StateController {
	return &StateController{updateStateInteractor, getStateInteractor, logger}
}

// Update calls the use case to update the state
func (sc *StateController) Update(w http.ResponseWriter, r *http.Request) {
	var err error
	var stateData StateData

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&stateData)
	if err != nil {
		sc.logger.Error("Invalid request format")
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	err = sc.updateStateInteractor.Execute(stateData.State)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(stateData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		return
	}
}

// Get calls the use case to get the current state
func (sc *StateController) Get(w http.ResponseWriter, r *http.Request) {
	state, err := sc.getStateInteractor.Execute()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	stateResponse := &StateData{state.Type}
	response, err := json.Marshal(stateResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		return
	}
}
