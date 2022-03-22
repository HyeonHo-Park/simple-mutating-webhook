# simple-mutating-webhook
### simple-mutating-webhook example은 아래와 같은 기능을 제공합니다
- test namesapce에 배포되는 deployment에 대한 mutating 및 validating 수행
- 배포되는 deployemnt의 container의 req, limit cpu validation (200m <= cpu <= 500m) 
- cpu < 200m -> 200m
- cpu > 500m -> inject
- 배포되는 container total req or limit cpu > 1000m -> inject

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
kubectl logs -f $(kubectl get pod -n simple-mutating-webhook | grep simple | awk '{print $1}')
```

### remove all resources in k8s
```shell
make remove
```