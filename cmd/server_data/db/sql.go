package db

// Lookup returns the named statement.
func Lookup(name string) string {
	return index[name]
}

var index = map[string]string{
	"account-find-passwd":         accountFindPasswd,
	"ban_machine-find-count":      banMachineFindCount,
	"account_join_history-insert": accountJoinHistoryInsert,
	"account_state-insert":        accountStateInsert,
	"account_state-update":        accountStateUpdate,
}

var accountFindPasswd = `
SELECT
 memb_guid
,memb__pwd
,bloc_code
FROM MEMB_INFO
WHERE
memb___id = @p1
`
var banMachineFindCount = `
SELECT
 count(*)
FROM IGC_MachineID_Banned
WHERE
HWID = @p1
`

var accountJoinHistoryInsert = `
INSERT INTO ConnectionHistory (
 AccountID
,ServerName
,IP
,Date
,State
,HWID
)
VALUES (
 :AccountID
,:ServerName
,:IP
,:Date
,:State
,:HWID
)
`

var accountStateInsert = `
INSERT INTO MEMB_STAT (
 memb___id
,ConnectStat
,ServerName
,IP
,ConnectTM
)
VALUES (
 :memb___id
,:ConnectStat
,:ServerName
,:IP
,:ConnectTM
)
`

var accountStateUpdate = `
UPDATE MEMB_STAT
SET
 ConnectStat = :ConnectStat
,ServerName = :ServerName
,IP = :IP
,ConnectTM = :ConnectTM
WHERE memb___id = ?
`
