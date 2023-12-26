#!/usr/bin/env bash

sudo apt-get update

# nginx
sudo apt-get -y install nginx
sudo service nginx start

# set up nginx server
sudo cp /workspace/.provision/nginx/nginx.conf /etc/nginx/sites-available/site.conf
sudo chmod 644 /etc/nginx/sites-available/site.conf
sudo ln -s /etc/nginx/sites-available/site.conf /etc/nginx/sites-enabled/site.conf
sudo service nginx restart

# make
sudo apt-get install make

# get updates
sudo apt-get update
sudo apt-get -y upgrade

# golang
wget  https://go.dev/dl/go1.20.2.linux-amd64.tar.gz 
sudo tar -xvf go1.20.2.linux-amd64.tar.gz
sudo mv go /usr/local

# setup golang env
export GOROOT=/usr/local/go 
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH 

# grpc and protobuf
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2