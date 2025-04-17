package discovery

import "github.com/google/wire"

// ProviderSet is discovery service providers.
var ProviderSet = wire.NewSet(
	NewAgentService, // example discovery remote grpc service.
)
