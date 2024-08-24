![Demo](./.github/gif-demo.gif)


# AlertManager -> Awtrix

Highly inspired by the [Blinky](https://www.getblinky.io/), I wanted to receive my monitoring alerts directly on my Awtrix matrix. That's all 🤷‍♂️

This program will expose a web server on port 5000 and must be configured as a webhook endpoint from AlertManager. 

## Install
### Using source code 

Install requirements : 
- Go (>=1.21.6)

Build binary : 

```bash
go build -o alertmanager2awtrix
```

You will obtain a binary `alertmanager2awtrix`.

### Using Docker

Sadly, no Dockerfile is provided in this project. I use [ko](https://ko.build) to generate a Dockerfile.

Use the following `docker-compose.yaml` :

```yaml

version: '3.8'
services:
  alertmanager2awtrix:
    image: ghcr.io/qjoly/alertmanager-awtrix:latest
    container_name: alertmanager2awtrix
    ports:
      - 5000:5000
    environment:
      - "AWTRIX_NOTIFY_URL=http://192.168.1.26/api/notify"
      - "AWTRIX_USERNAME=awtrix" # optional
      - "AWTRIX_PASSWORD=awtrix" # optional
      - "FIRING_ALERT_ICON=555"
      - "FIRING_TEXT_COLOR=#FF0000"
      - "RESOLVED_ALERT_ICON=138"
      - "RESOLVED_TEXT_COLOR=#C0C78C"
      - "HOLD_ALERT=true"
```

# Configuration

## Environment variables
- `AWTRIX_NOTIFY_URL` (e.g. `http://192.168.1.26/api/notify`)
- `AWTRIX_USERNAME` and `AWTRIX_PASSWORD` (only if needed)
- `FIRING_ALERT_ICON` (and `RESOLVED_ALERT_ICON`) with the icon id (it should be downloaded in the awtrix)
- `FIRING_TEXT_COLOR` (and `RESOLVED_TEXT_COLOR`) with RGB Hexadecimal color (e.g. `#FF0000`)
- `HOLD_ALERT` if you want to have to press the middle button to dismiss the alert

## Add webhook in alertmanager

```yaml
receivers:
  - name: 'web.hook'
    webhook_configs:
      - url: 'http://<your_ip>:5000/webhook'
```

# ROADMAP

- Support basic auth
- Use alert labels to choose the icon displayed

feel free to give me your ideas




