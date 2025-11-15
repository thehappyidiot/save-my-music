# Commands I regulary use

# Create environment variables from a .env file
export $(grep -v '^#' .env | xargs)

# Database Migrations:
goose up
goose down 

# Generate SQLC:
sqlc generate

# Go:
go test ./...
go build -o smm cmd/api/main.go && ./smm