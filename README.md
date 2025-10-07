## connectcheck

The aim of this service is to check how we can perform network call while k8s pod terminates

### Steps to run the service

1. `docker build -t connectcheck:latest .`
2. `kind load docker-image connectcheck:latest`
3. `kubectl apply -f deployment.yaml`
4. `kubectl apply -f service.yaml`

---

On pod termination through `kubectl rollout restart deployment connectcheck` we see the following logs while the pod is in `terminating` state

```
2025/10/07 05:20:33 Received signal: terminated
2025/10/07 05:20:33 ========================================
2025/10/07 05:20:33 Shutdown signal received. Making call to httpbin.org/delay/10 (this will take 10 seconds)...
2025/10/07 05:20:33 ========================================
2025/10/07 05:20:45 ========================================
2025/10/07 05:20:45 SUCCESS: Called httpbin.org/delay/10
2025/10/07 05:20:45   Status Code: 200
2025/10/07 05:20:45   Response Length: 323 bytes
2025/10/07 05:20:45   Time taken: 11.382627422s
2025/10/07 05:20:45 ========================================
2025/10/07 05:20:45 Shutting down server gracefully...
2025/10/07 05:20:45 Server stopped
```