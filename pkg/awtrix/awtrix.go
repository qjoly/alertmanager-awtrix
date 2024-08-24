package awtrix

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"alertmanager-awtrix/pkg/types"
)

type (
	AwtrixClient struct {
		resolvedAlertIcon string
		resolvedTextColor [3]int
		firingAlertIcon   string
		firingTextColor   [3]int
		awtrixApiURL      string
		awtrixUser        string
		awtrixPass        string
		holdAlert         bool
		logger            *slog.Logger
	}
)

func NewClient(awtrixUser string, awtrixPass string, awtrixApiURL string, firingTextColorHex string, firingAlertIcon string, resolvedAlertIcon string, resolvedTextColorHex string, holdAlert bool, logger *slog.Logger) (*AwtrixClient, error) {

	firingTextColor, err := hexToRGB(firingTextColorHex)
	if err != nil {
		return nil, fmt.Errorf("Can't parse %s, it should be in format '#FF00FF' ", firingTextColorHex)
	}

	resolvedTextColor, err := hexToRGB(resolvedTextColorHex)
	if err != nil {
		return nil, fmt.Errorf("Can't parse %s, it should be in format '#FF00FF' ", resolvedTextColorHex)
	}

	client := &AwtrixClient{
		awtrixUser:        awtrixUser,
		awtrixPass:        awtrixPass,
		awtrixApiURL:      awtrixApiURL,
		firingTextColor:   firingTextColor,
		firingAlertIcon:   firingAlertIcon,
		resolvedAlertIcon: resolvedAlertIcon,
		resolvedTextColor: resolvedTextColor,
		holdAlert:         holdAlert,
		logger:            logger,
	}

	return client, nil
}

func (ac *AwtrixClient) SendAwtrixNotification(alert types.Alert) error {

	var icon string
	var color [3]int

	if alert.Status == "firing" {
		color = ac.firingTextColor
		icon = ac.firingAlertIcon

	} else {

		color = ac.resolvedTextColor
		icon = ac.resolvedAlertIcon
	}

	// Define the payload to send to Awtrix
	payload := map[string]interface{}{
		"color":    color, // La couleur est maintenant un [3]int
		"repeat":   5,
		"duration": 2000,
		"hold":     ac.holdAlert,
		"icon":     icon,
		"text":     fmt.Sprintf("Alert: %s", alert.Labels["alertname"]),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", ac.awtrixApiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	if ac.awtrixUser != "" {
		auth := base64.StdEncoding.EncodeToString([]byte(ac.awtrixUser + ":" + ac.awtrixPass))
		req.Header.Set("Authorization", "Basic "+auth)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send notification to Awtrix: %s", resp.Status)
	}

	return nil
}

func hexToRGB(hex string) ([3]int, error) {
	var rgb [3]int

	if hex[0] != '#' {
		return rgb, fmt.Errorf("can't handle this color : %s", hex)
	}

	for i := 0; i < 3; i++ {
		val, err := strconv.ParseInt(hex[1+i*2:3+i*2], 16, 64)
		if err != nil {
			return rgb, fmt.Errorf("error during conversion : %v", err)
		}
		rgb[i] = int(val)
	}

	return rgb, nil
}
