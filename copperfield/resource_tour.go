package copperfield

import (
	"context"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &resourceTour{}

type resourceTour struct{}

func (o *resourceTour) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tour"
}

func (o *resourceTour) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"cities": schema.MapNestedAttribute{
				PlanModifiers: []planmodifier.Map{mapplanmodifier.UseStateForUnknown()},
				Required:      true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"season": schema.StringAttribute{
							Required: true,
						},
						"promoter": schema.StringAttribute{
							Computed:      true,
							PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
						},
						"candidates": schema.ListAttribute{
							Computed:      true,
							ElementType:   types.StringType,
							PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
						},
						"venues": schema.ListNestedAttribute{
							Computed:      true,
							PlanModifiers: []planmodifier.List{listplanmodifier.UseStateForUnknown()},
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"capacity": schema.Int64Attribute{
										Computed:      true,
										PlanModifiers: []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
									},
									"coordinates": schema.StringAttribute{
										Computed:      true,
										PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (o *resourceTour) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// retrieve values from plan
	var plan tour
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(plan.Cities.Elements()) != 1 {
		resp.Diagnostics.AddError("configuration error", "sorry, this provider only supports a single city")
		return
	}

	cities := make(map[string]city, 1)
	d := plan.Cities.ElementsAs(ctx, &cities, false)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, ok := cities["new_york"]; !ok {
		resp.Diagnostics.AddError("configuration error", "sorry, this provider only supports the city of 'new_york'")
		return
	}

	venueObj, d := types.ObjectValueFrom(ctx, venue{}.attrTypes(), &venue{
		Capacity:    types.Int64Value(500),
		Coordinates: types.StringValue("40.690574, -74.045564"),
	})
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	venueList, d := types.ListValueFrom(ctx, venue{}.attrType(), []attr.Value{venueObj})

	stringList, d := types.ListValueFrom(ctx, types.StringType, []string{"airplane", "statue"})
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	ny := cities["new_york"]
	ny.Promoter = types.StringValue("reagan")
	ny.Candidates = stringList
	ny.Venues = venueList

	cities["new_york"] = ny
	cityMap, d := types.MapValueFrom(ctx, city{}.attrType(), cities)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state tour
	state.Id = types.StringValue(uuid.New().String())
	state.Cities = cityMap

	// set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (o *resourceTour) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// retrieve values from state
	var state tour
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	return
}

func (o *resourceTour) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// retrieve values from plan
	var plan tour
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

}

func (o *resourceTour) Delete(_ context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
}

type tour struct {
	Id     types.String `tfsdk:"id"`
	Cities types.Map    `tfsdk:"cities"`
}

type city struct {
	Season     types.String `tfsdk:"season"`
	Promoter   types.String `tfsdk:"promoter"`
	Candidates types.List   `tfsdk:"candidates"`
	Venues     types.List   `tfsdk:"venues"`
}

func (o city) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"season":     types.StringType,
		"promoter":   types.StringType,
		"candidates": types.ListType{ElemType: types.StringType},
		"venues":     types.ListType{ElemType: venue{}.attrType()},
	}
}

func (o city) attrType() attr.Type {
	return types.ObjectType{
		AttrTypes: o.attrTypes(),
	}
}

type venue struct {
	Capacity    types.Int64  `tfsdk:"capacity"`
	Coordinates types.String `tfsdk:"coordinates"`
}

func (o venue) attrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"capacity":    types.Int64Type,
		"coordinates": types.StringType,
	}
}

func (o venue) attrType() attr.Type {
	return types.ObjectType{
		AttrTypes: o.attrTypes(),
	}
}
