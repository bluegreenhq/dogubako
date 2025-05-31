package context

type ContextKey string

const (
	// userinterface layer
	ContextKeyRequestID   ContextKey = "requestId"
	ContextKeyRequestTime ContextKey = "requestTime"
	// application layer
	ContextKeyTransaction ContextKey = "transaction"
)
