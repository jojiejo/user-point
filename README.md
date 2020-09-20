## About This Application (Telunjuk User Point)

Telunjuk User Point is a application to manage user point. This application is developed to complete the technical back-end test at Telunjuk. It has several features, such as:

- Create user
- Manage user point
- Delete user

While developing this application, I use Gin as Go Framework. 

## Installation

### System Requirement

- Go >= 1.1.0 (Developed in 1.1.3)
- MySQL >= 8.0.0 (Developed in 10.4.8-MariaDB)

### Steps

Install [golang](https://golang.org/dl/) on your server.

After you clone, a new folder should be created in your current location.

```bash
cd user-point-app
```

Install the project dependencies using go mod.

```bash
go mod tidy
```

Create a copy of .env file. Afterwards, you have to fill up .env with your environment variables. (For example, database server, port and name)

```bash
cp .env.example .env
```

Create an empty database with name user_point and do the following command to migrate the tables needed.

After successfully migrating the tables, you can simply run using this command.

```bash
go run main.go
```

Or you can build the application by using this command below.

```bash
go build main.go
```
