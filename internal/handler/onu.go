package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/achyar10/snmp-olt-zte/internal/model"
	"github.com/achyar10/snmp-olt-zte/internal/usecase"
	"github.com/achyar10/snmp-olt-zte/internal/utils"
	"github.com/achyar10/snmp-olt-zte/pkg/pagination"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type OnuHandlerInterface interface {
	GetByBoardIDAndPonID(w http.ResponseWriter, r *http.Request)
	GetByBoardIDPonIDAndOnuID(w http.ResponseWriter, r *http.Request)
	GetEmptyOnuID(w http.ResponseWriter, r *http.Request)
	GetOnuIDAndSerialNumber(w http.ResponseWriter, r *http.Request)
	UpdateEmptyOnuID(w http.ResponseWriter, r *http.Request)
	GetByBoardIDAndPonIDWithPaginate(w http.ResponseWriter, r *http.Request)
	ActivateONU(w http.ResponseWriter, r *http.Request)
	GetUnactivatedONU(w http.ResponseWriter, r *http.Request)
}

type OnuHandler struct {
	ponUsecase usecase.OnuUseCaseInterface
}

func NewOnuHandler(ponUsecase usecase.OnuUseCaseInterface) *OnuHandler {
	return &OnuHandler{ponUsecase: ponUsecase}
}

func (o *OnuHandler) GetByBoardIDAndPonID(w http.ResponseWriter, r *http.Request) {

	boardID := chi.URLParam(r, "board_id") // 1 or 2
	ponID := chi.URLParam(r, "pon_id")     // 1 - 8

	boardIDInt, err := strconv.Atoi(boardID) // convert string to int

	log.Info().Msg("Received a request to GetByBoardIDAndPonID")

	// Validate boardIDInt value and return error 400 if boardIDInt is not 1 or 2
	if err != nil || (boardIDInt != 1 && boardIDInt != 2) {
		log.Error().Err(err).Msg("Invalid 'board_id' parameter")
		utils.ErrorBadRequest(w, fmt.Errorf("invalid 'board_id' parameter. It must be 1 or 2")) // error 400
		return
	}

	ponIDInt, err := strconv.Atoi(ponID) // convert string to int

	// Validate ponIDInt value and return error 400 if ponIDInt is not between 1 and 8
	if err != nil || ponIDInt < 1 || ponIDInt > 16 {
		log.Error().Err(err).Msg("Invalid 'pon_id' parameter")
		utils.ErrorBadRequest(w, fmt.Errorf("invalid 'pon_id' parameter. It must be between 1 and 16")) // error 400
		return
	}

	query := r.URL.Query() // Get query parameters from the request

	log.Debug().Interface("query_parameters", query).Msg("Received query parameters")

	//Validate query parameters and return error 400 if query parameters is not "onu_id" or empty query parameters
	if len(query) > 0 && query["onu_id"] == nil {
		log.Error().Msg("Invalid query parameter")
		utils.ErrorBadRequest(w, fmt.Errorf("invalid query parameter")) // error 400
		return
	}

	// Call usecase to get data from SNMP
	onuInfoList, err := o.ponUsecase.GetByBoardIDAndPonID(r.Context(), boardIDInt, ponIDInt)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get data from SNMP")
		utils.ErrorInternalServerError(w, fmt.Errorf("cannot get data from snmp")) // error 500
		return
	}

	log.Info().Msg("Successfully retrieved data from SNMP")

	/*
		Validate onuInfoList value
		If onuInfoList is empty, return error 404
	*/

	if len(onuInfoList) == 0 {
		log.Warn().Msg("Data not found")
		utils.ErrorNotFound(w, fmt.Errorf("data not found")) // error 404
		return
	}

	// Convert result to JSON format according to WebResponse structure
	response := utils.WebResponse{
		Code:   http.StatusOK, // 200
		Status: "OK",          // "OK"
		Data:   onuInfoList,   // data
	}

	utils.SendJSONResponse(w, http.StatusOK, response) // 200

}

