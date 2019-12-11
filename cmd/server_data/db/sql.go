package db

// Lookup returns the named statement.
func Lookup(name string) string {
	return index[name]
}

var index = map[string]string{
	"account-find-passwd":    accountFindPasswd,
	"ban_machine-find-count": banMachineFindCount,
}

var accountFindPasswd = `
SELECT
 memb__pwd
,bloc_code
FROM MEMB_INFO
WHERE
memb___id = ?
`
var banMachineFindCount = `
SELECT
 count(*)
FROM IGC_MachineID_Banned
WHERE
HWID = ?
`
