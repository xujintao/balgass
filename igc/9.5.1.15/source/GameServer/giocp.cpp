// Giocp.cpp
//------------------------------------------
// Decompiled by Deathway
// Date : 2007-03-09
//------------------------------------------

// GS-N 0.99.60T 0x00473020 - Status : Completed :)
//	GS-N	1.00.18	JPN	0x00489FD0	-	Completed
#include "stdafx.h"
#include "spe.h"
#include "giocp.h"
#include "TLog.h"
#include "GameMain.h"
#include "GameServer.h"
#include "user.h"
#include "configread.h"
#include "PacketEncrypt.h"

CIOCP IOCP;

void CIOCP::GiocpInit()
{
	ExSendBuf = new unsigned char[MAX_EXSENDBUF_SIZE];
}

void CIOCP::GiocpDelete()
{
	delete[] ExSendBuf;
}

bool CIOCP::CreateGIocp(int server_port)
{
	g_ServerPort = server_port;

	InitializeCriticalSection(&criti);
	unsigned long ThreadID;
	g_IocpThreadHandle=CreateThread( NULL, 0, (LPTHREAD_START_ROUTINE)IocpServerWorkerEP, this, 0, &ThreadID);
	if ( g_IocpThreadHandle == 0 )
	{
		g_Log.Add("CreateThread() failed with error %d", GetLastError());
		return 0;
	}
	return 1;
}

void CIOCP::DestroyGIocp()
{
	closesocket(g_Listen);

	for (DWORD dwCPU=0; dwCPU < g_dwThreadCount;dwCPU++ )
	{
		TerminateThread( g_ThreadHandles[dwCPU] , 0);
	}

	TerminateThread(g_IocpThreadHandle, 0);

	if ( g_CompletionPort != NULL )
	{
		CloseHandle(g_CompletionPort);
		g_CompletionPort=NULL;
	}
}

bool CIOCP::CreateListenSocket()
{
	// create socket
	g_Listen = WSASocket(AF_INET, SOCK_STREAM, IPPROTO_IP, NULL, 0, WSA_FLAG_OVERLAPPED);
	if ( g_Listen == -1 )
	{
		g_Log.Add("WSASocket() failed with error %d", WSAGetLastError() );
		return 0;
	}

	// bind
	sockaddr_in InternetAddr;
	InternetAddr.sin_family=AF_INET;
	InternetAddr.sin_addr.S_un.S_addr=htonl(0);
	InternetAddr.sin_port=htons(g_ServerPort);
	int nRet = ::bind(g_Listen, (sockaddr*)&InternetAddr, 16);	
	if ( nRet == -1 )
	{
		g_Log.MsgBox("Bind error - server cannot be launched twice. Please terminate game server and restart.");
		SendMessage(ghWnd, WM_CLOSE, 0,0);	// Kill aplication
		return 0 ;
	}

	// listen
	nRet = listen(g_Listen, 5);
	if (nRet == -1)
	{
		g_Log.Add("listen() failed with error %d", WSAGetLastError());
		return 0;
	}
	return 1;
}

