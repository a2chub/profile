#!/bin/bash

useradd -m -s /bin/bash atusi
usermod -aG sudo atusi
chown -R atusi:atusi /home/atusi
# Note: パスワードは手動で設定してください: passwd atusi
su - atusi
