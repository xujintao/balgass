#include "stdafx.h"
#include "offset.h"
#include "aes.h"
#include "xor.h"
#include "model.h"
#include "c1c2.h"

void ParsePacket(void* PackStream, int unk1, int unk2)
{
	BYTE* buff;
	while (true)
	{
		__asm {
			MOV ECX, PackStream;
			MOV EDX, PARSE_PACKET_STREAM;
			CALL EDX;
			MOV buff, EAX;
		}
		if (!buff)
			break;
		BYTE DecBuff[7024];
		unsigned int DecSize;
		int proto;
		int size;
		int enc;
		switch (buff[0])
		{
		case 0xC1:
			proto = buff[2];
			size = buff[1];
			enc = 0;
			break;
		case 0xC2:
			proto = buff[3];
			size = *(WORD*)&buff[1];
			enc = 0;
			break;
		case 0xC3:
			enc = 1;
			size = buff[1];
			DecSize = CAes::Decrypt(&DecBuff[2], &buff[2], size - 2);
			DecBuff[0] = 0xC1;
			DecBuff[1] = DecSize + 2;
			size = DecSize + 2;
			buff = DecBuff;
			proto = DecBuff[2];
			break;
		case 0xC4:
			enc = 1;
			size = MAKEWORD(buff[2], buff[1]);
			DecSize = CAes::Decrypt(&DecBuff[3], &buff[3], size - 3);
			DecBuff[0] = 0xC2;
			DecBuff[2] = LOBYTE(DecSize + 3);
			DecBuff[1] = HIBYTE(DecSize + 3);
			size = DecSize + 3;
			buff = DecBuff;
			proto = buff[3];
			break;
		}
		if (unk1 == 1)
		{
			typedef int(*tProtocolCore2)(int, int, BYTE*, int, int);
			tProtocolCore2 ProtocolCore2 = (tProtocolCore2)PROTOCOL_CORE1;
			ProtocolCore2(unk2, proto, buff, size, enc);
		}
		else
		{
			typedef int(*tProtocolCore)(int, BYTE*, int, int);
			tProtocolCore ProtocolCore = (tProtocolCore)PROTOCOL_CORE2;
			ProtocolCore(proto, buff, size, enc); // Main.exe protocolcore
		}
	}
}

void __declspec(naked) muSendPacket(BYTE* buff, int len)
{
	__asm
	{
		PUSH EBP;
		MOV EBP, ESP;
		MOV EAX, len;
		PUSH EAX;
		PUSH buff;
		MOV ECX, DWORD PTR DS : [MU_SENDER_CLASS];
		MOV EDX, MU_SEND_PACKET;
		CALL EDX;
		MOV ESP, EBP;
		POP EBP;
		RETN;
	}
}

void SendPacket(BYTE* buff, int len, int enc, int unk1)
{
	if (buff[0] == 0xC1){
		CXor::Dec(&buff[3], 3, len - 3);
	}
	else if (buff[0] == 0xC2){
		CXor::Dec(&buff[4], 4, len - 4);
	}
	if (buff[2] == 0xF1 && buff[3] == 0x01) // ban by hwid
	{
		BYTE NewIDPass[200];
		memset(NewIDPass, 0x00, sizeof(NewIDPass));

		PMSG_IDPASS_OLD * lpMsg1 = (PMSG_IDPASS_OLD *)buff;
		PMSG_IDPASS_NEW * lpMsg2 = (PMSG_IDPASS_NEW *)NewIDPass;

		lpMsg2->h.c = 0xC1;
		lpMsg2->h.size = sizeof(PMSG_IDPASS_NEW);
		lpMsg2->h.headcode = 0xF1;
		lpMsg2->subcode = 0x01;
		lpMsg2->TickCount = lpMsg1->TickCount;

		memcpy(lpMsg2->Id, lpMsg1->Id, 10);
		memcpy(lpMsg2->Pass, lpMsg1->Pass, 20);
		memcpy(lpMsg2->CliSerial, lpMsg1->CliSerial, 16);
		memcpy(lpMsg2->CliVersion, lpMsg1->CliVersion, 8);
		lpMsg2->ServerSeason = 0x40;

		//memcpy(Login, buff, buff[1]);
		buff = NewIDPass;
		len = lpMsg2->h.size;
		//CXor::Enc(&Login[3], 3, len - 3);

		// Fix account id save for s9
		//char AccountID[11];
		//AccountID[10] = 0;
		//memcpy(AccountID, lpMsg2->Id, 10);
		//CXor::BuxConvert(AccountID, 10);
		//MemCpy(ACCOUNT_ID, AccountID, sizeof(AccountID));

		enc = 1;
	}
	else if (buff[2] == 0x0E)
	{
		BYTE NewClientTime[200];
		memset(NewClientTime, 0x00, sizeof(NewClientTime));

		PMSG_CLIENTTIME_OLD * lpMsg1 = (PMSG_CLIENTTIME_OLD *)buff;
		PMSG_CLIENTTIME_NEW * lpMsg2 = (PMSG_CLIENTTIME_NEW *)NewClientTime;

		lpMsg2->h.c = 0xC1;
		lpMsg2->h.size = sizeof(PMSG_CLIENTTIME_NEW);
		lpMsg2->h.headcode = 0x0E;

		lpMsg2->TimeH = HIBYTE(lpMsg1->Time);
		lpMsg2->TimeL = LOBYTE(lpMsg1->Time);
		//lpMsg2->Agility = *(WORD*)((*(DWORD*)IGC_STAT) + 0x11A) + *(WORD*)((*(DWORD*)IGC_STAT) + 0x134);
		lpMsg2->AttackSpeed = lpMsg1->AttackSpeed;
		lpMsg2->MagicSpeed = lpMsg1->MagicSpeed;
		//memcpy(lpMsg2->Version, DLL_VERSION, strlen(DLL_VERSION));

		buff = NewClientTime;
		len = lpMsg2->h.size;
		enc = 1;
	}

	if (buff[0] == 0xC1)
	{
		CXor::Enc(&buff[3], 3, len - 3);
	}
	else if (buff[0] == 0xC2)
	{
		CXor::Enc(&buff[4], 4, len - 4);
	}

	if (enc)
	{
		BYTE EncBuff[7024];
		unsigned int EncSize;
		if (buff[0] == 0xC1)
		{
			EncSize = CAes::Encrypt(&EncBuff[2], &buff[2], len - 2);
			EncBuff[0] = 0xC3;
			EncBuff[1] = EncSize + 2;
			buff = EncBuff;
			len = EncSize + 2;
		}
		else if (buff[0] == 0xC2)
		{
			EncSize = CAes::Encrypt(&EncBuff[3], &buff[3], len - 3);
			EncBuff[0] = 0xC4;
			EncBuff[1] = HIBYTE(EncSize + 3);
			EncBuff[2] = LOBYTE(EncSize + 3);
			buff = EncBuff;
			len = EncSize + 3;
		}
	}

	muSendPacket(buff, len);
}