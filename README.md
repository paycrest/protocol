# Paycrest Protocol
A P2P powered crypto-to-cash exchange protocol.

This is a work in progress. Follow [this guide](https://www.freecodecamp.org/news/how-to-write-a-good-readme-file/) if you want to contribute to this README.

## What is Paycrest Protocol?
T.B.D.

## Architecture
T.B.D.

## Development Setup
To set up your development environment for Paycrest Protocol, follow these steps:

1. Setup Paycrest Protocol repo in your local machine.
```bash
# clone the repo
git clone https://github.com/paycrest/protocol.git
cd protocol

# copy and update enviroment variables
cp .env.example .env
```

2. Start the development environment:

**With Docker**
```bash
# build the image
docker-compose build

# run containers
docker-compose up -d
```

*OR*

**Without docker**
*(Ensure you have PostgreSQL 13 or higher installed and running with a database created and the credentials updated in the .env file.)*

```bash
# Tap brew core repository. Might take some time.
brew tap homebrew/core

# Install PostgreSQL via homebrew
brew install postgresql

# start postgresql service
brew services start postgresql

# Install Redis via homebrew
brew install redis

# start redis service
brew services start redis
```

Setup database

```bash
createdb paycrest
createuser postgres

# connect to paycrest database
psql paycrest

# enable extended display
\x

```

Setup redis

```bash
# Run redis on the current TTY
redis-server

# connect to redis instance
redis-cli

# set new password
CONFIG SET requirepass "password"
```

```bash
# install go dependencies
go mod download

# install Air for live reloading
curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s

# Or 

go install github.com/cosmtrek/air@latest

# start server
air -v
```

That's it! The server will now be running at http://localhost:8000. You can use an API testing tool like Postman or cURL to interact with the API.

**Seed database**
```bash
# run seed-db script
go run 
```

## Usage
TODO: Add API documentation.

## Roadmap
T.B.D.

## Contributing
We welcome contributions to the Paycrest Protocol! To get started, follow these steps:

1. Fork the repository by clicking the "Fork" button on the top right corner of the repository page.

2. Clone your forked repository to your local machine:
```bash
$ git clone https://github.com/your-username/protocol.git
```
3. Setup the precommit hook by running 
```bash
sh init.sh
```

4. Make changes to the codebase and commit them.

5. Push your changes to your forked repository:
```bash
$ git push origin main
```

5. Create a pull request by visiting the [repository page](https://github.com/paycrest/protocol) and clicking the "New pull request" button.

Our team will review your pull request and work with you to get it merged into the main branch of the repository. 

If you encounter any issues or have questions, feel free to open an issue on the repository or contact us via email on support@paycrest.io

## Testing
Paycrest Protocol uses a combination of unit tests and integration tests to ensure the reliability of the codebase.

To run the tests, run the following command:
```bash
# run all tests
go test ./...

# run a specific test
go test ./path/to/test/file
```

It is mandatory that you write tests for any new features or changes you make to the codebase. Only PRs that include passing tests will be accepted.

## License
[Affero General Public License v3.0](https://choosealicense.com/licenses/agpl-3.0/)
