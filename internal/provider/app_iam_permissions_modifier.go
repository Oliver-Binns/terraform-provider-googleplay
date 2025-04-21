package provider

import (
	"context"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/oliver-binns/googleplay-go/users"
)

func expandAppPermissionsPlanModifier(permissions path.Path) planmodifier.Set {
	return &appPermissionsExpansionModifier{
		permissions: permissions,
	}
}

type appPermissionsExpansionModifier struct {
	permissions path.Path
}

func (m *appPermissionsExpansionModifier) Description(ctx context.Context) string {
	return "Expands Google Play Permissions to include explicitly granted permissions"
}

func (m *appPermissionsExpansionModifier) MarkdownDescription(ctx context.Context) string {
	return "Expands Google Play Permissions to include explicitly granted permissions"
}

func (m *appPermissionsExpansionModifier) PlanModifySet(
	ctx context.Context,
	req planmodifier.SetRequest,
	resp *planmodifier.SetResponse,
) {
	// fetch declared permissions:
	var permissions types.Set // contrived example is using string attributes
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, m.permissions, &permissions)...)

	// expand the permissions
	expanded_permissions := []users.AppLevelPermission{}
	diag := req.PlanValue.ElementsAs(ctx, &permissions, false)
	resp.Diagnostics.Append(diag...)

	for _, permission := range expanded_permissions {
		for _, inherited := range permission.Expand() {
			if !slices.Contains(expanded_permissions, inherited) {
				// missing from the plan:
				expanded_permissions = append(expanded_permissions, inherited)
			}
		}
	}

	// write the expanded permissions back to the plan
	planValue, diag := types.SetValueFrom(ctx, req.PlanValue.ElementType(ctx), permissions)
	resp.Diagnostics.Append(diag...)
	resp.PlanValue = planValue
}
