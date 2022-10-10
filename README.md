# Soteria Blokchain Functions

## Dependencies

- Google cloud SDK
- Go (1.11 or higher)

## Prerequisites

Make sure to be logged in to GCP in your cli
```
gcloud auth login
```
In order to run the function locally, you need to have the default credentials set.
The easiest way to go is via gcloud cli:
```
gcloud auth application-default login
```
Now you should be able to boot a local server that can run your functions like this:
```
make serve
```
## Usage

Use make file for easier commands

### Functions management

#### 1. Deploy a function:

------------
Deploy a function into GCP `Cloud Functions`
```
# for dev/staging environment
make deploy-fn fn=<YOUR_FUNC_NAME>

# for prod
make deploy-fn-prod fn=<YOUR_FUNC_NAME>
```
#### 2. Call a function
----------------
Calls a function hosted in GCP `Cloud Functions` from your local
```
# for dev/staging environment
make call-fn fn=<YOUR_FUNC_NAME>

# for prod
make call-fn-prod fn=<YOUR_FUNC_NAME>
```

### 3. Run a function from local server
-----------------
After booting up a local server, you can use curl or postman to run any function defined locally even if it is not deployed yet.
Make sure to have default credentials and the server running first.
```
curl http://localhost:8080/<YOUR_FUNC_NAME>
```

