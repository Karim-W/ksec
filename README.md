# ksec

KSec is a command line tool to manage secrets in Kubernetes with the following functionallities:
- [x] Create a secret from an env file
- [x] Append a secret to an existing secret
- [x] Get a secret from Kubernetes secrets 
- [x] Delete a secret from Kubernetes secrets
- [x] List all secrets in a namespace
- [x] Fill a file with a secret from Kubernetes secrets

## Installation
to install Ksec use the following command:
` go install github.com/karim-w/ksec `

## Usage
### Create a secret from an env file
` ksec -e <.env file path> -n <namespace> -s <secret name> `
this command will :
- create a secret from the env file and will add it to the kubernetes secrets
- create a yaml file with the env config map 
### Append a secret to an existing secret
` ksec -w <.env file path> -n <namespace> -s <secret name> -a `
this commmand will add a secret to a existing secret in kubernetes secrets
### List all secrets in a Kubernetes secret
` ksec -l -n <namespace> -s <secret name> `
this command will retrieve the secrets embedde in a kubernetes secrets
### Get a secret from Kubernetes secrets
` ksec -g -n <namespace> -s <secret name> -k <key> `
this command will retrieve a the value of secret within an existing kubernetes secret
### Delete a secret from Kubernetes secrets
` ksec -d -n <namespace> -s <secret name> -k <key> `
this command will delete a secret from an existing kubernetes secret

### Fill a file with secrets from Kubernetes secrets
` ksec -f <file path> -n <namespace> -s <secret name> `

