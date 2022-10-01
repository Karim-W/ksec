# ksec

KSec is a command line tool to manage secrets in Kubernetes with the following functionallities:
- [x] Create a secret from an env file
- [x] Append a secret to an existing secret
- [x] Get a secret from Kubernetes secrets 
- [x] Delete a secret from Kubernetes secrets
- [x] List all secrets in a namespace

## Installation
to install Ksec use the following command:
` go install github.com/karim-w/ksec `

## Usage
### Create a secret from an env file
` ksec -e <.env file path> -n <namespace> -s <secret name> `
### Append a secret to an existing secret
` ksec -w <.env file path> -n <namespace> -s <secret name> -a `
### List all secrets in a Kubernetes secret
` ksec -l -n <namespace> -s <secret name> `
### Get a secret from Kubernetes secrets
` ksec -g -n <namespace> -s <secret name> -k <key> `
### Delete a secret from Kubernetes secrets
` ksec -d -n <namespace> -s <secret name> -k <key> `

