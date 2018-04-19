#sudo rm -rf paic_data/

if [ "$1" = "up" ]; then
docker start $(docker ps -aq)
elif [ "$1" = "down" ]; then
docker stop $(docker ps -q)
elif [ "$1" = "rma" ]; then
docker rm $(docker ps -aq)
docker network prune
else
 echo "please start with <up|down|rma>"
fi
