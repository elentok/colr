package ui

// EditMode represents the current editor mode.
type EditMode int

const (
	ModeHSV EditMode = iota
	ModeRGB
)

// Field indices — shared Opacity is always the last field (index 3) in both modes.
const (
	FieldHue        = 0
	FieldSaturation = 1
	FieldValue      = 2

	FieldRed   = 0
	FieldGreen = 1
	FieldBlue  = 2

	FieldOpacity = 3
	FieldCount   = 4
)
