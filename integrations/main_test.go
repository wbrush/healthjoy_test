//  +build integration

package integrations

import (
	"bitbucket.org/optiisolutions/go-template-service/setup"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Printf("\n\nRunning TestMain(): Setting up test environment\n")

	//  set up integration test specific environment vars here
	os.Setenv("DB_DATABASE", "template_service")
	os.Setenv("DB_USER", "app")
	os.Setenv("DB_PASSWORD", "app")
	os.Setenv("DB_PORT", "5521")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_CONNS", "4")
	os.Setenv("DB_MIGRATION_PATH", "../dao/postgres")

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "") //  disable pubsub for now
	//os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "../../../../gcp_keys/optii-v3-pre-prod-fbe887520f60.json")
	os.Setenv("GCP_PROJECT_ID", "integration_test")
	//os.Setenv("PUBSUB_EMULATOR_HOST", "localhost:8921")

	//  set up regular environment vars here
	filename := "../files/.env"
	if _, err := os.Stat(filename); err == nil {
		_ = godotenv.Load(filename)
	}
	//fmt.Println("env.GCP_PROJECT_ID:", os.Getenv("GCP_PROJECT_ID"))
	//fmt.Println("env.PUBSUB_EMULATOR_HOST:", os.Getenv("PUBSUB_EMULATOR_HOST"))

	//  initialize service
	setup.SetupAndRun(false, "commit", "builtAt", "../docs/")

	//  perform various initializations
	InitializeDBs()

	fmt.Printf("\n\nPerforming Integration Tests\n")
	os.Exit(m.Run())
}
