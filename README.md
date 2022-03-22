# simple-mutating-webhook


### install simple mutating webhook in k8s
```shell
make install
```

### deploy sample deployment to k8s
```shell
make deploy CASE=ex1
make deploy CASE=ex2
make deploy CASE=ex3
```

### check webhook logs
```shell
kubectl logs $(kubectl get pod -n simple-mutating-webhook | grep simple | awk '{print $1}')
```

### remove all resources in k8s
```shell
make remove
```