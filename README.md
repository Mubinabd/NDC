# 📜 Post Service

Post Service is a microservice responsible for handling posts. It allows users to create, retrieve, update, and delete posts. The service communicates via gRPC and is accessible through an API Gateway.

## 🚀 Features
- Create, retrieve, update, and delete posts.
- gRPC-based service communication.
- JWT-based authentication and authorization.
- Logging of all requests and responses.
- API Gateway integration for HTTP access.

## 🛠️ Technologies Used
- **Programming Language**: Go
- **Frameworks**: gRPC, Gin
- **Database**: PostgreSQL
- **Authentication**: JWT
- **Logging**: Logrus or Zerolog
- **Containerization**: Docker

## 📌 API Endpoints (through API Gateway)
| Method | Endpoint         | Description                      |
|--------|-----------------|----------------------------------|
| POST   | `/api/posts`    | Create a new post               |
| GET    | `/api/posts`    | Get all posts                   |
| GET    | `/api/posts/{id}` | Get a specific post by ID     |
| PUT    | `/api/posts/{id}` | Update a post                 |
| DELETE | `/api/posts/{id}` | Delete a post                 |

Getting Started

To get started with Post-Service, follow these steps:
Prerequisites

Make sure you have the following tools installed:

    Golang (version 1.23 or higher)
    Docker (for containerization)
    PostgreSQL (for database management)
    Swagger (for API documentation)

Installation

Clone the repository:

git clone https://github.com/Mubinabd/NDC.git

Navigate to the project directory:

cd NDC

Start the project using Docker Compose:

docker-compose up --build

Once the project is running, you can access the API documentation at:

http://localhost:8080/api/swagger/index.html

    API Documentation: Access and manage project-specific API documentation through Swagger URLs.

Contributing

We welcome contributions! To get involved:

    Fork the repository.
    Create a new feature branch: git checkout -b feature/new-feature.
    Commit your changes: git commit -m 'Add some feature'.
    Push to the branch: git push origin feature/new-feature.
    Open a Pull Request for review.

License

This project is licensed under the MIT License – see the LICENSE file for details.
