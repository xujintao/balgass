#1, pi.ca.key pi.ca.crt
openssl genrsa -out pi.ca.key 2048
openssl req -new -x509 -days 3650 -key pi.ca.key -out pi.ca.crt -subj "/C=CN/ST=Shanghai/L=Shanghai/O=pi/CN=pi.ca"

#2, mu.com.key mu.com.csr
openssl genrsa -out mu.com.key 2048
openssl req -new -key mu.com.key -out mu.com.csr -config openssl.cnf -subj "/C=CN/ST=Shanghai/L=Shanghai/O=mu/CN=mu.com"

#3, sign
openssl x509 -req -days 1460 -in mu.com.csr -CAkey pi.ca.key -CA pi.ca.crt -set_serial 01 -out mu.com.crt -extensions v3_req -extfile openssl.cnf
rm -rf mu.com.csr

#4, dists
sudo mv mu.com.key mu.com.crt ~/volumes/nginx-mu/ssl/
sudo cp pi.ca.crt ~/volumes/nginx-mu/dns/

sudo cp pi.ca.crt /usr/local/share/ca-certificates
sudo update-ca-certificates -f

#5, test
# curl -v https://mu.com
