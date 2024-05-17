# k0pt

k0pt is a command-line tool for managing Kubernetes clusters and namespaces with ease.

## Features

- Use a specific namespace
- Delete a specific namespace
- Start a specific cluster
- Stop a specific cluster
- Get the status of a specific cluster
- Get admin credentials for a specific cluster
- Analyze resource usage across the cluster
- Calculate potential cost savings
- Optimize resource allocation across the cluster

## Installation

### Using Docker

1. Build the Docker image:

    ```sh
    docker build -t k0pt .
    ```

2. Run the Docker container:

    ```sh
    docker run --rm k0pt help
    ```

### From Source

1. Clone the repository:

    ```sh
    git clone https://github.com/yourusername/k0pt.git
    cd k0pt
    ```

2. Build the project:

    ```sh
    go build -o k0pt cmd/main.go
    ```

3. Run the binary:

    ```sh
    ./k0pt help
    ```

## Usage

Here are some usage examples:

- Display help:

    ```sh
    k0pt help
    ```

- Use a specific namespace:

    ```sh
    k0pt use-namespace <namespace>
    ```

- Delete a specific namespace:

    ```sh
    k0pt delete-namespace <namespace>
    ```

## Contributing

We welcome contributions! Please read `CONTRIBUTING.md` for guidelines on how to contribute to this project.

## Code of Conduct

This project adheres to the Contributor Covenant [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
