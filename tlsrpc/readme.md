# Generate a CA private key and certificate

    openssl genrsa -out ca.key 2048

    openssl req -new -x509 -key ca.key -out ca.crt

# Generate CSR private key and certificate

    openssl genrsa -out server.key 2048

    openssl req -new -key server.key -out server.csr

    openssl genrsa -out client1.key 2048

    openssl req -new -key client1.key -out client1.csr

# Sign the CSR using the CA

    openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.pem

    openssl x509 -req -in client1.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client1.pem 
