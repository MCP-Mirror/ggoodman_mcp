package integrations

type Instructions struct {
	Runtime string   `json:"runtime"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
}
