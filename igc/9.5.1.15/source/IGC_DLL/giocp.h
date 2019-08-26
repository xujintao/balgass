#ifndef GIOCP_H
#define GIOCP_H

// function headers

void ParsePacket(void* PackStream, int unk1, int unk2);
void DataSend(void* pData, int size);
void SendPacket(BYTE* buff, int len, int enc, int unk1);
extern SOCKET* g_MuSocket;

#endif