# Bookstore users API - Golang (Gin + MVC pattern + MySQL)

_**!!! This is not a production ready application. Developed while studying !!!**_

### Docker
1. Move .env.dist to .env & update the configs, if required
2. Move dbconfig.yml.dist to dbconfig.yml & update the configs, if required
    * In local, MySQL host should be name of the database container
3. `$ docker-compose up -d --build`
    * `api` may not be running as it requires the database. DB creation takes a while, so it will attempt an auto restart.
4. SSH into the api container and apply db migrations
    * `$ docker exec -it bookstore_users-api sh`
    * Apply the migrations `$ sql-migrate up`
    * Verify the migrations `$ sql-migrate status`

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
curl --location --request PATCH 'http://localhost:8000/user/1' \
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
