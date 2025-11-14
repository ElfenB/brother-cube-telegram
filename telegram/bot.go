package telegram

import (
	"brother-cube-telegram/logger"
	"context"
	"os"

	"github.com/go-telegram/bot"
)

// CommandInfo holds information about a bot command
type CommandInfo struct {
	Command     string
	Description string
	Usage       string
	Example     string
}

// Global command registry
var commandRegistry = map[string]CommandInfo{
	"help": {
		Command:     "/help",
		Description: "Show this help message with all available commands",
		Usage:       "/help",
		Example:     "/help",
	},
	"status": {
		Command:     "/status",
		Description: "Get current printer status and information",
		Usage:       "/status",
		Example:     "/status",
	},
	"preview": {
		Command:     "/preview",
		Description: "Generate a preview image of your label before printing",
		Usage:       "/preview <text>",
		Example:     "/preview Kitchen Labels",
	},
	"size": {
		Command:     "/size",
		Description: "Print a label with custom font size",
		Usage:       "/size <font_size> <text>",
		Example:     "/size 32 My Custom Label",
	},
	"preset": {
		Command:     "/preset",
		Description: "Print a label using a predefined preset with specific font settings",
		Usage:       "/preset [preset_name] [text] | /preset (to list presets)",
		Example:     "/preset kitchen Container A",
	},
}

// GetRegisteredCommands returns all registered commands
func GetRegisteredCommands() []CommandInfo {
	commands := make([]CommandInfo, 0, len(commandRegistry))
	for _, cmd := range commandRegistry {
		commands = append(commands, cmd)
	}

	// Add the default text message handler info
	commands = append(commands, CommandInfo{
		Command:     "Text Message",
		Description: "Send any text message (not starting with /) to print it with default settings",
		Usage:       "Just type your text",
		Example:     "Hello World",
	})

	return commands
}

func GetBot(ctx context.Context) *bot.Bot {
	// Get bot token from environment variable
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		logger.Error("TELEGRAM_BOT_TOKEN environment variable is required")
		os.Exit(1)
	}

	opts := []bot.Option{
		bot.WithDefaultHandler(defaultHandler),
		bot.WithMiddlewares(
			recoveryMiddleware,
			authorizationMiddleware,
			createMiddlewareWithCtxFactory(ctx, printerMiddlewareHandler),
			loggingMiddleware,
		),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		logger.Error("Failed to create bot: %v", err)
		os.Exit(1)
	}

	// Register handlers using the command registry
	registerCommandHandler(b, "help", bot.MatchTypeCommandStartOnly, helpHandler)
	registerCommandHandler(b, "status", bot.MatchTypeCommandStartOnly, statusHandler)
	registerCommandHandler(b, "preview", bot.MatchTypeCommand, previewHandler)
	registerCommandHandler(b, "size", bot.MatchTypeCommandStartOnly, sizeHandler)
	registerCommandHandler(b, "preset", bot.MatchTypeCommandStartOnly, presetHandler)

	// Register handler for unknown commands (any command that starts with /)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/", bot.MatchTypePrefix, unknownCommandHandler)

	return b
}

// registerCommandHandler registers a command handler and ensures it exists in the registry
func registerCommandHandler(b *bot.Bot, command string, matchType bot.MatchType, handler bot.HandlerFunc) {
	// Verify that the command exists in our registry
	if _, exists := commandRegistry[command]; !exists {
		logger.Warn("Command '%s' is not registered in commandRegistry - help will not include it", command)
	}

	// Register the handler
	b.RegisterHandler(bot.HandlerTypeMessageText, command, matchType, handler)
}
