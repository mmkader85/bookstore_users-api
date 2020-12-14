# Bookstore users API - Golang (Gin + MVC pattern + MySQL)

### SQL Migrations
Required package: https://github.com/rubenv/sql-migrate

```
# Most useful commands
$ sql-migrate status
$ sql-migrate up -dryrun
$ sql-migrate up
$ sql-migrate down
```

### Curl for end-points

##### Create User
```
curl --location --request POST 'http://localhost:8000/user' \
--header 'Content-Type: application/json' \
--data-raw '{
    "first_name": "F.Name",
    "last_name": "L.Name",
    "email": "email@example.com",
    "status": "active",
    "password": "password123"
}'
```

##### Get User
```
curl --location --request GET 'http://localhost:8000/user/1' \
--header 'X-Private: true'
```

##### Search User
```
curl --location --request GET 'http://localhost:8000/internal/user/search?status=active' \
--header 'X-Private: true'
```

##### Update User
```
curl --location --request PUT 'http://localhost:8000/user/1' \
--header 'Content-Type: application/json' \
--data-raw '{
    "first_name": "F.Name",
    "last_name": "L.Name",
    "email": "email@example.com",
    "status": "active"
}'
```

##### Patch User
```
curl --location --request PATCH 'http://localhost:8000/user/12' \
--header 'Content-Type: application/json' \
--data-raw '{
    "first_name": "F.Name",
    "email": "email@example.com"
}'
```

##### Delete User
```
curl --location --request DELETE 'http://localhost:8000/user/1'
```
