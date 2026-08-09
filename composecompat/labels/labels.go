// Package labels defines the Compose labels required by compatibility operations.
package labels

// Label is a Docker label name.
type Label string

// HASH identifies the Compose configuration hash stored on a service.
const HASH = Label("com.docker.compose.config-hash")

// Str returns the label name.
func (label Label) Str() string { return string(label) }
