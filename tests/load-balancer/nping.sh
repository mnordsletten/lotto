# Basically: nping --tcp-connect -p 90 10.100.0.30

# Prerequisite:
# Run 'docker run -dit --name my-apache-app -p 8080:80 -v "$PWD":/usr/local/apache2/htdocs/ httpd:2.4' on
# lotto-client3 (10.100.0.150) and lotto-client4 (10.100.0.160)

sent=1000
target=950
rate=100 # Requests pr second, higher than 5 requires sudo
mode="--tcp-connect"
port=90
# delay=

raw=$(nping -c $sent $mode -p $port --rate $rate 10.100.0.30)
res=$(printf "%s" "$raw" | grep "Successful connections:")
# attempts=$(printf "%s" "$res" | cut -d ' ' -f 4)
received=$(printf "%s" "$res" | cut -d ' ' -f 8)
# successful=$(printf "%s" "$res" | cut -d ' ' -f 8)
# failed=$(printf "%s" "$res" | cut -d ' ' -f 11)

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
