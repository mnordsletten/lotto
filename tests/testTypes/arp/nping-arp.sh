# Nping arp
# Script used for arp pinging the instance at a fixed rate
# Returns rate and average response time

TARGET={{index .Template "target"}}

# Input values to cmd
sent=500
target=480
rate=50 # Requests pr second, higher than 5 requires sudo
raw=$(sudo nping --arp -q --count $sent --rate $rate $TARGET)

# Parse output
received=$(printf "%s" "$raw" | grep Rcvd | cut -d ' ' -f 8)

# If we receive more than 450 then the test passes
if [ "$received" -gt "$target" ]; then
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