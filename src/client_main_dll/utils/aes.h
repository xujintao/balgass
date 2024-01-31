#pragma once
class CAes
{
public:
	static int Encrypt(BYTE * lpDest, BYTE * lpSource, int iSize);
	static int Decrypt(BYTE * lpDest, BYTE * lpSource, int iSize);
};
