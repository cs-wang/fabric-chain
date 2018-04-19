#sudo rm -rf paic_data/

if [ "$1" = "up" ]; then
#docker start $(docker ps -aq)
docker start zookeeper0
docker start zookeeper1
docker start zookeeper2
docker start kafka0
docker start kafka1
docker start kafka2
docker start kafka3

docker start orderer0.example.com

docker start peer0.org1.example.couchdb
docker start peer0.org1.example.com
docker start peer0.org2.example.couchdb
docker start peer0.org2.example.com
docker start peer0.org3.example.couchdb
docker start peer0.org3.example.com
docker start peer0.org4.example.couchdb
docker start peer0.org4.example.com

elif [ "$1" = "down" ]; then
#docker stop $(docker ps -q)

docker stop zookeeper0
docker stop zookeeper1
docker stop zookeeper2
docker stop kafka0
docker stop kafka1
docker stop kafka2
docker stop kafka3

docker stop orderer0.example.com
exit
docker stop peer0.org1.example.couchdb
docker stop peer0.org1.example.com
docker stop peer0.org2.example.couchdb
docker stop peer0.org2.example.com
docker stop peer0.org3.example.couchdb
docker stop peer0.org3.example.com
docker stop peer0.org4.example.couchdb
docker stop peer0.org4.example.com

else
 echo "please start with <up|down>"
fi
