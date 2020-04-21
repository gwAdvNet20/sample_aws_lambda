# ETL 

This application will take in a base64 encoded file and return back the browser count


To run this do the following on a mac:
```
curl -H "content-type: application/json" https://uwi10uh8wd.execute-api.us-east-1.amazonaws.com/default/etl2 -d '{"filename":"log_b.txt","contents":"'$(cat log_a.txt | base64)'"}'
```

