/*
	api_test.go

		A demo implementation of a CRUD API

		! Important:
			Never test on a production database. Testing will empty the mongo collection completely. Bad things will happen.
*/

package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

// Create the router we will use for the tests
func Router() *mux.Router {
	// Initialize router
	r := mux.NewRouter()

	// Route Handlers / Endpoints
	r.HandleFunc("/api/members", getMembers).Methods("GET")
	r.HandleFunc("/api/members/{clid}", getMember).Methods("GET")
	r.HandleFunc("/api/members", createMember).Methods("POST")
	r.HandleFunc("/api/members/{clid}", updateMember).Methods("PATCH")
	r.HandleFunc("/api/members/{clid}", deleteMember).Methods("DELETE")
	r.HandleFunc("/api/members", deleteMembers).Methods("DELETE")
	return r
}

// Try to empty the DB Collection
func TestEmptyDB(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing emptying collection")

	req, _ := http.NewRequest("DELETE", "/api/members", nil)
	recorder := httptest.NewRecorder()
	Router().ServeHTTP(recorder, req)
	expected := "Successfully deleted all members"
	received := recorder.Body.String()

	ok := assert.Equal(t, expected, received, "They should be the same")
	if ok {
		fmt.Println("Successfully deleted all members from the collection")
	}
}

// Try to retrieve data from the DB
func TestGetMembersEmpty(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing getting all documents in an empty collection")

	req, _ := http.NewRequest("GET", "/api/members", nil)
	recorder := httptest.NewRecorder()
	Router().ServeHTTP(recorder, req)
	expected := "The collection currently has no members."
	received := recorder.Body.String()

	ok := assert.Equal(t, expected, received, "They should be the same")
	if ok {
		fmt.Println("Successfully read empty collection")
	}
}

// Try adding a member with a ID
func TestAddMemberID(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing creating a member with a preselected ID")

	testData := []byte(`{"clid": "1", "firstname": "Julius", "lastname": "Caesar","jobtype": "Employee", "role": "Imperator", "tags": ["He wasn't actually an emperor"]}`)

	req, _ := http.NewRequest("POST", "/api/members", bytes.NewBuffer(testData))
	recorder := httptest.NewRecorder()
	Router().ServeHTTP(recorder, req)

	expected := "Created a new member"
	received := recorder.Body.String()

	ok := assert.Equal(t, expected, received, "They should be the same")
	if ok {
		fmt.Println("Successfully created an employee with a provided ID")
	}
}

// Try getting a member by a known ID
func TestGetMemberByID(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing getting a member with by ID")

	req, _ := http.NewRequest("GET", "/api/members/1", nil)
	recorder := httptest.NewRecorder()
	Router().ServeHTTP(recorder, req)

	expected := `{"clid":"1","firstname":"Julius","lastname":"Caesar","jobtype":"Employee","role":"Imperator","tags":["He wasn't actually an emperor"]}`
	received := recorder.Body.String()

	/*
		client := &http.Client{}
		req, _ := http.NewRequest("GET", "/api/members/1", nil)
		resp, _ := client.Do(req)
		defer resp.Body.Close()

		respBody, _ := ioutil.ReadAll(resp.Body)

		expected := `{"clid":"1","firstname":"Julius","lastname":"Caesar","jobtype":"Employee","role":"Imperator","tags":["He wasn't actually an emperor"]}`
		received := string(respBody)
	*/

	// For some reason, the response Body returns with a new line tag, so to make sure the values are the same we have to clear it.
	received = strings.Trim(received, "\n")

	ok := assert.Equal(t, expected, received, "They should be the same")
	if ok {
		fmt.Println("Successfully read member with provided ID")
	}
}

// Try getting all members from the DB again
func TestGetAllMembers(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing getting all members from a non-empty collection")

	req, _ := http.NewRequest("GET", "/api/members", nil)
	recorder := httptest.NewRecorder()
	Router().ServeHTTP(recorder, req)

	expected := `[{"clid":"1","firstname":"Julius","lastname":"Caesar","jobtype":"Employee","role":"Imperator","tags":["He wasn't actually an emperor"]}]`
	received := recorder.Body.String()
	received = strings.Trim(received, "\n")

	ok := assert.Equal(t, expected, received, "They should be the same")
	if ok {
		fmt.Println("Successfully read all members in the collection")
	}
}

