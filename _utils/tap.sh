source_socket="/run/nordlayer/nordlayer.sock"
socket="nordlayer.sock"
port=80
dump="/tmp/capture2.pcap"

socat -t100 "TCP-LISTEN:${port},reuseaddr,fork" "UNIX-CONNECT:${source_socket}" &
socat -t100 "UNIX-LISTEN:${socket},mode=777,reuseaddr,fork" "TCP:localhost:${port}" &

# Record traffic
tshark -i lo -w "${dump}" -F pcapng "dst port ${port} or src port ${port}"
