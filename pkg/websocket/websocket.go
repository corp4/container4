package websocket

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os/exec"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type Websocket struct {
	// Channels
	output chan []byte
	add    chan *websocket.Conn
	remove chan *websocket.Conn
	close  chan bool

	// Boolean to know if the websocket is open
	open bool

	// Server properties
	host   string
	port   int
	server *http.Server
}

func CreateWebsocket(host string, port int) *Websocket {
	ws := &Websocket{
		output: make(chan []byte),
		add:    make(chan *websocket.Conn),
		remove: make(chan *websocket.Conn),
		close:  make(chan bool),
		open:   false,
		host:   host,
		port:   port,
	}

	// Start the "broadcaster" goroutine in the background to update connections
	go broadcaster(ws)

	return ws
}

// Start the websocket server asynchronously
func (ws *Websocket) Start() {
	if ws.open {
		return
	}

	ws.open = true
	go startWebsocket(ws)
}

func (ws *Websocket) Stop() {
	if !ws.open {
		return
	}

	// Close all connections and Stop the server
	ws.close <- true
}

func startWebsocket(ws *Websocket) {
	// Start the WebSocket server
	ws.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%d", ws.host, ws.port),
		Handler: http.HandlerFunc(handleStream(ws)),
	}

	if err := ws.server.ListenAndServe(); err != nil {
		log.Warnln("Error while starting the server: ", err)
	}
}

func (ws *Websocket) Send(message []byte) error {
	if !ws.open {
		return fmt.Errorf("Websocket is not open")
	}

	ws.output <- message
	return nil
}

// Start a command and get the output (stdout and stderr) to send it to the websocket
// Will wait for the command to finish
func (ws *Websocket) RedirectCommandOutput(cmd *exec.Cmd, stdout io.ReadCloser, stderr io.ReadCloser) error {
	if !ws.open {
		return fmt.Errorf("Websocket is not open")
	}

	// Read the output of the command and store it in the channel
	if stdout != nil {
		go func() {
			buf := make([]byte, 1024)
			for {
				n, err := stdout.Read(buf)
				if err != nil {
					log.Debugln("End of output stdout: ", err)
					return
				}

				log.Debugln("Read: [stdout]", string(buf[:n]))
				ws.output <- []byte("[stdout]" + string(buf[:n]))
			}
		}()
	}

	// stderr
	if stderr != nil {
		go func() {
			buf := make([]byte, 1024)
			for {
				n, err := stderr.Read(buf)
				if err != nil {
					log.Debugln("End of output stderr: ", err)
					return
				}

				log.Debugln("Read: [stderr]", string(buf[:n]))
				ws.output <- []byte("[stderr]" + string(buf[:n]))
			}
		}()
	}

	return nil
}

func handleStream(ws *Websocket) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Upgrade la connexion HTTP à une connexion WebSocket
		var upgrader = websocket.Upgrader{}
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "Failed to upgrade connection", http.StatusBadRequest)
			log.Warnln("Failed to upgrade connection:", err)
			return
		}

		// Add the connection to the "add" channel
		ws.add <- conn

		// Wait for the connection to be closed
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}

		// Remove the connection from the "remove" channel
		ws.remove <- conn

		// Close the connection
		conn.Close()
	}
}

func broadcaster(ws *Websocket) {
	// List of connections to update
	connections := make([]*websocket.Conn, 0)
	buff := make([]byte, 0)

	loop := true
	for loop {
		select {
		// Add a connection to the list
		case conn := <-ws.add:
			log.Debugln("New connection")
			connections = append(connections, conn)

			// Send the content stored in the buffer to the new connection
			if err := conn.WriteMessage(websocket.TextMessage, buff); err != nil {
				ws.remove <- conn
			}

		// Remove a connection from the list
		case conn := <-ws.remove:
			log.Debugln("End of a connection")
			for i, c := range connections {
				if c == conn {
					connections = append(connections[:i], connections[i+1:]...)
					break
				}
			}
		// Update all connections with the content stored in the channel
		case data := <-ws.output:
			for _, conn := range connections {
				if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
					// Remove the connection if an error occurred while sending
					ws.remove <- conn
				}
			}

			// Append the new data to the buffer
			buff = append(buff, data...)

		// Close all connections
		case close := <-ws.close:
			if !close {
				continue
			}

			log.Infoln("Closing all connections and the server")
			for _, conn := range connections {
				conn.Close()
			}

			ws.server.Shutdown(context.Background())
			loop = false
		}
	}
}
