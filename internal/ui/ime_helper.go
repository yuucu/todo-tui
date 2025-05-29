package ui

import (
	"os"
	"strings"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
)

// IMEHelper provides enhanced support for Japanese input
type IMEHelper struct {
	// Track composition state for IME input
	isComposing     bool
	compositionText string
}

// NewIMEHelper creates a new IME helper instance
func NewIMEHelper() *IMEHelper {
	return &IMEHelper{
		isComposing:     false,
		compositionText: "",
	}
}

// ProcessKeyMsg processes key messages with IME support
func (ime *IMEHelper) ProcessKeyMsg(msg tea.KeyMsg) (isIME bool, text string) {
	// Check if this is likely an IME input
	if ime.isLikelyIMEInput(msg) {
		ime.handleIMEInput(msg)
		return true, ime.compositionText
	}

	// Handle final character commit
	if ime.isComposing && ime.isCommitKey(msg) {
		ime.isComposing = false
		result := ime.compositionText
		ime.compositionText = ""
		return false, result
	}

	return false, ""
}

// isLikelyIMEInput checks if the input is likely from an IME
func (ime *IMEHelper) isLikelyIMEInput(msg tea.KeyMsg) bool {
	// Check for multi-byte characters or IME composition
	str := msg.String()

	// If it's a single ASCII character, it's not IME input
	if len(str) == 1 && str[0] <= 127 {
		return false
	}

	// If it contains multi-byte characters, it's likely IME input
	if !utf8.ValidString(str) || utf8.RuneCountInString(str) != len(str) {
		return true
	}

	// Check for Japanese characters (Hiragana, Katakana, Kanji)
	for _, r := range str {
		if (r >= 0x3040 && r <= 0x309F) || // Hiragana
			(r >= 0x30A0 && r <= 0x30FF) || // Katakana
			(r >= 0x4E00 && r <= 0x9FAF) { // CJK Unified Ideographs
			return true
		}
	}

	return false
}

// handleIMEInput handles IME input processing
func (ime *IMEHelper) handleIMEInput(msg tea.KeyMsg) {
	ime.isComposing = true
	ime.compositionText = msg.String()
}

// isCommitKey checks if the key commits the current composition
func (ime *IMEHelper) isCommitKey(msg tea.KeyMsg) bool {
	key := msg.String()
	return key == enterKeyStr || key == " " || key == "tab"
}

// SetupIMEEnvironment sets up environment variables for better IME support
func SetupIMEEnvironment() {
	// Set environment variables for better Japanese input support
	if os.Getenv("LANG") == "" {
		os.Setenv("LANG", "ja_JP.UTF-8")
	}

	// Ensure proper locale for IME
	if os.Getenv("LC_CTYPE") == "" {
		os.Setenv("LC_CTYPE", "ja_JP.UTF-8")
	}

	// Enable XIM for better compatibility
	if os.Getenv("XMODIFIERS") == "" {
		// Try to detect and set appropriate input method
		if ime := detectInputMethod(); ime != "" {
			os.Setenv("XMODIFIERS", "@im="+ime)
		}
	}
}

// detectInputMethod attempts to detect available input methods
func detectInputMethod() string {
	// Common Japanese input methods on macOS and Linux
	methods := []string{"fcitx5", "fcitx", "ibus", "uim"}

	for _, method := range methods {
		if isInputMethodAvailable(method) {
			return method
		}
	}

	return ""
}

// isInputMethodAvailable checks if an input method is available
func isInputMethodAvailable(method string) bool {
	// This is a simplified check - in a real implementation,
	// you might want to check for running processes or installed packages
	switch method {
	case "fcitx5", "fcitx":
		return checkEnvOrProcess("FCITX", "")
	case "ibus":
		return checkEnvOrProcess("IBUS", "")
	case "uim":
		return checkEnvOrProcess("UIM", "")
	}
	return false
}

// checkEnvOrProcess checks for environment variables or running processes
func checkEnvOrProcess(envVar, _ string) bool {
	// Check environment variable
	if os.Getenv(envVar+"_SOCKET") != "" || os.Getenv(envVar+"_DAEMON") != "" {
		return true
	}

	// This is a simplified check - in production, you might want to
	// actually check for running processes
	return false
}

// ValidateUTF8Input ensures the input is valid UTF-8
func ValidateUTF8Input(input string) string {
	if !utf8.ValidString(input) {
		// Clean up invalid UTF-8 sequences
		return strings.ToValidUTF8(input, "")
	}
	return input
}
