# Commands I regulary use

# Create environment variables from a .env file
export $(grep -v '^#' ../.env | xargs)

# Database Migrations:
goose -dir ./database/schema postgres <connection_string> up
