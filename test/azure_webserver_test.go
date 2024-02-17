package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// You normally want to run this under a separate "Testing" subscription
// For lab purposes you will use your assigned subscription under the Cloud Dev/Ops program tenant
var subscriptionID string = "3a57a2b6-bb02-44a8-b964-575752c42004"

func TestAzureLinuxVMCreation(t *testing.T) {
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../",
		// Override the default terraform variables
		Vars: map[string]interface{}{
			"labelPrefix": "enca0002",
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of output variable
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	nicName := terraform.Output(t, terraformOptions, "nic_name")
	imageOffer := terraform.Output(t, terraformOptions, "image_offer")
	imageSKU := terraform.Output(t, terraformOptions, "image_sku")
	
	// Confirm VM exists
	assert.True(t, azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID))

	//confirm nic exists and is connected to vm
	nicList := azure.GetVirtualMachineNics(t, vmName, resourceGroupName, subscriptionID)
	assert.Contains(t, nicList, nicName)

	//confirm vm is running the correct version, values taken fron main.tf
	correctOffer := "0001-com-ubuntu-server-jammy"
    correctSKU := "22_04-lts-gen2"
    assert.Equal(t, correctOffer, imageOffer)
    assert.Equal(t, correctSKU, imageSKU)
}
