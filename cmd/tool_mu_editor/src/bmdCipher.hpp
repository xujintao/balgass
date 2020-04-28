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
	virtual void SetKey(const unsigned char* key, size_t len) = 0;
	//virtual void encrypt(unsigned char* out, const unsigned char* in, int len) = 0;
	//virtual void decrypt(unsigned char* out, const unsigned char* in, int len) = 0;
};

template <class E, class D>
class blockCipher : public blockCipherInterface {
public:
	void SetKey(const unsigned char* key, size_t len) {
		enc.SetKey(key, 0x10);
		dec.SetKey(key, 0x10);
	}
private:
	E enc;
	D dec;
};


class bmdCipher
{
public:
	bmdCipher() {
		this->bc = nullptr;
	}
	~bmdCipher() {
		if (this->bc != nullptr) {
			delete this->bc;
		}
	}
	void setkey(int alg, const unsigned char* key, size_t len) {
		switch (alg & 7)
		{
		case 0: this->bc = new(blockCipher<CryptoPP::TEAEncryption, CryptoPP::TEADecryption>);break;
		case 1: this->bc = new(blockCipher<CryptoPP::ThreeWayEncryption, CryptoPP::ThreeWayDecryption>);break;
		case 2: this->bc = new(blockCipher<CryptoPP::CAST128Encryption, CryptoPP::CAST128Decryption>);break;
		case 3: this->bc = new(blockCipher<CryptoPP::RC5Encryption, CryptoPP::RC5Decryption>);break;
		case 4: this->bc = new(blockCipher<CryptoPP::RC6Encryption, CryptoPP::RC6Decryption>);break;
		case 5: this->bc = new(blockCipher<CryptoPP::MARSEncryption, CryptoPP::MARSDecryption>);break;
		case 6: this->bc = new(blockCipher<CryptoPP::IDEAEncryption, CryptoPP::IDEADecryption>);break;
		case 7: this->bc = new(blockCipher<CryptoPP::GOSTEncryption, CryptoPP::GOSTDecryption>);break;
		}
		this->bc->SetKey(key, len);
	}
	void encrypt(unsigned char* out, const unsigned char* in, size_t len) {

	}
	void decrypt(unsigned char* out, const unsigned char* in, size_t len) {

	}

private:
	blockCipherInterface* bc;
};

size_t bmdenc(unsigned char* out, const unsigned char* in, size_t len) {
	if (out == nullptr) {
		return len + 34;
	}
	return len;
}

size_t bmddec(unsigned char* out, const unsigned char* in, size_t len) {
	if (out == nullptr) {
		return len - 34;
	}
	return len;
}