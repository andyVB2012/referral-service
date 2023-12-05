package http

import (
	"errors"
	"net/http"

	"github.com/andyVB2012/referral-service/repository"
	"github.com/gin-gonic/gin"
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

func (s Server) HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
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
	response := StatsResponse{
		Address:      address,
		ReferralCode: code,
		Stats:        stats,
	}
	ctx.JSON(http.StatusOK, gin.H{"data": response})

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

func (s Server) CreateReferralCode(ctx *gin.Context) {
	address := ctx.Param("address")
	if address == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidAddress})
		return
	}

	code, err := s.repository.CreateReferralCode(ctx, address)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyAdded) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, ErrCodeError)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": code})
}

func (s Server) AddAttributor(ctx *gin.Context) {

	var req AttributionRequest

	// Bind the JSON payload to the struct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Validate the request data
	if req.Address == "" || req.Message == "" || req.Signature == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "address, message, and signature are required"})
		return
	}
	code := ctx.Param("code")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidAddress})
		return
	}

	check, err := VerifySignature(req.Message, req.Signature, req.Address)
	if err != nil {
		if err.Error() == "invalid signature" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, "signature verification failed")
		return
	}
	if check == false {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
		return
	}

	err1 := s.repository.AddAttributor(ctx, code, req.Address)
	if err1 != nil {
		ctx.JSON(http.StatusInternalServerError, ErrAttributionFailed)
		return
	}

}
