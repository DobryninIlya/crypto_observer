package handlers

import (
	"cryptoObserver/internal/app/store/sqlstore"
	"cryptoObserver/internal/app/store/sqlstore/utils"
	worker "cryptoObserver/internal/app/workers"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// NewAddCurrencyHandler godoc
//
// @Summary Добавление валюты
// @Description Добавление валюты в список валют для отслеживания.
// @Tags currency
//
//	@Accept			multipart/form-data
//
// @Produce json
// @Param currencyID	formData	string	true	"Логин"
// @Success 200 {object} string "OK - Currency added successfully"
// @Failure 400 {object} string "Bad Request - Currency ID is required"
// @Router /currency/add [post]
func NewAddCurrencyHandler(log *logrus.Logger, store sqlstore.CurrencyInterface, pool *worker.WorkerPool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.addCurrency.NewAddCurrencyHandler"
		currencyID := strings.TrimSpace(r.FormValue("currencyID"))
		if currencyID == "" {
			log.WithFields(logrus.Fields{
				"path": path,
			}).Error("Currency ID is required")
			utils.Respond(w, r, http.StatusBadRequest, "Currency ID is required")
			return
		}
		pool.AddCurrency(currencyID)
		err := store.AddCurrency(currencyID)
		if err != nil {
			log.WithFields(logrus.Fields{
				"path":  path,
				"error": err.Error(),
			}).Error("Failed to add currency to store")
			pool.RemoveCurrency(currencyID)
			utils.Respond(w, r, http.StatusInternalServerError, "Failed to add currency: "+err.Error())
			return
		}
		log.WithFields(logrus.Fields{
			"path":       path,
			"currencyID": currencyID,
		}).Info("Currency added successfully")
		utils.Respond(w, r, http.StatusOK, "Currency added successfully")

	}
}
