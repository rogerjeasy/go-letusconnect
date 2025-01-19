// websocket/handler.go
package websocket

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/coder/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/rogerjeasy/go-letusconnect/handlers"
	"github.com/rogerjeasy/go-letusconnect/models"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

type Handler struct {
	manager *models.Manager
}

func NewHandler(manager *models.Manager) *Handler {
	return &Handler{
		manager: manager,
	}
}

type responseWriter struct {
	*fiber.Response
	written int
	conn    net.Conn
	writer  *bufio.Writer
	headers http.Header
}

func (w *responseWriter) Header() http.Header {
	if w.headers == nil {
		w.headers = make(http.Header)
	}
	return w.headers
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.Response.SetBodyRaw(b)
	w.written += len(b)
	return len(b), nil
}

func (w *responseWriter) WriteHeader(statusCode int) {
	if w.written > 0 {
		return
	}
	w.Response.SetStatusCode(statusCode)
	for key, values := range w.headers {
		for _, value := range values {
			w.Response.Header.Set(key, value)
		}
	}
	w.written = 1
}

// Implement http.Hijacker interface
func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if w.conn == nil {
		return nil, nil, fmt.Errorf("failed to get underlying connection")
	}

	// Create buffered reader/writer for the connection
	br := bufio.NewReader(w.conn)
	bw := bufio.NewWriter(w.conn)
	return w.conn, bufio.NewReadWriter(br, bw), nil
}

func (h *Handler) HandleWebSocket(c *fiber.Ctx) error {
	// Extract and validate token
	token := c.Query("token")
	if token == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Missing token")
	}
	userID, err := handlers.ValidateToken(token)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	// Set headers before upgrading
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Credentials", "true")
	c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, Sec-WebSocket-Key, Sec-WebSocket-Version, Sec-WebSocket-Extensions")
	c.Set("Connection", "Upgrade")
	c.Set("Upgrade", "websocket")

	// Create a done channel to coordinate shutdown
	done := make(chan struct{})

	// Convert Fiber context to http.ResponseWriter and *http.Request
	w := &responseWriter{
		Response: c.Response(),
		conn:     c.Context().Conn(),
		writer:   bufio.NewWriter(c.Context().Conn()),
	}

	var httpRequest http.Request
	if err := fasthttpadaptor.ConvertRequest(c.Context(), &httpRequest, true); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Could not convert request")
	}

	// Configure WebSocket options
	opts := &websocket.AcceptOptions{
		InsecureSkipVerify: true,
		OriginPatterns:     []string{"*"},
		CompressionMode:    websocket.CompressionDisabled,
	}

	// Perform the upgrade
	conn, err := websocket.Accept(w, &httpRequest, opts)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Could not upgrade connection")
	}

	// Create and register client first
	client := NewClient(conn, h.manager, userID)
	h.manager.Register <- client

	// Start client routines with proper shutdown coordination
	go func() {
		ReadPump(client)
		close(done)
	}()
	go WritePump(client)

	// Wait for connection to close
	<-done

	return nil
}
