export CLOUDSDK_CORE_PROJECT=nail-chiodo

dev:
	/usr/local/Caskroom/google-cloud-sdk/latest/google-cloud-sdk/bin/dev_appserver.py .

deploy:
	gcloud app deploy --no-promote

.PHONY: dev deploy
