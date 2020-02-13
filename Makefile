
generate_crds:
	operator-sdk generate crds

operator:
	operator-sdk build inventory

update: operator
	kubectl rollout restart deployment inventory

update_crd:
	kubectl delete -f deploy/crds/eranco74.redhat_v1alpha1_machinehealth_cr.yaml
	kubectl create -f deploy/crds/eranco74.redhat_v1alpha1_machinehealth_cr.yaml

register_machine:
	kubectl apply -f deploy/crds/eranco74.redhat_v1alpha1_machinehealth_cr.yaml

delete_machine:
	kubectl delete -f deploy/crds/eranco74.redhat_v1alpha1_machinehealth_cr.yaml

.PHONY: healthchecker
healthchecker:
	CGO_ENABLED=0 go build -o health_checker/build/health_checker health_checker/main.go
	docker build -t health_checker health_checker/

.PHONY: machine
machine:
	CGO_ENABLED=0  go build -o machine/build/machine machine/main.go
	docker build -t machine ./machine

push:
	docker tag machine quay.io/eranco/machine:latest
	docker push quay.io/eranco/machine:latest

run_machine:
	docker run --rm --name machine -d --hostname machine --net=host machine

kill_machine:
	docker kill machine