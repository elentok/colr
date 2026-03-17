package app

import "github.com/elentok/colr/ui"

// Re-export ui types for use within the app package.
type EditMode = ui.EditMode

const (
	ModeHSV = ui.ModeHSV
	ModeRGB = ui.ModeRGB

	FieldHue        = ui.FieldHue
	FieldSaturation = ui.FieldSaturation
	FieldValue      = ui.FieldValue
	FieldRed        = ui.FieldRed
	FieldGreen      = ui.FieldGreen
	FieldBlue       = ui.FieldBlue
	FieldOpacity    = ui.FieldOpacity
	FieldCount      = ui.FieldCount
)
