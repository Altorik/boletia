package currency

import (
	"boletia/internal/obtain"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ObtainHandler returns an HTTP handler for currency creation.
func ObtainHandler(getCurrency obtain.ICurrencyService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := getCurrency.GetCurrency(ctx)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		ctx.Status(http.StatusCreated)
	}
}
