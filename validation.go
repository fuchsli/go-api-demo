/*
	validation.go
		Provides the validation functions for the API
*/

package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// Validate data provided when creating a member
func validateMemberData(w http.ResponseWriter, m Member) bool {
	// If we're going to throw an error, we need the proper Content-Type
	w.Header().Set("Content-Type", "text/html")

	// Did the user provide a first name?
	if m.FirstName == "" {
		fmt.Fprintf(w, "The member must have a first name")
		return false
	}

	// Did the user provide a last name?
	if m.LastName == "" {
		fmt.Fprintf(w, "The member must have a last name")
		return false
	}

	// Is the provided JobType valid?
	lcJobType := strings.ToLower(m.JobType)
	if lcJobType != "contractor" && lcJobType != "employee" {
		fmt.Fprintf(w, "The job type provided is not valid. Please provide either 'Employee' or 'Contractor'")
		return false
	}

	// Did the user provide both a role and a duration for a member?
	if m.Role != "" && m.Duration != "" {
		fmt.Fprintf(w, "A member cannot have both a duration and a role")
		return false
	}

	// Did the user provide the right values for a contractor?
	if lcJobType == "contractor" && m.Role != "" {
		fmt.Fprintf(w, "A contractor cannot have a role")
		return false
	}

	// Did the user provide a duration for a contractor?
	if lcJobType == "contractor" && m.Duration == "" {
		fmt.Fprintf(w, "A contractor must have a duration")
		return false
	}

	// Did the user provide a duration for an employee?
	if lcJobType == "employee" && m.Duration != "" {
		fmt.Fprintf(w, "An employee cannot have a duration")
		return false
	}

	// Did the user provide a role for an employee?
	if lcJobType == "employee" && m.Role == "" {
		fmt.Fprintf(w, "An employee must have a role")
		return false
	}

	// Data is valid. Return Content-Type to application/json
	w.Header().Set("Content-Type", "application/json")
	return true
}

func validateUpdate(w http.ResponseWriter, filter bson.D, member Member, outcome string) (string, bool) {
	lcJobType := strings.ToLower(member.JobType)

	// Did the user update the first name?
	if member.FirstName != "" {
		update := bson.D{
			{"$set", bson.D{{"firstname", member.FirstName}}},
		}
		_, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			printErrorMessage(w, err)
			return "", false
		}
		outcome += "Updated first name successfully. "
	}

	// Did the user update the last name?
	if member.LastName != "" {
		update := bson.D{
			{"$set", bson.D{{"lastname", member.LastName}}},
		}
		_, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			printErrorMessage(w, err)
			return "", false
		}
		outcome += "Updated last name successfully. "
	}

	// Did the user update the job type?
	// If so, we need to make sure duration and role are handled accordingly
	if member.JobType != "" {
		if lcJobType != "contractor" && lcJobType != "employee" {
			fmt.Fprintf(w, "The job type provided is not valid. Please provide either 'Employee' or 'Contractor'.")
			return "", false
		}

		// If the job type is contractor, a duration must also be specified
		if lcJobType == "contractor" {
			if member.Duration == "" {
				fmt.Fprintf(w, "The contractor job type must have a specified duration.")
				return "", false
			}
			update := bson.D{
				{"$set", bson.D{
					{"jobtype", member.JobType},
					{"duration", member.Duration},
					{"role", ""},
				}},
			}
			_, err := collection.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				printErrorMessage(w, err)
				return "", false
			}
			outcome += "Updated member job type to contractor. "
		}

		// If the job type is employee, a role must also be specified
		if lcJobType == "employee" {
			if member.Role == "" {
				fmt.Fprintf(w, "The employee job type must have a specified role.")
				return "", false
			}
			update := bson.D{
				{"$set", bson.D{
					{"jobtype", member.JobType},
					{"role", member.Role},
					{"duration", ""},
				}},
			}
			_, err := collection.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				printErrorMessage(w, err)
				return "", false
			}
			outcome += "Updated member job type to employee. "
		}
	}

	// Did the user update the role?
	if member.Role != "" {
		if lcJobType == "" {
			fmt.Fprintf(w, "To set a role, please also include a job type of employee.")
			return "", false
		}

		if lcJobType == "contractor" {
			fmt.Fprintf(w, "A contractor cannot have a role.")
			return "", false
		}
		update := bson.D{
			{"$set", bson.D{
				{"jobtype", member.JobType},
				{"role", member.Role},
			}},
		}
		_, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			printErrorMessage(w, err)
			return "", false
		}
		outcome += "Updated role successfully. "
	}

	// Did the user update the duration?
	if member.Duration != "" {
		if lcJobType == "" {
			fmt.Fprintf(w, "To set a duration, please also include a job type of contractor.")
			return "", false
		}

		if lcJobType == "employee" {
			fmt.Fprintf(w, "An employee cannot have a duration.")
			return "", false
		}
		update := bson.D{
			{"$set", bson.D{
				{"jobtype", member.JobType},
				{"duration", member.Duration},
				{"role", ""},
			}},
		}
		_, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			printErrorMessage(w, err)
			return "", false
		}
		outcome += "Updated duration successfully. "
	}

	// Did the user update or remove the tags?
	if member.Tags != nil {
		update := bson.D{
			{"$set", bson.D{
				{"tags", member.Tags},
			}},
		}
		_, err := collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			printErrorMessage(w, err)
			return "", false
		}
		outcome += "Updated tags successfully. "
	}

	return outcome, true
}

// Ensure the  ID doesn't match an existing ID
func verifyUniqueID(clid string, outcome string) (string, string) {
	newID := clid
	filter := bson.D{{"clid", clid}}
	var testMember Member

	// Test whether or not the given ID matches a member
	err := collection.FindOne(context.TODO(), filter).Decode(&testMember)
	if err == nil {
		rand.Seed(time.Now().UTC().UnixNano())
		newID = strconv.Itoa(rand.Intn(99999999))
		outcome = "The provided ID was not unique, so a unique one with number " + newID + " was created. "
		verifyUniqueID(newID, outcome)
	}

	return newID, outcome
}