DWORD CIOCP::IocpServerWorker(void * p)
{
	SYSTEM_INFO SystemInfo;
	GetSystemInfo(&SystemInfo);
	g_dwThreadCount = SystemInfo.dwNumberOfProcessors * 2;
	if ( g_dwThreadCount > MAX_IO_THREAD_HANDLES )
	{
		g_dwThreadCount = MAX_IO_THREAD_HANDLES;
		g_Log.AddC(TColor::Orange, "[WARNING]: IOCP Thread Handles set to 16");
	}

	__try
	{
		g_CompletionPort=CreateIoCompletionPort(INVALID_HANDLE_VALUE, 0, 0, 0);

		if ( g_CompletionPort == NULL )
		{
			g_Log.Add("CreateIoCompletionPort failed with error: %d", GetLastError());
			__leave;
		}

		for (int n = 0; n<g_dwThreadCount; n++)
		{
			DWORD ThreadID;
			HANDLE hThread = CreateThread(NULL, 0, (LPTHREAD_START_ROUTINE)ServerWorkerThreadEP, this, 0, &ThreadID);
			if ( hThread == 0 )
			{
				g_Log.Add("CreateThread() failed with error %d", GetLastError() );
				__leave;
			}
			g_ThreadHandles[n] = hThread;
			CloseHandle(hThread);
		}

		if (CreateListenSocket() == false) {
			__leave;
		}

		while ( true )
		{
			sockaddr_in cAddr;
			int cAddrlen = 16;
			SOCKET fd = WSAAccept(g_Listen, (sockaddr*)&cAddr, &cAddrlen, NULL, 0 );
			if ( fd == -1 )
			{
				EnterCriticalSection(&criti);
				g_Log.Add("WSAAccept() failed with error %d", WSAGetLastError() );
				LeaveCriticalSection(&criti);
				continue;
			}

			EnterCriticalSection(&criti);

			in_addr cInAddr;
			memcpy(&cInAddr, &cAddr.sin_addr  , sizeof(cInAddr) );
			int ClientIndex = gObjAddSearch(fd, inet_ntoa(cInAddr) );
			if ( ClientIndex == -1 || !ObjectMaxRange(ClientIndex) )
			{
				g_Log.AddC(TColor::Red,  "error-L2 : ClientIndex = -1");
				closesocket(fd);
				LeaveCriticalSection(&criti);
				continue;
			}

			if (UpdateCompletionPort(fd, ClientIndex, 1) == 0 )
			{
				g_Log.AddC(TColor::Azure, "error-L1 : %d %d CreateIoCompletionPort failed with error %d", fd, ClientIndex, GetLastError() );
				closesocket(fd);
				LeaveCriticalSection(&criti);
				continue;
			}

			if (gObjAdd(fd, inet_ntoa(cInAddr), ClientIndex) == -1 )
			{
				g_Log.AddC(TColor::Azure, "error-L1 : %d %d gObjAdd() failed with error %d", fd, ClientIndex, GetLastError() );
				LeaveCriticalSection(&criti);
				closesocket(fd);
				continue;
			}
				
			memset(&gObj[ClientIndex].PerSocketContext->IOContext[0].Overlapped, 0, sizeof(WSAOVERLAPPED));
			memset(&gObj[ClientIndex].PerSocketContext->IOContext[1].Overlapped, 0, sizeof(WSAOVERLAPPED));

			gObj[ClientIndex].PerSocketContext->IOContext[0].wsabuf.buf = (char*)&gObj[ClientIndex].PerSocketContext->IOContext[0].Buffer;
			gObj[ClientIndex].PerSocketContext->IOContext[0].wsabuf.len = MAX_IO_BUFFER_SIZE;
			gObj[ClientIndex].PerSocketContext->IOContext[0].nTotalBytes = 0;
			gObj[ClientIndex].PerSocketContext->IOContext[0].nSentBytes = 0;
			gObj[ClientIndex].PerSocketContext->IOContext[0].nWaitIO = 0;
			gObj[ClientIndex].PerSocketContext->IOContext[0].nSecondOfs = 0;
			gObj[ClientIndex].PerSocketContext->IOContext[0].IOOperation = 0;
			gObj[ClientIndex].PerSocketContext->IOContext[1].wsabuf.buf = (char*)gObj[ClientIndex].PerSocketContext->IOContext[1].Buffer;
			gObj[ClientIndex].PerSocketContext->IOContext[1].wsabuf.len = MAX_IO_BUFFER_SIZE;
			gObj[ClientIndex].PerSocketContext->IOContext[1].nTotalBytes = 0;
			gObj[ClientIndex].PerSocketContext->IOContext[1].nSentBytes = 0;
			gObj[ClientIndex].PerSocketContext->IOContext[1].nWaitIO = 0;
			gObj[ClientIndex].PerSocketContext->IOContext[1].nSecondOfs = 0;
			gObj[ClientIndex].PerSocketContext->IOContext[1].IOOperation = 1;
			gObj[ClientIndex].PerSocketContext->m_socket = fd;
			gObj[ClientIndex].PerSocketContext->nIndex = ClientIndex;

			int RecvBytes;
			unsigned long Flags = 0;
			int nRet = WSARecv(fd, &gObj[ClientIndex].PerSocketContext->IOContext[0].wsabuf , 1, (unsigned long*)&RecvBytes, &Flags, 
								&gObj[ClientIndex].PerSocketContext->IOContext[0].Overlapped, NULL);
			if (nRet == SOCKET_ERROR)
			{
				int err = WSAGetLastError();
				if (err != WSA_IO_PENDING) {
					g_Log.AddC(TColor::Azure, "WSARecv() failed with error:%d [%d][%s]",
						err, ClientIndex, gObj[ClientIndex].m_PlayerData->Ip_addr);
					close_no_lock(ClientIndex, err);
					LeaveCriticalSection(&criti);
					continue;
				}
			}

			gObj[ClientIndex].PerSocketContext->IOContext[0].nWaitIO  = 1;
			gObj[ClientIndex].PerSocketContext->dwIOCount++;

			LeaveCriticalSection(&criti);
			GSProtocol.SCPJoinResultSend(ClientIndex, 1);
		}
	}
	__finally
	{
		
		if ( g_CompletionPort != NULL )
		{
			for ( int i = 0 ; i < g_dwThreadCount ; i++ )
			{
				PostQueuedCompletionStatus( g_CompletionPort , 0, 0, 0);
			}
		}

		if ( g_CompletionPort != NULL )
		{
			CloseHandle(g_CompletionPort);
			g_CompletionPort = NULL;
		}
		if ( g_Listen != INVALID_SOCKET )
		{
			closesocket( g_Listen);
			g_Listen = INVALID_SOCKET;
		}
	}

	return 1;
}

