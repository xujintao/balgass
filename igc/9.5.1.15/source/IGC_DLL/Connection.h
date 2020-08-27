#ifndef CONNECTION_H
#define CONNECTION_H

#define MU_CONNECT_FUNC 0x006BF89A // 1.04R
#define CONNECT_HOOK1 0x0043EDC5 // 1.04R
#define CONNECT_HOOK2 0x004E1D36 // 1.04R
#define GSCONNECT_HOOK1 0x004FF0B2 // 1.04R
#define GSCONNECT_HOOK2 0x0A440A47 // 1.04R
#define CHATSERVER_PORT 0x004B8F30+1 // 1.04R

enum CONNECTION_TYPE
{
	NOT_CONNECTED = 0,
	CS_CONNECTED,
	GS_CONNECTED,
};

typedef void (*tmuConnectToCS)(char* hostname, int port);
extern tmuConnectToCS muConnectToCS;

void OnConnect();
void OnGsConnect();

extern CONNECTION_TYPE g_Connection;
#endif