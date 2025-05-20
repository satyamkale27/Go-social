# Go-Social

**Go-Social** is a social networking platform built with Go. It allows users to register, create posts, follow/unfollow other users, comment on posts, and manage user roles. The project includes features like user authentication, email invitations, and role-based access control. PostgreSQL is used as the database, and SendGrid is integrated for sending email notifications.

---

## 🚀 Features

- User registration, login, and authentication
- Role-based access control
- Follow/unfollow users
- Create, update, and delete posts
- Comment on posts
- Email activation and invitations via SendGrid
- User feed with pagination, search, and tag filtering

---

## 🛠 Installation Guide

### Prerequisites

- **Go (>=1.20)** – [Install Go](https://golang.org/dl/)
- **Docker & Docker Compose** – [Install Docker](https://www.docker.com/)
- **SendGrid API Key** – [Create a SendGrid account](https://sendgrid.com/)

---

### 🔧 Setup Instructions

1. **Clone the Repository**
   ```bash
   git clone https://github.com/satyamkale27/Go-social.git
   cd Go-social
   ```

2. **Set Up Environment Variables**

   Create a `.env` file in the root directory and add the following:

   ```env
   ADDR=:8080
   API_URL=http://localhost:8080
   FRONTEND_URL=http://localhost:4000
   DB_ADDR=postgres://admin:adminpassword@localhost/socialnetwork?sslmode=disable
   FROM_EMAIL=your-email@example.com
   SENDGRID_API_KEY=your-sendgrid-api-key
   TOKEN_SECRET=example
   ```

3. **Start PostgreSQL Database**

   ```bash
   docker-compose up -d
   ```

4. **Run Database Migrations**

   ```bash
   cd cmd/migrate
   go run main.go up
   ```

5. **Build and Run the Application**

   ```bash
   cd ../..
   go build -o bin/main ./cmd/api
   ./bin/main
   ```

6. **Access the API**

   The backend API is available at: `http://localhost:8080`

---

## 📡 API Endpoints

### 🧑 Authentication

- **POST** `/api/v1/users/register` – Register a new user
  ```json
  {
    "username": "example",
    "email": "example@example.com",
    "password": "password123"
  }
  ```

- **POST** `/api/v1/users/login` – Login and get token
  ```json
  {
    "email": "example@example.com",
    "password": "password123"
  }
  ```

### 👤 User Management

- **GET** `/api/v1/users/{userId}` – Get user details
- **POST** `/api/v1/users/{userId}/follow` – Follow a user
- **POST** `/api/v1/users/{userId}/unfollow` – Unfollow a user
- **GET** `/api/v1/users/activate/{token}` – Activate a user account

### 📝 Posts

- **POST** `/api/v1/posts` – Create a new post
  ```json
  {
    "title": "Post Title",
    "content": "Post Content",
    "tags": ["tag1", "tag2"]
  }
  ```

- **GET** `/api/v1/posts/{postId}` – Get post by ID

- **PUT** `/api/v1/posts/{postId}` – Update post
  ```json
  {
    "title": "Updated Title",
    "content": "Updated Content"
  }
  ```

- **DELETE** `/api/v1/posts/{postId}` – Delete post

### 💬 Comments

- **POST** `/api/v1/posts/{postId}/comments` – Add a comment
  ```json
  {
    "content": "This is a comment."
  }
  ```

- **GET** `/api/v1/posts/{postId}/comments` – Get all comments on a post

### 📰 Feed

- **GET** `/api/v1/feed` – Get user feed
    - Query Parameters:
        - `limit`: Number of posts to fetch (default: 20)
        - `offset`: Offset for pagination (default: 0)
        - `sort`: Sorting order (`asc` or `desc`)
        - `tags`: Filter by tags (comma-separated)
        - `search`: Search by title or content

---


## 🧠 Conclusion

Through this project, I learned how to build a backend using Go and explored Go concepts in depth. I applied these concepts to create APIs and implement the repository pattern for clean and maintainable code. I also gained hands-on experience with PostgreSQL and database migrations. Throughout development, I wrote SQL queries, debugged various issues, and dealt with complex database errors — especially the tricky ones that occur when migrations or queries fail. This project was a valuable learning experience that strengthened my understanding of backend development in Go.