DWORD CIOCP::ServerWorkerThread()
{
	DWORD dwIoSize;
#ifdef _WIN64
	ULONGLONG ClientIndex;
#else
	ULONG ClientIndex;
#endif
	LPOVERLAPPED lpOverlapped = 0;
	
	while ( true )
	{
		BOOL bSuccess = GetQueuedCompletionStatus(g_CompletionPort, &dwIoSize, &ClientIndex, &lpOverlapped, -1); // WAIT_FOREVER
		EnterCriticalSection(&criti);
		// https://docs.microsoft.com/en-us/windows/win32/api/ioapiset/nf-ioapiset-getqueuedcompletionstatus#remarks
		if(bSuccess == FALSE)
		{
			// ERROR_NETNAME_DELETED 64
			// ERROR_CONNECTION_ABORTED 1236
			// ERROR_OPERATION_ABORTED 995
			// ERROR_SEM_TIMEOUT 121
			// ERROR_HOST_UNREACHABLE 1232
			int aError = GetLastError();
			if (lpOverlapped == NULL) {
				// If *lpOverlapped is NULL, the function did not dequeue
				// a completion packet from the completion port. In this
				// case, the function does not store information in the
				// variables pointed to by the lpNumberOfBytes and
				// lpCompletionKey parameters, and their values are indeterminate.
				g_Log.AddC(TColor::Orange, "GetQueuedCompletionStatus error:%d, needs end current thread",
					aError);
				LeaveCriticalSection(&criti);
				return 0;
			}
			else {
				// If *lpOverlapped is not NULL and the function dequeues a
				// completion packet for a failed I/O operation from the
				// completion port, the function stores information about
				// the failed operation in the variables pointed to by
				// lpNumberOfBytes, lpCompletionKey, and lpOverlapped.
				// To get extended error information, call GetLastError.
				_PER_IO_CONTEXT* lpIOContext = (_PER_IO_CONTEXT *)lpOverlapped;
				g_Log.AddC(TColor::Orange, "GetQueuedCompletionStatus error:%d [device:%d][op:%d]",
					aError, ClientIndex, lpIOContext->IOOperation);
				// https://stackoverflow.com/a/28770841/11508308
			}
		}

		if (dwIoSize == 0)
		{
			// handle remote close connection
			close_no_lock(ClientIndex, 0);
			LeaveCriticalSection(&criti);
			continue;
		}

		_PER_SOCKET_CONTEXT* lpPerSocketContext = gObj[ClientIndex].PerSocketContext;
		// if(lpPerSocketContext->m_socket == INVALID_SOCKET)
		// {
		// 	LeaveCriticalSection(&criti);
		// 	continue;
		// }

		// handle io count
		lpPerSocketContext->dwIOCount--;
		if(lpPerSocketContext->dwIOCount <= 0){
			lpPerSocketContext->dwIOCount = 0;
		}

		_PER_IO_CONTEXT* lpIOContext = (_PER_IO_CONTEXT *)lpOverlapped;
		if (lpIOContext->IOOperation == 1) // write
		{
			lpIOContext->nWaitIO = 0;
			lpIOContext->nSentBytes += dwIoSize;
			if(lpIOContext->nSentBytes < lpIOContext->nTotalBytes)
			{
				IoSendMore(lpPerSocketContext);
			}
			else if(lpIOContext->nSecondOfs > 0)
			{
				IoSendSecond(lpPerSocketContext);
			}
		}
		else if (lpIOContext->IOOperation == 0) // read
		{
			lpIOContext->nWaitIO = 0;
			lpIOContext->nSentBytes += dwIoSize;
			if (RecvDataParse(lpIOContext, lpPerSocketContext->nIndex ) == 0)
			{
				g_Log.Add("RecvDataParse error [%d][%s]",
					ClientIndex, gObj[ClientIndex].m_PlayerData->Ip_addr);
				close_no_lock(ClientIndex, 2);
				LeaveCriticalSection(&criti);
				continue;
			}

			unsigned long Flags = 0;
			memset(&lpIOContext->Overlapped, 0, sizeof (WSAOVERLAPPED));
			lpIOContext->wsabuf.len = MAX_IO_BUFFER_SIZE - lpIOContext->nSentBytes;
			lpIOContext->wsabuf.buf = (char*)&lpIOContext->Buffer[lpIOContext->nSentBytes];
			lpIOContext->IOOperation = 0;
			unsigned long RecvBytes = 0;
			int nRet = WSARecv(lpPerSocketContext->m_socket, &lpIOContext->wsabuf, 1, &RecvBytes, &Flags, &lpIOContext->Overlapped, NULL);
			if (nRet == SOCKET_ERROR)
			{
				int err = WSAGetLastError();
				// The error code WSA_IO_PENDING indicates that the overlapped
				// operation has been successfully initiated and that completion
				// will be indicated at a later time.
				if (err != WSA_IO_PENDING) {
					// Any other error code indicates that the overlapped operation
					// was not successfully initiated and no completion indication
					// will occur.
					// So we have to close client here.
					g_Log.Add("WSARecv() failed with error:%d [%d][%s]",
						err, ClientIndex, gObj[ClientIndex].m_PlayerData->Ip_addr);
					close_no_lock(ClientIndex, err);
					LeaveCriticalSection(&criti);
					continue;
				}
			}
			lpPerSocketContext->dwIOCount++;
			lpIOContext->nWaitIO = 1;
		}
		LeaveCriticalSection(&criti);
	}
	return 1;
}

