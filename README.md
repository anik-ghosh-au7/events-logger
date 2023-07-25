# Events-Logger (Hasura)

This application logs events triggered by Hasura, a popular GraphQL engine. These events, triggered by insert, update, or delete operations in the connected database, are sent to a webhook endpoint provided by this application. The application logs these events in an in-memory database and provides HTTP endpoints to retrieve these events.

## File Structure and Application Logic

The main codebase is split into two main files, `main.go` and `db.go`, which reside in the `main` and `database` packages respectively.

### `main.go`

This is the entry point of the application. It sets up HTTP routes and handles incoming requests. There are three main endpoints:

1. `/webhook`: This endpoint accepts `POST` requests. It expects a JSON payload in the request body, which it decodes into a `Payload` object (defined in the `database` package). The payload is then stored in the in-memory database with its `ID` as the key.

2. `/events`: This endpoint accepts `GET` requests and returns a list of all event IDs currently stored in the database.

3. `/events/{id}`: This endpoint accepts `GET` requests and returns the event corresponding to the provided `id`. If no such event is found in the database, it returns an error.

### `database/db.go`

This file contains code related to the in-memory database used for storing events. It defines the `Payload` structure and the `InMemoryDB` type.

The `Payload` structure represents an event. It includes various fields like `Event`, `CreatedAt`, `ID`, `Trigger`, and `Table`, which capture different aspects of the event. This structure is designed to parse the request payload as per the Hasura specification.

## Integration with Hasura

To use this application with Hasura, follow these steps:

1. Run the application, ensuring it is accessible over the internet. If you're running it locally, you might need to use a tool like ngrok to expose your local server to the internet.

2. Log in to your Hasura console. This is usually found at the URL where you're hosting Hasura, followed by `/console`.

3. Navigate to the "Events" tab, then "Create Trigger".

4. Give the trigger a name, select the table you want to set the trigger on, and choose the operations that should cause the trigger to fire (insert, update, delete, or any combination of the three).

5. In the "Webhook" field, enter the URL where your events-logger application is running, followed by `/webhook`. For example, if your application is running at `http://myapp.com`, you would enter `http://myapp.com/webhook`.

6. Click "Create".

Now, whenever the chosen operation happens on your table, Hasura will send an event to your events-logger application. The event will be logged in the in-memory database, and you can retrieve it using the `/events` and `/events/{id}` endpoints.

## Testing the Application

To test the application, you can use the following `curl` commands:

1. To retrieve all events:

```sh
curl -X GET http://localhost:8080/events
```

2. To retrieve a specific event by its id:

```sh
curl -X GET http://localhost:8080/events/{id}
```

> Note: Replace {id} with the actual ID of the event.

## Running the Application

To run the application, you need to have Go installed on your machine. You can then run the application using the following command:

```sh
go run main.go
```

The application will start a server and listen on port 8080.

> Note: If you're testing this locally and want to connect it to Hasura, you'll need to expose your local server to the internet. Tools like [ngrok](https://ngrok.com/) can help with this.

## Notes

Remember, this is a sample application and the in-memory database used here does not persist if the application restarts. In a real-world scenario, you would likely want to replace the in-memory database with a more persistent database solution.
