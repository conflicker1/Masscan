@echo off
masscan -oL output.txt -p443 0.0.0.0/0 --exclude 255.255.255.255 --rate 2147483647