// Try updating first name
func TestUpdateFirstName(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing updating first name")

	testData := []byte(`{"firstname":"NotJulius"}`)
	req, _ := http.NewRequest("PATCH", "/api/members/1", bytes.NewBuffer(testData))
	recorder := httptest.NewRecorder()
	Router().ServeHTTP(recorder, req)

	expected := "Updated first name successfully. "
	received := recorder.Body.String()

	ok := assert.Equal(t, expected, received, "They should be the same")
	if ok {
		fmt.Println("Successfully updated member first name")
	}
}

// Try updating first and last name
func TestUpdateFirstAndLastName(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing updating first and last name")

	testData := []byte(`{"firstname":"NotJulius","lastname":"NotCaesar"}`)
	req, _ := http.NewRequest("PATCH", "/api/members/1", bytes.NewBuffer(testData))
	recorder := httptest.NewRecorder()
	Router().ServeHTTP(recorder, req)

	expected := "Updated first name successfully. Updated last name successfully. "
	received := recorder.Body.String()

	ok := assert.Equal(t, expected, received, "They should be the same")
	if ok {
		fmt.Println("Successfully updated member first and last name")
	}
}

// Try updating last name
func TestUpdateLastName(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing updating last name")

	testData := []byte(`{"lastname":"CaesarAgain"}`)
	req, _ := http.NewRequest("PATCH", "/api/members/1", bytes.NewBuffer(testData))
	recorder := httptest.NewRecorder()
	Router().ServeHTTP(recorder, req)

	expected := "Updated last name successfully. "
	received := recorder.Body.String()

	ok := assert.Equal(t, expected, received, "They should be the same")
	if ok {
		fmt.Println("Successfully updated member last name")
	}
}

// Try changing a job type to 'contractor' without providing a duration
func TestUpdateJobTypeContractorNoDuration(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing updating job type to contractor without duration")

	testData := []byte(`{"jobtype":"Contractor"}`)
	req, _ := http.NewRequest("PATCH", "/api/members/1", bytes.NewBuffer(testData))
	recorder := httptest.NewRecorder()
	Router().ServeHTTP(recorder, req)

	expected := "The contractor job type must have a specified duration."
	received := recorder.Body.String()

	ok := assert.Equal(t, expected, received, "They should be the same")
	if ok {
		fmt.Println("Failed to update job type contractor without a duration")
	}
}

// Try updating a member to 'contractor' while also including a role
func TestUpdateJobTypeContractorDurationRole(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing updating a contractor with a role")

	testData := []byte(`{"jobtype":"Contractor","duration":"6 weeks","role":"Executive"}`)
	req, _ := http.NewRequest("PATCH", "/api/members/1", bytes.NewBuffer(testData))
	recorder := httptest.NewRecorder()
	Router().ServeHTTP(recorder, req)

	expected := "A contractor cannot have a role."
	received := recorder.Body.String()

	ok := assert.Equal(t, expected, received, "They should be the same")
	if ok {
		fmt.Println("Failed to update job type contractor with a role")
	}
}

// Try adding a contractor with a duration
func TestUpdateJobTypeContractorDuration(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing updating a contractor with a duration")

	testData := []byte(`{"jobtype":"Contractor","duration":"6 months"}`)
	req, _ := http.NewRequest("PATCH", "/api/members/1", bytes.NewBuffer(testData))
	recorder := httptest.NewRecorder()
	Router().ServeHTTP(recorder, req)

	expected := "Updated member job type to contractor. Updated duration successfully. "
	received := recorder.Body.String()

	ok := assert.Equal(t, expected, received, "They should be the same")
	if ok {
		fmt.Println("Successfully updated contract duration")
	}
}

