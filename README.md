# Komodo REST API

Golang Ecommerce API using MySQL DB and Fiber

## How to install

### For local

1. Install dependencies.

```golang
go get ./...
```

2. Add Config at `.env` file in root folder.

```golang
PORT = 9000
JWT_SECRET = "secret"
MYSQL_USER = "root"
MYSQL_PASSWORD = ""
MYSQL_HOST = "localhost"
MYSQL_PORT = 3306
MYSQL_DATABASE = "ecommerce_go"
```

3. Import table and data using `schema.sql` and `data.sql` at `./scripts` folder.

### Using Docker Compose

1. Setting enviroment variable using `docker-compose.yml` at root folder.
2. If MYSQL_HOST and MYSQL_PORT are changed, make sure to change `mysql:3306` part in CMD at dockerfile at root folder.

```golang
CMD ["./wait-for.sh" , "mysql:3306" , "--timeout=300" , "--" , "./komodo-backend"]
```

## How to use

### For local

1. Use this command at root folder to start the app.

```golang
go run main.go
```

2. Run test using this command at root folder.

```golang
go test ./...
```

### Using Docker Compose

1. Clear docker compose

```golang
docker-compose rm -v
```

2. Run this command to start docker compose.

```golang
docker-compose up --build
```

## App Endpoints

| No  | Path                | Method | Payload                                                                                                                                                                                                                                                                                                                     | Description                                        |
| --- | ------------------- | ------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------- |
| 1   | /buyers/register    | POST   | <pre lang="json">{<br> "email":"buyer@mail.com",<br> "name":"john buyer",<br> "password":"12345",<br> "sendingAddress":"Jl jalan"<br>}</pre>                                                                                                                                                                                | Buyer register                                     |
| 2   | /buyers/login       | POST   | <pre lang="json">{<br> "email":"buyer@mail.com",<br> "password":"12345"<br>}</pre>                                                                                                                                                                                                                                          | Buyer login                                        |
| 3   | /sellers/register   | POST   | <pre lang="json">{<br> "email":"seller@mail.com",<br> "name":"john seller",<br> "password":"12345",<br> "pickupAddress":"Jl jalan"<br>}</pre>                                                                                                                                                                               | Seller register                                    |
| 4   | /sellers/login      | POST   | <pre lang="json">{<br> "email":"seller@mail.com",<br> "password":"12345"<br>}</pre>                                                                                                                                                                                                                                         | Seller login                                       |
| 5   | /products           | GET    |                                                                                                                                                                                                                                                                                                                             | Get all products                                   |
| 6   | /products           | POST   | <pre lang="json">{<br> "name":"pro1",<br> "description":"check",<br> "price":91051551.13,<br> "sellerId":1<br>}</pre>                                                                                                                                                                                                       | Create a product                                   |
| 7   | /orders/find/byuser | GET    |                                                                                                                                                                                                                                                                                                                             | Get all orders by buyer/seller id inside JWT token |
| 8   | /orders             | POST   | <pre lang="json">{<br>"buyerId": 6,<br>"sellerId": 3,<br>"deliverySourceAddress": "source",<br>"deliveryDestinationAddress": "destination",<br>"items": [<br>{<br>"productId": 1,<br>"quantity": 12<br>},<br>{<br>"productId": 2,<br>"quantity": 8<br>},<br>{<br>"productId": 3,<br>"quantity": 10<br>}<br>]<br><br>}</pre> | Create a order                                     |
| 9   | /orders/:id/accept  | PUT    |                                                                                                                                                                                                                                                                                                                             | Accept order                                       |

## Endpoints security

| No  | Path                | Method | Need Login? | Access by |
| --- | ------------------- | ------ | ----------- | --------- |
| 1   | /buyers/register    | POST   | no          | all       |
| 2   | /buyers/login       | POST   | no          | all       |
| 3   | /sellers/register   | POST   | no          | all       |
| 4   | /sellers/login      | POST   | no          | all       |
| 5   | /products           | GET    | no          | all       |
| 6   | /products           | POST   | yes         | seller    |
| 7   | /orders/find/byuser | POST   | yes         | all       |
| 8   | /orders             | POST   | yes         | buyer     |
| 9   | /orders/:id/accept  | PUT    | yes         | seller    |

## Database Design
![db_design](https://user-images.githubusercontent.com/28037175/116769487-b67d8e80-aa66-11eb-8820-cfac90be9eeb.png)
