# Go Blog Aggregator (gator)

Go Blog Aggregator, or `gator`, is a command-line application written in Go for fetching, storing, and managing RSS feeds. It allows users to register, follow their favorite feeds, and view aggregated content directly in the terminal.

This project serves as a practical example of building a Go application with a PostgreSQL backend, database migrations, and a clean, command-driven architecture.

## Features

*   **User Management:** Register new users.
*   **Feed Management:** Add new RSS feeds to the system.
*   **Feed Following:** Users can follow and unfollow specific feeds.
*   **Content Aggregation:** A worker process fetches the latest posts from all registered feeds.
*   **CLI Interface:** All functionality is exposed through a simple command-line interface.

## Prerequisites

Before you begin, ensure you have the following installed on your system:
*   **Go:** Version 1.21 or later.
*   **PostgreSQL:** A running PostgreSQL database instance.
*   **Goose:** The database migration tool. You can install it with:
    ```bash
    go install github.com/pressly/goose/v3/cmd/goose@latest
    ```

## Setup Instructions

Follow these steps to get the application running locally.

### 1. Clone the Repository

```bash
git clone https://github.com/arnicfil/go_blog_aggregator.git
cd go_blog_aggregator
```

### 2. Set Up the PostgreSQL Database

You need to create a database and a user for the application.

## Linux

```sql

-- Connect to PostgreSQL as a superuser (e.g., 'postgres')
psql -U postgres

-- Create a user (replace 'your_password' with a secure password)
CREATE USER gator_user WITH PASSWORD 'your_password';

-- Create the database
CREATE DATABASE gator OWNER gator_user;

-- Exit psql
\q
```

### 3. Configure the Application

The application requires a configuration file to store the database connection string and other settings. Create a configuration file (e.g., `config.json`) in your home directory.


Example `config.json`:
```json
{
  "db_url": "postgres://gator_user:your_password@localhost:5432/gator?sslmode=disable",
  "name": "arnicfil"
}
```
*The `"name"` field is used by commands that require a logged-in user context.*

### 4. Run Database Migrations

Use `goose` to apply the database schema. The migration files are located in the `sql/schema` directory.

```bash
# From the root of the project directory
goose -dir "sql/schema" postgres "postgres://gator_user:your_password@localhost:5432/gator?sslmode=disable" up
```
This will create the `users`, `feeds`, `posts`, and `feeds_follow` tables.

## Usage

### Build the Executable

Build the Go program to create an executable.

```bash
go build -o gator .
```

### Running Commands

All commands are run through the `gator` executable.

**General Syntax:**
```bash
./gator <command> [arguments...]
```

**Available Commands:**

*   **Register a new user:**
    ```bash
    ./gator register
    ```

*   **Add a new feed to the system:**
    ```bash
    # ./gator addFeed <feed_name> <feed_url>
    ./gator addFeed "My favourite blog" "https://www.myfavouriteblog.com/index.xml"
    ```

*   **Follow a feed:**
    ```bash
    # ./gator follow <feed_name>
    ./gator follow "My favourite blog"
    ```

*   **View feeds you are following:**
    ```bash
    ./gator following
    ```

*   **List all feeds in the system:**
    ```bash
    ./gator feeds
    ```

*   **Run the aggregator to fetch posts (for a single feed):**
    *This command can be expanded into a long-running worker.*
    ```bash
    ./gator agg
    ```
