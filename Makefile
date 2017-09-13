build:
	go build
	docker build . -t imeoer/faas-hyper
	docker build prometheus/ -t imeoer/faas-prometheus
