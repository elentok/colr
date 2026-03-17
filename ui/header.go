package ui

import (
	"fmt"
	"strings"
)

// RenderHeader renders the header region.
func RenderHeader(clipText, normalizedText, toastMsg string, width int) string {
	inner := width - 2 // account for border padding

	clip := HeaderLabelStyle.Render("Clipboard:") + " " + HeaderValueStyle.Render(truncate(clipText, inner-14))
	norm := HeaderLabelStyle.Render("Normalized:") + " " + HeaderValueStyle.Render(normalizedText)

	var status string
	if toastMsg != "" {
		if strings.HasPrefix(toastMsg, "Failed") {
			status = HeaderLabelStyle.Render("Status:") + " " + ToastErrStyle.Render(toastMsg)
		} else {
			status = HeaderLabelStyle.Render("Status:") + " " + ToastOKStyle.Render(toastMsg)
		}
	} else {
		status = HeaderLabelStyle.Render("Status:") + " " + ToastOKStyle.Render("OK")
	}

	return fmt.Sprintf("%s\n%s\n%s", clip, norm, status)
}

func truncate(s string, max int) string {
	if max < 3 {
		return s
	}
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max-1]) + "…"
}
