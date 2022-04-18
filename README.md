# 🏦 Banking System

- A backend API to imitate how a bank transfer transaction takes place

- Handled data inconsistency and potential deadlocks for concurrent transactions

## Features 
- **Create and manage account**

  - Owner, balance and currency

- **Record all balance changes**

  - Create an account entry for each change

- **Money transfer transaction**
  - Perform money transaction between 2 accounts consistently within a transaction

## DB Schema
![Banking-System](https://user-images.githubusercontent.com/43776315/163681485-499ea22d-b2fd-49d9-acd6-0d23792cc164.png)

## 🏃‍♂️ Setup & How to Run

### Using docker-compose
- **Start the services**
```bash
make dc-up
```

- **Stop the servcies**
```bash
make dc-down
```

<hr>

### Using Docker (a bit more manual work)

- **Clone the repository**
```bash
git clone https://github.com/skamranahmed/banking-system.git
```

- **Create a `banking-system-network` in Docker**
```bash
make create-bank-network
```

- **Setup Postgres via Docker**
```bash
make setup-postgres
```

- **Create `bank` db and `bank_test` db in Postgres**
```bash
make create-db
```

- **Download the project dependencies via the below command**
```bash
make download
```

- **Migrate the DB Schema using the SQL Script**
For this I have used the `golang-migrate` library. You will first have to install this package. 

For `MacOS`
```bash
brew install golang-migrate
```

For other OS or kernels, refer this: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

Once the installation is done, run the below commands:

1. To migrate the database schema for the `bank` db:
```bash
make migrate-up
```

2. To migrate the database schema for the `bank_test` db:
```bash
make migrate-up-test
```

- **Configuring the environment variables**

Make sure you are in the `banking-system` root directory.
```bash
cp config/localConfigSample.yaml config/localConfig.yaml
```

- **Run the server**
1. Non-dockerized:
```bash
make run
```

2. Dockerized:

- Build the docker image for the backend service:
```bash
make build
```

- Run the docker image for the backend service:
```bash
make dockerized-server-run
```

<hr>