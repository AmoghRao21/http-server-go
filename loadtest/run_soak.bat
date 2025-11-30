locust -f locustfile.py ^
  --host=http://localhost:8080 ^
  --users=300 ^
  --spawn-rate=20 ^
  --run-time=30m ^
  --headless ^
  --html=soak_report.html
