.PHONY: build clean

build:
	go build -o dnsupdater .

deploy:
	GOOS=linux GOARCH=arm64 go build -o dnsupdater .
	scp -o StrictHostKeyChecking=no dnsupdater ruben@homeserver:/home/ruben/bin/dnsupdater

clean:
	rm -rf dnsupdater
