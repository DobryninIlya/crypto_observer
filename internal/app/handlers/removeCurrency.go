package handlers

import (
	"cryptoObserver/internal/app/store/sqlstore"
	"cryptoObserver/internal/app/store/sqlstore/utils"
	worker "cryptoObserver/internal/app/workers"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// NewRemoveCurrencyHandler godoc
//
// @Summary Удаление валюты
// @Description Удаление валюты в список валют для отслеживания.
// @Tags currency
// @Accept multipart/form-data
// @Produce json
// @Param currencyID formData string true "ID валюты"
// @Success 200 {object} string "OK - Currency removed successfully"
// @Failure 400 {object} string "Bad Request - Currency ID is required"
// @Router /currency/remove [delete]
func NewRemoveCurrencyHandler(log *logrus.Logger, store sqlstore.CurrencyInterface, pool *worker.WorkerPool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const path = "handlers.removeCurrency.NewRemoveCurrencyHandler"
		currencyID := strings.TrimSpace(r.FormValue("currencyID"))
		if currencyID == "" {
			log.WithFields(logrus.Fields{
				"path": path,
			}).Error("Currency ID is required")
			utils.Respond(w, r, http.StatusBadRequest, "Currency ID is required")
			return
		}
		pool.RemoveCurrency(currencyID)
		err := store.RemoveCurrency(currencyID)
		if err != nil {
			log.WithFields(logrus.Fields{
				"path":  path,
				"error": err.Error(),
			}).Error("Failed to remove currency from store")
			//pool.AddCurrency(currencyID) // Может быть, стоит вернуть валюту в пул, если удаление не удалось?
			utils.Respond(w, r, http.StatusInternalServerError, "Failed to add currency: "+err.Error())
			return
		}
		log.WithFields(logrus.Fields{
			"path":       path,
			"currencyID": currencyID,
		}).Info("Currency removed successfully")
		utils.Respond(w, r, http.StatusOK, "Currency removed successfully")

	}
}
