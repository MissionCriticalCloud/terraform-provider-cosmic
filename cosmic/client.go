package cosmic

import "github.com/MissionCriticalCloud/go-cosmic/v6/cosmic"

// CosmicClient wraps CosmicClient to instantiate a new
// Cosmic client whilst adding some bonus fields and methods.
type CosmicClient struct {
	*cosmic.CosmicClient

	ZoneName string
}
