# tempmonitor
Uses the linux sys file to monitor your servers temperature in a docker container!

# About
This project is designed to monitor 1 CPU's temperature on a linux server. Many have a file named "/sys/class/thermal/thermal_zone0/temp" which lets you monitor the physical temperature of your CPU. Doing this is important to make sure your server isn't overheating, but is also hard to automate, I want to help fix that with this small project. If you have multiple cpu's, you would change "thermal_zone0" to "thermal_zone1" or so in the path.

The code is (hopefully) pretty simple. You GET a simple HTTP server, it returns the temperature if you authenticate correctly, and optionally (mostly for integration with Uptime Kuma) lets you specify the "X-Temp-Expect" header, which I will explain in another section, but what that does is if the expectation you specify there is not matching the actual temperature, it returns a 417 status code, so you can easily check for that with uptime monitoring tools.

# Dependencies
Docker (recommend docker-compose too) or Go
A Linux computer/server with a /sys/class/thermal/thermal_zone0/temp file, or thermal_zone(x)
For installation:
An internet connection
(wget and unzip) or (git)

# How to run
Make a new folder somewhere, then go into it.
wget and unzip:
```
wget https://github.com/ByteAfterlife/tempmonitor/archive/refs/heads/main.zip && unzip main.zip && cd tempmonitor-main
```
git:
```
git clone https://github.com/ByteAfterlife/tempmonitor.git && cd tempmonitor
```
Then, open main.go with your favorite text editor, and change
```
const (
    username = "ut"
    password = "1Xbq59yFD6toT5Y3HSLPU8kB4R88c95JHnKw0kpN3cxbML5VGSwTSiOqz6qEZuFH"
)
```
To have a different username and password for security. By default, those are the credentials.
Next part depends if you want to use Docker or not. If yes:
```
docker build -t tempmonitor .
docker-compose up -d || docker compose up -d
```
And if you do not want to use thermal_zone0, open docker-compose.yml and change the line:
```
      - /sys/class/thermal/thermal_zone0/temp:/temp:ro
```
Where it says thermal_zone0 to whatever CPU number you want
If you do not have docker compose, you should, it's pretty good, but you can run this instead if you do not have compose:
```
docker build -t tempmonitor .
docker run --name tempmonitor -p 8080:8080 --restart unless-stopped -v /sys/class/thermal/thermal_zone0/temp:/temp:ro tempmonitor
```

If you do not want to use docker:
(replace thermal_zone0 with whatever other CPU you want to monitor in the first command)
```
sed -i 's|ioutil.ReadFile("/temp")|ioutil.ReadFile("/sys/class/thermal/thermal_zone0/temp")|' main.go && go build -o server main.go && ./server
```
To make the server run on system startup, run:
```
sudo tee /etc/systemd/system/server.service > /dev/null <<EOF
[Unit]
Description=tempmonitor
After=network.target

[Service]
Type=simple
ExecStart=$(pwd)/server
Restart=on-failure
WorkingDirectory=$(pwd)
User=$(whoami)

[Install]
WantedBy=multi-user.target
EOF
```

You have now finished installation! Next, to test, you can run in the same terminal:
```
curl -u (username):(password) localhost:8080
```
To test X-Temp-Expect, try:
```
curl -v -u (username):(password) -H "X-Temp-Expect: <75" localhost:8080
```
If you see "< HTTP/1.1 200 OK" Then your server's temperature is less than 75 degrees celsius (which is good!)<br>
If you see "< HTTP/1.1 417 Expectation Failed" Then your server's temperature is over 75 celsius...

You can also test if you can access it externally (not just in your terminal on the same server) like this:
```bash
echo "http://$(echo $SSH_CONNECTION | awk '{print $3}'):8080"
```
Which should print something like "http://192.168.56.10:8080" or "http://192.168.1.19:8080"
Open that link in a browser (in many terminals, you can just hold ctrl and left click it) and you should see your server temperature!

# X-Temp-Expect Format
X-Temp-Expect gives you a 417 if your condition is not matched, and a 200 if it is. Examples:
```
X-Temp-Expect: <75 : Will only give you a 200 status code if your server's temperature is below 75. If it is exactly 75 or more, it will return the code 417
X-Temp-Expect: <=75 : Will only give you a 200 status code if your server's temperature is below 75 or exactly 75. If it is more than 75, it will return the code 417
X-Temp-Expect: >75 : Will only give you a 200 status code if your server's temperature is above 75. If it is exactly 75 or less, it will return the code 417
X-Temp-Expect: >=75 : Will only give you a 200 status code if your server's temperature is above 75 or exactly 75. If it is less than 75, it will return the code 417
```
