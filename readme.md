export RABBITMQ_SERVER=amqp://test:test@127.0.0.1:5672
export ES_SERVER=127.0.0.1:9200

RabbitMQ:

brew services start rabbitmq

DATA:

LISTEN_ADDRESS=127.0.0.1:8081 STORAGE_ROOT=./tmp/3 go run dataServer/dataServer.go


API:

LISTEN_ADDRESS=127.0.0.1:4396 go run apiServer/apiServer.go

curl 127.0.0.1:9200/metadata -XPUT -H 'content-Type:application/json' -d'{
"mappings":{
        "properties":{
            "name":{"type":"text","index":"false"},
            "version":{"type":"integer"},
            "size":{"type":"integer"},
            "hash":{"type":"text"}
        }
}
}'
echo -n "nmsl" | openssl dgst -sha256 -binary | base64

h(v2): cAPvsxZe1PR54zIESQy0BaxC1pYJIvaHSF3qEOZYYIo=

curl -v 127.0.0.1:4396/objects/test5 -XPUT -d"nmsl" -H "Digest:SHA-256=BvGiORYPpUbgKcmtFViPCiSuysp8X5Lhi65jMJX9w74="


chap4.

step1.
echo -n "chap4 xxx" | openssl dgst -sha256 -binary | base64
op6UprjPcQScCGU3MOg81OKO54eb0lfMS97dfnoxKtQ=
step2.
curl -v 127.0.0.1:4396/objects/test4_1 -XPUT -d"chap4 xxx" -H "Digest:SHA-256=HlR1AxKd9IPG+6WAPZk6285+FuB5PmAVeBQ9HbWRmeQ="

curl http://127.0.0.1:4396/locate/HlR1AxKd9IPG+6WAPZk6285+FuB5PmAVeBQ9HbWRmeQ= 