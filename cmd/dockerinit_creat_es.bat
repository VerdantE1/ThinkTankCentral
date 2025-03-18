docker run --name es ^
  -p 127.0.0.1:9200:9200 ^
  -e "discovery.type=single-node" ^
  -e "xpack.security.http.ssl.enabled=false" ^
  -e "xpack.license.self_generated.type=trial" ^
  -e "xpack.security.enabled=false" ^
  -e ES_JAVA_OPTS="-Xms84m -Xmx512m" ^
  --network elastic ^
  -d elasticsearch:8.17.0