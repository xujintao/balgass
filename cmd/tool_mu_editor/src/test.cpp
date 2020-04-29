#include <string.h>
#include <stdio.h>
#include "bmdCipher.hpp"

size_t hex2bytes(unsigned char* dst, const char* src) {
	if (src == nullptr) {
		return 0;
	}
	size_t len = strlen(src);
	if (len == 0 || len % 2 == 1) {
		return 0;
	}
	if (dst == nullptr) {
		return len >> 1;
	}
	unsigned char byteH, byteL;
	for (int i = 0; i < len; i++)
	{
		unsigned char b = src[i];
		if (b >= '0' && b <= '9')
			b -= '0';
		else if (b >= 'A' && b <= 'F')
			b = b - 'A' + 10;
		else if (b >= 'a' && b <= 'f')
			b = b - 'a' + 10;
		else
			b = 0xF;
		if (i % 2 == 0) {
			byteH = b;
		}
		else {
			byteL = b;
			dst[i>>1] = (byteH << 4) | byteL;
		}
	}
	return len >> 1;
}

size_t bytes2hex(char* dst, const unsigned char* src, size_t len) {
	if (src == nullptr || len == 0) {
		return 0;
	}
	if (dst == nullptr) {
		return len << 1;
	}
	for (int i = 0; i < len; i++)
	{
		unsigned char byteH = src[i] >> 4 & 0x0F;
		unsigned char byteL = src[i] & 0x0F;
		unsigned char bytes[2] = { byteH,byteL };
		for (int j = 0;j < 2; j++) {
			unsigned char b = bytes[j];
			if (b >= 0 && b <= 9)
				dst[(i<<1) + j] = b + '0';
			else if (b >= 10 && b <= 15) {
				dst[(i << 1) + j] = b + 'A';
			}
		}
	}
	return len << 1;
}

struct text {
	const char* plain;
	const char* cipher;
};

text texts[] = {
	{
		"12AFab",
		""
	},
	{
		"",
		""
	}
};

void main() {
	// encrypt
	for (auto v : texts) {
		unsigned char* streamPlain = nullptr;
		unsigned char* streamCipher = nullptr;
		char* cipher = nullptr;

		__try {
			size_t lenPlain = hex2bytes(0, v.plain);
			if (lenPlain == 0) {
				printf("hex2bytes failed, invalid plain string: %s\r\n", v.plain);
				return;
			}
			streamPlain = new unsigned char[lenPlain];
			hex2bytes(streamPlain, v.plain);

			size_t lenCipher = bmdenc(nullptr, streamPlain, lenPlain);
			if (lenCipher == 0) {
				printf("bmdenc failed, invalid plain stream: %s\r\n", v.plain);
				return;
			}
			streamCipher = new unsigned char[lenCipher];
			bmdenc(streamCipher, streamPlain, lenPlain);
			
			cipher = new char[(lenCipher << 1) + 1];
			cipher[lenCipher << 1] = 0;
			bytes2hex(cipher, streamCipher, lenCipher);
			if (strcmp(cipher, v.cipher) != 0) {
				printf("bmdenc failed, invalid plain stream: %s\r\n", v.plain);
				return;
			}
		}
		__finally {
			if (streamPlain != nullptr) {
				delete[]streamPlain;
			}
			if (streamCipher != nullptr) {
				delete[]streamCipher;
			}
			if (cipher != nullptr) {
				delete[]cipher;
			}
		}
	}
}
