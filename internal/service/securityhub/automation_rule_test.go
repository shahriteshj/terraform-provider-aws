// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package securityhub_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/securityhub/types"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfsecurityhub "github.com/hashicorp/terraform-provider-aws/internal/service/securityhub"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func testAccAutomationRule_basic(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_securityhub_automation_rule.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SecurityHubEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAutomationRuleDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAutomationRuleConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutomationRuleExists(ctx, resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAutomationRule_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_securityhub_automation_rule.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SecurityHubEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAutomationRuleDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAutomationRuleConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutomationRuleExists(ctx, resourceName),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfsecurityhub.ResourceAutomationRule(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccAutomationRule_stringFilters(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_securityhub_automation_rule.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SecurityHubEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAutomationRuleDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAutomationRuleConfig_stringFilters(rName, string(types.StringFilterComparisonEquals), "1234567890"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutomationRuleExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.aws_account_id.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.aws_account_id.0.comparison", string(types.StringFilterComparisonEquals)),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.aws_account_id.0.value", "1234567890"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAutomationRuleConfig_stringFilters(rName, string(types.StringFilterComparisonContains), "0987654321"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutomationRuleExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.aws_account_id.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.aws_account_id.0.comparison", string(types.StringFilterComparisonContains)),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.aws_account_id.0.value", "0987654321"),
				),
			},
		},
	})
}

func testAccAutomationRule_numberFilters(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_securityhub_automation_rule.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SecurityHubEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAutomationRuleDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAutomationRuleConfig_numberFilters(rName, "eq = 5"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutomationRuleExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.confidence.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.confidence.0.eq", "5"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAutomationRuleConfig_numberFilters(rName, "lte = 50"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutomationRuleExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.confidence.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.confidence.0.lte", "50"),
				),
			},
		},
	})
}

func testAccAutomationRule_dateFilters(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_securityhub_automation_rule.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	endDate := time.Now().Add(5 * time.Minute).Format(time.RFC3339)
	startDate := time.Now().Format(time.RFC3339)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SecurityHubEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAutomationRuleDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAutomationRuleConfig_dateFiltersAbsoluteRange(rName, startDate, endDate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutomationRuleExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.created_at.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.created_at.0.end", endDate),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.created_at.0.start", startDate),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAutomationRuleConfig_dateFiltersRelativeRange(rName, string(types.DateRangeUnitDays), 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutomationRuleExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.created_at.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.created_at.0.date_range.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.created_at.0.date_range.0.unit", string(types.DateRangeUnitDays)),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.created_at.0.date_range.0.value", "10"),
				),
			},
		},
	})
}

func testAccAutomationRule_mapFilters(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_securityhub_automation_rule.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SecurityHubEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAutomationRuleDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAutomationRuleConfig_mapFilters(rName, string(types.MapFilterComparisonEquals), "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutomationRuleExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.resource_details_other.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.resource_details_other.0.comparison", string(types.MapFilterComparisonEquals)),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.resource_details_other.0.key", "key1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.resource_details_other.0.value", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAutomationRuleConfig_mapFilters(rName, string(types.MapFilterComparisonContains), "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutomationRuleExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.resource_details_other.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.resource_details_other.0.comparison", string(types.MapFilterComparisonContains)),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.resource_details_other.0.key", "key2"),
					resource.TestCheckResourceAttr(resourceName, "criteria.0.resource_details_other.0.value", "value2"),
				),
			},
		},
	})
}

func testAccAutomationRule_tags(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_securityhub_automation_rule.test"
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, names.SecurityHubEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckAutomationRuleDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccAutomationRuleConfig_tags(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutomationRuleExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAutomationRuleConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutomationRuleExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccAutomationRuleConfig_tags(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAutomationRuleExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
		},
	})
}

func testAccCheckAutomationRuleExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).SecurityHubClient(ctx)

		_, err := tfsecurityhub.FindAutomationRuleByARN(ctx, conn, rs.Primary.ID)

		return err
	}
}

func testAccCheckAutomationRuleDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).SecurityHubClient(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_securityhub_automation_rule" {
				continue
			}

			_, err := tfsecurityhub.FindAutomationRuleByARN(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("Security Hub Automation Rule (%s) still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccAutomationRuleConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_securityhub_automation_rule" "test" {
  description = "test description"
  rule_name   = %[1]q
  rule_order  = 1

  actions {
    finding_fields_update {
      severity {
        label   = "LOW"
        product = "0.0"
      }

      types = ["Software and Configuration Checks/Industry and Regulatory Standards"]

      user_defined_fields = {
        key = "value"
      }
    }
    type = "FINDING_FIELDS_UPDATE"
  }

  criteria {
    aws_account_id {
      comparison = "EQUALS"
      value      = "1234567890"
    }
  }
}
`, rName)
}

func testAccAutomationRuleConfig_stringFilters(rName, comparison, value string) string {
	return fmt.Sprintf(`
resource "aws_securityhub_automation_rule" "test" {
  description = "test description"
  rule_name   = %[1]q
  rule_order  = 1

  actions {
    finding_fields_update {
      user_defined_fields = {
        key = "value"
      }
    }
    type = "FINDING_FIELDS_UPDATE"
  }

  criteria {
    aws_account_id {
      comparison = %[2]q
      value      = %[3]q
    }
  }
}
`, rName, comparison, value)
}

func testAccAutomationRuleConfig_numberFilters(rName, value string) string {
	return fmt.Sprintf(`
resource "aws_securityhub_automation_rule" "test" {
  description = "test description"
  rule_name   = %[1]q
  rule_order  = 1

  actions {
    finding_fields_update {
      user_defined_fields = {
        key = "value"
      }
    }
    type = "FINDING_FIELDS_UPDATE"
  }

  criteria {
    confidence {
      %[2]s
    }
  }
}
`, rName, value)
}

func testAccAutomationRuleConfig_dateFiltersAbsoluteRange(rName, start, end string) string {
	return fmt.Sprintf(`
resource "aws_securityhub_automation_rule" "test" {
  description = "test description"
  rule_name   = %[1]q
  rule_order  = 1

  actions {
    finding_fields_update {
      user_defined_fields = {
        key = "value"
      }
    }
    type = "FINDING_FIELDS_UPDATE"
  }

  criteria {
    created_at {
      end   = %[3]q
      start = %[2]q
    }
  }
}
`, rName, start, end)
}

func testAccAutomationRuleConfig_dateFiltersRelativeRange(rName, unit string, value int) string {
	return fmt.Sprintf(`
resource "aws_securityhub_automation_rule" "test" {
  description = "test description"
  rule_name   = %[1]q
  rule_order  = 1

  actions {
    finding_fields_update {
      user_defined_fields = {
        key = "value"
      }
    }
    type = "FINDING_FIELDS_UPDATE"
  }

  criteria {
    created_at {
      date_range {
        unit  = %[2]q
        value = %[3]d
      }
    }
  }
}
`, rName, unit, value)
}

func testAccAutomationRuleConfig_mapFilters(rName, comparison, key, value string) string {
	return fmt.Sprintf(`
resource "aws_securityhub_automation_rule" "test" {
  description = "test description"
  rule_name   = %[1]q
  rule_order  = 1

  actions {
    finding_fields_update {
      user_defined_fields = {
        key = "value"
      }
    }
    type = "FINDING_FIELDS_UPDATE"
  }

  criteria {
    resource_details_other {
      comparison = %[2]q
      key        = %[3]q
      value      = %[4]q
    }
  }
}
`, rName, comparison, key, value)
}

func testAccAutomationRuleConfig_tags(rName, key, value string) string {
	return fmt.Sprintf(`
resource "aws_securityhub_automation_rule" "test" {
  description = "test description"
  rule_name   = %[1]q
  rule_order  = 1

  actions {
    finding_fields_update {
      user_defined_fields = {
        key = "value"
      }
    }
    type = "FINDING_FIELDS_UPDATE"
  }

  criteria {
    aws_account_id {
      comparison = "EQUALS"
      value      = "1234567890"
    }
  }
  tags = {
    %[2]q = %[3]q
  }
}
`, rName, key, value)
}

func testAccAutomationRuleConfig_tags2(rName, key, value, key2, value2 string) string {
	return fmt.Sprintf(`
resource "aws_securityhub_automation_rule" "test" {
  description = "test description"
  rule_name   = %[1]q
  rule_order  = 1

  actions {
    finding_fields_update {
      user_defined_fields = {
        key = "value"
      }
    }
    type = "FINDING_FIELDS_UPDATE"
  }

  criteria {
    aws_account_id {
      comparison = "EQUALS"
      value      = "1234567890"
    }
  }
  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, key, value, key2, value2)
}
