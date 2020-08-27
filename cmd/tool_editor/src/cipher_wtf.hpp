#pragma once
#include "config.h"
#include <boost/archive/xml_iarchive.hpp>
#include <boost/archive/xml_oarchive.hpp>
#include <boost/serialization/vector.hpp>
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

	// unmarshal to std::vector<msg>
	std::vector<msg> msgs;
	std::stringstream ss(std::string(in, in + len));
/*	boost::archive::xml_iarchive ia(ss);
	ia >> BOOST_SERIALIZATION_NVP(msgs);*/
	try {
		boost::archive::xml_iarchive ia(ss);
		ia >> BOOST_SERIALIZATION_NVP(msgs);
	}
	catch (const boost::archive::archive_exception& e) {
		if (out != nullptr) {
			strcpy((char*)out, e.what());
		}
		return strlen(e.what()) + 1;
	}
	catch (...) {
		return 0;
	}

	// encode and marshal to wtf
	size_t outLen = sizeof(head);
	for (auto m : msgs) {
		outLen += 4 + m.value.length();
	}
	if (out == nullptr) {
		return outLen;
	}
	unsigned char* buf = new unsigned char[outLen];
	head* h = (head*)buf;
	h->flag = 0xCC;
	h->ver = 1;
	strcpy(h->unk, "WTF File");
	h->size = msgs.size();
	unsigned char* p = buf + sizeof(head);
	for (auto m : msgs) {
		*(short*)p = m.index;
		p += 2;
		size_t size = m.value.length();
		*(short*)p = size;
		p += 2;
		memcpy(p, m.value.data(), size);
		for (int i = 0; i < size; i++) {
			p[i] ^= 0xCA;
		}
		p += size;
	}
	memcpy(out, buf, outLen);
	delete[] buf;
	return outLen;
}

size_t decrypt(unsigned char* out, const unsigned char* in, size_t len) {
	// validate parameter
	if (in == nullptr || len < sizeof(head)) {
		return 0;
	}

	// validate head
	head* h = (head*)in;
	if (h->flag != 0xCC || h->ver != 1 || h->size <= 0) {
		return 0;
	}

	// index to element array
	unsigned char* p = (unsigned char*)(in + sizeof(head));
	len -= sizeof(head);

	// decode and unmarshal to std::vector<msg>
	std::vector<msg> msgs;
	for (int i = 0; i < h->size; i++) {
		if (len <= 4) {
			return 0;
		}
		int index = *(short*)p;
		p += 2;
		len -= 2;
		int size = *(short*)p;
		p += 2;
		len -= 2;
		if (len < size) {
			return 0;
		}
		unsigned char* buf = new unsigned char[size];
		memcpy(buf, p, size);
		for (int j = 0; j < size; j++) {
			buf[j] ^= 0xCA;
		}
		// std::vector<unsigned char> v(buf, buf + size);
		std::vector<unsigned char> v(buf, buf + size);
		delete[] buf;
		msgs.push_back(msg(index, std::string(v.begin(), v.end())));
		p += size;
		len -= size;
	}

	// marshal to xml
	std::stringstream ss;
	try {
		boost::archive::xml_oarchive oa(ss);
		oa << BOOST_SERIALIZATION_NVP(msgs); // oa << BOOST_SERIALIZATION_NVP(m);
	}
	catch (const boost::archive::archive_exception& e) {
		if (out != nullptr) {
			strcpy((char*)out, e.what());
		}
		return strlen(e.what()) + 1;
	}
	catch (...) {
		return 0;
	}
	
	// return
	std::string outStr = ss.str();
	size_t outLen = outStr.size();
	if (out != nullptr) {
		memcpy(out, outStr.c_str(), outLen);
	}
	return outLen;
}

NAMESPACE_END