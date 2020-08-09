#pragma once
#include "config.h"
#include <boost/archive/xml_iarchive.hpp>
#include <boost/archive/xml_oarchive.hpp>
#include <sstream>

NAMESPACE_BEGIN(WTFCipher)

class msg {
	friend class boost::serialization::access;
	template<class Archive>
	void serialize(Archive& ar, const unsigned int version) {
		ar & BOOST_SERIALIZATION_NVP(index);
		ar & BOOST_SERIALIZATION_NVP(value);
	}
public:
	int index;
	std::string value;
	msg(){}
	msg(int index, std::string value) :
		index(index), value(value) {}
};

class msgs {
	friend class boost::serialization::access;
	template<class Archive>
	void serialize(Archive& ar, const unsigned int version) {
		ar & boost::serialization::make_nvp("count", size);
		if (msgList == nullptr) {
			msgList = new msg[size];
		}
		for (int i = 0; i < size; i++) {
			ar & boost::serialization::make_nvp("msg", msgList[i]);
		}
	}
public:
	msg* msgList;
	int size;
	msgs() {
		msgList = nullptr;
		size = 0;
	}
	~msgs() {
		if (msgList != nullptr) {
			delete[] msgList;
		}
	}
};

#pragma pack (1)
typedef struct {
	unsigned char flag;
	unsigned char ver;
	char unk[22];
	int size;
}head;
#pragma pack ()

size_t encrypt(unsigned char* out, const unsigned char* in, size_t len) {
	// validate parameter
	if (in == nullptr || len == 0) {
		return 0;
	}

	// unmarshal
	msgs msgs;
	std::stringstream ss(std::string(in, in + len));
	boost::archive::xml_iarchive ia(ss);
	ia >> boost::serialization::make_nvp("msgs", msgs);

	// marshal
	size_t outLen = sizeof(head);
	for (int i = 0; i < msgs.size; i++) {
		outLen += 4 + msgs.msgList[i].value.length();
	}
	if (out == nullptr) {
		return outLen;
	}
	unsigned char* buf = new unsigned char[outLen];
	head* h = (head*)buf;
	h->flag = 0xCC;
	h->ver = 1;
	strcpy(h->unk, "WTF File");
	h->size = msgs.size;
	unsigned char* p = buf + sizeof(head);
	for (int i = 0; i < msgs.size; i++) {
		msg& m = msgs.msgList[i];
		*(short*)p = m.index;
		p += 2;
		size_t size = m.value.length();
		*(short*)p = size;
		p += 2;
		while (size--)
			p[size] ^= 0xCA;
		p += size;
	}
	memcpy(out, buf, outLen);
	delete[] buf;
	return outLen;
}

size_t decrypt(unsigned char* out, const unsigned char* in, size_t len) {
	// validate parameter
	if (in == nullptr || len == 0) {
		return 0;
	}

	// validate len
	if (len < sizeof(head)) {
		return 0;
	}

	// validate head
	head* h = (head*)in;
	if (h->flag != 0xCC || h->ver != 1 || h->size <= 0) {
		return 0;
	}

	msgs msgs;
	msgs.size = h->size;
	auto msgList = std::make_unique<msg[]>(h->size);

	// index to element array
	unsigned char* p = (unsigned char*)(in + sizeof(head));
	len -= sizeof(head);

	// range elements
	for (int i = 0; i < h->size; i++) {
		if (len <= 4) {
			return 0;
		}
		int index = *(short*)p;
		p += 2;
		len -= 2;
		int bufSize = *(short*)p;
		p += 2;
		len -= 2;
		if (len < bufSize) {
			return 0;
		}
		unsigned char* buf = new unsigned char[bufSize];
		memcpy(buf, p, bufSize);
		for (int j = 0; j < bufSize; j++) {
			buf[j] ^= 0xCA;
		}
		// std::vector<unsigned char> v(buf, buf + bufSize);
		std::vector<unsigned char> v(buf, buf + bufSize);
		delete[] buf;
		msgList[i] = msg(index, std::string(v.begin(), v.end()));
		p += bufSize;
		len -= bufSize;
	}
	msgs.msgList = msgList.release();

	// serialization
	std::stringstream ss;
	boost::archive::xml_oarchive oa(ss);
	oa << boost::serialization::make_nvp("msgs", msgs); // oa << BOOST_SERIALIZATION_NVP(m);
	
	// return
	std::string outStr = ss.str();
	size_t outLen = outStr.size();
	if (out != nullptr) {
		memcpy(out, outStr.c_str(), outLen);
	}
	return outLen;
}

NAMESPACE_END