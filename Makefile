build-docker:
	docker build -t crawler -f infra/Dockerfiles/Dockerfile .

run-docker:
	docker run -it crawler
