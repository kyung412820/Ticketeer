# Ticketeer Kubernetes Manifests

## Included
- namespace
- backend / frontend Deployment + Service
- postgres / redis Deployment + Service
- backend / frontend ConfigMap
- backend Secret
- postgres / redis PVC
- ingress

## Notes
- Backend image: `ticketeer-backend:latest`
- Frontend image: `ticketeer-frontend:latest`
- For Minikube / Docker Desktop Kubernetes, build or load these images into the cluster first.
- `frontend-service.yaml` uses NodePort `30080`
- ingress host: `ticketeer.local`

## Apply order
```bash
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/configmap-backend.yaml
kubectl apply -f k8s/secret-backend.yaml
kubectl apply -f k8s/configmap-frontend.yaml
kubectl apply -f k8s/postgres-pvc.yaml
kubectl apply -f k8s/redis-pvc.yaml
kubectl apply -f k8s/postgres-deployment.yaml
kubectl apply -f k8s/postgres-service.yaml
kubectl apply -f k8s/redis-deployment.yaml
kubectl apply -f k8s/redis-service.yaml
kubectl apply -f k8s/backend-deployment.yaml
kubectl apply -f k8s/backend-service.yaml
kubectl apply -f k8s/frontend-deployment.yaml
kubectl apply -f k8s/frontend-service.yaml
kubectl apply -f k8s/ingress.yaml
```

## Verify
```bash
kubectl get all -n ticketeer
kubectl get ingress -n ticketeer
kubectl describe pod -n ticketeer
```
