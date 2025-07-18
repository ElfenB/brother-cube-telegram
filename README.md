# Brother Cube Telegram Bot

A simple Telegram bot built with Go. It interacts with Brother printers to print messages sent by users.

The interaction with the printer is enabled by the ptouch-print tool by Dominic Radermacher.

- <https://git.familie-radermacher.ch/linux/ptouch-print.git>
- <https://dominic.familie-radermacher.ch/projekte/ptouch-print>

## Setup

1. Create a new bot by messaging [@BotFather](https://t.me/botfather) on Telegram
2. Follow the instructions to get your bot token
3. Set the environment variable:

```bash
export TELEGRAM_BOT_TOKEN="your_actual_bot_token_here"
```

Or create a `.env` file (copy from `.env.example`) and load it before running

## Running

```bash
go run main.go
```
