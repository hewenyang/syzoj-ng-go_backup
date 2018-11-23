package group

type standardGroupUserRole struct {
	Role int `json:"role"`
}
type standardGroupPolicy struct {
	MemberCreateProblemset bool `json:"member_create_problemset"`
}
type standardGroupProvider struct{}

func (standardGroupProvider) GetDefaultGroupPolicy() GroupPolicy {
	return &standardGroupPolicy{}
}

func (g *standardGroupPolicy) GetDefaultRole() GroupUserRole {
	return &standardGroupUserRole{}
}

func (g *standardGroupPolicy) GetGuestRole() GroupUserRole {
	return &standardGroupUserRole{}
}

func (g *standardGroupPolicy) GetRegisteredUserRole() GroupUserRole {
	return &standardGroupUserRole{}
}

func (g *standardGroupPolicy) GetCreatorRole() GroupUserRole {
	return &standardGroupUserRole{Role: 3}
}

func (g *standardGroupPolicy) CheckPrivilege(u_ GroupUserRole, p GroupPrivilege) error {
	u := u_.(*standardGroupUserRole)
	switch u.Role {
	case 0: // Guest
		return GroupPermissionDeniedError
	case 1: // Member
		switch p {
		case GroupViewProblemsetPrivilege:
			return nil
		case GroupCreateProblemsetPrivilege:
			if g.MemberCreateProblemset {
				return nil
			}
			return GroupPermissionDeniedError
		}
		return GroupPermissionDeniedError
	case 2: // Admin
		switch p {
		case GroupViewProblemsetPrivilege:
			return nil
		case GroupCreateProblemsetPrivilege:
			return nil
		case GroupManageProblemsetPrivilege:
			return nil
		}
		return GroupPermissionDeniedError
	case 3: // Owner
		return nil
	}
	panic(GroupPermissionInvalidError)
}
