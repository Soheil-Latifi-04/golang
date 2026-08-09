# Slack Age Bot

A Slack bot written in Go that calculates your age from your year of birth (YOB).

## Prerequisites

* Go 1.21 or later
* A Slack app with [Socket Mode](https://api.slack.com/apis/connections/socket) enabled
* Bot token (`xoxb-...`) and app-level token (`xapp-...`) with `connections:write` scope

## Setup

1. Clone the repository:

```bash
   git clone https://github.com/soheil/slack-age-bot.git
   cd slack-age-bot
   ```

2. Copy the example environment file and fill in your tokens:

```bash
   cp .env.example .env
   ```

3. Install dependencies:

```bash
   go mod download
   ```

## Running locally

The bot loads a `.env` file automatically when present (via [godotenv](https://github.com/joho/godotenv)). Copy `.env.example` to `.env`, add your tokens, then run:

```bash
go run .
```

You can also set variables in your shell instead of using a `.env` file:

**Linux / macOS:**

```bash
export SLACK\_BOT\_TOKEN="xoxb-your-bot-token"
export SLACK\_APP\_TOKEN="xapp-your-app-token"
go run .
```

**Windows (PowerShell):**

```powershell
$env:SLACK\_BOT\_TOKEN = "xoxb-your-bot-token"
$env:SLACK\_APP\_TOKEN = "xapp-your-app-token"
go run .
```

Shell variables take precedence over values in `.env` if both are set.

## Usage

Message the bot in Slack:

```
my yob is 1990
```

The bot replies with your calculated age and year of birth.

## Build

```bash
go build -o slack-age-bot .
```

## Environment variables

|Variable|Required|Description|
|-|-|-|
|`SLACK\_BOT\_TOKEN`|Yes|Bot User OAuth Token (`xoxb-...`)|
|`SLACK\_APP\_TOKEN`|Yes|App-level token for Socket Mode|



