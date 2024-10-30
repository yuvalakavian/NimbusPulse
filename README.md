Here’s an updated README with the project repo name as **NimbusPulse**.

---

# NimbusPulse

**NimbusPulse** is a real-time environmental monitoring system built with Go, leveraging WebSocket for live data streaming and HTTP for sensor data ingestion. This project serves as a foundation for expanding into a full IoT solution, tracking environmental metrics such as temperature and humidity in real-time.

## Features

- **Receive Sensor Data**: Collects temperature, humidity, and other environmental data via HTTP POST requests.
- **Real-Time Updates**: Streams the latest sensor data to connected clients through WebSocket.
- **Expandable Architecture**: Designed for easy expansion to include more sensors, data processing, and persistent storage options.

## Technologies Used

- **Go**: For server-side logic and handling concurrency.
- **WebSocket**: For real-time data streaming to clients.
- **HTTP**: For receiving sensor data from IoT devices.

## Project Structure

- **`main.go`**: Main server setup, HTTP and WebSocket endpoints.
- **`sensor.go`**: Manages sensor data storage and parsing.
- **`ws.go`**: Handles WebSocket connections.

## Getting Started

### Prerequisites

- Go 1.16 or higher
- Basic understanding of HTTP and WebSocket protocols
- IoT device(s) capable of sending data (or a tool like `curl` for testing)

### Installation

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/your-username/NimbusPulse.git
   cd NimbusPulse
   ```

2. **Install Dependencies**:
   This project uses the `gorilla/websocket` package for WebSocket handling. Install it with:
   ```bash
   go get -u github.com/gorilla/websocket
   ```

3. **Run the Server**:
   ```bash
   go run main.go
   ```
   The server will start on `http://localhost:8080`.

### Testing the Server

1. **Send Mock Sensor Data**:
   Use `curl` to simulate sending sensor data:

   ```bash
   curl -X POST http://localhost:8080/sensor -H "Content-Type: application/json" -d '{"id":"sensor1", "temperature":25.5, "humidity":60.0}'
   ```

2. **Connect to WebSocket Endpoint**:
   Open the browser console and connect to the WebSocket to receive real-time updates:

   ```javascript
   const ws = new WebSocket("ws://localhost:8080/ws");
   ws.onmessage = function(event) {
       console.log("Real-time data:", event.data);
   };
   ```

## API Endpoints

### 1. POST `/sensor`
Accepts JSON data for environmental readings.

**Request Body**:
```json
{
    "id": "sensor1",
    "temperature": 25.5,
    "humidity": 60.0
}
```

### 2. WebSocket `/ws`
A WebSocket endpoint that streams real-time data to clients.

## Future Improvements

Here are some ideas for expanding the project:

- **Persistent Storage**: Integrate with a time-series database like InfluxDB to store historical sensor data.
- **Alerting System**: Configure threshold-based alerts (e.g., high temperature warnings) and send notifications via email, SMS, or push notifications.
- **Dashboard**: Build a frontend dashboard to visualize real-time and historical data using JavaScript libraries like Chart.js or D3.js.
- **Support Additional Sensors**: Add support for more environmental metrics (e.g., CO₂ levels, air quality).
  
## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