func (o *OnuHandler) GetByBoardIDPonIDAndOnuID(w http.ResponseWriter, r *http.Request) {

	boardID := chi.URLParam(r, "board_id") // 1 or 2
	ponID := chi.URLParam(r, "pon_id")     // 1 - 8
	onuID := chi.URLParam(r, "onu_id")     // 1 - 128

	boardIDInt, err := strconv.Atoi(boardID) // convert string to int

	log.Info().Msg("Received a request to GetByBoardIDPonIDAndOnuID")

	// Validate boardIDInt value and return error 400 if boardIDInt is not 1 or 2
	if err != nil || (boardIDInt != 1 && boardIDInt != 2) {
		log.Error().Err(err).Msg("Invalid 'board_id' parameter")
		utils.ErrorBadRequest(w, fmt.Errorf("invalid 'board_id' parameter. It must be 1 or 2")) // error 400
		return
	}

	ponIDInt, err := strconv.Atoi(ponID) // convert string to int

	// Validate ponIDInt value and return error 400 if ponIDInt is not between 1 and 8
	if err != nil || ponIDInt < 1 || ponIDInt > 16 {
		log.Error().Err(err).Msg("Invalid 'pon_id' parameter")
		utils.ErrorBadRequest(w, fmt.Errorf("invalid 'pon_id' parameter. It must be between 1 and 16")) // error 400
		return
	}

	onuIDInt, err := strconv.Atoi(onuID) // convert string to int

	// Validate onuIDInt value and return error 400 if onuIDInt is not between 1 and 128
	if err != nil || onuIDInt < 1 || onuIDInt > 128 {
		log.Error().Err(err).Msg("Invalid 'onu_id' parameter")
		utils.ErrorBadRequest(w, fmt.Errorf("invalid 'onu_id' parameter. It must be between 1 and 128")) // error 400
		return
	}

	// Call usecase to get data from SNMP
	onuInfoList, err := o.ponUsecase.GetByBoardIDPonIDAndOnuID(boardIDInt, ponIDInt, onuIDInt)

	if err != nil {
		log.Error().Err(err).Msg("Failed to get data from SNMP")
		utils.ErrorInternalServerError(w, fmt.Errorf("cannot get data from snmp")) // error 500
		return
	}

	log.Info().Msg("Successfully retrieved data from SNMP")

	/*
		Validate onuInfoList value
		If onuInfoList.Board, onuInfoList.PON, and onuInfoList.ID is 0, return error 404
		example: http://localhost:8080/board/1/pon/1/onu/129
	*/

	if onuInfoList.Board == 0 && onuInfoList.PON == 0 && onuInfoList.ID == 0 {
		log.Error().Msg("Data not found")
		utils.ErrorNotFound(w, fmt.Errorf("data not found")) // error 404
		return
	}

	// Convert a result to JSON format according to WebResponse structure
	response := utils.WebResponse{
		Code:   http.StatusOK, // 200
		Status: "OK",          // "OK"
		Data:   onuInfoList,   // data
	}

	utils.SendJSONResponse(w, http.StatusOK, response) // 200
}

func (o *OnuHandler) GetEmptyOnuID(w http.ResponseWriter, r *http.Request) {

	boardID := chi.URLParam(r, "board_id") // 1 or 2
	ponID := chi.URLParam(r, "pon_id")     // 1 - 8

	boardIDInt, err := strconv.Atoi(boardID) // convert string to int

	log.Info().Msg("Received a request to GetEmptyOnuID")

	// Validate boardIDInt value and return error 400 if boardIDInt is not 1 or 2
	if err != nil || (boardIDInt != 1 && boardIDInt != 2) {
		log.Error().Err(err).Msg("Invalid 'board_id' parameter")
		utils.ErrorBadRequest(w, fmt.Errorf("invalid 'board_id' parameter. It must be 1 or 2")) // error 400
		return
	}

	ponIDInt, err := strconv.Atoi(ponID) // convert string to int

	// Validate ponIDInt value and return error 400 if ponIDInt is not between 1 and 8
	if err != nil || ponIDInt < 1 || ponIDInt > 16 {
		log.Error().Err(err).Msg("Invalid 'pon_id' parameter")
		utils.ErrorBadRequest(w, fmt.Errorf("invalid 'pon_id' parameter. It must be between 1 and 16")) // error 400
		return
	}

	// Call usecase to get data from SNMP
	onuIDEmptyList, err := o.ponUsecase.GetEmptyOnuID(r.Context(), boardIDInt, ponIDInt)

	if err != nil {
		log.Error().Err(err).Msg("Failed to get data from SNMP")
		utils.ErrorInternalServerError(w, fmt.Errorf("cannot get data from snmp")) // error 500
		return
	}

	log.Info().Msg("Successfully retrieved data from SNMP")

	// Convert result to JSON format according to WebResponse structure
	response := utils.WebResponse{
		Code:   http.StatusOK,  // 200
		Status: "OK",           // "OK"
		Data:   onuIDEmptyList, // data
	}

	utils.SendJSONResponse(w, http.StatusOK, response) // 200
}

