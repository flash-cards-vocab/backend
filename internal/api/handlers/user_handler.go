package handlers

import (
	"errors"
	"net/http"

	user_usecase "github.com/flash-cards-vocab/backend/app/usecase/user"
	"github.com/flash-cards-vocab/backend/entity"
	"github.com/flash-cards-vocab/backend/internal/api/handler_interfaces"
	"github.com/gin-gonic/gin"
)

type handlerUser struct {
	user_uc user_usecase.UseCase
}

func NewUserHandler(user_uc user_usecase.UseCase) handler_interfaces.RestUserHandler {
	return &handlerUser{user_uc: user_uc}
}

func (h *handlerUser) Register(c *gin.Context) {
	// paramId := c.Param("user_id")
	// user_id, err := uuid.Parse(paramId)
	// if err == nil {
	// 	c.JSON(http.StatusInternalServerError, handler_interfaces.ErrorResponse{Message: err.Error()})
	// }
	var newUserData entity.User
	err := c.ShouldBindJSON(&newUserData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handler_interfaces.ErrorResponse{Message: err.Error()})
	}

	data, err := h.user_uc.Register(newUserData)
	if err != nil {
		if errors.Is(err, user_usecase.ErrNotFound) {
			c.JSON(http.StatusNotFound, handler_interfaces.ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, handler_interfaces.ErrorResponse{Message: err.Error()})
		}
	}
	c.JSON(http.StatusOK, handler_interfaces.SuccessResponse{Data: data})
}

func (h *handlerUser) Login(c *gin.Context) {
	// paramId := c.Param("user_id")
	// user_id, err := uuid.Parse(paramId)
	// if err == nil {
	// 	c.JSON(http.StatusInternalServerError, handler_interfaces.ErrorResponse{Message: err.Error()})
	// }
	////////////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////////
	// We can obtain the session token from the requests cookies, which come with every request
	// cookie, ok := c.Request.Header["Authorization"]
	// if !ok {
	// 	// if err == http.ErrNoCookie {
	// 	// 	// If the cookie is not set, return an unauthorized status
	// 	// 	// w.WriteHeader(http.StatusUnauthorized)
	// 	// 	// return
	// 	// 	c.JSON(http.StatusInternalServerError, handler_interfaces.ErrorResponse{Message: err.Error()})
	// 	// }
	// 	// For any other type of error, return a bad request status
	// 	// w.WriteHeader(http.StatusBadRequest)
	// 	// return
	// 	c.JSON(http.StatusInternalServerError, handler_interfaces.ErrorResponse{Message: "unaithorized"})
	// }

	// // Get the JWT string from the cookie
	// tknStr := cookie[0]

	// // Initialize a new instance of `Claims`
	// // claims := &jwt.Claims{}

	// // Parse the JWT string and store the result in `claims`.
	// // Note that we are passing the key in this method as well. This method will return an error
	// // if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// // or if the signature does not match
	// tkn, err := jwt.Parse(tknStr, func(token *jwt.Token) (interface{}, error) {
	// 	return "TODO: REPLACE THIS STRING WITH LEGIT SECRET CODE", nil
	// })
	// // ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
	// // 	return "TODO: REPLACE THIS STRING WITH LEGIT SECRET CODE", nil
	// // })
	// if err != nil {
	// 	if err == jwt.ErrSignatureInvalid {
	// 		c.JSON(http.StatusInternalServerError, handler_interfaces.ErrorResponse{Message: "Unauthorized"})
	// 		// w.WriteHeader(http.StatusUnauthorized)
	// 		// return
	// 	}
	// 	c.JSON(http.StatusInternalServerError, handler_interfaces.ErrorResponse{Message: "Bad request"})
	// 	// w.WriteHeader(http.StatusBadRequest)
	// 	// return
	// }
	// if !tkn.Valid {
	// 	c.JSON(http.StatusInternalServerError, handler_interfaces.ErrorResponse{Message: "Unauthorized"})
	// 	// w.WriteHeader(http.StatusUnauthorized)
	// 	// return
	// }

	////////////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////////

	////////////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////////
	////////////////////////////////////////////////////////////////////////////////

	var loginUserData entity.UserLogin
	err := c.ShouldBindJSON(&loginUserData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handler_interfaces.ErrorResponse{Message: err.Error()})
	}

	data, err := h.user_uc.Login(loginUserData)
	if err != nil {
		if errors.Is(err, user_usecase.ErrNotFound) {
			c.JSON(http.StatusNotFound, handler_interfaces.ErrorResponse{Message: err.Error()})
		} else if errors.Is(err, user_usecase.ErrUserPasswordMismatch) {
			c.JSON(http.StatusForbidden, handler_interfaces.ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, handler_interfaces.ErrorResponse{Message: err.Error()})
		}
	}
	c.JSON(http.StatusOK, handler_interfaces.SuccessResponse{Data: data})
}
