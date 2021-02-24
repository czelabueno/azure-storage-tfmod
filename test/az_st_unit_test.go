package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resourcehealth/mgmt/2017-07-01/resourcehealth"
	ct "github.com/daviddengcn/go-colortext"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func Test_AzStorageExample(t *testing.T) {

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../",
	})

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)
	ValidateModule(t, terraformOptions)
}

func ValidateModule(t *testing.T, terraformOptions *terraform.Options) (result bool, err error) {
	subscriptionID := terraform.OutputRequired(t, terraformOptions, "subscriptionId")
	resourceID := terraform.OutputRequired(t, terraformOptions, "resourceId")
	resourceName := terraform.OutputRequired(t, terraformOptions, "resourceName")

	time.Sleep(10 * time.Second)
	// Create availabilityStatusesClient
	availabilityStatusesClient := resourcehealth.NewAvailabilityStatusesClient(subscriptionID)

	// Get azure token access from Environment(Client credential or client certificate) or File json or Az CLI
	authorizer, err := azure.NewAuthorizer()
	if err == nil {
		availabilityStatusesClient.Authorizer = *authorizer
		availabilityStatusesClient.RetryAttempts = 5
		availabilityStatusesClient.RetryDuration = time.Second * 5
	} else {
		t.Fatalf("Authorization is Failed: %s", err.Error())
		t.Fail()
	}

	availabilityStatus, err := availabilityStatusesClient.GetByResource(context.Background(), resourceID, "", "")

	if err != nil {
		ct.Foreground(ct.Red, false)
		t.Fatalf("Cant connect with azure resourcehealth api service: %s", err.Error())
		ct.Foreground(ct.White, false)
		t.Fail() // So if error is not null the test must be fail
	}

	// Checking that resource validated is healthy available
	if availabilityStatus.Properties.AvailabilityState != resourcehealth.Available {
		ct.Foreground(ct.Red, false)
		logger.Logf(t, "Resource %s is unhealthy status: \t%s", resourceName, fmt.Sprint(availabilityStatus.Properties.AvailabilityState))
		ct.Foreground(ct.White, false)
		t.Error("Resource " + resourceName + " is unhealthy :( . Please check resource config")
		t.Fail() // So if resource is unhealthy the test should be fail
	} else {
		result = true
		ct.Foreground(ct.Green, true)
		logger.Logf(t, "Validation complete! Resource "+resourceName+" is: "+fmt.Sprint(resourcehealth.Available))
		ct.Foreground(ct.White, false)
	}
	return
}
