# terraform-provider-copperfield

This is a small, single-resource Terraform provider structured to illustrate a plan
problem I've been experiencing.

### What happens

The resource schema includes two `Computed` lists side-by-side in a nested resource
hierarchy:
- `candidates` (list of strings)
- `venues` (list of objects)

When modifying an unrelated attribute within the resource, members of the `venues` 
list vanish from the plan without a trace:

```
Terraform will perform the following actions:

  # copperfield_tour.sol_19830408 will be updated in-place
  ~ resource "copperfield_tour" "sol_19830408" {
      ~ cities = {
          ~ "new_york" = {
              ~ season     = "spring" -> "summer"                          <--- this value was deliberately changed
              ~ venues     = [
                  - {                                                      <+-- this list element has vanished (ooh! ahh!)
                      - capacity    = 500 -> null                           |
                      - coordinates = "40.690574, -74.045564" -> null       |
                    },                                                     <+
                ]
                # (2 unchanged attributes hidden)                          <--- the "candidates" list is in here, unmodified
            },
        }
        id     = "2185e2b4-bd0f-46de-a7f2-58fa967c6c39"
    }

```

Debugs indicate that the computed objects, including the `venue` details are present in the plan and marked as unknown
before plan modifiers are invoked:

```
Marking Computed attributes with null configuration values as unknown (known after apply) in the plan to prevent potential Terraform errors:
marking computed attribute that is null in the config as unknown: tf_attribute_path=AttributeName("id")
marking computed attribute that is null in the config as unknown: tf_attribute_path=AttributeName("cities").ElementKeyString("new_york").AttributeName("candidates")
marking computed attribute that is null in the config as unknown: tf_attribute_path=AttributeName("cities").ElementKeyString("new_york").AttributeName("promoter")
marking computed attribute that is null in the config as unknown: tf_attribute_path=AttributeName("cities").ElementKeyString("new_york").AttributeName("venues").ElementKeyInt(0).AttributeName("capacity")
marking computed attribute that is null in the config as unknown: tf_attribute_path=AttributeName("cities").ElementKeyString("new_york").AttributeName("venues").ElementKeyInt(0).AttributeName("coordinates")
marking computed attribute that is null in the config as unknown: tf_attribute_path=AttributeName("cities").ElementKeyString("new_york").AttributeName("venues")
```

But, when the plan modifiers are invoked, they're not run against attributes of objects inthe `venues` list:

```
Calling provider defined planmodifier.String: description="Once set, the value of this attribute in state will not change." tf_attribute_path=id
Calling provider defined planmodifier.Map:    description="Once set, the value of this attribute in state will not change." tf_attribute_path=cities
Calling provider defined planmodifier.List:   description="Once set, the value of this attribute in state will not change." tf_attribute_path=cities["new_york"].candidates
Calling provider defined planmodifier.List:   description="Once set, the value of this attribute in state will not change." tf_attribute_path=cities["new_york"].venues
Calling provider defined planmodifier.String: description="Once set, the value of this attribute in state will not change." tf_attribute_path=cities["new_york"].promoter
```

Does element `0` of the `venues` list still exist at this point? Not sure how to find out.

### What I expected to happen

The `venues` list should still be populated, plan modifiers should run against attributes of its members.
