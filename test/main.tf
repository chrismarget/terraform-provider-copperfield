terraform {
  required_providers {
    copperfield = {
      source = "example.com/chrismarget/copperfield"
    }
  }
}

provider "copperfield" {}

resource "copperfield_tour" "sol_19830408" {
  cities = {
    "new_york" = {      // keys other than 'new_york' are not supported
      season = "spring" // change this value to experience the problem
    }
  }
}

output "x" {
  value = copperfield_tour.sol_19830408
}