#define USERIJ 
bool CIOCP::RecvDataParse(_PER_IO_CONTEXT * lpIOContext, int uIndex)	
{
	int lOfs = 0;
	int size = 0;
	BYTE headcode;
	BYTE xcode = 0;
	unsigned char* recvbuf = lpIOContext->Buffer;
	unsigned char byDec[9216];
	while ( true )
	{
		if ( recvbuf[lOfs] == 0xC1 || recvbuf[lOfs] == 0xC3 )
		{
			if (lpIOContext->nSentBytes < 3) {
				return true;
			}
			PBMSG_HEAD* lphead = (PBMSG_HEAD*)(recvbuf + lOfs);
			size = lphead->size;
			headcode = lphead->headcode;
			xcode = lphead->c;
		}
		else if ( recvbuf[lOfs] == 0xC2 || recvbuf[lOfs] == 0xC4 )
		{
			if (lpIOContext->nSentBytes < 4) {
				return true;
			}
			PWMSG_HEAD* lphead = (PWMSG_HEAD*)(recvbuf + lOfs);
			size = lphead->sizeH<<8 | lphead->sizeL;
			headcode = lphead->headcode;
			xcode = lphead->c;
		}
		else
		{
			g_Log.Add("Header error lOfs:%d size:%d", lOfs, lpIOContext->nSentBytes);
			return false;
		}

		// validate frame
		if (size <= 0)
		{
			g_Log.Add("error-L1: size: %d [%02X][%02X][%02X]", size, recvbuf[lOfs], recvbuf[lOfs+1], recvbuf[lOfs+2], recvbuf[lOfs+3]);
			return false;
		}

		if (lpIOContext->nSentBytes >= size)
		{
			if ( xcode == 0xC3 )
			{
				int ret = g_PacketEncrypt.Decrypt(&byDec[2], &recvbuf[lOfs + 2], size - 2);
				if ( ret < 0 )
				{
					g_Log.AddC(TColor::Red, "[%s][Packet-Decrypt BYTE] Error: ret < 0 %x/%x/%x)", gObj[uIndex].m_PlayerData->Ip_addr,recvbuf[lOfs], recvbuf[lOfs+1], recvbuf[lOfs+2]);
					return false;
				}
				else
				{
					unsigned char* pDecBuf = byDec;
					headcode = pDecBuf[2];
					byDec[0] = 0xC1;
					byDec[1] = ret + 2;
					gObj[uIndex].m_PlayerData->PacketsPerSecond++;
					if(gObj[uIndex].m_PlayerData->PacketsPerSecond >= g_ConfigRead.data.common.PacketLimit)
					{
						g_Log.AddC(TColor::Red, "[ANTI-HACK] Packets Per Second: %d / %d", gObj[uIndex].m_PlayerData->PacketsPerSecond, g_ConfigRead.data.common.PacketLimit);
						return false;
					}

					CStreamPacketEngine_Server PacketStream;
					PacketStream.Clear();
					if (PacketStream.AddData(byDec, ret+2) == 0)
					{
						g_Log.AddC(TColor::Red,  "error-L1 : CStreamPacketEngine Adding Error : ip = %s account:%s name:%s HEAD:%x (%s,%d) State:%d",
							gObj[uIndex].m_PlayerData->Ip_addr,
							gObj[uIndex].AccountID, gObj[uIndex].Name,
							headcode,
							__FILE__, __LINE__,
							gObj[uIndex].Connected);
						return false;
					}
					if ( PacketStream.ExtractPacket(byDec) != 0 )
					{
						g_Log.AddC(TColor::Red,  "error-L1 : CStreamPacketEngine ExtractPacket Error : ip = %s account:%s name:%s HEAD:%x (%s,%d) State:%d",
							gObj[uIndex].m_PlayerData->Ip_addr, gObj[uIndex].AccountID, gObj[uIndex].Name, headcode, __FILE__, __LINE__, gObj[uIndex].Connected);
						return false;
					}
					GSProtocol.ProtocolCore(headcode, byDec, ret, uIndex, 1);
				}
			}
			else if ( xcode == 0xC4 )
			{
				int ret = g_PacketEncrypt.Decrypt(&byDec[3], &recvbuf[lOfs + 3], size - 3);
				if ( ret < 0 )
				{
					g_Log.AddC(TColor::Red, "[Packet-Decrypt WORD] Error: ret < 0 %x/%x/%x)", recvbuf[lOfs], recvbuf[lOfs+1], recvbuf[lOfs+2]);
					return false;
				}
				else
				{
					unsigned char* pDecBuf = byDec;
					headcode = pDecBuf[3];
					byDec[0] = 0xC2;
					WORD size = (ret &0xFFFF)+3;
					byDec[1] = SET_NUMBERH(size);
					byDec[2] = SET_NUMBERL(size);

					gObj[uIndex].m_PlayerData->PacketsPerSecond++;
					if(gObj[uIndex].m_PlayerData->PacketsPerSecond >= g_ConfigRead.data.common.PacketLimit)
					{
						g_Log.AddC(TColor::Red, "[ANTI-HACK] Packets Per Second: %d / %d", gObj[uIndex].m_PlayerData->PacketsPerSecond, g_ConfigRead.data.common.PacketLimit);
						return false;
					}

					CStreamPacketEngine_Server PacketStream;
					PacketStream.Clear();
					if ( PacketStream.AddData(byDec, ret+3) == 0 )
					{
						g_Log.AddC(TColor::Red,  "error-L1 : CStreamPacketEngine Adding Error : ip = %s account:%s name:%s HEAD:%x (%s,%d) State:%d",
							gObj[uIndex].m_PlayerData->Ip_addr, gObj[uIndex].AccountID, gObj[uIndex].Name, headcode, __FILE__, __LINE__, gObj[uIndex].Connected);
						return false;
					}
					if ( PacketStream.ExtractPacket(byDec) != 0 )
					{
						g_Log.AddC(TColor::Red,  "error-L1 : CStreamPacketEngine ExtractPacket Error : ip = %s account:%s name:%s HEAD:%x (%s,%d) State:%d",
							gObj[uIndex].m_PlayerData->Ip_addr, gObj[uIndex].AccountID, gObj[uIndex].Name, headcode, __FILE__, __LINE__, gObj[uIndex].Connected);
						return false;
					}
					GSProtocol.ProtocolCore(headcode, byDec, ret, uIndex, 1);
				}
			}
			else
			{
				CStreamPacketEngine_Server ps;
				ps.Clear();

				if ( ps.AddData(&recvbuf[lOfs], size) == 0 )
				{
					g_Log.AddC(TColor::Red,  "error-L1 : CStreamPacketEngine Adding Error : ip = %s account:%s name:%s HEAD:%x (%s,%d) State:%d",
						gObj[uIndex].m_PlayerData->Ip_addr, gObj[uIndex].AccountID, gObj[uIndex].Name, headcode, __FILE__, __LINE__, gObj[uIndex].Connected);
					return false;
				}

				if ( ps.ExtractPacket(byDec) != 0 )
				{
					g_Log.AddC(TColor::Red,  "error-L1 : CStreamPacketEngine ExtractPacket Error : ip = %s account:%s name:%s HEAD:%x (%s,%d) State:%d",
						gObj[uIndex].m_PlayerData->Ip_addr, gObj[uIndex].AccountID, gObj[uIndex].Name, headcode, __FILE__, __LINE__, gObj[uIndex].Connected);
					return false;
				}

				gObj[uIndex].m_PlayerData->PacketsPerSecond++;
				if(gObj[uIndex].m_PlayerData->PacketsPerSecond >= g_ConfigRead.data.common.PacketLimit)
				{
					g_Log.AddC(TColor::Red, "[ANTI-HACK] Packets Per Second: %d / %d", gObj[uIndex].m_PlayerData->PacketsPerSecond, g_ConfigRead.data.common.PacketLimit);
					return false;
				}
				GSProtocol.ProtocolCore(headcode, byDec, size, uIndex, 0); // here
			}

			lOfs += size; // wait
			lpIOContext->nSentBytes -= size;
			if ( lpIOContext->nSentBytes <= 0 )
			{
				lpIOContext->nSentBytes = 0;
				break;
			}
		}
		else if ( lOfs > 0 )
		{
			// the n frame needs more bytes
			if ( lpIOContext->nSentBytes < 1 )
			{
				g_Log.Add("error-L1 : recvbuflen 1 %s %d", __FILE__, __LINE__);
				break;
			}

			if ( lpIOContext->nSentBytes < MAX_IO_BUFFER_SIZE ) 
			{
				memcpy(recvbuf, &recvbuf[lOfs], lpIOContext->nSentBytes);
				g_Log.Add("Message copy %d", lpIOContext->nSentBytes);
			}
			break;
		}
		else
		{
			// the first frame needs more bytes
			break;
		}
	}

	return true;
}

