# Connect to apache server running on client4 (TCP)
# Basically: nping --tcp-connect -p 8080 10.100.0.160

# Prerequisite:
# Run 'docker run -dit --name my-apache-app -p 8080:80 -v "$PWD":/usr/local/apache2/htdocs/ httpd:2.4' on
# lotto-client4 (10.100.0.160)

SERVER_ADDRESS={{index .Template "serverAddress"}}
SERVER_PORT={{index .Template "serverPort"}}

sent=600 # The apache server we use can only process up to 657 requests
rate=100 # Requests pr second, higher than 5 requires sudo
mode="--tcp-connect"

raw=$(nping -c $sent $mode -p $SERVER_PORT --rate $rate $SERVER_ADDRESS)

res=$(printf "%s" "$raw" | grep "Successful connections:")

received=$(printf "%s" "$res" | cut -d ' ' -f 8)

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
