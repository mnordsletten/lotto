# http testing (downloads file from load balancer)
# Script used for performing http requests (get file.txt from lotto-client3 10.100.0.150 and lotto-client4 10.100.0.160)

SERVER_ADDRESS={{index .Template "serverAddress"}}

# Input values to cmd
sent=1
received=0

# Loop and attempt to download file 5 times. Controlling the size of the download for success
for i in $(seq 1 $sent); do
  out=$(wget -q -T 10 --tries 1 -O - http://$SERVER_ADDRESS/1GB_file.txt | wc -c - | cut -d' ' -f 1)
  if [ "$out" -eq "1073741824" ]; then
    received=$((received + 1))
    raw="$raw\
    run: $i: received: $out SUCCESS"
  else
    raw="$raw\
    run: $i: received: $out FAIL"
  fi
done

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