bool CIOCP::DataSend(int aIndex, unsigned char* lpMsg, DWORD dwSize, bool Encrypt)
{
	unsigned char * SendBuf;
	if (aIndex < g_ConfigRead.server.GetObjectStartUserIndex())
	{
		return true;
	}

	EnterCriticalSection(&criti);
	if (ObjectMaxRange(aIndex) == FALSE)
	{
		g_Log.Add("error-L2 : Index(%d) %x %x %x ", dwSize, lpMsg[0], lpMsg[1], lpMsg[2]);
		LeaveCriticalSection(&criti);
		return false;
	}

	_PER_SOCKET_CONTEXT* lpPerSocketContext = gObj[aIndex].PerSocketContext;
	if(lpPerSocketContext->m_socket == INVALID_SOCKET)
	{
		LeaveCriticalSection(&criti);
		return false;
	}

	if (gStalkProtocol)
	{
		if (gStalkProtocolId[0] == gObj[aIndex].AccountID[0])
		{
			if (gStalkProtocolId[1] == gObj[aIndex].AccountID[1])
			{
				if (!strcmp(gStalkProtocolId, gObj[aIndex].AccountID))
				{
					g_Log.AddHeadHex("DATA SEND", lpMsg, dwSize);
				}
			}
		}
	}

	if (!(lpMsg[0] == 0xC1 || lpMsg[0] == 0xC2 || lpMsg[0] == 0xC3 || lpMsg[0] == 0xC4))
	{
		g_Log.AddC(TColor::Red, "[ERROR] Trying to send packet without HEADER (%s)(%s)", gObj[aIndex].AccountID, gObj[aIndex].Name);
		LeaveCriticalSection(&criti);
		return FALSE;
	}

#ifdef EMU_NOCRYPT
	if (lpMsg[0] == 0xC3)
		lpMsg[0] = 0xC1;
	if (lpMsg[0] == 0xC4)
		lpMsg[0] = 0xC2;
#endif

#ifdef C3C4_DISABLECRYPT
	if (lpMsg[0] == 0xC3)
		lpMsg[0] = 0xC1;
	if (lpMsg[0] == 0xC4)
		lpMsg[0] = 0xC2;
#endif

	if (lpMsg[0] == 0xC3 || lpMsg[0] == 0xC4)
	{
		int ret;
		BYTE btsize;

		if (lpMsg[0] == 0xC3)
		{
			btsize = lpMsg[1];
			ret = g_PacketEncrypt.Encrypt(&ExSendBuf[2], &lpMsg[1], dwSize - 1);
			ExSendBuf[0] = 0xC3;
			ExSendBuf[1] = ret + 2;
			SendBuf = ExSendBuf;
			dwSize = ret + 2;
			lpMsg[1] = btsize;
		}
		else
		{
			btsize = lpMsg[2];
			ret = g_PacketEncrypt.Encrypt(&ExSendBuf[3], &lpMsg[2], dwSize - 2);
			ExSendBuf[0] = 0xC4;
			ExSendBuf[1] = SET_NUMBERH(ret + 3);
			ExSendBuf[2] = SET_NUMBERL(ret + 3);
			SendBuf = ExSendBuf;
			dwSize = ret + 3;
		}
	}
	else
	{
		SendBuf = lpMsg;
	}

	if ( gObj[aIndex].Connected < PLAYER_CONNECTED )
	{
		LeaveCriticalSection(&criti);
		return false;
	}

	if ( dwSize > sizeof(lpPerSocketContext->IOContext[0].Buffer))
	{
		g_Log.Add("Error : Max msg(%d) %s %d", dwSize, __FILE__, __LINE__);
		shutdown_no_lock(aIndex);
		LeaveCriticalSection(&criti);
		return false;
	}

	_PER_IO_CONTEXT* lpIoCtxt = &lpPerSocketContext->IOContext[1];
	if (lpIoCtxt->nWaitIO > 0)
	{
		if ((lpIoCtxt->nSecondOfs + dwSize ) > MAX_IO_BUFFER_SIZE-1)
		{
			g_Log.Add("the second send buffer is full %d + %d [%d][%s]",
				lpIoCtxt->nSecondOfs, dwSize, aIndex, gObj[aIndex].AccountID);
			// in this case, we use closesocket instead of shutdown for the network
			// current is unreliable, such as underlying tcp re-transmission too
			// many times
			close_no_lock(aIndex, 1);
			LeaveCriticalSection(&criti);
			return false;
		}

		memcpy( &lpIoCtxt->BufferSecond[lpIoCtxt->nSecondOfs], SendBuf, dwSize);
		lpIoCtxt->nSecondOfs += dwSize;
		LeaveCriticalSection(&criti);
		return true;
	}

	lpIoCtxt->nTotalBytes = 0;
	
	if ( lpIoCtxt->nSecondOfs > 0 )
	{
		memcpy(lpIoCtxt->Buffer, lpIoCtxt->BufferSecond, lpIoCtxt->nSecondOfs);
		lpIoCtxt->nTotalBytes = lpIoCtxt->nSecondOfs;
		lpIoCtxt->nSecondOfs = 0;
	}

	if ((lpIoCtxt->nTotalBytes+dwSize) > MAX_IO_BUFFER_SIZE-1)
	{
		g_Log.Add("the first send buffer is full %d + %d [%d][%s]",
			lpIoCtxt->nTotalBytes, dwSize, aIndex, gObj[aIndex].AccountID);
		shutdown_no_lock(aIndex);
		LeaveCriticalSection(&criti);
		return false;
	}
	memcpy(&lpIoCtxt->Buffer[lpIoCtxt->nTotalBytes], SendBuf, dwSize);
	lpIoCtxt->nTotalBytes += dwSize;
	lpIoCtxt->wsabuf.buf = (char*)&lpIoCtxt->Buffer;
	lpIoCtxt->wsabuf.len = lpIoCtxt->nTotalBytes;
	lpIoCtxt->nSentBytes = 0;
	lpIoCtxt->IOOperation = 1;
	
	unsigned long SendBytes;
	int nRet = WSASend(lpPerSocketContext->m_socket, &lpIoCtxt->wsabuf , 1, &SendBytes, 0, &lpIoCtxt->Overlapped, NULL);
	if (nRet == SOCKET_ERROR) {
		int err = WSAGetLastError();
		// The error code WSA_IO_PENDING indicates that the overlapped
		// operation has been successfully initiated and that completion
		// will be indicated at a later time.
		if (err != WSA_IO_PENDING) {
			// Any other error code indicates that the overlapped operation
			// was not successfully initiated and no completion indication
			// will occur.
			g_Log.AddC(TColor::Red, "WSASend [line:%d] error:%d [%d][%s] [%02X][%02X][%02X]", __LINE__,
				err, aIndex, gObj[aIndex].m_PlayerData->Ip_addr,
				(BYTE)lpIoCtxt->wsabuf.buf[0], (BYTE)lpIoCtxt->wsabuf.buf[1], (BYTE)lpIoCtxt->wsabuf.buf[2]);
			LeaveCriticalSection(&criti);
			return false;
		}
	}
	lpPerSocketContext->dwIOCount++;
	lpIoCtxt->nWaitIO = 1;
	LeaveCriticalSection(&criti);
	return true;
}