// Try updating a member to 'employee' without a role
func TestUpdateJobTypeEmployeeNoRole(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing updating job type to employee with no role")

	testData := []byte(`{"jobtype":"Employee"}`)
	req, _ := http.NewRequest("PATCH", "/api/members/1", bytes.NewBuffer(testData))
	recorder := httptest.NewRecorder()
	Router().ServeHTTP(recorder, req)

	expected := "The employee job type must have a specified role."
	received := recorder.Body.String()

	ok := assert.Equal(t, expected, received, "They should be the same")
	if ok {
		fmt.Println("Failed to update contractor to employee with no role")
	}
}

// Try updating a member to 'employee' with a role
func TestUpdateJobTypeEmployeeRole(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing updating job type to employee with a provided role")

	testData := []byte(`{"jobtype":"Employee","role":"Mastermind"}`)
	req, _ := http.NewRequest("PATCH", "/api/members/1", bytes.NewBuffer(testData))
	recorder := httptest.NewRecorder()
	Router().ServeHTTP(recorder, req)

	expected := "Updated member job type to employee. Updated role successfully. "
	received := recorder.Body.String()

	ok := assert.Equal(t, expected, received, "They should be the same")
	if ok {
		fmt.Println("Successfully updated contractor to employee with role")
	}
}

// Try updating only the role field
func TestUpdateEmployeeNoJobTypeOnlyRole(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing updating only role")

	testData := []byte(`{"role":"Dishwasher"}`)
	req, _ := http.NewRequest("PATCH", "/api/members/1", bytes.NewBuffer(testData))
	recorder := httptest.NewRecorder()
	Router().ServeHTTP(recorder, req)

	expected := "To set a role, please also include a job type of employee."
	received := recorder.Body.String()

	ok := assert.Equal(t, expected, received, "They should be the same")
	if ok {
		fmt.Println("Failed to update role without specified job type employee")
	}
}

// Try updating a member to 'employee' with a duration
func TestUpdateEmployeeDuration(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing updating an employee with a duration")

	testData := []byte(`{"duration":"2 years"}`)
	req, _ := http.NewRequest("PATCH", "/api/members/1", bytes.NewBuffer(testData))
	recorder := httptest.NewRecorder()
	Router().ServeHTTP(recorder, req)

	expected := "To set a duration, please also include a job type of contractor."
	received := recorder.Body.String()

	ok := assert.Equal(t, expected, received, "They should be the same")
	if ok {
		fmt.Println("Failed to update employee with a duration")
	}
}

// Try updating first name, last name, job type to contractor, and duration
func TestUpdateFnameLnameJobTypeDuration(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing updating first name, last name, job type contarctor, and duration")

	testData := []byte(`{"firstname":"Johnny","lastname":"Nolastname","jobtype":"contractor","duration":"2 years"}`)
	req, _ := http.NewRequest("PATCH", "/api/members/1", bytes.NewBuffer(testData))
	recorder := httptest.NewRecorder()
	Router().ServeHTTP(recorder, req)

	expected := "Updated first name successfully. Updated last name successfully. Updated member job type to contractor. Updated duration successfully. "
	received := recorder.Body.String()

	ok := assert.Equal(t, expected, received, "They should be the same")
	if ok {
		fmt.Println("Successfully updated first name, last name, job type to contractor, and duration")
	}
}

// Try updating tags
func TestUpdateTags(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing updating tags")

	testData := []byte(`{"tags": ["Emperor", "Really Cool Guy"]}`)
	req, _ := http.NewRequest("PATCH", "/api/members/1", bytes.NewBuffer(testData))
	recorder := httptest.NewRecorder()
	Router().ServeHTTP(recorder, req)

	expected := "Updated tags successfully. "
	received := recorder.Body.String()

	ok := assert.Equal(t, expected, received, "They should be the same")
	if ok {
		fmt.Println("Successfully updated tags")
	}
}

// Try deleting a member by ID
func TestDeleteMemberByID(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing deleting a member by ID")

	req, _ := http.NewRequest("DELETE", "/api/members/1", nil)
	recorder := httptest.NewRecorder()
	Router().ServeHTTP(recorder, req)

	expected := "Member successfully deleted"
	received := recorder.Body.String()

	ok := assert.Equal(t, expected, received, "They should be the same")
	if ok {
		fmt.Println("Successfully deleted member with provided ID")
	}
}

