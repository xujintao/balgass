#pragma once
#include "tea.h"  // 0
#include "3way.h" // 1
#include "cast.h" // 2
#include "rc5.h"  // 3
#include "rc6.h"  // 4
#include "mars.h" // 5
#include "idea.h" // 6
#include "gost.h" // 7

#pragma comment(lib, "cryptlib.lib")

class blockCipherInterface {
public:
	virtual void setkey(const unsigned char* key, size_t len) = 0;
	virtual int encrypt(const unsigned char* in, size_t len, unsigned char* out) = 0;
	virtual int decrypt(const unsigned char* in, size_t len, unsigned char* out) = 0;
};

template <class E, class D>
class blockCipher : public blockCipherInterface {
public:
	void setkey(const unsigned char* key, size_t len) {
		enc.SetKey(key, 0x10);
		dec.SetKey(key, 0x10);
	}
	int encrypt( const unsigned char* in, size_t len, unsigned char* out) {
		if (out == nullptr || len == 0) {
			return len;
		}
		len /= 8;
		for (;len > 0;len--) {
			enc.ProcessAndXorBlock(in, nullptr, out);
			in += 8;
			out += 8;
		}
		return len;
	}
	int decrypt(const unsigned char* in, size_t len, unsigned char* out) {
		if (out == nullptr || len == 0) {
			return len;
		}
		len /= 8;
		for (;len > 0;len--) {
			dec.ProcessAndXorBlock(in, nullptr, out);
			in += 8;
			out += 8;
		}
		return len;
	}
private:
	E enc;
	D dec;
};

class bmdCipher
{
public:
	bmdCipher() {
		bc = nullptr;
	}
	~bmdCipher() {
		if (bc != nullptr) {
			delete bc;
		}
	}
	void setkey(int alg, const unsigned char* key, size_t len) {
		switch (alg & 7)
		{
		case 0: bc = new(blockCipher<CryptoPP::TEAEncryption, CryptoPP::TEADecryption>);break;
		case 1: bc = new(blockCipher<CryptoPP::ThreeWayEncryption, CryptoPP::ThreeWayDecryption>);break;
		case 2: bc = new(blockCipher<CryptoPP::CAST128Encryption, CryptoPP::CAST128Decryption>);break;
		case 3: bc = new(blockCipher<CryptoPP::RC5Encryption, CryptoPP::RC5Decryption>);break;
		case 4: bc = new(blockCipher<CryptoPP::RC6Encryption, CryptoPP::RC6Decryption>);break;
		case 5: bc = new(blockCipher<CryptoPP::MARSEncryption, CryptoPP::MARSDecryption>);break;
		case 6: bc = new(blockCipher<CryptoPP::IDEAEncryption, CryptoPP::IDEADecryption>);break;
		case 7: bc = new(blockCipher<CryptoPP::GOSTEncryption, CryptoPP::GOSTDecryption>);break;
		}
		bc->setkey(key, len);
	}
	int encrypt(unsigned char* out, const unsigned char* in, size_t len) {
		if (bc == nullptr) {
			return -1;
		}
		return bc->encrypt(in, len, out);
	}
	int decrypt(unsigned char* out, const unsigned char* in, size_t len) {
		if (bc == nullptr) {
			return -1;
		}
		return bc->decrypt(in, len, out);
	}

private:
	blockCipherInterface* bc;
};

char* key = "webzen#@!01webzen#@!01webzen#@!0";

size_t bmdenc(unsigned char* out, const unsigned char* in, size_t len) {
	if (in==nullptr || len==0) {
		return 0;
	}
	if (out == nullptr) {
		return len + 34;
	}
	srand(time(nullptr));
	int alg1 = rand() % 11;
	int alg2 = rand() % 7;
	out[0] = (unsigned char)(alg1);
	out[1] = (unsigned char)(alg2);
	out = out + 2;
	unsigned char key1[32]; // never zero init
	memcpy(out, key1, 32);

	// encrypt with key1
	bmdCipher bc;
	bc.setkey(alg1, key1, 32);
	bc.encrypt(out+32, in, len);

	// encrypt with key2
	bmdCipher bc2;
	bc2.setkey(alg2, (const unsigned char*)key, 32);
	size_t len1024 = 1024;
	if (len > len1024) {
		// encrypt the head part with key1 included
		bc2.encrypt(out, out, len1024);
		// encrypt the tail part
		unsigned char* pos = out + 32 + len - len1024;
		bc2.encrypt(pos, pos, len1024);
	}
	if (len > len1024<<2) {
		// encrypt the middle part
		unsigned char* pos = out + (len>>1);
		bc2.encrypt(pos, pos, len1024);
	}
	return len+34;
}

size_t bmddec(unsigned char* out, const unsigned char* in, size_t len) {
	if (in==nullptr || len<=34) {
		return 0;
	}
	if (out == nullptr) {
		return len - 34;
	}
	int alg1 = in[0];
	int alg2 = in[1];
	in = in + 2;
	len -= 2;
	unsigned char* tmp = new unsigned char[len];
	memcpy(tmp, in, len);

	// decrypt with key2
	bmdCipher bc2;
	bc2.setkey(alg2, (const unsigned char*)key, 32);
	size_t len1024 = 1024;
	size_t lenOrigin = len - 32;
	if (lenOrigin > len1024<<2) {
		// decrypt the middle part
		unsigned char* pos = tmp + (lenOrigin>>1);
		bc2.decrypt(pos, pos, len1024);
	}
	if (lenOrigin > len1024) {
		// decrypt the tail part
		unsigned char* pos = tmp + len - len1024;
		bc2.decrypt(pos, pos, len1024);
		// decrypt the head part with key1 included
		bc2.decrypt(tmp, tmp, len1024);
	}

	// decrypt with key1
	bmdCipher bc1;
	bc1.setkey(alg1, tmp, 32);
	tmp = tmp + 32;
	len -= 32;
	bc1.decrypt(out, tmp, len);
	delete[]tmp;
	return len;
}