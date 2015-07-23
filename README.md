Docker Volume Management for Springpath
=======================================

springpath-volume-driver is a storage driver for
the clustered storage solution provided by
http://springpathinc.com

How it Works
============

Container Run -> Docker -> VolumeDriver Mount -> Container Started.
Container Stop -> Docker -> VolumeDriver Unmount -> Container Stopped.


