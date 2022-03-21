# Create Key Pair
NS="simple-mutating-webhook"
DNS="simple-mutating-webhook.simple-mutating-webhook.svc"

# have to use openssl@3 on mac
/usr/local/opt/openssl/bin/openssl req -x509 -newkey rsa:2048 -keyout tls.key -out tls.crt -days 365  \
  -addext "subjectAltName = DNS:${DNS}" \
  -nodes -subj "/CN=${DNS}"

# Create NS
kubectl create ns ${NS}

# Create TLS secret for service
kubectl create secret tls webhook-certs -n ${NS} \
    --cert "tls.crt" \
    --key "tls.key"

# Create Deployment, Service
kubectl apply -f ./manifests/webhook/deployment.yaml
kubectl apply -f ./manifests/webhook/service.yaml

# Register Mutator Server as a Mutate Webhook to Kubernetes
export CA_PEM_BASE64="$(openssl base64 -A <"tls.crt")"
cat ./manifests/webhook/mutating-webhook-configuration.yaml | sed "s/{{CA_PEM_BASE64}}/$CA_PEM_BASE64/g" | kubectl apply -n ${NS} -f -

# Clean files
rm -rf tls.*