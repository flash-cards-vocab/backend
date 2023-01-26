package handlers

import (
	"errors"
	"net/http"

	userUC "github.com/flash-cards-vocab/backend/app/usecase/user"
	"github.com/flash-cards-vocab/backend/entity"
	handlerIntf "github.com/flash-cards-vocab/backend/internal/api/handler_interfaces"
	"github.com/flash-cards-vocab/backend/pkg/helpers"
	"github.com/gin-gonic/gin"
)

type handlerUser struct {
	userUsecase userUC.UseCase
}

func NewUserHandler(userUsecase userUC.UseCase) handlerIntf.RestUserHandler {
	return &handlerUser{userUsecase: userUsecase}
}

func (h *handlerUser) GetProfile(c *gin.Context) {
	userCtx, err := helpers.GetAuthContext(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "User id not found"})
	}

	data, err := h.userUsecase.GetProfile(userCtx.UserId)
	if err != nil {
		if errors.Is(err, userUC.ErrNotFound) {
			c.JSON(http.StatusNotFound, handlerIntf.ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
		}
	}
	c.JSON(http.StatusOK, handlerIntf.SuccessResponse{Result: data})
}

func (h *handlerUser) UsernameExists(c *gin.Context) {
	username := c.Param("username")

	exists, err := h.userUsecase.UsernameExists(username)
	if err != nil {
		if errors.Is(err, userUC.ErrNotFound) {
			c.JSON(http.StatusNotFound, handlerIntf.ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
		}
	}
	if exists {
		c.JSON(http.StatusNotAcceptable, handlerIntf.SuccessResponse{Result: exists})
	} else {
		c.JSON(http.StatusOK, handlerIntf.SuccessResponse{Result: exists})
	}
}

func (h *handlerUser) Register(c *gin.Context) {
	var newUserData entity.UserRegistration
	err := c.ShouldBindJSON(&newUserData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
	}

	data, err := h.userUsecase.Register(newUserData)
	if err != nil {
		if errors.Is(err, userUC.ErrNotFound) {
			c.JSON(http.StatusNotFound, handlerIntf.ErrorResponse{Message: err.Error()})
		} else if errors.Is(err, userUC.ErrUserExistsAlready) {
			c.JSON(400, handlerIntf.ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
		}
	}
	c.JSON(http.StatusOK, handlerIntf.SuccessResponse{Result: data})
}

func (h *handlerUser) Login(c *gin.Context) {
	// paramId := c.Param("user_id")
	// user_id, err := uuid.Parse(paramId)
	// if err == nil {
	// 	c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
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
	// 	// 	c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
	// 	// }
	// 	// For any other type of error, return a bad request status
	// 	// w.WriteHeader(http.StatusBadRequest)
	// 	// return
	// 	c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "unaithorized"})
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
	// 		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "Unauthorized"})
	// 		// w.WriteHeader(http.StatusUnauthorized)
	// 		// return
	// 	}
	// 	c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "Bad request"})
	// 	// w.WriteHeader(http.StatusBadRequest)
	// 	// return
	// }
	// if !tkn.Valid {
	// 	c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: "Unauthorized"})
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
		c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
	}

	data, err := h.userUsecase.Login(loginUserData)
	if err != nil {
		if errors.Is(err, userUC.ErrNotFound) {
			c.JSON(http.StatusNotFound, handlerIntf.ErrorResponse{Message: err.Error()})
		} else if errors.Is(err, userUC.ErrUserPasswordMismatch) {
			c.JSON(http.StatusForbidden, handlerIntf.ErrorResponse{Message: err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, handlerIntf.ErrorResponse{Message: err.Error()})
		}
	}
	c.JSON(http.StatusOK, handlerIntf.SuccessResponse{Result: data})
}
