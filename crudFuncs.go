package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

// Get a list of all members
func getMembers(w http.ResponseWriter, r *http.Request) {
	var members []*Member

	// Passing nil as the filter matches all documents in the collection
	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		printErrorMessage(w, err)
		return
	}
	// Close the cursor
	defer cur.Close(context.TODO())

	// Iterating through the cursor allows us to find one document at a time
	for cur.Next(context.TODO()) {
		// Create a value into which a single document can be decoded

		var member Member
		err := cur.Decode(&member)
		if err != nil {
			printErrorMessage(w, err)
			return
		}
		members = append(members, &member)
	}

	if err := cur.Err(); err != nil {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "The following error occurred: %v", err)
	}

	if len(members) == 0 {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "The collection currently has no members.")
	} else {
		w.Header().Set("Content-Type", "application/json")
		// Encode as JSON to display in browser
		json.NewEncoder(w).Encode(members)
	}
}

// Get a member by ID
func getMember(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	// Create a variable into which the resulting member data can be encoded
	var resultMember Member
	filter := bson.D{{"clid", params["clid"]}}
	err := collection.FindOne(context.TODO(), filter).Decode(&resultMember)
	if err != nil {
		printErrorMessage(w, err)
		return
	} else {
		// Encode resultMember as JSON
		json.NewEncoder(w).Encode(&resultMember)
	}
}

// Create a new member
func createMember(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	outcome := ""

	var member Member
	_ = json.NewDecoder(r.Body).Decode(&member)

	// The user can provide a custom ID as long as it's unique
	if member.ID == "" {
		rand.Seed(time.Now().UTC().UnixNano())
		member.ID = strconv.Itoa(rand.Intn(99999999))
	}

	// Ensure the ID is unique and validate provided information
	member.ID, outcome = verifyUniqueID(member.ID, outcome)
	isValidData := validateMemberData(w, member)

	// If the data is valid, insert it into the database
	if isValidData {
		_, err := collection.InsertOne(context.TODO(), member)
		if err != nil {
			printErrorMessage(w, err)
			return
		}

		outcome += "Created a new member"
		// members = append(members, member)
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, outcome)
	}

}

// Update member information
func updateMember(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	params := mux.Vars(r)
	var member Member
	var testMember Member

	filter := bson.D{{"clid", params["clid"]}}

	// Test whether or not the given ID matches a member
	err := collection.FindOne(context.TODO(), filter).Decode(&testMember)
	if err != nil {
		fmt.Fprintf(w, "No member for the provided ID could be found")
		return
	}

	_ = json.NewDecoder(r.Body).Decode(&member)
	var outcome string

	// End the update function without updating if a validation error
	outcome, ok := validateUpdate(w, filter, member, outcome)
	if !ok {
		return
	}

	if outcome == "" {
		outcome = "No data was changed"
	}
	fmt.Fprintf(w, outcome)
}

// Deletes a member
func deleteMember(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	params := mux.Vars(r)

	filter := bson.D{{"clid", params["clid"]}}

	var testMember Member

	// Test whether or not the given ID matches a member
	err := collection.FindOne(context.TODO(), filter).Decode(&testMember)
	if err != nil {
		fmt.Fprintf(w, "No member for the provided ID could be found")
		return
	}

	// Finds the matching ID and deletes the document
	_, err = collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		printErrorMessage(w, err)
		return
	}
	fmt.Fprintf(w, "Member successfully deleted")
}

// Deletes all members
func deleteMembers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	_, err := collection.DeleteMany(context.Background(), bson.D{})
	if err != nil {
		printErrorMessage(w, err)
		return
	}
	fmt.Fprintf(w, "Successfully deleted all members")
}
