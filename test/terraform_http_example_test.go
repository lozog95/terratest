package test

import (
	"testing"
	"github.com/gruntwork-io/terratest/terraform"
	"fmt"
	"github.com/gruntwork-io/terratest/util"
	"github.com/gruntwork-io/terratest/aws"
	"github.com/gruntwork-io/terratest/http"
	"time"
)

// An example of how to test the Terraform module in examples/terraform-http-example using Terratest.
func TerraformHttpExampleTest(t *testing.T) {
	t.Parallel()

	// Give this EC2 Instance and other resources in the Terraform code a name with a unique ID so it doesn't clash
	// with anything else in the AWS account.
	instanceName := fmt.Sprintf("terratest-http-example-%s", util.UniqueId())

	// Specify the text the EC2 Instance will return when we make HTTP requests to it.
	instanceText := fmt.Sprintf("Hello, %s!", util.UniqueId())

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	awsRegion := aws.PickRandomRegion(t)

	terraformOptions := terraform.Options {
		// The path to where our Terraform code is located
		TerraformDir: "../examples/terraform-http-example",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]string {
			"aws_region":    awsRegion,
			"instance_name": instanceName,
			"instance_text": instanceText,
		},
	}

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.Apply(t, terraformOptions)

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform output` to get the value of an output variable
	instanceUrl := terraform.Output(t, terraformOptions, "instance_url")

	// It can take a minute or so for the Instance to boot up, so retry a few times
	maxRetries := 15
	timeBetweenRetries := 5 * time.Second

	// Verify that we get back a 200 OK with the expected instanceText
	http_helper.HttpGetWithRetry(t, instanceUrl, 200, instanceText, maxRetries, timeBetweenRetries)
}