// Try adding a new member
func TestAddMember(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing creating a new valid member")

	testData := []byte(`{"firstname": "pirate","lastname": "Booty","jobtype": "Contractor","duration": "4 minutes","tags": ["Has scurvy","Needs Vitamin C"]}`)
	req, _ := http.NewRequest("POST", "/api/members", bytes.NewBuffer(testData))
	recorder := httptest.NewRecorder()
	Router().ServeHTTP(recorder, req)

	expected := "Created a new member"
	received := recorder.Body.String()

	ok := assert.Equal(t, expected, received, "They should be the same")
	if ok {
		fmt.Println("Successfully created a new member")
	}
}

// Try adding a member without a first name
func TestAddMemberNoFirstName(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing creating a member with no first name")

	testData := []byte(`{"lastname": "Lincoln","jobtype": "Employee","role": "President","tags": ["Was President","Is Tall"]}`)

	req, _ := http.NewRequest("POST", "/api/members", bytes.NewBuffer(testData))
	recorder := httptest.NewRecorder()
	Router().ServeHTTP(recorder, req)

	expected := "The member must have a first name"
	received := recorder.Body.String()

	ok := assert.Equal(t, expected, received, "They should be the same")
	if ok {
		fmt.Println("Failed to create member with no first name")
	}
}

// Try adding a member with no last name
func TestAddMemberNoLastName(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing creating a member no last name")

	testData := []byte(`{"firstname": "Steven","jobtype": "Contractor","duration": "1 Year","tags": ["Is Contractor","Is not the boss"]}`)

	req, _ := http.NewRequest("POST", "/api/members", bytes.NewBuffer(testData))
	recorder := httptest.NewRecorder()
	createMember(recorder, req)

	expected := "The member must have a last name"
	received := recorder.Body.String()

	ok := assert.Equal(t, expected, received, "They should be the same")
	if ok {
		fmt.Println("Failed to create member with no last name")
	}
}

// Try adding a member with an invalid job type
func TestAddMemberWrongJobType(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing creating a member with incorrect job type")

	testData := []byte(`{"firstname": "Baberaham", "lastname": "Lincoln","jobtype": "President","duration": "Forever","tags": ["Is a babe","Is not Bugs Bunny when he puts on a dress and plays girl bunny"]}`)

	req, _ := http.NewRequest("POST", "/api/members", bytes.NewBuffer(testData))
	recorder := httptest.NewRecorder()
	createMember(recorder, req)

	expected := "The job type provided is not valid. Please provide either 'Employee' or 'Contractor'"
	received := recorder.Body.String()

	ok := assert.Equal(t, expected, received, "They should be the same")
	if ok {
		fmt.Println("Failed to create member with invalid job type")
	}
}

// Try adding a contractor with a role
func TestAddMemberContractorRole(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing creating a contractor with a role")

	testData := []byte(`{"firstname": "Baberaham", "lastname": "Lincoln","jobtype": "Contractor","role": "CEO","tags": ["Is a babe","Is not Bugs Bunny when he puts on a dress and plays girl bunny"]}`)

	req, _ := http.NewRequest("POST", "/api/members", bytes.NewBuffer(testData))
	recorder := httptest.NewRecorder()
	createMember(recorder, req)

	expected := "A contractor cannot have a role"
	received := recorder.Body.String()

	ok := assert.Equal(t, expected, received, "They should be the same")
	if ok {
		fmt.Println("Failed to create a contractor with a role")
	}
}

// Try adding a contractor with no duration
func TestAddContractorNoDuration(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing creating a contractor with no duration")

	testData := []byte(`{"firstname": "Baberaham", "lastname": "Lincoln","jobtype": "Contractor","tags": ["Is a babe","Is not Bugs Bunny when he puts on a dress and plays girl bunny"]}`)

	req, _ := http.NewRequest("POST", "/api/members", bytes.NewBuffer(testData))
	recorder := httptest.NewRecorder()
	createMember(recorder, req)

	expected := "A contractor must have a duration"
	received := recorder.Body.String()

	ok := assert.Equal(t, expected, received, "They should be the same")
	if ok {
		fmt.Println("Failed to create a contractor without a duration")
	}
}

