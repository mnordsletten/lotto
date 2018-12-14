# http testing (downloads file from load balancer)
# Script used for performing http requests (get file.txt from lotto-client3 10.100.0.150 and lotto-client4 10.100.0.160)

# Input values to cmd
sent=5
received=0

# first stress the load-balancer
sudo docker run --rm rcmorano/docker-hey -n 100 -c 50 http://10.100.0.30:90/1GB_file.txt

# Loop and attempt to download file 5 times. Controlling the size of the download for success
for i in $(seq 1 $sent); do
  out=$(wget -q -T 10 --tries 1 -O - http://10.100.0.30:90/1GB_file.txt | wc -c - | cut -d' ' -f 1)
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

if [ -z $success ]; then success=false; fi
if [ -z $sent ]; then sent=0; fi
if [ -z $received ]; then received=0; fi
if [ -z $rate ]; then rate=0; fi
if [ -z $raw ]; then raw=""; fi
jq \
  --argjson success $success \
  --argjson sent $sent \
  --argjson received $received \
  --argjson rate $rate \
  --arg raw "$raw" \
  '. |
  .["success"]=$success |
  .["sent"]=$sent |
  .["received"]=$received |
  .["rate"]=$rate |
  .["raw"]=$raw
  '<<<'{}'
