package locker

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

var ErrInvalidPort error = errors.New("invalid port format, port must be in the format ':port' and in the range 1-65535")

// LockerServer is a WebSocket server that manages locks for different IDs.
// Each ID can have multiple clients connected, but only one client can hold the lock at a time.
type LockerServer struct {
	upgrader websocket.Upgrader
	locker   *lockManager
	port     string
}

// NewLockServer creates a new LockerServer instance.
// port must follow the format ":port".
// Example: ":8080"
func NewLockServer(port string) (*LockerServer, error) {

	isValidPort := isValidPortFormat(port)

	if !isValidPort {
		return nil, ErrInvalidPort
	}

	return &LockerServer{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		locker: newLockManager(),
		port:   port,
	}, nil
}

// Start starts the WebSocket server and listens for incoming connections on the specified port.
func (ls *LockerServer) Start() error {
	http.HandleFunc("/ws", ls.sequenceManager)

	err := http.ListenAndServe(ls.port, nil)
	if err != nil {
		return err
	}
	return nil
}

// lockManager is a simple lock manager that uses a map to store locks for different IDs.
type lockManager struct {
	mu    sync.Mutex
	locks map[string]chan struct{}
}

// newLockManager creates a new lock manager instance.
func newLockManager() *lockManager {
	return &lockManager{
		locks: make(map[string]chan struct{}),
	}
}

// getLock returns a channel that acts as a lock for the given ID.
func (lm *lockManager) getLock(id string) chan struct{} {
	lm.mu.Lock()
	defer lm.mu.Unlock()

	lock, exists := lm.locks[id]
	if !exists {
		lock = make(chan struct{}, 1)
		lock <- struct{}{} // Inicializa com uma posição livre
		lm.locks[id] = lock
	}
	return lock
}

// ServeHTTP handles incoming WebSocket connections and manages locks for different IDs.

func (ls *LockerServer) sequenceManager(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	conn, err := ls.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "could not upgrade connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	lock := ls.locker.getLock(id)

	<-lock

	err = conn.WriteMessage(websocket.TextMessage, []byte("go"))
	if err != nil {
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			lock <- struct{}{}
			return
		}

		if strings.EqualFold(string(message), "unlock") {
			lock <- struct{}{}
			return
		}
	}
}

// IsValidPortFormat checks if the string follows the ":port" format and the port is in the valid range (1-65535).
func isValidPortFormat(s string) bool {
	re := regexp.MustCompile(`^:(\d{1,5})$`)
	matches := re.FindStringSubmatch(s)
	if len(matches) != 2 {
		return false
	}

	port, err := strconv.Atoi(matches[1])
	if err != nil || port < 1 || port > 65535 {
		return false
	}

	return true
}
