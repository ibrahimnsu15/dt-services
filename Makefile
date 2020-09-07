service := dt-services
version := 0.0.11
docker_org := pineappleworkshop
gcloud_proj := pineappleworkshop
region := us-central1-c
cluster := pw
docker-image := ${docker_org}/${service}:${version}
root := $(abspath $(shell pwd))
port := 7004

list:
	@grep '^[^#[:space:]].*:' Makefile | grep -v ':=' | grep -v '^\.' | sed 's/:.*//g' | sed 's/://g' | sort

bootstrap:
	go mod init $(service)
	make init

init:
	go mod tidy

port-forward:
    @kubectl port-forward pw-mongodb-replicaset-0 27017:27017
    @kubectl port-forward consul-server-0 8500:8500

dev:
	go run main.go

docker-build:
	docker build -t $(docker-image) .

docker-dev:
	make docker-build
	make docker-run

docker-push:
	docker push $(docker-image)

docker-run:
	@docker run -itp $(port):$(port)  $(docker-image)

bumpversion-patch:
	bumpversion patch --allow-dirty

bumpversion-minor:
	bumpversion minor --allow-dirty

bumpversion-major:
	bumpversion major --allow-dirty

bootstrap-deploy:
	gcloud container clusters get-credentials ${cluster} --zone ${region} --project ${gcloud_proj}
	make docker-build
	make docker-push
	kubectl create -f deployments/k8s/service.yml
	kubectl create -f deployments/k8s/deploy.yml

ci-deploy:
	./google-cloud-sdk/bin/gcloud container clusters get-credentials ${cluster} --zone ${region} --project ${gcloud_proj}
	./kubectl apply -f deployments/k8s/deploy.yml

deploy:
	git add .
	git commit -m "${version}"
	git push origin master
	gcloud container clusters get-credentials ${cluster} --zone us-central1-c --project ${gcloud_proj}
	make docker-build
	make docker-push
	kubectl apply -f deployments/k8s/deploy.yml

purge:
	go clean
	rm -rf $(root)/vendor
