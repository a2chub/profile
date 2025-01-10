FROM ubuntu:22.04

RUN apt-get update && apt-get install -y curl

RUN apt-get install -y sudo
RUN apt-get install -y vim
RUN apt-get install -y git

COPY 1_init_user.sh /1_init_user.sh
RUN chmod +x /1_init_user.sh
ENTRYPOINT ["/1_init_user.sh"]

COPY 2_git_clone.sh /home/atusi/2_git_clone.sh
WORKDIR /home/atusi
RUN chmod +x /home/atusi/2_git_clone.sh
RUN /home/atusi/2_git_clone.sh

CMD ["/bin/bash"]
