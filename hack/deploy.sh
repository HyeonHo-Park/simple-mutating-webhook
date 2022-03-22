CASE=$1

kubectl apply -f ./manifests/examples/namespace.yaml
sleep 1

kubectl apply -f ./manifests/examples/${CASE}-deployment.yaml