bool CIOCP::IoSendSecond(_PER_SOCKET_CONTEXT * lpPerSocketContext)
{
	int aIndex = lpPerSocketContext->nIndex;
	_PER_IO_CONTEXT * lpIoCtxt = &lpPerSocketContext->IOContext[1];
	memcpy(lpIoCtxt->Buffer, lpIoCtxt->BufferSecond, lpIoCtxt->nSecondOfs);
	lpIoCtxt->nTotalBytes = lpIoCtxt->nSecondOfs;
	lpIoCtxt->nSecondOfs = 0;
	lpIoCtxt->wsabuf.buf = (char*)&lpIoCtxt->Buffer;
	lpIoCtxt->wsabuf.len = lpIoCtxt->nTotalBytes;
	lpIoCtxt->nSentBytes = 0;
	lpIoCtxt->IOOperation = 1;

	unsigned long SendBytes;
	int nRet = WSASend(lpPerSocketContext->m_socket, &lpIoCtxt->wsabuf, 1, &SendBytes, 0, &lpIoCtxt->Overlapped, NULL);
	if (nRet == SOCKET_ERROR) {
		int err = WSAGetLastError();
		// The error code WSA_IO_PENDING indicates that the overlapped
		// operation has been successfully initiated and that completion
		// will be indicated at a later time.
		if (err != WSA_IO_PENDING) {
			// Any other error code indicates that the overlapped operation
			// was not successfully initiated and no completion indication
			// will occur.
			g_Log.AddC(TColor::Red, "WSASend [line:%d] error:%d [%d][%s]", __LINE__,
				err, aIndex, gObj[aIndex].m_PlayerData->Ip_addr);
			return false;
		}
	}
	lpPerSocketContext->dwIOCount++;
	lpIoCtxt->nWaitIO = 1;
	return true;
}

