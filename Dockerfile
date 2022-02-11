FROM golang:1.17.6-alpine as dev

# アップデートとgit, gccのインストール！！
RUN apk update && apk add git && apk add build-base
# appディレクトリの作成
RUN mkdir /go/app
# ワーキングディレクトリの設定
WORKDIR /go/app
# ホストのファイルをコンテナの作業ディレクトリに移行
ADD . /go/app