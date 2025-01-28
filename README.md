# Bookstore API

## Overview
A Go-based service providing basic CRUD operations for Books, Authors, Customers, Orders, and automatic Sales Reports.

## How to Run
1. Install Go 1.23.4 or higher.
2. Run `go mod tidy` to install dependencies.
3. Run `go run api/bookstore.go`.
4. The server listens on port 8080 by default.

## Endpoints

### Books
- **POST /books** — Create a book.  
- **GET /books** — List/search books.  
- **GET /books/{id}** — Get a single book.  
- **PUT /books/{id}** — Update a book.  
- **DELETE /books/{id}** — Delete a book.

### Authors
- **POST /authors** — Create an author.  
- **GET /authors** — List/search authors.  
- **GET /authors/{id}** — Get a single author.  
- **PUT /authors/{id}** — Update an author.  
- **DELETE /authors/{id}** — Delete an author.

### Customers
- **POST /customers** — Create customer.  
- **GET /customers** — List/search customers.  
- **GET /customers/{id}** — Get a single customer.  
- **PUT /customers/{id}** — Update a customer.  
- **DELETE /customers/{id}** — Delete a customer.

### Orders
- **POST /orders** — Create an order.  
- **GET /orders** — List/search orders.  
- **GET /orders/{id}** — Get a single order.  
- **PUT /orders/{id}** — Update an order.  
- **DELETE /orders/{id}** — Delete an order.

### Reports
- **GET /reports** — Aggregate and return all JSON sales reports.

## Usage Steps
1. Start the server.  
2. Use any REST client (e.g., cURL or Postman).  
3. Send HTTP requests to the endpoints above.  
4. Check generated JSON files in the `./reports` folder for daily sales reports.

### Additional Features

#### 1. **Periodic Sales Report Generator**
- The application includes a periodic background task that runs every 24 hours.
- This task aggregates sales data, generating a JSON report with the following details:
  - **Total Revenue**: The sum of all sales within the last 24 hours.
  - **Total Orders**: The total number of orders placed.
  - **Total Books Sold**: A cumulative count of books sold.
  - **Top-Selling Books**: A list of books with the highest sales during the period.
- Each report is saved in the `reports` directory with filenames in the format `report_YYYYMMDD_HHMM.json`.
- The sales report generation runs in the background, ensuring it doesn’t interfere with the main API responsiveness.

#### 2. **Logging**
- A comprehensive logging mechanism has been implemented to:
  - Record API requests and responses.
  - Log significant events such as order placements and the execution of background tasks.
  - Capture errors, including failed requests and system anomalies.
- Logs are stored in the `api.log` file with timestamps for easy debugging and monitoring.

#### 3. **Manual Testing (Postman as a client)**
Below are some examples of tests I have done using Postman
- **Create a Book**
  - **Endpoint**: `POST /books`
    - ![Posting a book](./assests/book_post.jpg)
    - ![Trying with an invalid author id](./assests/book-invalid_author_id.JPG)
  - **Endpoint**: `GET /books?genre=romance`  
    - ![Getting books by genre](./assests/books_by_genre.JPG)
  - **Endpoint**: `PUT /authors/1`  
    - ![Updating author](./assests/updating_author.JPG)
  - **Endpoint**: `DELETE /authors/1`  
    - ![Deleting author](./assests/deleting_author.JPG)
  - **Endpoint**: `GET /reports`  
    - ![Getting reports](./assests/report.JPG)


 

- **Test Logging**

  - ![Logging requests](./assests/log.JPG)


### Additional Notes 🛠️✨

This project is just the beginning of an exciting journey! 🚀 For future updates, I plan to take the following steps to make the application even more robust and feature-rich:

- **Enhanced Logging**: Include response status logging to provide better visibility into API behavior and streamline debugging. 📝✅
- **Expanding Concurrency**: Leverage concurrency in other areas of the API to further optimize performance and scalability. ⚡🧵
- **Authentication & Authorization**: Implement robust user authentication (e.g., JWT) and role-based access control to secure endpoints. 🔐🔏
- **Database Migration**: Transition to a powerful database engine like PostgreSQL for reliable, scalable, and persistent data storage. 🗄️🐘
- **Containerization**: Embrace Docker to containerize the application, ensuring consistent deployment across environments and paving the way for seamless CI/CD integration. 🐳✨

Stay tuned for these updates as the project evolves into a production-ready masterpiece! 💡🚀

