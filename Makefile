IMG=pcrofts/awsranchhand
BUILD_IMG=pcrofts/build-awsranchhand
BUILD_NAME=build-awsranchhand

default: builddocker

buildgo:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o main ./go/src/github.com/margic/awsranchhand

cleandocker:
	-rm awsranchhand
	-docker rm -v $(BUILD_NAME)
	-docker rmi -f $(IMG)


builddocker: cleandocker
	docker build -t $(BUILD_IMG) -f Dockerfile.build .
	docker run --name $(BUILD_NAME) -t $(BUILD_IMG) /bin/true
	docker cp $(BUILD_NAME):/go/bin/awsranchhand ./awsranchhand
	chmod 755 ./awsranchhand
	docker build --no-cache --rm=true -t $(IMG) -f Dockerfile.static .
	docker rm -v $(BUILD_NAME)
	rm ./awsranchhand

tagdocker: builddocker
	docker tag -f $(IMG)

run: builddocker
	docker run -it --rm $(IMG) --help

alpine-build:
	GOOS=linux GOARCH=amd64 go build -o awsranchhand
	docker build -t $(IMG) -f Dockerfile.static .
