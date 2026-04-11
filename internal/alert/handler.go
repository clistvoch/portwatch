package alert

// HandlerFunc is a function adapter that implements Handler.
type HandlerFunc func(Change)

// Handle calls the underlying function.
func (f HandlerFunc) Handle(c Change) {
	f(c)
}

// NopHandler is a Handler that discards all changes. Useful in tests.
type NopHandler struct{}

func (NopHandler) Handle(Change) {}

// MultiHandler fans a change out to multiple handlers in order.
type MultiHandler struct {
	handlers []Handler
}

// NewMultiHandler creates a MultiHandler from the provided handlers.
func NewMultiHandler(handlers ...Handler) *MultiHandler {
	h := make([]Handler, len(handlers))
	copy(h, handlers)
	return &MultiHandler{handlers: h}
}

// Handle delivers the change to every registered handler.
func (m *MultiHandler) Handle(c Change) {
	for _, h := range m.handlers {
		h.Handle(c)
	}
}

// Add appends a handler to the MultiHandler.
func (m *MultiHandler) Add(h Handler) {
	m.handlers = append(m.handlers, h)
}

// Len returns the number of handlers registered in the MultiHandler.
func (m *MultiHandler) Len() int {
	return len(m.handlers)
}
