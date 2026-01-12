package domain

const (
	PermissionAdministrator int64 = 1 << 62
)

var AllPermissionsList = []int64{

	PermissionAdministrator,
}

var ValidPermissionsMask int64

func init() {
	for _, p := range AllPermissionsList {
		ValidPermissionsMask |= p
	}
}

func IsValidPermission(p int64) bool {
	return (p & ^ValidPermissionsMask) == 0
}
