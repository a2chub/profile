#!/bin/bash

useradd -m -s /bin/bash atusi
usermod -aG sudo atusi
chown -R atusi:atusi /home/atusi
passwd -d atusi
su - atusi
