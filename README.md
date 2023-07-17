# E-Voting RDBMS API

This API provides endpoints for user registration, login, admin actions, and voting in an electronic voting system.

## Endpoints

### Register a User

Register a new user in the system.

- Endpoint: `POST /register`
- Body: JSON object containing user registration data
- Example Body:
    ```json
    {
        "full_name": "John Doe",
        "mothers_name": "Jane Doe",
        "email": "john.doe@example.com",
        "phone_number": "1234567890",
        "jmbg": "1234567890123",
        "password": "password123"
    }
    ```

### User Login

Authenticate a user and retrieve their login information.

- Endpoint: `POST /login`
- Body: JSON object containing user login credentials
- Example Body:
    ```json
    {
        "email": "john.doe@example.com",
        "password": "password123"
    }
    ```

### Admin Actions

The following endpoints require admin authentication.

#### Add Candidate

Add a new candidate to the database.

- Endpoint: `POST /admin/candidate`
- Body: JSON object containing candidate information
- Example Body:
    ```json
    {
        "full_name": "Jane Smith",
        "district": "District 1",
        "short_bio": "Experienced leader with a passion for change."
    }
    ```

#### Get All Data

Retrieve all user, candidate and voting data from the database.

- Endpoint: `GET /admin/alldata`

#### Delete Candidate

Delete a candidate from the database by their ID.

- Endpoint: `DELETE /admin/candidate/{id}`
- Example: `DELETE /admin/candidate/1`

### Vote

Vote for a candidate in the election.

- Endpoint: `POST /vote`
- Body: JSON object containing vote information
- Example Body:
    ```json
    {
        "user_id": 1,
        "candidate_id": 2,
        "num_votes": 1
    }
    ```

## Testing Endpoints

You can test the API endpoints using Postman. Here's how to test each endpoint:

1. Register a User:
   - Set the request method to POST.
   - Enter the URL as `http://localhost:8000/register`.
   - Set the request body as a JSON object containing the user registration data.
   - Click the "Send" button to register the user.

2. User Login:
   - Set the request method to POST.
   - Enter the URL as `http://localhost:8000/login`.
   - Set the request body as a JSON object containing the user login credentials.
   - Click the "Send" button to authenticate the user.

3. Admin Actions:
   - Before testing the admin actions, make sure you have the admin credentials and pass them in the request headers or as parameters as required by your implementation.
   - Add Candidate:
     - Set the request method to POST.
     - Enter the URL as `http://localhost:8000/admin/candidate`.
     - Set the request body as a JSON object containing the candidate information.
     - Click the "Send" button to add the candidate.
   - Get All Data:
     - Set the request method to GET.
     - Enter the URL as `http://localhost:8000/admin/alldata`.
     - Click the "Send" button to retrieve all user and candidate data.
   - Delete Candidate:
     - Set the request method to DELETE.
     - Enter the URL as `http://localhost:8000/admin/candidate/{id}` where `{id}` is the ID of the candidate you want to delete.
     - Click the "Send" button to delete the candidate.

4. Vote:
   - Set the request method to POST.
   - Enter the URL as `http://localhost:8000/vote`.
   - Set the request body as a JSON object containing the vote information.
   - Click the "Send" button to cast the vote.


