package identitytypeenum

const (
	Account = 0
	IamUser = 3
	IamRole = 4
)

var identityText = map[string]int{
	"root-account": Account,
	"iam-user":     IamUser,
	"assumed-role": IamRole,
}

func IdentityCode(code string) int {
	return identityText[code]
}
