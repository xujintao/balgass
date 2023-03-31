package db

// Lookup returns the named statement.
func Lookup(name string) string {
	return index[name]
}

var index = map[string]string{
	"users-find-user":            usersFindUser,
	"users-update-block":         usersUpdateBlock,
	"ban_machines-find-count":    banMachinesFindCount,
	"user_login_historys-insert": userLoginHistorysInsert,
	"user_states-find-user":      userStatesFindUser,
	"user_states-insert":         userStatesInsert,
	"user_states-update-connect": userStatesUpdateConnect,
	"user_states-update-disconn": userStatesUpdateDisconn,
	"user_chars-find-count":      userCharsFindCount,
	"user_chars-insert":          userCharsInsert,
	"vips-find-user":             vipsFindUser,
	"vips-find-user-update":      vipsFindUserUpdate,
	"vips-update-renew":          vipsUpdateRenew,
	"vips-update-upgrade":        vipsUpdateUpgrade,
	"chars-find-chars":           charsFindChars,
	"guild_members-find-status":  guildMembersFindStatus,
	"battlecores-find-bfname":    battlecoreFindBFName,
}

var usersFindUser = `
SELECT
 memb_guid
,memb__pwd
,bloc_code
FROM MEMB_INFO
WHERE
memb___id = @p1
`

var usersUpdateBlock = `
UPDATE MEMB_INFO
SET
 bloc_code = :bloc_code
WHERE memb___id = :memb___id
`

var banMachinesFindCount = `
SELECT
 count(*)
FROM IGC_MachineID_Banned
WHERE
HWID = @p1
`

var userLoginHistorysInsert = `
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

var userStatesFindUser = `
SELECT
 count(*)
FROM MEMB_STAT S
INNER JOIN MEMB_INFO I ON S.memb___id = I.memb___id
WHERE
I.memb___id = @p1
`

var userStatesInsert = `
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

var userStatesUpdateConnect = `
UPDATE MEMB_STAT
SET
 ConnectStat = :ConnectStat
,ServerName = :ServerName
,IP = :IP
,ConnectTM = :ConnectTM
WHERE memb___id = :memb___id
`

var userStatesUpdateDisconn = `
UPDATE MEMB_STAT
SET
 ConnectStat = :ConnectStat
,DisConnectTM = :DisConnectTM
WHERE memb___id = :memb___id
`

var userCharsFindCount = `
SELECT
 count(*)
FROM AccountCharacter
WHERE id = @p1
`

var userCharsInsert = `
INSERT INTO AccountCharacter (
 Id
)
VALUES (
 @p1
)
`

var userCharsFindchar = `
SELECT
 GameID1
,GameID2
,GameID3
,GameID4
,GameID5
,MoveCnt
,WarehouseExpansion
,Summoner
,RageFighter
,SecCode
FROM AccountCharacter
WHERE Id = @p1
`

var vipsFindUser = `
SELECT
 Date
,Type
FROM T_VIPList
WHERE AccountID = @p1
  AND Date > @p2
`

var vipsFindUserUpdate = `
SELECT
 count(*)
,(SELECT count(*) FROM T_VIPList WHERE AccountID = @p1 AND Date > @p2)
FROM T_VIPList
WHERE AccountID = @p1
`

var vipsUpdateRenew = `
UPDATE T_VIPList
SET
 Date = :Date
,Type = :Type
WHERE AccountID = :AccountID
`

var vipsUpdateUpgrade = `
UPDATE T_VIPList
SET
 Date = :Date
,Type = :Type
WHERE AccountID = :AccountID
  AND Type < :Type
`

var charsFindChars = `
SELECT
 Name
,Class
,cLevel
,mLevel
,CtlCode
,RESETS
,Inventory
FROM Character
WHERE AccountID = @p1
`

var guildMembersFindStatus = `
SELECT
 G_Status
FROM GuildMember
WHERE NAME = @p1
`

var battlecoreFindBFName = `
SELECT
 UBFName
,ServerCode
FROM BattleCoreInfo
WHERE Name = @p1
`