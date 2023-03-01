# Kubernetes controller - Scaler

## Description
This is a Kubernetes controller that will scale up deployments based on the time defined in the CRD. 
```
apiVersion: api.framsouza.co/v1alpha1
kind: Scaler
metadata:
  name: scaler-sample
spec:
  start: "9:00"
  end: "10:15"
  replicas: 5
  deployments:
    - name: nginx
      namespace: default
```
