# sqs2gl
This application reads Graylog GELF-format messages from an AWS SQS queue and forwards them to a 
Graylog GELF UDP or GELF HTTP listener.

For security and reliability reasons, it is generally preferable to configure a 
Graylog GELF HTTP listener on localhost (127.0.0.1) and install sqs2gl on the
Graylog host. If the HTTP POST fails, the message will be left in the queue.

### Development Status
This is a beta release.

### Installation
1) Clone the repo and compile using "go build"
2) Copy the binary (sqs2gl) and config file (sqs2gl.conf) to /opt/sqs2gl. If you put it elsewhere you will need to update the .service file.
3) Copy sqs2gl.service to /etc/systemd/system/
4) Update the User and Group in sqs2gl.service if you do not wish to run as root
5) Update the configuration file (sqs2gl.conf)
6) Run 'systemctl daemon-reload'
7) Run 'systemctl start sqs2gl' to start the application

### Copyright
Copyright (c) 2021 Tenebris Technologies Inc. All rights reserved.

Please see the LICENSE file for additional information.