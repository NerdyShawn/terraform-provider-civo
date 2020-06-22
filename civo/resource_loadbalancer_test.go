package civo

import (
	"fmt"
	"testing"

	"github.com/civo/civogo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// example.Widget represents a concrete Go type that represents an API resource
func TestAccCivoLoadBalancer_basic(t *testing.T) {
	var loadBalancer civogo.LoadBalancer

	// generate a random name for each test run
	resName := "civo_loadbalancer.foobar"
	var domainName = acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCivoLoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				// use a dynamic configuration with the random name from above
				Config: testAccCheckCivoLoadBalancerConfigBasic(domainName),
				// compose a basic test, checking both remote and local values
				Check: resource.ComposeTestCheckFunc(
					// query the API to retrieve the widget object
					testAccCheckCivoLoadBalancerResourceExists(resName, &loadBalancer),
					// verify remote values
					testAccCheckCivoLoadBalancerValues(&loadBalancer, domainName),
					// verify local values
					resource.TestCheckResourceAttr(resName, "protocol", "http"),
					resource.TestCheckResourceAttr(resName, "port", "80"),
				),
			},
		},
	})
}

// func TestAccCivoLoadBalancer_update(t *testing.T) {
// 	var firewallRule civogo.FirewallRule

// 	// generate a random name for each test run
// 	resName := "civo_firewall_rule.testrule"
// 	var firewallRuleName = acctest.RandomWithPrefix("rename-fw-rule")

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckCivoFirewallRuleDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: testAccCheckCivoFirewallRuleConfigUpdates(firewallRuleName),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckCivoFirewallRuleResourceExists(resName, &firewallRule),
// 					resource.TestCheckResourceAttr(resName, "protocol", "tcp"),
// 					resource.TestCheckResourceAttr(resName, "start_port", "443"),
// 				),
// 			},
// 			{
// 				// use a dynamic configuration with the random name from above
// 				Config: testAccCheckCivoFirewallRuleConfigUpdates(firewallRuleName),
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckCivoFirewallRuleResourceExists(resName, &firewallRule),
// 					testAccCheckCivoFirewallRuleUpdated(&firewallRule),
// 					resource.TestCheckResourceAttr(resName, "protocol", "tcp"),
// 					resource.TestCheckResourceAttr(resName, "start_port", "443"),
// 				),
// 			},
// 		},
// 	})
// }

func testAccCheckCivoLoadBalancerValues(loadBalancer *civogo.LoadBalancer, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if loadBalancer.Hostname != name {
			return fmt.Errorf("bad protocol, expected \"%s\", got: %#v", name, loadBalancer.Hostname)
		}
		return nil
	}
}

// testAccCheckExampleResourceExists queries the API and retrieves the matching Widget.
func testAccCheckCivoLoadBalancerResourceExists(n string, loadBalancer *civogo.LoadBalancer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		// retrieve the configured client from the test setup
		client := testAccProvider.Meta().(*civogo.Client)
		resp, err := client.FindLoadBalancer(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("LoadBalancer not found: (%s) %s", rs.Primary.ID, err)
		}

		// If no error, assign the response Widget attribute to the widget pointer
		*loadBalancer = *resp

		// return fmt.Errorf("Domain (%s) not found", rs.Primary.ID)
		return nil
	}
}

// func testAccCheckCivoLoadBalancerUpdated(firewall *civogo.FirewallRule) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		if firewall.Protocol != "tcp" {
// 			return fmt.Errorf("bad protocol, expected \"%s\", got: %#v", "tcp", firewall.Protocol)
// 		}
// 		if firewall.StartPort != "443" {
// 			return fmt.Errorf("bad port, expected \"%s\", got: %#v", "443", firewall.StartPort)
// 		}
// 		return nil
// 	}
// }

func testAccCheckCivoLoadBalancerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*civogo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "civo_loadbalancer" {
			continue
		}

		_, err := client.FindLoadBalancer(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("LoadBlanacer still exists")
		}
	}

	return nil
}

func testAccCheckCivoLoadBalancerConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "civo_instance" "vm" {
	hostname = "instance-%s"
}

resource "civo_loadbalancer" "foobar" {
	hostname = "%s"
	protocol = "http"
	port = 80
	max_request_size = 30
	policy = "round_robin"
	health_check_path = "/"
	max_conns = 10
	fail_timeout = 40
	depends_on = [civo_instance.vm]

	backend {
		instance_id = civo_instance.vm.id
		protocol =  "http"
		port = 80
	}
}
`, name, name)
}

// func testAccCheckCivoLoadBalancerConfigUpdates(name string) string {
// 	return fmt.Sprintf(`
// resource "civo_firewall" "foobar" {
// 	name = "%s"
// }

// resource "civo_firewall_rule" "testrule" {
// 	firewall_id = civo_firewall.foobar.id
// 	protocol = "tcp"
// 	start_port = "443"
// 	end_port = "443"
// 	cidr = ["192.168.1.2/32"]
// 	direction = "inbound"
// 	label = "web"
// }
// `, name)
// }
