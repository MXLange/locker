# locker  
A WebSocket-based implementation to ensure ordered execution of actions across multiple clients.

## 🚀 Getting Started

### 🖥️ Run the Server

To start the demo server, run the following command:

```bash
go run ./cmd/server/main.go
```

### 🧪 Run the Clients

In separate terminals, run the clients to test interaction and execution ordering:

**Terminal 1:**
```bash
go run cmd/client_1/main.go
```

**Terminal 2:**
```bash
go run cmd/client_2/main.go
```

## 📦 Project Structure

```
cmd/
├── server/     # WebSocket server implementation
├── client_1/   # Example client 1
└── client_2/   # Example client 2
```

## 🛠️ Description

This project showcases how to coordinate the execution of actions between multiple WebSocket clients while preserving the order of operations. It can be used as a base for distributed locking, task queues, or real-time collaboration tools.

---

Feel free to clone, experiment, and adapt the project to your needs!