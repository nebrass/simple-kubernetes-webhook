.PHONY: test
test:
	@echo "\n🛠️  Running unit tests..."
	go test ./...

.PHONY: build
build:
	@echo "\n🔧  Building Go binaries..."
	GOOS=darwin GOARCH=amd64 go build -o bin/admission-webhook-darwin-amd64 .
	GOOS=linux GOARCH=amd64 go build -o bin/admission-webhook-linux-amd64 .

.PHONY: generate-certs
generate-certs:
	@echo "\n🔐 Generating certificates..."
	./dev/gen-certs.sh

.PHONY: docker-build
docker-build:
	@echo "\n📦 Building simple-kubernetes-webhook Docker image..."
	docker build -t simple-kubernetes-webhook:latest .

# From this point `kind` is required
.PHONY: cluster
cluster:
	@echo "\n🔧 Creating Kubernetes cluster..."
	kind create cluster --config dev/manifests/kind/kind.cluster.yaml

.PHONY: delete-cluster
delete-cluster:
	@echo "\n♻️  Deleting Kubernetes cluster..."
	kind delete cluster

.PHONY: push
push: docker-build
	@echo "\n📦 Pushing admission-webhook image into Kind's Docker daemon..."
	kind load docker-image simple-kubernetes-webhook:latest

.PHONY: deploy-config
deploy-config:
	@echo "\n⚙️  Applying cluster config..."
	kubectl apply -f dev/manifests/cluster-config/

.PHONY: delete-config
delete-config:
	@echo "\n♻️  Deleting Kubernetes cluster config..."
	kubectl delete -f dev/manifests/cluster-config/ || true

.PHONY: deploy
deploy: push delete generate-certs deploy-config
	@echo "\n🚀 Deploying simple-kubernetes-webhook..."
	kubectl apply -f dev/manifests/webhook/

.PHONY: delete
delete:
	@echo "\n♻️  Deleting simple-kubernetes-webhook deployment if existing..."
	kubectl delete -f dev/manifests/webhook/ || true

.PHONY: pod
pod:
	@echo "\n🚀 Deploying test pod..."
	kubectl apply -f dev/manifests/pods/busybox.yaml

.PHONY: pod-image
pod-image:
	@echo "\n🔍 Deployed test pod image:" $(shell kubectl get pod busybox -n apps -o jsonpath='{.spec.containers[0].image}')

.PHONY: delete-pod
delete-pod:
	@echo "\n♻️ Deleting test pod..."
	kubectl delete -f dev/manifests/pods/busybox.yaml || true

.PHONY: bad-pod
bad-pod:
	@echo "\n🚀 Deploying \"bad\" pod..."
	kubectl apply -f dev/manifests/pods/bad-name.pod.yaml

.PHONY: delete-bad-pod
delete-bad-pod:
	@echo "\n🚀 Deleting \"bad\" pod..."
	kubectl delete -f dev/manifests/pods/bad-name.pod.yaml || true

.PHONY: logs
logs:
	@echo "\n🔍 Streaming simple-kubernetes-webhook logs..."
	kubectl logs -l app=simple-kubernetes-webhook -f

.PHONY: delete-all
delete-all: delete delete-config delete-pod delete-bad-pod delete-cluster

.PHONY: verify
verify: 
	@echo "\n🕵️‍♂️  Running verification checks..."

	@echo "\n🕵️‍♂️  Check Validating webhook"
	@if $(MAKE) bad-pod > /dev/null 2>&1; then \
		echo "❌  Validating webhook might not be running correctly"; \
		exit 1; \
	else \
		echo "✅  Validating webhook seems to be functioning"; \
		output=$$(kubectl apply -f dev/manifests/pods/bad-name.pod.yaml 2>&1); \
		echo "$$output" | grep -o 'admission webhook .* denied the request: .*' || { \
			echo "❌  Failed to extract the expected error message."; \
			echo "$$output"; \
			exit 1; \
		}; \
	fi

	@echo "\n🕵️‍♂️  Check Mutating webhook"
	@image=$$(kubectl get pod busybox -n apps -o jsonpath='{.spec.containers[0].image}'); \
	if echo "$$image" | grep -q "^contoso.acr.io"; then \
		echo "✅  Image deployed from expected registry (contoso.acr.io)"; \
	else \
		echo "❌  Image deployed from unexpected registry. Expected 'contoso.acr.io'"; \
		exit 1; \
	fi
