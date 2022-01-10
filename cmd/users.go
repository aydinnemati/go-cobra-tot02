package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var reader = bufio.NewReader(os.Stdin)
var MONGO_IP = os.Getenv("MONGO_IP")
var MONGO_PORT = os.Getenv("MONGO_PORT")
var MONGO_COL = os.Getenv("MONGO_COL")
var MONGO_DB = os.Getenv("MONGO_DB")

// usersCmd represents the users command
var usersCmd = &cobra.Command{
	Use:   "users",
	Short: "fo create or get users",
	Long:  `create or get users by flags from mongodb database`,
	Run: func(cmd *cobra.Command, args []string) {
		// connect to mongodb
		client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://"+MONGO_IP+":"+MONGO_PORT))
		if err != nil {
			panic(err)
		}
		if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
			panic(err)
		}
		usersCollection := client.Database(MONGO_DB).Collection(MONGO_COL)

		inputflag, _ := cmd.Flags().GetBool("input")
		argsflag, _ := cmd.Flags().GetBool("args")

		if inputflag {
			createUserInput(reader, usersCollection)
		} else if argsflag {
			firstname, _ := cmd.Flags().GetString("firstname")
			lastname, _ := cmd.Flags().GetString("lastname")
			if firstname == "" {
				fmt.Println("please use right arguments (firstname)")
			} else if lastname == "" {
				fmt.Println("please use right arguments (lastname)")
			} else {
				createUserArgs(firstname, lastname, usersCollection)
			}

		} else {
			getUsers(usersCollection)
		}

	},
}

func init() {
	rootCmd.AddCommand(usersCmd)

	// load env vars
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// flags
	usersCmd.Flags().BoolP("input", "i", false, "read user data from stdin")
	usersCmd.Flags().BoolP("args", "a", false, "read user data from args")

	// args
	usersCmd.PersistentFlags().String("firstname", "", "user's firstname")
	usersCmd.PersistentFlags().String("lastname", "", "user's lastname")

}

func getUsers(dbcol *mongo.Collection) {
	// retrieve all the documents in a collection
	cursor, err := dbcol.Find(context.TODO(), bson.D{})
	// check for errors in the finding
	if err != nil {
		panic(err)
	}

	// convert the cursor result to bson
	var results []bson.M
	// check for errors in the conversion
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	// display the documents retrieved
	fmt.Println("displaying all users:")
	for _, result := range results {
		fmt.Println(result)
	}
}

func createUserInput(r *bufio.Reader, dbcol *mongo.Collection) {
	fmt.Println("pleas enter user's first name:")
	fname, _ := reader.ReadString('\n')
	fmt.Println("pleas enter user's last name:")
	lname, _ := reader.ReadString('\n')
	fmt.Println(lname)
	fmt.Println(fname)
	addinguser := bson.D{{"firstname", fname}, {"lastname", lname}}
	// insert the bson object using InsertOne()
	result, err := dbcol.InsertOne(context.TODO(), addinguser)
	// check for errors in the insertion
	if err != nil {
		panic(err)
	}
	// display the id of the newly inserted object
	fmt.Println(result.InsertedID)
}

func createUserArgs(fname string, lname string, dbcol *mongo.Collection) {
	fmt.Println(fname)
	fmt.Println(lname)
	addinguser := bson.D{{"firstname", fname}, {"lastname", lname}}
	// insert the bson object using InsertOne()
	result, err := dbcol.InsertOne(context.TODO(), addinguser)
	// check for errors in the insertion
	if err != nil {
		panic(err)
	}
	// display the id of the newly inserted object
	fmt.Println(result.InsertedID)
}
