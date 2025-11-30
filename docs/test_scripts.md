# 🧪 Test Scripts & Verification Procedures

This document contains **complete testing instructions** for verifying correctness, performance, and reliability of the Minimalist High‑Performance Go HTTP Server.

---

# ✅ 1. Basic Connectivity Tests

### **Test: Root Endpoint**
```bash
curl -v http://localhost:8080/
```

### **Expected Response**
```
hello
```

---

# ✅ 2. Echo Endpoint Test
```bash
curl -v "http://localhost:8080/echo?message=hello"
```

Expected:
```json
{ "message": "hello" }
```

---

# 🔐 3. Authentication Tests

### **Auth Header**
```
Authorization: Basic YWRtaW46c2VjcmV0
```

### **Test Unauthorized Request**
```bash
curl -v http://localhost:8080/data
```
Expected:
```
401 unauthorized
```

### **Test Authorized**
```bash
curl -v -H "Authorization: Basic YWRtaW46c2VjcmV0" http://localhost:8080/data
```

---

# 📥 4. POST /data Test

```bash
curl -v -X POST http://localhost:8080/data \
 -H "Content-Type: application/json" \
 -H "Authorization: Basic YWRtaW46c2VjcmV0" \
 -d '{"name": "alpha"}'
```

Expected:
```json
{ "id": 1, "name": "alpha" }
```

---

# 📤 5. GET /data Test

```bash
curl -v -H "Authorization: Basic YWRtaW46c2VjcmV0" http://localhost:8080/data
```

---

# 🎯 6. PUT /data/:id

```bash
curl -v -X PUT http://localhost:8080/data/1 \
 -H "Authorization: Basic YWRtaW46c2VjcmV0" \
 -H "Content-Type: application/json" \
 -d '{"name": "updated"}'
```

---

# ✂️ 7. PATCH /data/:id

```bash
curl -v -X PATCH http://localhost:8080/data/1 \
 -H "Authorization: Basic YWRtaW46c2VjcmV0" \
 -H "Content-Type: application/json" \
 -d '{"flag": true}'
```

---

# 🗑 8. DELETE /data/:id

```bash
curl -v -X DELETE http://localhost:8080/data/1 \
 -H "Authorization: Basic YWRtaW46c2VjcmV0"
```

---

# 🧱 9. Static File Test

Place a file:
```
static/test.txt
```

Then test:
```bash
curl -v http://localhost:8080/static/test.txt
```

---

# 🧨 10. Body Size Limit Test (Expect 413)

```bash
python3 - << 'EOF'
import requests
big = "x" * (2_000_000)
r = requests.post(
    "http://localhost:8080/data",
    headers={"Authorization": "Basic YWRtaW46c2VjcmV0"},
    data=big
)
print(r.status_code, r.text)
EOF
```

Expected:
```
413 body too large
```

---

# 🔥 11. Stress Testing (Locust)

File: `locustfile.py`
```python
from locust import HttpUser, task, between

class Load(HttpUser):
    wait_time = between(1, 2)

    @task
    def root(self):
        self.client.get("/")
```

Run:
```bash
locust
```

Then open:
```
http://localhost:8089
```

---

# 🧩 12. Keep‑Alive Connection Test

```bash
curl -v http://localhost:8080/ -H "Connection: keep-alive"
curl -v http://localhost:8080/echo?message=ok -H "Connection: keep-alive"
```

Expected:
- No reconnection between requests.

---

# 🧹 13. Cleanup Test (DELETE + Validate)

```bash
curl -v -X DELETE http://localhost:8080/data/1 -H "Authorization: Basic YWRtaW46c2VjcmV0"
curl -v http://localhost:8080/data/1 -H "Authorization: Basic YWRtaW46c2VjcmV0"
```

Expected:
```
404 not found
```

---

# 🏁 Summary

This testing suite verifies:
- Protocol correctness  
- CRUD behavior  
- Storage correctness  
- Authentication  
- Static file serving  
- Error handling  
- Resource limits  
- Performance characteristics  
 
