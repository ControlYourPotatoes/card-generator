# Database Configuration Guide

## Setting Up Your Environment

The card generator application requires a PostgreSQL database connection. We've updated the configuration to support both local PostgreSQL instances and Supabase.

### Create a .env File

Create a file named `.env` in the root directory of the project with the following content:

```
# Database Configuration for Supabase
DB_HOST=qkkixxahqhuhwvtnqokb.supabase.co
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=postgres
DB_SSLMODE=require
USE_DB_PREFIX=true

# Test settings
RUN_DB_TESTS=true
```

Replace `your_password_here` with your actual database password.

### Configuration Options

| Variable | Description | Default |
|----------|-------------|---------|
| DB_HOST | Database hostname | localhost |
| DB_PORT | Database port | 5432 |
| DB_USER | Database username | postgres |
| DB_PASSWORD | Database password | postgres |
| DB_NAME | Database name | postgres |
| DB_SSLMODE | SSL mode (disable, require, etc.) | require |
| USE_DB_PREFIX | Whether to prefix host with "db." (for Supabase) | true |
| RUN_DB_TESTS | Whether to run database tests | - |

### Local Development Setup

For local development with PostgreSQL:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=card_generator
DB_SSLMODE=disable
USE_DB_PREFIX=false

RUN_DB_TESTS=true
```

### Running Tests

To run tests that use the database:

1. Ensure your `.env` file is properly set up
2. Make sure `RUN_DB_TESTS=true` is set
3. Run tests with: `go test -v ./internal/storage/database`

To skip database tests:

```
RUN_DB_TESTS=false
```

Or run without the database tests:

```
go test -v -short ./...
```

## Troubleshooting

If you encounter connection issues:

1. Verify your database credentials
2. Check if your IP is allowed to connect to the database
3. Ensure the correct SSL mode is set (cloud databases typically require `require`)
4. For Supabase, verify that `USE_DB_PREFIX=true` is set

## Security Note

Never commit your `.env` file or database credentials to version control. The `.env` file is included in `.gitignore` by default.
