# Input values to cmd

SERVER_ADDRESS={{index .Template "serverAddress"}}
SERVER_PORT={{index .Template "serverPort"}}

sent=1
raw=$(iperf3 -c $SERVER_ADDRESS -p $SERVER_PORT -n 1024M -f m 2>&1)
if [ "$?" -eq "0" ]; then
  received=$[$received + 1]
fi
rate=$(printf "%s" "$raw" | grep sender | perl -n -e'/([0-9]+) Mbits/ && print $1')

sent=$[$sent + 1]
raw+=$(iperf3 -c $SERVER_ADDRESS -p $SERVER_PORT -n 1024M -f m -R 2>&1)
if [ "$?" -eq "0" ]; then
  received=$[$received + 1]
fi

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
