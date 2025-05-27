# Fileserver
A simple http server for file upload and download from a server it's running on, access through it IP address and port 8080 , i.e `http://<IP Addresss>:8080`

## FileServer v0.1.0 Features
- **Large file support**: Up to 32GB uploads
- **File download capability**
- **Compact binary**: Small executable size
- **Cross-platform**: Linux, macOS, and Windows support

## Quick Setup Commands

**For Linux:**
```bash
curl -L -C - -o fileserver https://github.com/imrany/fileserver/releases/download/v0.1.0/fileserver-linux
chmod +x fileserver
./fileserver
```

**For macOS:**
```bash
curl -L -C - -o fileserver https://github.com/imrany/fileserver/releases/download/v0.1.0/fileserver-darwin
chmod +x fileserver
./fileserver
```

**For Windows:**
```bash
curl -L -C - -o fileserver.exe https://github.com/imrany/fileserver/releases/download/v0.1.0/fileserver-windows.exe
fileserver.exe
```

## Notes
This is a handy tool to quickly spin up a file server for transfers. Once you run it, it should start a local web server where you can upload/download files through a web interface.