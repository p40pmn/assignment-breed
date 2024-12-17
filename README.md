# Breed Inquiry Service

This repository provides a backend service for managing and querying breed data, built with Go (Golang) using the **Echo** framework and PostgreSQL as the database.

---

## API Specification

The full API documentation and specifications are available on Postman.  
You can view or test the APIs using the link below:

🔗 [Postman API Documentation](<https://documenter.getpostman.com/view/25020935/2sAYJ1jMqK>)

---

### Steps to Run Locally

1. **Clone the repository:**
    ```bash
      git clone git@github.com:p40pmn/assignment-breed.git
      cd assignment-breed
    ```

2. **Install dependencies:**
    ```bash
      go mod tidy
    ```
  
3. **Set up environment variables:**
    ```env
      PORT=
      DATABASE_URL=
    ```

4. **Run the application:**
    ```bash
      go run cmd/main.go
    ```

1. **Test the API:** Use Postman or curl to verify the endpoints:
    ```bash
      curl -X POST "http://localhost:8080/breed-inquiry" -H "Content-Type: application/json" -d '{"keyword": "example"}'
    ```

