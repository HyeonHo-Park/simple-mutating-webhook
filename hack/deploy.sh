# Create Key Pair
NS="simple-mutating-webhook"
DNS="simple-mutating-webhook.${NS}.svc"
KEYNAME="webhook-server-tls"
openssl req -nodes -new -x509 -keyout ca.key -out ca.crt -subj "/CN=Simple Mutating Webhook CA"
openssl genrsa -out ${KEYNAME}.key 2048
openssl req -new -key ${KEYNAME}.key -subj "/CN=${DNS}" \
    | openssl x509 -req -CA ca.crt -CAkey ca.key -CAcreateserial -out ${KEYNAME}.crt

# Register Mutator Server as a Mutate Webhook to Kubernetes
export CA_PEM_BASE64="$(openssl base64 -A <"ca.crt")"
cat ./manifests/webhook/mutating-webhook-configuration.yaml | sed "s/{{CA_PEM_BASE64}}/$CA_PEM_BASE64/g" | kubectl apply -n ${NS} -f -

# Create NS
kubectl create ns ${NS}

# Create TLS secret for service
kubectl create secret tls webhook-certs -n ${NS} \
    --cert "${KEYNAME}.crt" \
    --key "${KEYNAME}.key"

# Create Deployment, Service
kubectl apply -f ./manifests/webhook/deployment.yaml
kubectl apply -f ./manifests/webhook/service.yaml

# Clean files
rm -rf ${KEYNAME}.* ca.*