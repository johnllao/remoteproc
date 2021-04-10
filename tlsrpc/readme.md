# Generate a CA private key and certificate

    openssl genrsa -out ca.key 2048

    openssl req -new -x509 -sha256 -key ca.key -out ca.crt -days 365

# Generate CSR private key and certificate

    openssl genrsa -out server.key 2048

    openssl req -new -key server.key -out server-csr.crt

# Sign the CSR using the CA

    openssl x509 -req -in server-csr.crt -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 365 -sha256


# Generate a server private key and certificate

    openssl genrsa -out server.key 2048

    openssl req -new -x509 -sha256 -key server.key -out server.crt -days 365