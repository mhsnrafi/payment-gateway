## Payment Gateway
1. Build an API that allows a merchant:
   a. To process a payment through your payment gateway.
   b. To retrieve details of a previously made payment.
2. Build a bank simulator to test your payment gateway API.

## Tech/framework used
1. Go programming language
2. GORM library for database operations
3. Redis (for cache)
4. Postgres DB
5. Docker

These technologies were chosen to provide a scalable, performant, and maintainable solution for building a payment gateway API and bank simulator

## Problem Statement and Solution
The objective of this challenge is to create a payment gateway API that
will allow merchants to collect payments from their online customers.
The payment gateway must have the capability to process payments, 
which includes capturing information such as the card number, 
expiry date, amount, currency, and CVV. Additionally, merchants 
should be able to retrieve information on previously processed payments. The payment gateway will handle validating requests, storing card information, and forwarding payment requests to the acquiring bank, which will be simulated using the CKO bank simulator for testing purposes. The response from the payment gateway will include a masked card number, card details, and a status code indicating the success or failure of the payment

### Assumtions and Explanation how the function works
There are couple of API's Designed to solve this challenge and to secure the solution
1. Process Payment:
   1. The function "ProcessPayment" is used to process the payment from the customer. It takes in a struct called "PaymentRequest" as input and returns two outputs, "ProcessPaymentResponse" and an error. 
   2. A unique identifier for the payment is generated using the "uuid.New().String()" method and it is assigned to the "PaymentID" field. This identifier will be used to identify the payment in the future. 
   3. The payment information is then stored in the database using the "DbConnection.Create" method. This method inserts a new record in the database with the information from the "paymentInfo" struct. If there is an error while inserting the record, the error message is logged using the "logger.Error" method and an error is returned with the message "failed to store payment details". 
   4. The "simulateAcquiringBank" function is called to simulate the acquiring bank component which is responsible for processing the payment. This function returns the status of the payment. 
   5. If the returned status from the "simulateAcquiringBank" function is "constants.SUCCESS", the status of the payment in the database is updated using the "DbConnection.Model" method. If there is an error while updating the record, the error message is logged using the "logger.Error" method and an error is returned with the message "failed to store payment details". 
   6. The status of the payment is logged using the "logger.Info" method to keep track of the payment processing. 
   7. Finally, the "ProcessPaymentResponse" struct is returned with the status of the payment and the payment identifier. If there were no errors, the error value will be nil.

2. Retrieve Payment:
   1. The RetrievePayment function is used to retrieve the details of a specific payment, given the payment ID. The function returns a PaymentResponse struct and an error. 
   2. The first step in this function is to retrieve the payment details from the database using the DbConnection's Where method. If there's an error, the error message is logged using the logger.Error method and an error is returned with the message "failed to retrieve payment details". 
   3. Next, the function attempts to retrieve the payment details from cache using the GetPaymentDetailsFromCache function. If the payment details are found in cache, the function returns the payment details and a nil error. 
   4. If the payment details are not found in cache, the function masks the card number by calling the utils.MaskCardNumber function and populates a PaymentResponse struct with the payment details, including the masked card number, amount, currency, and status. 
   5. Finally, the function sets the payment details in cache using the SetPaymentDetails function and returns the PaymentResponse struct and a nil error
   
3. Take JWT authentication approach to make API's secure, How we generate access token and refresh that tokens how implement below is the approach
   1. The GenerateAccessTokens function creates two types of tokens for a given email.
   2. These tokens are an access token and a refresh token. 
   3. The access token has a set time to expire after a certain number of minutes and the refresh token has a set time to expire after a certain number of days. 
   4. The function calls a CreateToken function twice to create both the access and refresh tokens and returns them. 
   5. If there is an error during the creation of either token, the function returns an error.

4. API Collection Json already included in the project

5. To make system resiliat, using the retries mechanism circuit breaker aprroach
6. Proper logging is added used zap logger
7. To secure api, used JWT authentication mechanism


## How to run the solution, follow these steps:

Clone the repository: Run the following command to clone the repository to your local machine: 
```bash
git clone https://github.com/mhsnrafi/checkout.git

```

### Install dependencies: 
Change into the project directory and run go mod download to install the required dependencies.

### Start the API: 
Use the command docker-compose up to start the API.

### Configure credentials: 
The credentials required to connect to the database and run the API are described in the .env.local file.

### Run migrations: 
Connect to the database using the credentials from the .env.local file and run the migrations located in the Migrations/payments.sql file. This will create the required tables for the card balance and fraud.

###  Generate access token: 
Call the "Generate access token" endpoint to obtain an access token, which is required to authorize the API calls. Add the header "Bearer-Token" to each API request, using the access token obtained in this step.

### Use the API: 
The Postman collection is attached for easy use of the API.


### Endpoints
- Process Payment endpoint: `http://localhost:8080/v1/process-payment`
- Retrieve Payment endpoint: `http://localhost:8080/v1/get-payment?payment_id=de5f78db-c618-433b-b3f5-23c8a2519ea6`
- Generate access token endpoint: `http://localhost:8080/v1/auth/generate_access_token`
- Refresh Token endpoint: `http://localhost:8080/v1/auth/refresh`
- Prometheus endpoint: `http://localhost:9090`
- Grafana Endpoint: `http://localhost:3000`
- Swagger Endpoint: `http://localhost:8080/swagger/index.html#/`


## API Endpoints explanation
The API has a single endpoint to process payments.


### POST /auth/generate_access_token

This endpoint used to authenticate and validate the used is verified and generate access token details.

### Request Payload
```json
{
  "Email": "mohsin@checkout.example"
}
```