func (o *OnuHandler) GetOnuIDAndSerialNumber(w http.ResponseWriter, r *http.Request) {

	boardID := chi.URLParam(r, "board_id") // 1 or 2
	ponID := chi.URLParam(r, "pon_id")     // 1 - 8

	boardIDInt, err := strconv.Atoi(boardID) // convert string to int

	log.Info().Msg("Received a request to GetOnuSerialNumber")

	// Validate boardIDInt value and return error 400 if boardIDInt is not 1 or 2
	if err != nil || (boardIDInt != 1 && boardIDInt != 2) {
		log.Error().Err(err).Msg("Invalid 'board_id' parameter")
		utils.ErrorBadRequest(w, fmt.Errorf("invalid 'board_id' parameter. It must be 1 or 2")) // error 400
		return
	}

	ponIDInt, err := strconv.Atoi(ponID) // convert string to int

	// Validate ponIDInt value and return error 400 if ponIDInt is not between 1 and 8
	if err != nil || ponIDInt < 1 || ponIDInt > 16 {
		log.Error().Err(err).Msg("Invalid 'pon_id' parameter")
		utils.ErrorBadRequest(w, fmt.Errorf("invalid 'pon_id' parameter. It must be between 1 and 16")) // error 400
		return
	}

	// Call usecase to get Serial Number from SNMP
	onuSerialNumber, err := o.ponUsecase.GetOnuIDAndSerialNumber(boardIDInt, ponIDInt)

	if err != nil {
		log.Error().Err(err).Msg("Failed to get data from SNMP")
		utils.ErrorInternalServerError(w, fmt.Errorf("cannot get data from snmp")) // error 500
		return
	}

	log.Info().Msg("Successfully retrieved data from SNMP")

	// Convert a result to JSON format according to WebResponse structure
	response := utils.WebResponse{
		Code:   http.StatusOK,   // 200
		Status: "OK",            // "OK"
		Data:   onuSerialNumber, // data
	}

	utils.SendJSONResponse(w, http.StatusOK, response) // 200
}

func (o *OnuHandler) UpdateEmptyOnuID(w http.ResponseWriter, r *http.Request) {
	boardID := chi.URLParam(r, "board_id") // 1 or 2
	ponID := chi.URLParam(r, "pon_id")     // 1 - 8

	boardIDInt, err := strconv.Atoi(boardID) // convert string to int

	log.Info().Msg("Received a request to UpdateEmptyOnuID")

	// Validate boardIDInt value and return error 400 if boardIDInt is not 1 or 2
	if err != nil || (boardIDInt != 1 && boardIDInt != 2) {
		log.Error().Err(err).Msg("Invalid 'board_id' parameter")
		utils.ErrorBadRequest(w, fmt.Errorf("invalid 'board_id' parameter. It must be 0 or 1")) // error 400
		return
	}

	ponIDInt, err := strconv.Atoi(ponID) // convert string to int

	// Validate ponIDInt value and return error 400 if ponIDInt is not between 1 and 8
	if err != nil || ponIDInt < 1 || ponIDInt > 16 {
		log.Error().Err(err).Msg("Invalid 'pon_id' parameter")
		utils.ErrorBadRequest(w, fmt.Errorf("invalid 'pon_id' parameter. It must be between 1 and 16")) // error 400
		return
	}

	// Call usecase to get data from SNMP
	err = o.ponUsecase.UpdateEmptyOnuID(r.Context(), boardIDInt, ponIDInt)

	if err != nil {
		log.Error().Err(err).Msg("Failed to get data from SNMP")
		utils.ErrorInternalServerError(w, fmt.Errorf("cannot get data from snmp")) // error 500
		return
	}

	log.Info().Msg("Successfully retrieved data from SNMP")

	// Convert result to JSON format according to WebResponse structure
	response := utils.WebResponse{
		Code:   http.StatusOK,                 // 200
		Status: "OK",                          // "OK"
		Data:   "Success Update Empty ONU_ID", // data
	}

	utils.SendJSONResponse(w, http.StatusOK, response) // 200
}

