# Commands I regulary use

# Create environment variables from a .env file
export $(grep -v '^#' ../.env | xargs)

# Database
sqitch status
sqitch deploy
sqitch revert
sqitch add deployment_name --requires other_deployment --requires another_deployment -n 'Description.'
