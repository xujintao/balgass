#include "stdafx.h"
#include "xor.h"
CXor* Xor = new CXor;

static BYTE key[32] = {
	0xAB, 0x11, 0xCD, 0xFE,
	0x18, 0x23, 0xC5, 0xA3,
	0xCA, 0x33, 0xC1, 0xCC,
	0x66, 0x67, 0x21, 0xF3,
	0x32, 0x12, 0x15, 0x35,
	0x29, 0xFF, 0xFE, 0x1D,
	0x44, 0xEF, 0xCD, 0x41,
	0x26, 0x3C, 0x4E, 0x4D
};

void CXor::Enc(unsigned char* buff, int headerSize, int len) {
	for (int i = 0; i < len; i++) {
		buff[i] ^= key[(i + headerSize) & 31] ^ buff[i - 1];
	}
}

void CXor::Dec(unsigned char* buff, int headerSize, int len) {
	for (int i = len - 1; i >= 0; i--) {
		buff[i] ^= key[(i + headerSize) & 31] ^ buff[i - 1];
	}
}

void CXor::BuxConvert(char* buf, int size)
{
	static unsigned char bBuxCode[3] = { 0xFC, 0xCF, 0xAB };
	for (int n = 0; n<size; n++)
	{
		buf[n] ^= bBuxCode[n % 3];		// Nice trick from WebZen
	}
}