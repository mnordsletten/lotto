# hey http testing
# Script used for performing http requests (get file.txt from lotto-client4 10.100.0.160)
# Returns rate and average response time

# Prerequisites to test:
# On lotto-client4 (10.100.0.160):
# - Produce a 1G text file with random content:
# 	base64 /dev/urandom | head -c 1G > file.txt
# - Start apache server on 8080 (returns content of home folder):
# 	docker run --rm -dit --name my-apache-app -p 8080:80 -v "$PWD":/usr/local/apache2/htdocs/ httpd:2.4

# Input values to cmd
sent=100
concurrency=50

raw=$(docker run --rm rcmorano/docker-hey -n $sent -c $concurrency http://10.100.0.30:1600/1GB_file.txt)

# Parse output
received=$(printf "%s" "$raw" | awk '/responses/ {print $2}' )
rate=$(printf "%s" "$raw" | awk '/Requests\/sec/ {print $2}' )

if [ "$sent" -eq "$received" ]; then
  success=true
fi

jq \
  --argjson success ${success:-false} \
  --argjson sent ${sent:-0} \
  --argjson received ${received:-0} \
  --argjson rate ${rate:-0} \
  --arg raw "${raw:-''}" \
  '. |
  .["success"]=$success |
  .["sent"]=$sent |
  .["received"]=$received |
  .["rate"]=$rate |
  .["raw"]=$raw
  '<<<'{}'