bool CIOCP::IoSendMore(_PER_SOCKET_CONTEXT * lpPerSocketContext)
{
	int aIndex = lpPerSocketContext->nIndex;
	_PER_IO_CONTEXT* lpIoCtxt = &lpPerSocketContext->IOContext[1];
	lpIoCtxt->wsabuf.buf = (char*)&lpIoCtxt->Buffer[lpIoCtxt->nSentBytes];
	lpIoCtxt->wsabuf.len = lpIoCtxt->nTotalBytes - lpIoCtxt->nSentBytes;
	lpIoCtxt->IOOperation = 1;

	unsigned long SendBytes;
	int nRet = WSASend(lpPerSocketContext->m_socket, &lpIoCtxt->wsabuf, 1, &SendBytes, 0, &lpIoCtxt->Overlapped, NULL);
	if (nRet == SOCKET_ERROR) {
		int err = WSAGetLastError();
		// The error code WSA_IO_PENDING indicates that the overlapped
		// operation has been successfully initiated and that completion
		// will be indicated at a later time.
		if (err != WSA_IO_PENDING) {
			// Any other error code indicates that the overlapped operation
			// was not successfully initiated and no completion indication
			// will occur.
			g_Log.AddC(TColor::Red, "WSASend [line:%d] error:%d [%d][%s]", __LINE__,
				err, aIndex, gObj[aIndex].m_PlayerData->Ip_addr);
			return false;
		}
	}
	lpPerSocketContext->dwIOCount++;
	lpIoCtxt->nWaitIO = 1;
	return true;
}

