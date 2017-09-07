build:
	go build
	docker build . -t imeoer/faas-hyper
