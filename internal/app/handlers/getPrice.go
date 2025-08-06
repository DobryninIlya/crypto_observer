package handlers

import (
	"cryptoObserver/internal/app/store/sqlstore"
	"cryptoObserver/internal/app/store/sqlstore/utils"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

// NewGetPriceHandler godoc
//
// @Summary Получение цены валюты
// @Description Получение цены валюты по ID
// @Tags currency
// @Accept multipart/form-data
// @Produce json
// @Param currencyID formData string true "ID валюты"
// @Param timestamp formData string true "timestamp"
// @Success 200 {object} model.Decimal "Price of the currency"
// @Failure 400 {object} string "Bad Request - Currency ID is required"
// @Failure 400 {object} string "Bad Request - Timestamp is required"
// @Failure 404 {object} string "Not Found - No price found for the given currency and timestamp"
// @Router /currency/price [post]
func NewGetPriceHandler(log *logrus.Logger, store sqlstore.CurrencyInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.getPrice.NewGetPriceHandler"
		currencyID := strings.TrimSpace(r.FormValue("currencyID"))
		if currencyID == "" {
			log.WithFields(logrus.Fields{
				"path": path,
			}).Error("Currency ID is required")
			utils.Respond(w, r, http.StatusBadRequest, "Currency ID is required")
			return
		}
		timestamp := strings.TrimSpace(r.FormValue("timestamp"))
		if timestamp == "" {
			log.WithFields(logrus.Fields{
				"path": path,
			}).Error("Timestamp is required")
			utils.Respond(w, r, http.StatusBadRequest, "Timestamp is required")
			return
		}
		timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			log.WithFields(logrus.Fields{
				"path":  path,
				"error": err.Error(),
			}).Error("Invalid timestamp format")
			utils.Respond(w, r, http.StatusBadRequest, "Invalid timestamp format: "+err.Error())
			return
		}
		result, err := store.GetPrice(currencyID, timestampInt)
		if err != nil {
			log.WithFields(logrus.Fields{
				"path":  path,
				"error": err.Error(),
			}).Error("Failed to get price from store")
			utils.Respond(w, r, http.StatusInternalServerError, "Failed to get price from store: "+err.Error())
			return
		}
		if result.IsZero() {
			log.WithFields(logrus.Fields{
				"path": path,
			}).Warn("No price found for the given currency and timestamp")
			utils.Respond(w, r, http.StatusNotFound, "No price found for the given currency and timestamp")
			return
		}
		utils.Respond(w, r, http.StatusOK, result)

	}
}