// Try adding a member with both a role and a duration
func TestAddBothRoleDuration(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing creating a member with both a role and a duration")

	testData := []byte(`{"firstname": "Baberaham", "lastname": "Lincoln","jobtype": "Contractor","role": "CEO", "duration": "5 Hours"}`)

	req, _ := http.NewRequest("POST", "/api/members", bytes.NewBuffer(testData))
	recorder := httptest.NewRecorder()
	createMember(recorder, req)

	expected := "A member cannot have both a duration and a role"
	received := recorder.Body.String()

	ok := assert.Equal(t, expected, received, "They should be the same")
	if ok {
		fmt.Println("Failed to create a member with both a role and a duration")
	}
}

// Try adding an employee without a role
func TestAddEmployeeNoRole(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing creating a employee with no role")

	testData := []byte(`{"firstname": "Baberaham", "lastname": "Lincoln","jobtype": "Employee", "tags": ["Maybe a president"]}`)

	req, _ := http.NewRequest("POST", "/api/members", bytes.NewBuffer(testData))
	recorder := httptest.NewRecorder()
	createMember(recorder, req)

	expected := "An employee must have a role"
	received := recorder.Body.String()

	ok := assert.Equal(t, expected, received, "They should be the same")
	if ok {
		fmt.Println("Failed to create an employee without a role")
	}
}

// Try adding an employee with a duration
func TestAddEmployeeDuration(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing creating an employee with a duration")

	testData := []byte(`{"firstname": "Baberaham", "lastname": "Lincoln","jobtype": "Employee", "duration": "1 year", "tags": ["Maybe a president"]}`)

	req, _ := http.NewRequest("POST", "/api/members", bytes.NewBuffer(testData))
	recorder := httptest.NewRecorder()
	createMember(recorder, req)

	expected := "An employee cannot have a duration"
	received := recorder.Body.String()

	ok := assert.Equal(t, expected, received, "They should be the same")
	if ok {
		fmt.Println("Failed to create an employee with a duration")
	}
}

// Try adding a valid contractor
func TestAddValidContractor(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing creating a valid contractor")

	testData := []byte(`{"firstname": "Baberaham", "lastname": "Lincoln","jobtype": "Contractor", "duration": "1 year", "tags": ["Maybe a president"]}`)

	req, _ := http.NewRequest("POST", "/api/members", bytes.NewBuffer(testData))
	recorder := httptest.NewRecorder()
	createMember(recorder, req)

	expected := "Created a new member"
	received := recorder.Body.String()

	ok := assert.Equal(t, expected, received, "They should be the same")
	if ok {
		fmt.Println("Successfully created a contractor")
	}
}

// Try adding a valid employee
func TestAddValidEmployee(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing creating a valid employee")

	testData := []byte(`{"firstname": "Julius", "lastname": "Caesar","jobtype": "Employee", "role": "Imperator", "tags": ["He wasn't actually an emperor"]}`)

	req, _ := http.NewRequest("POST", "/api/members", bytes.NewBuffer(testData))
	recorder := httptest.NewRecorder()
	createMember(recorder, req)

	expected := "Created a new member"
	received := recorder.Body.String()

	ok := assert.Equal(t, expected, received, "They should be the same")
	if ok {
		fmt.Println("Successfully created an employee")
	}
}

// Try to empty the Collection again
func TestEmptyDBAgain(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing emptying collection")

	_, err := collection.DeleteMany(context.Background(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully deleted all members from the collection again")
}

// Test getting a wrong status code
func TestBad(t *testing.T) {
	fmt.Println("----------------")
	fmt.Println("Testing a bad request; one that 405s")

	testData := []byte(`{"firstname":"Bill"}`)
	req, _ := http.NewRequest("PUT", "/api/members", bytes.NewBuffer(testData))
	recorder := httptest.NewRecorder()
	Router().ServeHTTP(recorder, req)

	ok := assert.Equal(t, 405, recorder.Code, "They should be the same")
	if ok {
		fmt.Println("Successfully failed to send a bad request")
		fmt.Println("----------------")
	}
}
