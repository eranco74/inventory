
generate_crds:
	operator-sdk generate crds

build_operator:
	operator-sdk build inventory

update: build_operator
	kubectl rollout restart deployment inventory

update_crd:
	kubectl delete -f deploy/crds/eranco74.redhat_v1alpha1_machinehealth_cr.yaml
	kubectl create -f deploy/crds/eranco74.redhat_v1alpha1_machinehealth_cr.yaml

register_machine:
	kubectl apply -f deploy/crds/eranco74.redhat_v1alpha1_machinehealth_cr.yaml

delete_machine:
	kubectl delete -f deploy/crds/eranco74.redhat_v1alpha1_machinehealth_cr.yaml

build_healthchecker:
	CGO_ENABLED=0 go build -o health_checker/build/health_checker health_checker/main.go

image: build_healthchecker
	docker build -t health_checker health_checker/