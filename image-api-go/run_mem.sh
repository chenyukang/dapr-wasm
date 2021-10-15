go build --tags "tensorflow image"; 

nohup ./image-api-go > infer.log &

sleep 3

pid=`pidof ./image-api-go`
echo "begin record $pid ..."
rm -rf ff.sqlite
/home/coder/.local/bin/procpath  record  -i 1 -r 60 -d ff.sqlite -p $pid
/home/coder/.local/bin/procpath plot -d ff.sqlite -q rss -p $pid -f rss.svg
echo "recorded $pid ...."
