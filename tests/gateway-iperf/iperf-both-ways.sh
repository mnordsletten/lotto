# Input values to cmd

sent=1
raw=$(iperf3 -c 10.100.0.150 -p 12345 -n 1024M -f m 2>&1)
if [ "$?" -eq "0" ]; then
  received=$[$received + 1]
fi
rate=$(printf "%s" "$raw" | grep sender | perl -n -e'/([0-9]+) Mbits/ && print $1')

sent=$[$sent + 1]
raw+=$(iperf3 -c 10.100.0.150 -p 12345 -n 1024M -f m -R 2>&1)
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
