locust -f locustfile.py ^
  --host=http://localhost:8080 ^
  --users=500 ^
  --spawn-rate=50 ^
  --run-time=2m ^
  --headless ^
  --html=report.html
