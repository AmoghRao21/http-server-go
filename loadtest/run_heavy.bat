locust -f locustfile.py ^
  --host=http://localhost:8080 ^
  --users=2000 ^
  --spawn-rate=200 ^
  --run-time=5m ^
  --headless ^
  --html=heavy_report.html
