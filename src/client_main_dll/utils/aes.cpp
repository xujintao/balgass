#include "stdafx.h"
#include "crypto.h"
#include "aes.h"

int CAes::Encrypt(BYTE * lpDest, BYTE * lpSource, int iSize) {
	return g_PacketEncrypt.Encrypt(lpDest, lpSource, iSize);
}

int CAes::Decrypt(BYTE * lpDest, BYTE * lpSource, int iSize) {
	return g_PacketEncrypt.Decrypt(lpDest, lpSource, iSize);
}