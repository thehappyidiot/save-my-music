# Commands I regurly use

# Create environment variables from a .env file
export $(grep -v '^#' .env | xargs)

# Database
sqitch status
sqitch deploy
sqitch revert

