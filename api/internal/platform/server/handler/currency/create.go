package currency

import (
	bol "boletia/api/internal"
	"boletia/api/internal/obtain"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CurrencyRequest struct {
	Currency string `uri:"currency" binding:"required"`
	Finit    string `form:"finit"`
	Fend     string `form:"fend"`
}

// ObtainHandler returns an HTTP handler for currency creation.
func ObtainHandler(getCurrency obtain.ICurrencyService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req CurrencyRequest
		if err := ctx.ShouldBindUri(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := ctx.ShouldBindQuery(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var criteria bol.Criteria
		newCriteria, err := criteria.NewCriteria(req.Currency, req.Finit, req.Fend)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		data, err := getCurrency.GetCurrency(ctx, newCriteria)
		if err != nil {
			switch {
			case errors.Is(err, bol.ErrorNotFound):
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			case errors.Is(err, bol.ErrorCurrencyCode):
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		ctx.JSON(http.StatusOK, data)
	}
}
