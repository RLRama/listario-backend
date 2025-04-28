# API documentation

## Base URL

`https://localhost:8080/` or wherever it's hosted

## Endpoints

### User endpoints

#### Register

- POST - Curl command:

   ```bash
   curl -X POST http://localhost:8080/register \
     -H "Content-Type: application/json" \
     -d '{
           "username": "testuser",
           "password": "Password123!",
           "email": "testuser@example.com"
         }'
   ```

#### Login

- POST - Curl command:

   ```bash
   curl -X POST http://localhost:8080/user/login \
     -H "Content-Type: application/json" \
     -d '{
           "identifier": "testuser",
           "password": "Password123!"
         }'
   ```

#### Update username and/or email

- PUT - Curl command:

   ```bash
   curl -X PUT http://localhost:8080/user/protected/update   -H "Content-Type: application/json"   -b "token=<token_here>"   -d '{"username":"john","email":"testuser10@example.com"}'
   ```

- > `username` and `email` in the payload can be left blank if they are not to be changed.

#### PUT - Update password

- Curl command:

   ```bash
   curl -X PUT http://localhost:8080/user/protected/update-password   -H "Content-Type: application/json"   -b "token=<token_here>"   -d '{"current_password":"old_password123","new_password":"new_password123"}'
   ```

#### POST - Refresh token

- Curl command:

   ```bash
   curl -X POST http://localhost:8080/user/protected/refresh   -H "Content-Type: application/json"   -b "token=<token_here>"
   ```

#### POST - Logout

- Curl command:

   ```bash
   curl -X POST http://localhost:8080/user/protected/logout   -H "Content-Type: application/json"   -b "token=<token_here>"
   ```

#### GET - User details (me)

- Curl command:

   ```bash
   curl -X GET http://localhost:8080/user/protected/me   -H "Content-Type: application/json"   -b "token=<token_here>"
   ```

### Task endpoints

#### POST - Create a task

- Curl command:

```bash
curl -X POST http://localhost:8080/tasks \
-H "Content-Type: application/json" \
-b "token=<token_here>" \
-d '{"title":"Write report","description":"Annual report for Q1","category_id":1,"due_date":"2025-05-01T00:00:00Z","tag_ids":[1,2]}'
```

#### GET - List tasks

- Curl command

```bash
curl -X GET "http://localhost:8080/tasks?category_id=1&tag_ids=1,2" \
-b "token=<token_here>"
```

#### GET - Get a specific task

- Curl command

```bash
curl -X PUT http://localhost:8080/tasks/1 \
-H "Content-Type: application/json" \
-b cookies.txt \
-d '{"title":"Updated report","description":"Revised Q1 report","category_id":1,"completed":true,"due_date":"2025-05-02T00:00:00Z","tag_ids":[1]}'
```

#### DELETE - Delete a task

- Curl command

```bash
curl -X DELETE http://localhost:8080/tasks/1 \
-b cookies.txt
```

#### PUT - Update a task

- Curl command

```bash
curl -X GET http://localhost:8080/tasks/1 \
-b "token=<token_here>"
```

### Test endpoints

#### GET - Hello world!

- Curl command:

   ```bash
   curl -X GET http://localhost:8080/test/hello
   ```

#### GET - Test DB connection

- Curl command:

   ```bash
   curl -X GET http://localhost:8080/test/db-connection
   ```

#### GET - Protected test endpoint

- Curl command:

   ```bash
   curl -X GET http://localhost:8080/test/protected -b "token=<token_here>"
   ```