run-cloud-proxy:
	sudo ./cloud_sql_proxy \
	-dir=/cloudsql -\
	-projects=${project_id} \
	-instances=${INSTANCE_CONNECTION_NAME}
psql-cloud-proxy:
	psql -U postgres -h /cloudsql/${INSTANCE_CONNECTION_NAME}
go-test:
	go test ./... -v
run:
	go run main.go
dev:
	CompileDaemon --build="go build main.go" --command=./main
go-mod-init:
	rm -rf vendor
	rm go.mod
	rm go.sum
	go mod init
	go mod vendor
docker-prune:
	docker image prune
docker-build-local:
	docker build -t goucher:${version} .
docker-run-local:
	docker run -p 3000:3000 goucher:${version}
publish-google-registry:
	gcloud builds submit --tag gcr.io/${project_id}/goucher:${version}
deploy-cloud-run:
	gcloud run deploy --image gcr.io/${project_id}/goucher:${version} --platform managed \
	--concurrency=40 --memory 128Mi --timeout=30s --max-instances 10 --cpu 1 \
	--add-cloudsql-instances ${INSTANCE_CONNECTION_NAME} \
	--set-env-vars INSTANCE_CONNECTION_NAME=${INSTANCE_CONNECTION_NAME} \
	--set-env-vars DB_USER=${DB_USER} \
    --set-env-vars DB_PASS=${DB_PASS} \
    --set-env-vars DB_NAME=${DB_NAME}