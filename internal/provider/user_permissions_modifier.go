package provider

import (
	"context"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/oliver-binns/googleplay-go/users"
)

func expandUserPermissionsPlanModifier(permissions path.Path) planmodifier.Set {
	return &userPermissionsExpansionModifier{
		permissions: permissions,
	}
}

type userPermissionsExpansionModifier struct {
	permissions path.Path
}

func (m *userPermissionsExpansionModifier) Description(ctx context.Context) string {
	return "Expands Google Play Permissions for users to include implicitly granted permissions"
}

func (m *userPermissionsExpansionModifier) MarkdownDescription(ctx context.Context) string {
	return "Expands Google Play Permissions for users to include implicitly granted permissions"
}

func (m *userPermissionsExpansionModifier) PlanModifySet(
	ctx context.Context,
	req planmodifier.SetRequest,
	resp *planmodifier.SetResponse,
) {
	// fetch declared permissions:
	var permissions types.Set // contrived example is using string attributes
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, m.permissions, &permissions)...)

	// expand the permissions
	expanded_permissions := []users.DeveloperLevelPermission{}
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
