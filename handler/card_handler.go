package handler

import (
	"errors"
	"fmt"
	"net/http"

	card_usecase "github.com/flash-cards-vocab/backend/app/usecase/card"
	collection_usecase "github.com/flash-cards-vocab/backend/app/usecase/collection"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type handlerCard struct {
	card_uc card_usecase.UseCase
	apikey  string
}

func NewCardHandler(card_uc card_usecase.UseCase, apikey string) RestCardHandler {
	return &handlerCard{card_uc: card_uc, apikey: apikey}
}

func (h *handlerCard) UploadCardImage(c *gin.Context) {
	// start_time := time.Now()

	req_apikey := c.Request.Header.Get("Apikey")
	fmt.Println(c.Request.Header)
	fmt.Println(req_apikey, "req_apikey")
	fmt.Println(h.apikey, "h.apikey")
	if req_apikey != h.apikey {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "API key is incorrect"})
		return
	}

	file, handler, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Println("Error here 1", file, handler, err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "ERROR HERE1" + err.Error()})
		return
	}
	defer file.Close()
	storage_location := c.Request.FormValue("location")
	fmt.Println(storage_location, "storage_location")
	if handler == nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "Invalid file"})
		return
	}
	if storage_location == "" {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "location is required"})
		return
	}

	cloud_filename, err := h.card_uc.UploadCardImage(file, storage_location, handler.Filename) // Upload(file, storage_location, handler.Filename, random_filename != "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "here" + err.Error()})
		return
	}

	// resp := fmt.Sprintf(`{
	// 	"http_status": 200,
	// 	"process_time": "%.6fms",
	// 	"data": {
	// 		"url": "%s"
	// 	}
	// }`, time.Since(start_time).Seconds()/1000, cloud_filename)

	c.Request.Header.Set("Content-Type", "application/json")
	// c.JSON(http.StatusOK, resp)
	c.JSON(http.StatusOK, SuccessResponse{Data: cloud_filename})

}

func (h *handlerCard) AddExistingCardToCollection(c *gin.Context) {
	paramCardId := c.Param("card_id")
	cardId, err := uuid.Parse(paramCardId)
	if err == nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}
	paramCollectionId := c.Param("collection_id")
	collectionId, err := uuid.Parse(paramCollectionId)
	if err == nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	// userId, err = helpers.GetAuthContext(c)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, ErrorResponse{Message: "User id not found"})
	// }

	err = h.card_uc.AddExistingCardToCollection(collectionId, cardId)
	if err == nil {
		c.JSON(http.StatusOK, SuccessResponse{"Collection Viewed"})
	} else {
		if errors.Is(err, collection_usecase.ErrNotFound) {
			c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		}
	}
}
