package entities

type Action = string

const (
	ActionResize    = "resize"
	ActionMiniature = "miniature"
	ActionWatermark = "watermark"
)

type Message struct {
	Action Action `json:"action"`
	Image  Image  `json:"image"`
}
