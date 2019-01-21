
TARGET={{index .Template "target"}}

sent=1000
rate=100 # Requests pr second, higher than 5 requires sudo
mode="--udp"
port=4242
# delay=
data="hi"

raw=$(nping -c $sent --rate $rate $mode -p $port --data-string $data $TARGET)
res=$(printf "%s" "$raw" | grep "UDP packets")

# Possible:
# attempts=$(printf "%s" "$res" | cut -d ' ' -f 4)
# if [ -z $attempts ]; then attempts=0; fi

received=$(printf "%s" "$res" | cut -d ' ' -f 7)
# Or:

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