### Response Payload
```json
{
  "success": true,
  "data": {
    "token": {
      "access": {
        "expires": "2023-02-07 16:45:51",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI3MTQ5MjA0MDAyODQwNTM5MDkxIiwiZXhwIjoxNjc1Nzg0NzUxLCJpYXQiOjE2NzU2OTIzNTEsImVtYWlsIjoibW9oc2luQGNoZWNrb3V0LmV4YW1wbGUiLCJ0eXBlIjoiYWNjZXNzIn0.868WMt96788_ooL_rnAJHFD7XkcXYjywC6rJhRQdonM"
      },
      "refresh": {
        "expires": "2023-02-13 15:05:51",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIyMzEyMTY1MTg0NDk2NjYxNTY1IiwiZXhwIjoxNjc2Mjk3MTUxLCJpYXQiOjE2NzU2OTIzNTEsImVtYWlsIjoibW9oc2luQGNoZWNrb3V0LmV4YW1wbGUiLCJ0eXBlIjoicmVmcmVzaCJ9.KDn2k0rV_iyj9eEFTVDNVFbOm0C_zlcGHAvx79Bh0CA"
      }
    }
  }
}
```

### POST /auth/refresh

This endpoint used refresh the access token

### Request Payload
```json
{
  "Token": "refresh_token",
  "Email": "user@example.com"
}
```

### Response Payload
```json
{
  "success": true,
  "data": {
    "Email": "mohsin@checkout.example",
    "token": {
      "access": {
        "expires": "2023-02-08 14:35:36",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIzMDQ4NTg1MjU5ODkxMzM0MTU2IiwiZXhwIjoxNjc1ODYzMzM2LCJpYXQiOjE2NzU3NzA5MzYsImVtYWlsIjoibW9oc2luQGNoZWNrb3V0LmV4YW1wbGUiLCJ0eXBlIjoiYWNjZXNzIn0._2SZ9qfI5sosf6k23updXJYMqjld5UaqeWJcOoIbTcg"
      },
      "refresh": {
        "expires": "2023-02-14 12:55:36",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIzNzE2MTAwOTY1ODQzMTM3MDEwIiwiZXhwIjoxNjc2Mzc1NzM2LCJpYXQiOjE2NzU3NzA5MzYsImVtYWlsIjoibW9oc2luQGNoZWNrb3V0LmV4YW1wbGUiLCJ0eXBlIjoicmVmcmVzaCJ9.meYb31SNahNFdolWhNTf3cBSHb9Zm3z2dtu1uYoIM5A"
      }
    }
  }
}
```

### POST /process-payment

This endpoint processes a payment by storing the payment details
 in a database and simulating the acquiring bank.

### Request Payload
```json
{
 "card_number": "1234 1234 1234 1234",
 "exp_month": 12,
 "exp_year": 24,
 "cvv": 123,
 "amount": 10.0,
 "currency": "USD"
}
```

### Response Payload
```json
{
  "success": false,
  "data": {
    "Message": "Payment processed successfully",
    "Payment ID": "de5f78db-c618-433b-b3f5-23c8a2519ea6",
    "Payment Status": "Success"
  }
}
```

### GET /get-payment?payment_id={id}
This endpoint retrieve the previously made payments with masked card number


### Response Payload
```json
{
  "success": true,
  "data": {
    "Payment Details": {
      "payment_id": "de5f78db-c618-433b-b3f5-23c8a2519ea6",
      "masked_card_number": "**** **** **** 4567",
      "amount": 34,
      "currency": "USD",
      "status": "Success"
    }
  }
}
```


## Tests
The API includes a set of unit tests to ensure proper functionality. To run the tests, use the following command.
```bash
go test -v ./...
```

## API Documentation
To test the API endpoints directly from the documentation, making it easier to ensure that the API is working as expected build swagger api documentationa  user-friendly interface to quickly understand the API’s capabilities and functions
```bash
http://localhost:8080/swagger/index.html#/
```



Login into grafana using the creds mentioned in the `.env.local` file and import the dashboards using `dashbords/*.json`


## why cloud technologies you’d use and why.
1. To run the solution in a cloud environment, one can use a cloud-based PostgreSQL database such as Amazon RDS
2. The Docker image can be deployed on a cloud-based container orchestration platform such as Amazon ECS
3. To scale the solution, one can use a load balancer such as Amazon ELB
4. To monitor the solution, one can use a cloud-based monitoring solution such as Amazon CloudWatch

###### Screenshots of the projects
![](https://i.postimg.cc/R3jDvFVH/Screenshot-2023-02-07-at-3-42-45-PM.png)
![](https://i.postimg.cc/c6M05qXm/Screenshot-2023-02-07-at-3-43-03-PM.png)
![](https://i.postimg.cc/mcHGYKP7/Screenshot-2023-02-07-at-3-46-03-PM.png)
![](https://i.postimg.cc/bZMh2qwz/Screenshot-2023-02-07-at-3-46-10-PM.png)
![](https://i.postimg.cc/GHLrWMf1/Screenshot-2023-02-07-at-3-46-18-PM.png)
![](https://i.postimg.cc/1fvPsGxV/Screenshot-2023-02-07-at-3-46-24-PM.png)
![](https://i.postimg.cc/68TK38wj/Screenshot-2023-02-07-at-3-49-32-PM.png)
![](https://i.postimg.cc/D4FMtp54/Screenshot-2023-02-08-at-1-54-16-PM.png)
![](https://i.postimg.cc/G9C3Lc5G/Screenshot-2023-02-08-at-1-54-24-PM.png)
![](https://i.postimg.cc/bD8Xhvwv/Screenshot-2023-02-08-at-1-55-28-PM.png)

