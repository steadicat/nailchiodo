export CLOUDSDK_CORE_PROJECT=nail-chiodo

dev:
	go run main.go

deploy:
	gcloud app deploy --no-promote

.PHONY: dev deploy
