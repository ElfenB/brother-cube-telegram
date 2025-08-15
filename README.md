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

## Deployment

For easier management, you can use `taskfile` to run tasks. To install it, follow the instructions in their [documentation](https://taskfile.dev/docs/installation).

If you don't want to use it, you can run the tasks manually.

The `Taskfile.yaml` defines tasks for building, uploading, and managing the service on the Raspberry Pi.

To configure the settings for your host, modify the vars on top of the `Taskfile.yaml`.

```bash
# See available tasks
task
```

To start the service immediately, run:

```bash
sudo systemctl start brother-cube-telegram-pi.service
```

To activate the service on boot, run:

```bash
task activate-service
```
