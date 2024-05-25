# Fintech Server-Side Project

Welcome to the Fintech Server-Side Project repository! This project is a robust and efficient backend system built using Go (Golang). It follows the repository design pattern and integrates MySQL for data management. The project also implements advanced Go features such as channels and asynchronous operations.

## Table of Contents
- [Features](#features)
- [Technologies Used](#technologies-used)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)
- [Acknowledgements](#acknowledgements)

## Features

- **Efficient Concurrency Handling**: Utilizes Go’s channels and goroutines for managing asynchronous operations and concurrent tasks.
- **Repository Design Pattern**: Ensures clean, maintainable, and testable code architecture.
- **User Authentication**: Implements secure user authentication and authorization mechanisms.
- **Wallet Management**: Features for creating, activating, refunding, and withdrawing from wallets.
- **Transaction Handling**: Robust handling of transactions, including logging and error management.
- **MySQL Integration**: Efficiently manages data storage and retrieval using MySQL.

## Technologies Used

- **Programming Language**: Go (Golang)
- **Database**: MySQL
- **Frameworks and Libraries**: 
  - Gin (for HTTP web framework)
  - go-playground/validator (for data validation)
- **Other Tools**: Docker (for containerization)

## Getting Started

### Prerequisites

Before you begin, ensure you have the following installed on your system:

- Go (version 1.16 or higher)
- MySQL
- Docker (optional, for containerization)

### Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/yourusername/fintech-server-side.git
    cd fintech-server-side
    ```

2. Install Go dependencies:

    ```sh
    go mod download
    ```

3. Set up your environment variables. Create a `.env` file in the root directory and add the following:

    ```sh
    DB_HOST=your_db_host
    DB_PORT=your_db_port
    DB_USER=your_db_user
    DB_PASSWORD=your_db_password
    DB_NAME=your_db_name
    ```

4. Initialize the database:

    ```sh
    go run scripts/init_db.go
    ```

5. Run the server:

    ```sh
    go run main.go
    ```

## Usage

Once the server is running, you can interact with it via API endpoints. The following are some of the key endpoints:

- **Create a new User**: `POST api/auth/register`
- **Activate a wallet**: `POST api/wallet/activate`
- **Refund a wallet**: `POST api/wallet/refund`
- **Withdraw from a wallet**: `POST api/wallet/withdraw`

You will find the rest of api's in swagger docs.

Happy coding! 🚀
