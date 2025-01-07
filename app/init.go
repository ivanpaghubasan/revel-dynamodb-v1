package app

import (
	"fmt"
	"log"
	"revel-dynamodb-v1/app/repositories"
	"revel-dynamodb-v1/app/services"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	_ "github.com/revel/modules"
	"github.com/revel/revel"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string

	Service *services.MovieService
)

func initDynamoDB() *dynamodb.DynamoDB {
	// Set aws session
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("us-west-2"),
		Endpoint: aws.String("http://localhost:8000"),
	})

	if err != nil {
		log.Fatalf("Error creating a new aws session: %v", err)
	}

	// Initialize dynamodb client
	svc := dynamodb.New(sess)
	fmt.Println("DynamoDB has been initialized...")

	tables, err := svc.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		log.Fatal(err)
	}

	isTableExisting := false
	tableName := "Movies"

	fmt.Println("Checking if the table is already existing...")

	for _, table := range tables.TableNames {
		if *table == tableName {
			isTableExisting = true
			break
		}
	}

	if !isTableExisting {
		_, err := svc.CreateTable(&dynamodb.CreateTableInput{
			TableName: aws.String(tableName),
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("Year"),
					AttributeType: aws.String("N"),
				},
				{
					AttributeName: aws.String("Title"),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("Year"),
					KeyType:       aws.String("HASH"),
				},
				{
					AttributeName: aws.String("Title"),
					KeyType:       aws.String("RANGE"),
				},
			},
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(100),
				WriteCapacityUnits: aws.Int64(100),
			},
		})

		if err != nil {
			log.Fatalf("Got error calling create table: %s", err)
		}

		fmt.Printf("Table %s has been created. ", tableName)
	}

	return svc
}

func initRepositories() repositories.Repository {
	svc := initDynamoDB()
	tableName := "Movies"
	repo := repositories.New(svc, tableName)
	fmt.Println("Repositories has been initialized...")
	return repo
}

func initServices() {
	repo := initRepositories()
	fmt.Println("Services has been initialized...")
	Service = services.New(repo)
}

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.BeforeAfterFilter,       // Call the before and after filter functions
		revel.ActionInvoker,           // Invoke the action.
	}

	// Register startup functions with OnAppStart
	// revel.DevMode and revel.RunMode only work inside of OnAppStart. See Example Startup Script
	// ( order dependent )
	// revel.OnAppStart(ExampleStartupScript)
	// revel.OnAppStart(InitDB)
	// revel.OnAppStart(FillCache)
	revel.OnAppStart(initServices)
}

// HeaderFilter adds common security headers
// There is a full implementation of a CSRF filter in
// https://github.com/revel/modules/tree/master/csrf
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")
	c.Response.Out.Header().Add("Referrer-Policy", "strict-origin-when-cross-origin")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

//func ExampleStartupScript() {
//	// revel.DevMod and revel.RunMode work here
//	// Use this script to check for dev mode and set dev/prod startup scripts here!
//	if revel.DevMode == true {
//		// Dev mode
//	}
//}
