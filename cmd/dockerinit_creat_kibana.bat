docker run --name kibana ^
-p 127.0.0.1:5601:5601 ^
-e "ELASTICSEARCH_HOSTS=http://127.0.0.1:9200" ^
-d kibana:8.17.0
  
pause