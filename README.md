# go-grpc-micro-bank-server

This is the server component of the go-grpc-micro-bank project. It provides the backend functionality for the banking application using gRPC.

## Prerequisites

1. **Go (Golang)**
   - Ensure Go is installed on your machine. You can download it from [Go's official site](https://golang.org/dl/).
   - Verify the installation by running:
     ```bash
     go version
     ```

2. **golang-migrate**
   - The project uses `golang-migrate` to manage database migrations.
   - To install `golang-migrate`, follow these steps:

     **Using Homebrew (macOS and Linux)**
     ```bash
     brew install golang-migrate
     ```

     **Using Docker**
     ```bash
     docker run -v $(pwd)/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database <DATABASE_URL> up
     ```

     **Download Pre-compiled Binaries**
     - Visit the [golang-migrate releases page](https://github.com/golang-migrate/migrate/releases).
     - Download the appropriate binary for your OS.
     - Add the binary to your system's PATH.

   - Verify the installation by running:
     ```bash
     migrate -version
     ```

## Installation

1. Clone the repository:

```bash
git clone https://github.com/fajaramaulana/go-grpc-micro-bank-server.git
```

2. Install the dependencies:

```bash
go mod download
```

## Usage

To start the server, run the following command:

```bash
go run main.go
```

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).
