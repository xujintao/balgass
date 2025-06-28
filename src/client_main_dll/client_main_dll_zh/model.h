#pragma once
struct PBMSG_HEAD	// Packet - Byte Type
{
public:
	void set(LPBYTE lpBuf, BYTE head, BYTE size)	// line : 18
	{
		lpBuf[0] = 0xC1;
		lpBuf[1] = size;
		lpBuf[2] = head;
	};	// line : 22

	void setE(LPBYTE lpBuf, BYTE head, BYTE size)	// line : 25
	{
		lpBuf[0] = 0xC3;
		lpBuf[1] = size;
		lpBuf[2] = head;
	};	// line : 29

	BYTE c;
	BYTE size;
	BYTE headcode;
};

struct PMSG_CLIENTTIME_OLD
{
	PBMSG_HEAD h;
	DWORD Time;	// 4
	WORD AttackSpeed;	// 8
	WORD MagicSpeed;	// A
};

#pragma pack (1)
struct PMSG_IDPASS_OLD
{
	PBMSG_HEAD h;
	BYTE subcode;	// 3
	char Id[10];	// 4
	char Pass[10];	// E
	DWORD TickCount;	// 18
	BYTE CliVersion[8];	// 1C
	BYTE CliSerial[16];	// 21
};

struct PMSG_IDPASS_NEW
{
	PBMSG_HEAD h;
	BYTE subcode;	// 3
	char Id[10];	// 4
	char Pass[20];	// E
	char HWID[100];
	DWORD TickCount;	// 18
	BYTE CliVersion[8];	// 1C
	BYTE CliSerial[16];	// 21
	DWORD ServerSeason;
};
struct PMSG_CLIENTTIME_NEW
{
	PBMSG_HEAD h;
	WORD TimeH;
	WORD TimeL;
	WORD AttackSpeed;	// 8
	WORD Agility;
	WORD MagicSpeed;	// A
	char Version[10];
};
#pragma pack ()