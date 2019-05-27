## GO API Demo

This document serves as a tutorial for the Go API Demo.

### Preface

To emulate the feel of an actual API, this program has a running version hosted at 155.138.234.238 on port 8081.

### Testing the Application

This application has been tested with Postman. Feel free to use any application you like to call the API. But, for me, this is just a shout out to Postman for having such an awesome tool.

### Overview

In short, this API serves as a way to perform CRUD actions on a Mongo database containing member information. The database is called "go-api" and the collection is called "members". 

Each member is composed of the following fields: 
- firstname
- lastname
- jobtype
    - The job type can be either 'Employee' or 'Contractor'
    - Any other inputs will be rejected
- duration
    - If the job type is 'contractor', a duration must be provided 
- role
    - If the job type is 'employee', a role must be provided
- tags
    - An array of strings that act as additional information for the members
    - In the future, these will likely behave as filters for the data
    

The following are the ways to access and manipulate the data. The ID is randomly and uniquely generated. A desired ID can be provided as long as it is not already in use. If the ID is already in use, a new random and unique value will be generated:

- GET     /api/members
- GET     /api/members/{id}
- POST    /api/members
- PATCH   /api/members/{id}
- DELETE  /api/members/{id}
- DELETE  /api/members

#### GET /api/members

Sending a GET request to /api/members retrieves all of the documents found in the collection. They will be returned as JSON data. 

If the collection does not contain any documents, a message reading "The collection currently has no members" as plain text is returned.

#### GET /api/members/{id}

Sending a GET request to api/members/{id}, where {id} is a provided ID, returns the data for that member in JSON form. 

If no member can be found for the specified ID, an error message as plain text is returned.

#### POST /api/members

Sending a POST request to /api/members will successfully create a new member if the raw JSON data has been passed correctly. 

The application will alert you to any errors that might exist in your request. Some examples of these include:

- A job type other than 'Employee' or 'Contractor'
- No first name
- No last name

Also note that any additional fields will not be stored in the database. For example, if you try to create a member with "rich": "very", the document will save without that information.

#### PATCH /api/members/{id}

Sending a PATCH request to /api/members/{id} will update the provided information for the given ID. If the ID does not match up with an existing member, a message alerting the user that no matching member was found. Some use cases for the update function include: 

- Change a member's last name.
- Changing a member from a contractor to an employee (also requires including a role)
- Changing a member from an employee to a contractor (also requires including a duration)
- Changing the tags associated with a member

If the PATCH request is successful, a message saying that the update was successful is returned.

#### DELETE /api/members/{id}

Sending a DELETE request to /api/members/{id} will delete the document for the given ID. A message saying that the member has been successfully deleted is returned. 

If the ID does not match an existing ID in the database, a message saying that no member for the provided ID could be found is returned.

#### DELETE /api/members

Sending a DELETE request to /api/members will delete all documents inside of the collection. The result is an empty collection.

### Technical Tutorial

This section governs the more technical use of the program. 

#### Dependencies

The following third party libraries are used in this program:

- github.com/gorilla/mux
- go.mongodb.org/mongo-driver/mongo
- go.mongodb.org/mongo-driver/bson
- go.mongodb.org/mongo-driver/mongo/options
- go.mongodb.org/mongo-driver/mongo/readpref

To run the tests for the application, an additional testing dependency is required:

- github.com/stretchr/testify/assert

#### Program Structure

To reduce the amount of code in a single file and to organize the code in an easier way, the main package has been broken into component files:

- api.go
- crudFuncs.go
- errorFuncs.go
- validation.go
- api_test.go

##### api.go

api.go is the main file for the program. It includes:

- Global variable declarations
- Member struct declaration
- An init function, which connects to the MongoDB database and collection
- The main function, which creates the router, route handlers, and endpoints 

##### crudFuncs.go

crudFuncs.go handles all of the CRUD operations. It includes:

- getMembers, a function to display all members in the collection
- getMember, a function to display a single member with a matching ID
- createMember, a function to add new members to the collection
- updateMember, a function to change information about a member with a matching ID
- deleteMember, a function to delete a member's document with a matching ID
- deleteMembers, a function to delete all members in the collection

##### errorFuncs.go 

errorFuncs.go handles the errors for the program. It includes:

- printErrorMessage, a function to read a non-nil error to the responseWriter. It is intended to serve as a message to the user, so it will not terminate the program. It handles error messages like "mongo: no documents returned."
- handleError, a function that handles more critical errors. Unlike printErrorMessage, these errors are critical. They cause the application log the error to the terminal and close the program. 

##### validation.go

validation.go ensures that the data provided by a user is actually usable information. It includes the following functions:

- validateMemberData, a function that checks whether a member that's being created matches up with expected input. Any errors will be returned to the browser as text alerting the user as to what went wrong. The program will continue to run, and the user can change input data and try again.
- validateUpdate, a function that checks the information that a user is trying to update. If the data successfully updates, a message saying that the member was successfully updated is displayed. If unsuccessful, a specific reason for why the update was unsuccessful is displayed. The program continues to run and the user can change input data and try again.
- verifyUniqueID, a function that ensures the provided ID is actually unique. If it's not, it will call itself recursively until a unique ID is found.

#### Running the Application

To run the application, enter the following into a terminal on a system that has Go installed:

go run api.go crudFuncs.go errorFuncs.go validation.go

Or you can build the executable with 'go build' and run the executable with './api'

#### Testing

The application comes with a pre-built test. It is not exhaustive. However, it does handle a good number of cases. The test file is api_test.go.

 To run the test, enter the following into a terminal on a system that has Go installed:

go test