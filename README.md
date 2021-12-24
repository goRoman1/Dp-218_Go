# Scooter Service

The application is similar to uber service. Where you as a:
* supplier  - can get your scooters for rent.
* user - can take a scooter for rent.

# Run

The application runs with docker.
```
docker-compose up --build -d
```
Server runs on http://localhost:8080/

# How to start the trip

On the page ```http://localhost:8080/customer/map``` you can choose a departure station.  

On click "show station" button you will move to the ```http://localhost:8080/start-trip/{station_id}``` page.
Which shows you all available scooters on the chosen station. Here you also choose the destination station.  

"Start trip" button will start your trip with chosen scooter to the chosen station.

Information about trips will be written to the database table - "Orders".
