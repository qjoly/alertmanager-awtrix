package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/qjoly/alertmanager-awtrix/pkg/awtrix"
	"github.com/qjoly/alertmanager-awtrix/pkg/types"
	"github.com/qjoly/alertmanager-awtrix/pkg/version"
)

var ac *awtrix.AwtrixClient
var logger *slog.Logger

type AlertManagerWebhook struct {
	Alerts []types.Alert `json:"alerts"`
}

func getEnvAsBool(key string, defaultVal bool) bool {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}

func createAwtrixClient() *awtrix.AwtrixClient {

	awtrixURL := os.Getenv("AWTRIX_NOTIFY_URL")
	if awtrixURL == "" {
		fmt.Println("No Awtrix API found, please add a AWTRIX_NOTIFY_URL environment variable")
		os.Exit(1)
	}

	awtrixUser := os.Getenv("AWTRIX_USERNAME")
	awtrixPass := os.Getenv("AWTRIX_PASSWORD")

	firingAlertIcon := os.Getenv("FIRING_ALERT_ICON")
	if firingAlertIcon == "" {
		firingAlertIcon = "555"
	}

	firingTextColor := os.Getenv("FIRING_TEXT_COLOR")
	if firingTextColor == "" {
		firingTextColor = "#C0C78C"
	}

	resolvedAlertIcon := os.Getenv("RESOLVED_ALERT_ICON")
	if resolvedAlertIcon == "" {
		resolvedAlertIcon = "138"
	}

	resolvedTextColor := os.Getenv("RESOLVED_TEXT_COLOR")
	if resolvedTextColor == "" {
		resolvedTextColor = "#C0C78C"
	}

	holdAlert := getEnvAsBool("HOLD_ALERT", true)

	client, err := awtrix.NewClient(awtrixUser, awtrixPass, awtrixURL, firingTextColor, firingAlertIcon, resolvedAlertIcon, resolvedTextColor, holdAlert, logger)
	if err != nil {
		logger.Error("Can't create awtrix client" + err.Error())
		os.Exit(1)
	}

	return client

}

// Webhook handler
func webhookHandler(w http.ResponseWriter, r *http.Request) {
	var webhook AlertManagerWebhook

	err := json.NewDecoder(r.Body).Decode(&webhook)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	for _, alert := range webhook.Alerts {
		logger.Info("Alert: %s, Status: %s, Description: %s, Started: %s", alert.Labels["alertname"], alert.Status, alert.Annotations["description"], alert.StartsAt)
		if err := ac.SendAwtrixNotification(alert); err != nil {
			logger.Error("Error sending notification to Awtrix: " + err.Error())
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}

func main() {
	ac = createAwtrixClient()
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	http.HandleFunc("/webhook", webhookHandler)

	logger.Info("Starting server on :5000 at version " + version.GoBuildVersion + " from sha " + version.GoBuildSHA)
	log.Fatal(http.ListenAndServe(":5000", nil))
}
