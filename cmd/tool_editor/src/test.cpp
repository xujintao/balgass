#include <assert.h>
#include <string.h>
#include <stdio.h>
#include "cipher_random.hpp"
#include <boost/archive/text_oarchive.hpp>
#include <iostream>

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
				dst[(i << 1) + j] = b + 'A' - 10;
			}
		}
	}
	return len << 1;
}

struct text {
	const char* name;
	const char* plain;
	const char* cipher;
};

text texts[] = {
	{
		"tea block cipher, keylen:16, blocksize:8",
		"0000456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF",
		"00000123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF"
		"C21C16DC48CBC1633056D3BDBF3B75D53056D3BDBF3B75D53056D3BDBF3B75D5"
	},
	{
		"3way block cipher, keylen:12, blocksize:12",
		"0100456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF",
		"01000123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF"
		"269E7585A21DF0A653292DFCF1BAFF797E0E2B4AEDA3ACDC0123456789ABCDEF"
	},
	{
		"cast128 block cipher, keylen:16, blocksize:8",
		"0200456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF",
		"02000123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF"
		"376D3EF44A723E6AE83B50A03AF412D7E83B50A03AF412D7E83B50A03AF412D7"
	},
	{
		"rc5 block cipher, keylen:16, blocksize:8",
		"0300456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF",
		"03000123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF"
		"EF15087701B4E0E5CC96815A45E5D31FCC96815A45E5D31FCC96815A45E5D31F"
	},
	{
		"rc6 block cipher, keylen:16, blocksize:16",
		"0400456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF",
		"04000123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF"
		"0C41D531208044C42B1624BC887E721B9A8C95E925D522D436E87CA2A2F1B7F3"
	},
	{
		"mars block cipher, keylen:16, blocksize:16",
		"0500567890ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF",
		"05000123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF"
		"4D0B73B7F416E219D2C133C7F244BBF80DC015315BC73DEEB115659678D53374"
	},
	{
		"idea block cipher, keylen:16, blocksize:8",
		"0600567890ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF",
		"06000123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF"
		"10934DDD3D81F0AC22378E5924B8507522378E5924B8507522378E5924B85075"
	},
	{
		"gost block cipher, keylen:32, blocksize:8",
		"0700567890ABCDEF1234567890ABCDEF1234567890ABCDEF1234567890ABCDEF",
		"07000123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF"
		"B7CBA48C21681C455CEB997B12DE2BB75CEB997B12DE2BB75CEB997B12DE2BB7"
	}
};

bool test_random_encrypt() {
	for (auto v : texts) {
		unsigned char* streamPlain = nullptr;
		unsigned char* streamCipher = nullptr;
		char* plain = nullptr;
		char* cipher = nullptr;

		__try {
			size_t lenPlain = hex2bytes(0, v.plain);
			if (lenPlain == 0) {
				printf("failed [%s].hex2bytes, invalid plain string: %s\r\n", v.name, v.plain);
				return false;
			}
			streamPlain = new unsigned char[lenPlain];
			hex2bytes(streamPlain, v.plain);

			size_t lenCipher = RandomCipher::encrypt(nullptr, streamPlain, lenPlain, false);
			if (lenCipher == 0) {
				printf("failed [%s].bmdenc, invalid plain stream: %s\r\n", v.name, v.plain);
				return false;
			}
			streamCipher = new unsigned char[lenCipher];
			RandomCipher::encrypt(streamCipher, streamPlain, lenPlain, false);

			cipher = new char[(lenCipher << 1) + 1];
			cipher[lenCipher << 1] = 0;
			bytes2hex(cipher, streamCipher, lenCipher);
			if (strcmp(cipher, v.cipher) != 0) {
				printf("failed [%s].bmdenc, got[%s] expect[%s]\r\n", v.name, cipher, v.cipher);
				return false;
			}
			printf("success [%s]\r\n", v.name);
		}
		__finally {
			if (streamPlain != nullptr) {
				delete[]streamPlain;
			}
			if (streamCipher != nullptr) {
				delete[]streamCipher;
			}
			if (plain != nullptr) {
				delete[]plain;
			}
			if (cipher != nullptr) {
				delete[]cipher;
			}
		}
	}
	return true;
}

bool test_random_decrypt() {
	for (auto v : texts) {
		unsigned char* streamPlain = nullptr;
		unsigned char* streamCipher = nullptr;
		char* plain = nullptr;
		char* cipher = nullptr;

		__try {
			size_t lenCipher = hex2bytes(0, v.cipher);
			if (lenCipher == 0) {
				printf("failed [%s].hex2bytes, invalid cipher string: %s\r\n", v.name, v.cipher);
				return false;
			}
			streamCipher = new unsigned char[lenCipher];
			hex2bytes(streamCipher, v.cipher);

			size_t lenPlain = RandomCipher::decrypt(nullptr, streamCipher, lenCipher);
			if (lenPlain == 0) {
				printf("failed [%s].bmddec, invalid plain stream: %s\r\n", v.name, v.plain);
				return false;
			}
			streamPlain = new unsigned char[lenPlain];
			RandomCipher::decrypt(streamPlain, streamCipher, lenCipher);

			plain = new char[(lenPlain << 1) + 1];
			plain[lenPlain << 1] = 0;
			bytes2hex(plain, streamPlain, lenPlain);
			if (strcmp(plain, v.plain) != 0) {
				printf("failed [%s].bmdenc, got[%s] expect[%s]\r\n", v.name, plain, v.plain);
				return false;
			}
			printf("success [%s]\r\n", v.name);
		}
		__finally {
			if (streamPlain != nullptr) {
				delete[]streamPlain;
			}
			if (streamCipher != nullptr) {
				delete[]streamCipher;
			}
			if (plain != nullptr) {
				delete[]plain;
			}
			if (cipher != nullptr) {
				delete[]cipher;
			}
		}
	}
	return true;
}

bool test_wtf_decrypt() {
	boost::archive::text_oarchive oa(std::cout);
	int i = 1;
	oa << i;
	return true;
}

typedef bool (*func)();
typedef struct {
	const char* name;
	func f;
}test_item;
test_item items[] = {
	{ "random_encrypt", test_random_encrypt },
	{ "random_decrypt", test_random_decrypt },
	{ "wtf_decrypt", test_wtf_decrypt },
};

void main() {
	for (auto i : items) {
		printf("test %s begin\r\n", i.name);
		if (i.f() != true) {
			printf("test %s failed\r\n\r\n", i.name);
			break;
		}
		printf("test %s ok\r\n\r\n", i.name);
	}
}
