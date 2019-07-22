/*
	Members API Assignment
	@author Eric Leohner

		Overview
			This program connects to a server and a Mongo database.

			It allows for CRUD actions on the database through the following routes and methods:
				/api/members         GET    - returns all members in the database
				/api/members/{id}  GET    - returns a specific member in the database with the provided ID
				/api/members         POST   - adds a new member to the database
				/api/members/{id}  PATCH  - updates information for a member with the provided clid
				/api/members/{id}  DELETE - deletes information for a member with the provided clid

			A working demonstration of this API is hosted at fuchsli.com on port 8081

			The README file serves as the tutorial for this assignment

*/

package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Member Struct
type Member struct {
	ID        string   `json:"clid" bson:"clid"`
	FirstName string   `json:"firstname" bson:"firstname"`
	LastName  string   `json:"lastname" bson:"lastname"`
	JobType   string   `json:"jobtype" bson:"jobtype"`
	Role      string   `json:"role,omitempty" bson:"role,omitempty"`
	Duration  string   `json:"duration,omitempty" bson:"duration,omitempty"`
	Tags      []string `json:"tags" bson:"tags"`
}

// Global variables
var (
	members    []Member
	collection *mongo.Collection
)

func redirectTLS(w http.ResponseWriter, r *http.Request) {
	domain := "fuchsli.com"
	http.Redirect(w, r, "https://"+domain+r.RequestURI, http.StatusMovedPermanently)
}

func redirectWWW(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Request host has www. prefix. Redirect to host with www. trimmed.
		host := strings.TrimPrefix(r.Host, "www.")

		if host != r.Host {
			// http.Redirect(w, r, u.String(), http.StatusFound)
			http.Redirect(w, r, "https://"+host+r.RequestURI, http.StatusFound)
			return
		}

		h.ServeHTTP(w, r)

	})
}

// Connect to MongoDB
func init() {
	// Create a MongoDB Connection on Port 27017
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	handleError(err)

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	handleError(err)

	fmt.Println("Connected to MongoDB")
	collection = client.Database("go-api").Collection("members")
}

func main() {

	certfile := "/etc/letsencrypt/live/fuchsli.com-0003/fullchain.pem"
	privkey := "/etc/letsencrypt/live/fuchsli.com-0003/privkey.pem"

	// Initialize router
	r := mux.NewRouter()

	// Route Handlers / Endpoints
	r.HandleFunc("/api/members", getMembers).Methods("GET")
	r.HandleFunc("/api/members/{id}", getMember).Methods("GET")
	r.HandleFunc("/api/members", createMember).Methods("POST")
	r.HandleFunc("/api/members/{id}", updateMember).Methods("PATCH")
	r.HandleFunc("/api/members/{id}", deleteMember).Methods("DELETE")
	r.HandleFunc("/api/members", deleteMembers).Methods("DELETE")

	go func() {
		if err := http.ListenAndServe(":8082", http.HandlerFunc(redirectTLS)); err != nil {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	config := &tls.Config{MinVersion: tls.VersionTLS10}
	server := &http.Server{
		Addr:      ":8081",
		Handler:   redirectWWW(r),
		TLSConfig: config,
	}
	log.Fatal(server.ListenAndServeTLS(certfile, privkey))
}
