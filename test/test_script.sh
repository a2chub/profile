#!/bin/bash

# イメージ名とタグを設定
IMAGE_NAME="test-image"
TAG="latest"

# Dockerイメージをビルド
echo "Dockerイメージをビルドしています..."
sudo docker build -t ${IMAGE_NAME}:${TAG} .

# ビルドが成功した場合のみ実行
if [ $? -eq 0 ]; then
    echo "コンテナを実行しています..."
    sudo docker run -it  -u 0 --rm ${IMAGE_NAME}:${TAG} /bin/bash
    
    # コンテナの実行が終了したら、イメージを削除
    echo "イメージを削除しています..."
    sudo docker rmi ${IMAGE_NAME}:${TAG}
else
    echo "ビルドに失敗しました"
    exit 1
fi
