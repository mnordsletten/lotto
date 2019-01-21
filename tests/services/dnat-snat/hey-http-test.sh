# hey http testing
# Script used for performing http requests
# Returns rate and average response time

# Prerequesites:
# Requires a web server to run on client3

# Input values to cmd
sent=1000
concurrency=200
timeout=2
raw=$(docker run --rm rcmorano/docker-hey -t $timeout -n $sent -c $concurrency http://10.100.0.30:1500)

# Parse output, important to set a default value if the command over fails
received=$(printf "%s" "$raw" | awk '/responses/ {print $2}' )
rate=$(printf "%s" "$raw" | awk '/Requests\/sec/ {print $2}' )

# Only passes if 100% of packets were received
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
