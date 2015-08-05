Docker Volume Management for Springpath
=======================================

springpath-volume-driver is a storage driver for
the clustered storage solution provided by
http://springpathinc.com

How it Works
============

Container Run -> Docker -> VolumeDriver Mount -> Container Started.
Container Stop -> Docker -> VolumeDriver Unmount -> Container Stopped.

Installation
------------

Installing the plugin requires docker being installed, and a golang
enviroment set up. On Ubuntu, this is provided by the golang-go
package. A version >= 1.4.2 is recommended.

This plugin requires the as yet unrelease version of the docker engine.

Install docker using the usual methods, and replace the docker engine
with the pre-release version.

`$ sudo wget https://test.docker.com/builds/Linux/x86_64/docker-1.8.0-rc2 -O /usr/bin/docker`

and restart the docker service.

Then install the docker plugin using

`$ go get github.com/springpath/springpath-docker-plugin`

and start it as root

`$ springpath-docker-plugin`

Currently the plugin does not daemonize.
