package handler

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"

	"github.com/smolse/fluffy-pancake/internal/models"
	"github.com/smolse/fluffy-pancake/internal/service"
)

// Handler is a struct for handling HTTP requests.
type Handler struct {
	Service *service.RiskService
}

// NewHandler is a constructor function that creates a new Handler instance with the given risk service instance.
func NewHandler(svc *service.RiskService) *Handler {
	return &Handler{
		Service: svc,
	}
}

// GetRisk is a handler function that gets a risk by ID.
func (h *Handler) GetRisk(c *gin.Context) {
	// Extract and validate the risk ID from the request
	id, err := models.NewRiskIdFromString(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the risk from the service
	risk, err := h.Service.GetRisk(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"risk": risk})
}

// CreateRisk is a handler function that creates a risk.
func (h *Handler) CreateRisk(c *gin.Context) {
	riskAttributes := models.RiskAttributes{}

	// Bind the request body to the risk attributes
	if err := c.BindJSON(&riskAttributes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request body: %s", err.Error())})
		return
	}

	// Make sure the risk state is valid
	if !slices.Contains(models.ValidRiskStates, riskAttributes.State) {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid risk state: %s", riskAttributes.State)})
		return
	}

	id := models.NewRiskId()
	if err := h.Service.CreateRisk(id, riskAttributes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id.String()})
}

// ListRisks is a handler function that lists all risks.
func (h *Handler) ListRisks(c *gin.Context) {
	risks, err := h.Service.ListRisks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"risks": risks})
}