func (o *OnuHandler) GetByBoardIDAndPonIDWithPaginate(w http.ResponseWriter, r *http.Request) {

	boardID := chi.URLParam(r, "board_id") // 1 or 2
	ponID := chi.URLParam(r, "pon_id")     // 1 - 8

	// Get page and page size parameters from the request
	pageIndex, pageSize := pagination.GetPaginationParametersFromRequest(r)

	boardIDInt, err := strconv.Atoi(boardID) // convert string to int

	log.Info().Msg("Received a request to GetByBoardIDAndPonIDWithPaginate")

	// Validate boardIDInt value and return error 400 if boardIDInt is not 1 or 2
	if err != nil || (boardIDInt != 1 && boardIDInt != 2) {
		log.Error().Err(err).Msg("Invalid 'board_id' parameter")
		utils.ErrorBadRequest(w, fmt.Errorf("invalid 'board_id' parameter. It must be 1 or 2")) // error 400
		return
	}

	ponIDInt, err := strconv.Atoi(ponID) // convert string to int

	// Validate ponIDInt value and return error 400 if ponIDInt is not between 1 and 8
	if err != nil || ponIDInt < 1 || ponIDInt > 16 {
		log.Error().Err(err).Msg("Invalid 'pon_id' parameter")
		utils.ErrorBadRequest(w, fmt.Errorf("invalid 'pon_id' parameter. It must be between 1 and 16")) // error 400
		return
	}

	item, count := o.ponUsecase.GetByBoardIDAndPonIDWithPagination(boardIDInt, ponIDInt, pageIndex,
		pageSize)

	/*
		Validate item value
		If item is empty, return error 404
	*/

	if len(item) == 0 {
		log.Error().Msg("Data not found")
		utils.ErrorNotFound(w, fmt.Errorf("data not found")) // error 404
		return
	}

	// Convert result to JSON format according to Pages structure
	pages := pagination.New(pageIndex, pageSize, count)

	// Convert result to JSON format according to WebResponse structure
	responsePagination := pagination.Pages{
		Code:      http.StatusOK,   // 200
		Status:    "OK",            // "OK"
		Page:      pages.Page,      // page
		PageSize:  pages.PageSize,  // page size
		PageCount: pages.PageCount, // page count
		TotalRows: pages.TotalRows, // total rows
		Data:      item,            // data
	}

	utils.SendJSONResponse(w, http.StatusOK, responsePagination) // 200
}

func (o *OnuHandler) ActivateONU(w http.ResponseWriter, r *http.Request) {
	var payload model.ActivateONURequest

	// Decode body
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Error().Err(err).Msg("Invalid JSON payload")
		utils.ErrorBadRequest(w, fmt.Errorf("invalid request payload"))
		return
	}

	// Validate required fields
	if payload.OLTIndex == "" || payload.SerialNumber == "" || payload.Region == "" || payload.Code == "" {
		log.Error().Msg("Missing required fields in payload")
		utils.ErrorBadRequest(w, fmt.Errorf("missing required fields"))
		return
	}

	slot, port, err := utils.ParseOltIndex(payload.OLTIndex)
	if err != nil {
		log.Error().Err(err).Msg("Invalid OLTIndex format")
		utils.ErrorBadRequest(w, fmt.Errorf("invalid olt_index format"))
		return
	}

	// Get ONU ID: dari payload atau dari fungsi available
	var onuID int
	if payload.Onu != nil {
		onuID = *payload.Onu
	} else {
		available, err := utils.GetAvailableONUOnly(payload.OLTIndex, 128)
		if err != nil || len(available) == 0 {
			log.Error().Err(err).Msg("No available ONU found")
			utils.ErrorInternalServerError(w, fmt.Errorf("no available ONU found"))
			return
		}
		onuID = available[0].ID
	}

	// Build command
	cmd := utils.BuildZTERegisterCommand(slot, port, payload.Region, payload.SerialNumber, payload.Code, onuID, payload.VlanID)

	// Run telnet command
	resp, err := utils.RunTelnetCommand(cmd)
	if err != nil {
		log.Error().Err(err).Msg("Activation failed via Telnet")
		utils.ErrorInternalServerError(w, fmt.Errorf("activation failed"))
		return
	}

	// Detect re-registration / duplicate
	isAlreadyExists := strings.Contains(resp, "entry is existed") ||
		strings.Contains(resp, "already exists") ||
		strings.Contains(resp, "The service is already existed")

	if isAlreadyExists {
		log.Warn().Msg("ONU already registered")
		utils.SendJSONResponse(w, http.StatusConflict, utils.WebResponse{
			Code:   http.StatusConflict,
			Status: "already_registered",
			Data: map[string]interface{}{
				"used_onu":       onuID,
				"olt_index":      payload.OLTIndex,
				"serial_number":  payload.SerialNumber,
				"command_output": resp,
			},
		})
		return
	}

	// Return success response
	response := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data: map[string]interface{}{
			"status":         "success",
			"used_onu":       onuID,
			"olt_index":      payload.OLTIndex,
			"serial_number":  payload.SerialNumber,
			"command_output": resp,
		},
	}

	utils.SendJSONResponse(w, http.StatusOK, response)
}

func (o *OnuHandler) GetUnactivatedONU(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	output, err := utils.RunTelnetCommand("show pon onu u")
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute Telnet command")
		utils.ErrorInternalServerError(w, fmt.Errorf("failed to execute telnet command: %v", err))
		return
	}

	onuItems := utils.ParseONULineOutput(output)
	duration := time.Since(start).Seconds()

	response := utils.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data: map[string]interface{}{
			"duration":     fmt.Sprintf("%.2fs", duration),
			"detected_onu": onuItems,
		},
	}

	utils.SendJSONResponse(w, http.StatusOK, response)
}
