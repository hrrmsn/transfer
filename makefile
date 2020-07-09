build:
	go build cmd/delivery-server/main.go

run:
	./main -port 8081

delivery:
	curl -X POST -H 'Content-Type: application/json' -d '{"lat": 17.986511, "lng": 63.441092}' 'localhost:8081/delivery'
