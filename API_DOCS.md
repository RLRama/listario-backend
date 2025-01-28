# API documentation

## Base URL

`https://localhost:8080/` or wherever it's hosted

## Endpoints

### User endpoints

#### Register

- Curl command:

   ```
   curl -X POST http://localhost:8080/register \
     -H "Content-Type: application/json" \
     -d '{
           "username": "testuser",
           "password": "Password123!",
           "email": "testuser@example.com"
         }'
   ```

#### Login

- Curl command:

   ```
   curl -X POST http://localhost:8080/user/login \
     -H "Content-Type: application/json" \
     -d '{
           "identifier": "testuser",
           "password": "Password123!"
         }'
   ```

#### Update username and/or email

- Curl command:

   ```
   curl -X PUT http://localhost:8080/user/protected/update   -H "Content-Type: application/json"   -b "token=<token_here>"   -d '{"username":"john","email":"testuser10@example.com"}'
   ```

- > `username` and `email` in the payload can be left blank if they are not to be changed.

#### Update password

- Curl command:

   ```
   curl -X PUT http://localhost:8080/user/protected/update-password   -H "Content-Type: application/json"   -b "token=<token_here>"   -d '{"current_password":"old_password123","new_password":"new_password123"}'
   ```

### Test endpoints

#### Hello world!

- Curl command:

   ```
   curl -X GET http://localhost:8080/test/hello
   ```

#### Test DB connection

- Curl command:

   ```
   curl -X GET http://localhost:8080/test/db-connection
   ```

#### Protected test endpoint

- Curl command:

   ```
   curl -X GET http://localhost:8080/test/protected -b "token=<token_here>"
   ```