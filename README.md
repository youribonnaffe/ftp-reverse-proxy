# ftp-reverse-proxy: HTTP reverse proxy to access files hosted on a FTP server.

Want to download files from a FTP server with HTTP? Here you go!
It can also be used to upload files 🤓.

## Usage:

> ./ftp-reverse-proxy -port 8080 -target ftp://user:pwd@host:port

## Test it:

Run a FTP server locally (here with Docker):

> docker  run -ti --rm --name ftp -p 21:21 -e FTP_USER_NAME=bob -e FTP_USER_PASS=azerty -p 30000-30009:30000-30009 -e "PUBLICHOST=localhost" stilliard/pure-ftpd

Start the reverse proxy:

> ./ftp-reverse-proxy -port 8080 -target ftp://bob:azerty@localhost:21

Upload a file to the FTP server:

> curl -v -XPOST http://localhost:8080/file.txt -d @file.txt

Download a file from the FTP server:

> curl -v http://localhost:8080/file.txt

## Roadmap

This is a tool I built mostly to practice a bit with Go. It can be useful for applications that support HTTP to retrieve
files and where you don't want to add FTP support to it.

Feel free to open an issue to share you use case if you find limitations and issues (there are many!).

For instance, a connection to the FTP is made for each request so you will quickly hit limits such as the number of parallel 
clients and such.

## FAQ

- Does it support SFTP?

_Nop, feel free to open an issue for that._

- How can I proxy to multiple FTP servers?

_You can run the same binary multiple times, on different ports and put an HTTP proxy in front._

- Does it support HTTPS?

_Nop, you can however put in behind a real proxy such as Nginx for instance that will expose it with
HTTPS._
