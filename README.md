# Deck Cards 

The Application for this challenge is written in Go Programming language using Gin framework. The workflow for running this challenge is written using a Makefile. The application is using on port 8000

### The application is created using: 

1. Go - Gin Framework 
2. Mongodb - database
3. It runs using Docker and Docker compose 
4. It uses a Makefile for running the applicatiomn 

### How to run the app. 
1. Using docker-compose

### Running the application using a makefile with docker-compose

1. Command to run / start the application. Run `make compose-up` 
2. To stop the app using docker-compose. Run `make compose-down`

### I created the following API for the following 
1. Create a new Deck 
2. Open a Deck 
3. Draw a Card

## Create a new Deck

This endpoints creates a deck. I set the remaining cards to be `52` by default but if the cards is supplied the remaining will be the number of cards supplied. 

| Endpoint                                    | Description                                                                    | Method
| -----------------                           | -----------------------------------------------------------------------------  | ---------
|`127.0.0.1:8000/create/deck`                | This end points creates the deck without shuffle. By default shuffle is false  | POST
|`127.0.0.1:8000/create/deck?shuffled=true`  | This end points creates the deck by shuffle. Shuffle is set to true            | POST
|`127.0.0.1:8000/create/deck?cards=AS,KD,AC,2C,KH` |This endpount creates a partial deck cards with supplied cards             | POST
|`127.0.0.1:8000/create/deck?cards=AS,KD,AC,2C,KH&?shuffled=true` | This endpount creates a partial deck cards with supplied cards and it suffles | POST

## Open a Deck 

This endpoints show the details of a deck card 

| Endpoint                                    | Description                                                                    | Method
| -----------------                           | -----------------------------------------------------------------------------  | ---------
| `127.0.0.1:8000/open/deck/:deck_id`         | This end points gets the details of the deck. Where deck_id is the deck id. Example 127.0.0.1:8000/open/deck/b1a7faaf-48b1-4ef1-890d-da9b6bd345dd | GET


## Draw a Deck 

This endpoints is used to draw a card. By default if no count is added in the query string. The default `count = 1`. 
The count selects the top cards from the database. If the count is more than the remaining cards in the database, user will not be able to draw cards in the database.

| Endpoint                                    | Description                                                                    | Method
| -----------------                           | -----------------------------------------------------------------------------  | ---------
|`127.0.0.1:8000/draw/deck/:deck_id`        | This end points is used to draw a card from the database. Example 127.0.0.1:8000/draw/deck/b1a7faaf-48b1-4ef1-890d-da9b6bd345dd. If count is not specified in the query string. The default count = 1  | GET
|`127.0.0.1:8000/draw/deck/:deck_id?count=10` | This end points is used to draw a card from the database. Example 127.0.0.1:8000/draw/deck/35721961-0685-4984-be16-52b29c30e15c?count=10. | GET


## Check Healthcheck

This endpoints show the health of the application. It just to check if the application is healthy 

| Endpoint                                    | Description                                                                    | Method
| -----------------                           | -----------------------------------------------------------------------------  | ---------
| `127.0.0.1:8000/health`         | This end points checks the health of the application | GET


Where `127.0.0.1` is the localhost and `8000` is the port

I also enabled CORS.

I used postman to verify that all the endpoints are working fine.


## Unit Testing

I have currently written a unit test. To check that the creation of deck works. 

To run the test, pls run command `go test`