# Pinger
# Script used for pinging the instance at a fixed rate
# Returns rate and average response time

# Input values to cmd
sent=5
rate=5 # Requests pr second, higher than 5 requires sudo
raw=$(ping -c $sent -i $(awk "BEGIN {print 1/$rate}") -q 10.100.0.150)

# Parse output
received=$(printf "%s" "$raw" | grep received | cut -d ' ' -f 4)

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
