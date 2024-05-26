run: weather-cache
	OPENWEATHERMAP_API_KEY=`cat OPENWEATHERMAP_API_KEY` ./weather-cache

weather-cache: main.go
	go build

server.crt: server.key
	openssl req -new -x509 -sha256 -key $< -out $@ -days 3650
	
server.key:
	openssl ecparam -genkey -name secp384r1 -out $@

