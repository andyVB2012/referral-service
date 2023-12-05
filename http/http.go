package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	// "github.com/parsaakbari1209/go-mongo-crud-rest-api/model"
	// "github.com/parsaakbari1209/go-mongo-crud-rest-api/repository"

	"github.com/andyVB2012/referral-service/repository"
)

var (
	ErrCodeNotFound      = errors.New("Referral code not found")
	ErrUserNotFound      = errors.New("User not found")
	ErrAttributionFailed = errors.New("Attribution failed")
	ErrCodeError         = errors.New("Referral code error")
	ErrNoTraderAddr      = errors.New("No trader address")
	ErrInvalidAddress    = errors.New("Invalid address")
)

type Server struct {
	repository repository.Repository
}

func NewServer(repository repository.Repository) *Server {
	return &Server{repository: repository}
}

func (s Server) GetAttributionStats(ctx *gin.Context) {
	address := ctx.Param("address")
	if address == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidAddress})
		return
	}
	code, err := s.repository.GetCode(ctx, address)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrCodeNotFound)
		return
	}
	stats, err := s.repository.GetStats(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": stats})

}

func (s Server) GetAllAttributions(ctx *gin.Context) {
	address := ctx.Param("address")
	if address == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidAddress})
		return
	}
	code, err := s.repository.GetCode(ctx, address)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrCodeNotFound)
		return
	}
	stats, err := s.repository.GetStats(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	attributors, err := s.repository.GetAllAttributors(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	response := AttributionResponse{
		Address:      address,
		ReferralCode: code,
		Stats:        stats,
		Attributors:  attributors,
	}

	// Wrap the response in an array
	results := []AttributionResponse{response}

	// Send the response
	ctx.JSON(http.StatusOK, gin.H{"results": results})

}

func (s Server) GetReferralCode(ctx *gin.Context) {
	address := ctx.Param("address")
	if address == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidAddress})
		return
	}
	code, err := s.repository.GetCode(ctx, address)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrCodeNotFound)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": code})
}

// func (s Server) GetFollow(ctx *gin.Context) {
// 	user1 := ctx.Param("user1")
// 	if user1 == "" {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid argument user1"})
// 		return
// 	}

// 	user2 := ctx.Param("user2")
// 	if user2 == "" {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid argument user2"})
// 		return
// 	}

// 	if user1 == user2 {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user1 == user2"})
// 		return
// 	}

// 	follow, err := s.repository.IsCodeInDb(ctx, user1, user2)
// 	if err != nil {
// 		if errors.Is(err, repository.ErrUserNotFound) {
// 			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{"data": follow})
// }

// func (s Server) GetFollowings(ctx *gin.Context) {
// 	user1 := ctx.Param("user1")
// 	if user1 == "" {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid argument user1"})
// 		return
// 	}

// 	follow, err := s.repository.GetFollowings(ctx, user1)
// 	if err != nil {
// 		if errors.Is(err, repository.ErrUserNotFound) {
// 			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{"data": follow})
// }

// func (s Server) GetFollowers(ctx *gin.Context) {
// 	user2 := ctx.Param("user2")
// 	if user2 == "" {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid argument user2"})
// 		return
// 	}

// 	follow, err := s.repository.GetFollowers(ctx, user2)
// 	if err != nil {
// 		if errors.Is(err, repository.ErrUserNotFound) {
// 			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{"data": follow})
// }

// func (s Server) GetAll(ctx *gin.Context) {
// 	follow, err := s.repository.GetAll(ctx)
// 	if err != nil {
// 		if errors.Is(err, repository.ErrUserNotFound) {
// 			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{"data": follow})
// }

// func (s Server) CreateFollow(ctx *gin.Context) {
// 	var follow model.Follow
// 	if err := ctx.ShouldBindJSON(&follow); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
// 		return
// 	}
// 	followOut, err := s.repository.CreateFollow(ctx, follow)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{"data": followOut})
// }

// func (s Server) CreateFollowBatch(ctx *gin.Context) {
// 	var follows []model.Follow
// 	var followsOut []model.Follow
// 	if err := ctx.ShouldBindJSON(&follows); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
// 		return
// 	}
// 	followsOut, err := s.repository.CreateFollowBatch(ctx, follows)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{"data": followsOut})
// }

// func (s Server) DeleteFollow(ctx *gin.Context) {
// 	user1 := ctx.Param("user1")
// 	if user1 == "" {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid argument user1"})
// 		return
// 	}

// 	user2 := ctx.Param("user2")
// 	if user2 == "" {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid argument user2"})
// 		return
// 	}

// 	if user1 == user2 {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user1 == user2"})
// 		return
// 	}

// 	if err := s.repository.DeleteFollow(ctx, user1, user2); err != nil {
// 		if errors.Is(err, repository.ErrUserNotFound) {
// 			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{})
// }
