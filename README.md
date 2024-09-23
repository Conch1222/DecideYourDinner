# DecideYourDinner

This project will help you decide what you want to eat for dinner today based on your inputs and your location.

## Usage

[Link](https://drive.google.com/file/d/1YLZE-OIB5QN8LX_CDQnP9RaMAdXSlU2c/view?usp=sharing)

## Setup

1. create a [GCP API Key](https://cloud.google.com/api-keys/docs/overview) and enable [Places API](https://developers.google.com/maps/documentation/places/web-service?hl=zh-tw) and [Geolocation API](https://developers.google.com/maps/documentation/geocoding/overview)

2. create Key.txt in the **File** folder and put the api key in it

3. deploy the project on Docker

```
docker-compose up --build
```

4. Once done, enter `http://localhost:8080/login` into the webpage to decide your dinner

**NOTE: Run all Docker commands in the root project directory**

### Default Account

- username: admin
- password: admin

## Test

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
