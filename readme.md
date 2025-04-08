# Chat Room Service

## Features

- **Multiple Client Support**: Handles multiple clients concurrently using **goroutines** and **channels**.
- **Thread-Safe Operations**: Access to shared resources is protected using **mutexes** to ensure thread-safety.
- **Timeout Handling**: `/messages` endpoint has a timeout to prevent blocking indefinitely.
- **Client Disconnect Handling**: When a client leaves, they will stop receiving messages.
- **Graceful Exit Handling**: Clients that leave the chat will automatically stop receiving messages.
- **Efficient Broadcasting**: Messages are sent to all connected clients using a central message distribution channel.

## Getting Started

### Prerequisites

- **Go** (for running the server)
- **Postman** or **cURL** (for testing the endpoints)

### Running the Application

1. **Clone the repository**:

   ```bash
   git clone https://github.com/akshaykajithkumar/chat-room-service.git
   cd chat-room-service
   ```

2. **Run the server**:

   ```bash
   go run main.go
   ```

## API Endpoints

### 1. Join the Chat Room

To join the chat room, send a POST request to `/join` with a client ID as a query parameter.

#### Request:

```bash
curl --location --request POST 'http://localhost:5000/join?id=client123'
```

### 2. Receive Messages

To receive messages, send a GET request to `/messages` with your client ID as a query parameter. This will keep the connection open and stream messages. If no new messages arrive within a set timeout period, the connection will close.

#### Request:

```bash
curl --location 'http://localhost:5000/messages?id=client123'
```

### 3. Send a Message

To send a message to the chat room, use the POST request to `/send` with the client ID and message query parameters.

#### Request:

```bash
curl --location --request POST 'http://localhost:5000/send?id=client123&message=Hi'
```

### 4. Leave the Chat Room

To leave the chat room, send a POST request to `/leave` with the client ID as a query parameter. Once a client leaves, they will no longer receive messages.

#### Request:

```bash
curl --location --request DELETE 'http://localhost:5000/leave?id=client123'
```


