# Go Hiring Challenge

This repository contains a Go application for managing products and their prices, including functionalities for CRUD operations and seeding the database with initial data.

## Project Structure

1. **cmd/**: Contains the main application and seed command entry points.

   - `server/main.go`: The main application entry point, serves the REST API.
   - `seed/main.go`: Command to seed the database with initial product data.

2. **app/**: Contains the application logic.
3. **sql/**: Contains a very simple database migration scripts setup.
4. **models/**: Contains the data models and repositories used in the application.
5. `.env`: Environment variables file for configuration.

## Setup Code Repository

1. Create a github/bitbucket/gitlab repository and push all this code as-is.
2. Create a new branch, and provide a pull-request against the main branch with your changes. Instructions to follow.

## Application Setup

- Ensure you have Go installed on your machine.
- Ensure you have Docker installed on your machine.
- Important makefile targets:
  - `make tidy`: will install all dependencies.
  - `make docker-up`: will start the required infrastructure services via docker containers.
  - `make seed`: ⚠️ Will destroy and re-create the database tables.
  - `make test`: Will run the tests.
  - `make run`: Will start the application.
  - `make docker-down`: Will stop the docker containers.

## API Testing with Postman

### Setup Instructions

1. **Install Postman**: Download and install Postman from [https://www.postman.com/downloads/](https://www.postman.com/downloads/)

2. **Import the Collection**:
   - Open Postman
   - Click **Import** in the top left
   - Select **Upload Files** and choose `postman tests/postman collection.json`
   - The collection will be imported with all endpoints and tests

3. **Configure Environment Variables**:
   - In the imported collection, you'll see variables defined:
     - `base_url`: Default is `http://localhost:8080` (adjust if your server runs on a different port)
     - `product_code`: Default is `PROD001` (adjust to test different products)
     - `test_category_code`: Default is `test-category` (used for creating test categories)

### Running Tests

1. **Start the Application**:
   ```bash
   make docker-up
   make seed
   make run
   ```

2. **Run Tests in Postman**:
   - Open the imported collection "Go Hiring Challenge API"
   - Click the **Run** button (arrow icon) or use **Runner** from the menu
   - Select all requests or specific test suites to run
   - View results in the Test Results panel

### Available Test Endpoints

#### Catalog Endpoints
- **GET /catalog** - Retrieve all products with pagination (default: offset=0, limit=10)
- **GET /catalog?offset=0&limit=5** - Retrieve products with custom pagination
- **GET /catalog/{code}** - Retrieve detailed information for a specific product
- **GET /catalog/INVALID_CODE** - Test 404 error handling

#### Categories Endpoints
- **GET /categories** - Retrieve all available categories
- **POST /categories** - Create a new category
- **POST /categories** (with missing fields) - Test validation error handling

### Test Coverage

Each request includes automated tests that verify:
- Correct HTTP status codes
- Proper response structure and format
- Required fields in responses
- Error handling for invalid inputs

Follow up for the assignemnt here: [ASSIGNMENT.md](ASSIGNMENT.md)
