# In-Memory Key-Value Store
A simple, thread-safe, in-memory key-value store implemented in Go. It exposes a basic HTTP API to perform CRUD (Create, Read, Update, Delete) operations on string-based keys and values

## Features
*   Thread-safe in-memory key-value store
*   Generic key-value store implementation using Go generics
*   HTTP-based CRUD interface.
*   TTL (Time-To-Live) support for keys

## Getting Started

### Installation and Running
1.  Clone the repository:
    ```sh
    git clone https://github.com/dhruv883/key-value-store.git
    cd key-value-store
    ```

2.  Run the application:
    ```sh
    go run .
    ```
    The server will start and listen for requests on port `:3000`.

3. Usage
- Set: ```curl http://localhost:3000/put/key/value```
- Get: ```curl http://localhost:3000/get/key```
- Update: ```curl http://localhost:3000/update/key/value```
- Delete: ```curl http://localhost:3000/delete/key```
