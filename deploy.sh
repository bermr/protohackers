#!/bin/bash

SERVER_IP=$1;
SOLUTION=$2;

mkdir -p bin
go build -o ./bin/$SOLUTION.out ./$SOLUTION/$SOLUTION.go

scp -i ../chave-oracle.key ./bin/$SOLUTION.out ubuntu@$SERVER_IP:/home/ubuntu/projects/protohackers/$SOLUTION.out

ssh -tt -i ../chave-oracle.key -l ubuntu $SERVER_IP <<ENDSSH
  cd /home/ubuntu/projects/protohackers
  ./$SOLUTION.out
  exit
ENDSSH
