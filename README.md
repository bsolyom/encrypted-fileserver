# encrypted-fileserver
Simple web fileserver with file encryption written in Go
Upload a file to the web server, then download it with the given code, with configurable expiration time (after which the files are deleted), maximum file size, etc. Enter a password for the uploaded files to encrypt them on the server (optional).
In the [env](./env) file you can edit some variables, for example, set a cert for the server

## Upload a file to the server
Just pick the file, give it a password (optional) and click on the upload button on the main page
![image](https://github.com/user-attachments/assets/011bca60-72ed-49a0-8c85-24722c3dbca9)

## Download a file from the server
If the file is not encrypted, you can download it like this: `http(s)://host.host/down/code` \
If the file is encrypted, you have to use the main page (or the terminal, see below) to decrypt the file and download it \
![image](https://github.com/user-attachments/assets/69a87be8-03f4-4b04-9467-9af1463f1711)

## Also, you can use terminal to do these, the commands are specified in the [index.html](./static/index.html)
