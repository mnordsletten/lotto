sudo apt install -y iperf3
nohup iperf3 -s -p 12345 &>/dev/null &
