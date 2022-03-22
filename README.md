# simple-mutating-webhook


### install simple mutating webhook in k8s
```shell
make install
```

### deploy sample deployment to k8s
```shell
make deploy ex1
make deploy ex2
make deploy ex3
```

### check webhook logs
```shell
kubectl logs $(kubectl get pod -n simple-mutating-webhook | grep simple | awk '{print $1}')
```

### remove all resources in k8s
```shell
make remove
```