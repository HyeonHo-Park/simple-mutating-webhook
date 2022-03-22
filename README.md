# simple-mutating-webhook

### install mutating webhook in k8s
```shell
make install
```

### deploy sample deployment
```shell
make deploy
```

### check webhook logs
```shell
kubectl logs $(kubectl get pod -n simple-mutating-webhook | grep simple | awk '{print $1}')
```
