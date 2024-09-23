# DecideYourDinner

**Run the command in the root project directory**

- To deploy the project on Docker

```
docker-compose up --build
```

- To manually enter the Docker commands to build and run

```
docker build -t dweb:ori .
docker run -d --name web --network mynetwork -e DB_HOST=mysql-db -e DB_PORT=3306 -e DB_USER=admin -e DB_PASSWORD=admin -e DB_NAME=web -p 8080:8080 dweb:ori
docker run -d --name mysql-db --network mynetwork -e MYSQL_ROOT_PASSWORD='Root2147' -e MYSQL_DATABASE=web -e MYSQL_USER=admin -e MYSQL_PASSWORD=admin -p 3307:3306 -v C:/Users/TimHe/Documents/Go/dinner/DecideYourDinner/sql/init_web.sql:/docker-entrypoint-initdb.d/init.sql mysql:5.7
```

- To enter mysql container

```
docker exec -it mysql-db mysql -uadmin -padmin
```
