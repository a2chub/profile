sudo apt update
sudo apt install docker.io
sudo systemctl start docker
sudo systemctl enable docker

# 例) 
#sudo docker run hello-world

#sudo docker run -d -p 80:80 --name my-nginx nginx
#sudo docker exec -it my-nginx /bin/bash
#sudo docker stop my-nginx
#sudo docker rm my-nginx

