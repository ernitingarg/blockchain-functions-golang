DEV=black-stream-292507
PROD=soteria-production
GOVERSION=go111
PATHTODEVKEY="/Users/thomasproust/Workplace/soteria/bitcoin-functions/.key/dev.json"
PATHTOPRODKEY="/Users/thomasproust/Workplace/soteria/bitcoin-functions/.key/prod.json"


.PHONY set-dev:
set-dev:
	export GOOGLE_APPLICATION_CREDENTIALS=$(PATHTODEVKEY); \
	export GCP_PROJECT=$(DEV); \
	gcloud config set project $(DEV)

.PHONY set-prod:
set-prod:
	export GOOGLE_APPLICATION_CREDENTIALS=$(PATHTOPRODKEY); \
	export GCP_PROJECT=$(PROD); \
	gcloud config set project $(PROD)


.PHONY: deploy-fn
deploy-fn: set-dev
	gcloud functions deploy $(fn) \
	--runtime $(GOVERSION) \
	--trigger-http \
	--allow-unauthenticated \


.PHONY: call-fn
call-fn: set-dev
	gcloud functions call $(fn)

.PHONY: deploy-fn-prod
deploy-fn-prod: set-prod
	gcloud functions deploy $(fn) \
	--runtime $(GOVERSION) \
	--trigger-http \
	--allow-unauthenticated \

.PHONY: call-fn-prod
call-fn-prod: set-prod
	gcloud functions call $(fn)

.PHONY: serve
serve: set-prod
	export GOOGLE_APPLICATION_CREDENTIALS=$(PATHTOPRODKEY); \
	export GCP_PROJECT=$(PROD); \
	go run cmd/main.go

.PHONY: deploy-pb-prod
deploy-pb-prod: set-prod
	gcloud functions deploy $(fn) \
	--runtime $(GOVERSION) \
	--trigger-topic $(topic) \