bool CIOCP::UpdateCompletionPort(SOCKET sd, int ClientIndex, BOOL bAddToList)
{
	HANDLE cp = CreateIoCompletionPort((HANDLE) sd, g_CompletionPort, ClientIndex, 0);
	if ( cp == 0 )
	{
		g_Log.Add("CreateIoCompletionPort: %d", GetLastError() );
		return FALSE;
	}

	gObj[ClientIndex].PerSocketContext->dwIOCount = 0;
	return TRUE;
}

// no lock
void CIOCP::close_no_lock(int index, int err)
{
	if (index < g_ConfigRead.server.GetObjectStartUserIndex()
	|| index >= g_ConfigRead.server.GetObjectMax()) {
		return;
	}

	_PER_SOCKET_CONTEXT* lpPerSocketContext = gObj[index].PerSocketContext;
	SOCKET fd = lpPerSocketContext->m_socket;
	if (fd != INVALID_SOCKET)
	{
		g_Log.AddC(TColor::Orange, "connection closed [err:%d] [%d][%s]",
			err, index, gObj[index].m_PlayerData->Ip_addr);
		if (err != 0) {
			// set linger
			LINGER l;
			l.l_onoff = 1;
			l.l_linger = 0;
			if (SOCKET_ERROR == setsockopt(fd, SOL_SOCKET, SO_LINGER, (char*)&l, sizeof(LINGER))) {
				g_Log.Add("setsockopt failed with error %d", WSAGetLastError() );
			}
		}
		if (closesocket(fd) == SOCKET_ERROR)
		{
			g_Log.AddC(TColor::Red, "close socket error:%d [%d][%s]",
				WSAGetLastError(), index, gObj[index].m_PlayerData->Ip_addr);
		}
		lpPerSocketContext->m_socket = INVALID_SOCKET;
	}
	gObjDel(index);
}

void CIOCP::shutdown_no_lock(int index) // go to game
{
	if (index < g_ConfigRead.server.GetObjectStartUserIndex()
	|| index >= g_ConfigRead.server.GetObjectMax()) {
		return;
	}

	_PER_SOCKET_CONTEXT* lpPerSocketContext = gObj[index].PerSocketContext;
	if (lpPerSocketContext->m_socket != INVALID_SOCKET)
	{
		// just shutdown
		g_Log.AddC(TColor::Orange, "shutdown connection [%d][%s]",
			index, gObj[index].m_PlayerData->Ip_addr);
		if (shutdown(lpPerSocketContext->m_socket, SD_SEND) == SOCKET_ERROR) {
			g_Log.AddC(TColor::Red, "shutdown connection error:%d [%d][%s]",
				WSAGetLastError(), index, gObj[index].m_PlayerData->Ip_addr);
		}
	}
}

void CIOCP::CloseClient(int index) // go to game
{
	EnterCriticalSection(&criti);
	shutdown_no_lock(index);
	LeaveCriticalSection(&criti);
}

void CIOCP::CloseClientHard(int index) // go to game
{
	EnterCriticalSection(&criti);
	close_no_lock(index, 3);
	LeaveCriticalSection(&criti);
}