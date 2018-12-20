# Pinging 10.100.0.160
# Script used for pinging the instance at a fixed rate
# Returns rate and average response time

sent=20
rate=5 # Requests pr second, higher than 5 requires sudo
cmdOut=$(ping -c $sent -i $(awk "BEGIN {print 1/$rate}") -q 10.100.0.160)
received=$(printf "%s" "$cmdOut" | grep received | cut -d ' ' -f 4)

# Fail test if there were received packets. This test should not succeed
if [ "$received" -eq "0" ]; then
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
