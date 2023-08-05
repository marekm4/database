# database
Serverless database server over HTTP websocket

## Usage
### Docker compose
```
services:
  web:
    build: .
    environment:
      DATABASE_SOCKET: ws://db:8080/database
  db:
    image: "marekm43/database"
```
By default, it uses port `8080`, this can be changed by setting up `PORT` env. Dumps are stored at `/app/database.txt`.
See [websocket_test.go](websocket_test.go) for more examples.

### Queries
#### Get value
```
select username
```

#### Set or update value
```
update username john
```

#### Set or increment counter
```
increment money 100
```

#### Append element to collection
```
append orders pizza
```

## FAQ
### What is it?
It's a database server that is adapted to run as a low cost "stateless" web service in the cloud.

### Why does it exist?
Most cloud providers allow you to run HTTP stateless service for free, but database servers, even cheapest ones, are at least a few $ per month.

### How does it work?
1. On startup, it downloads database files from cost-efficient storage, like S3
2. Then it runs database engine on it
3. When service is going to be redeployed files are send back to S3

### Which cloud providers are supported?
For now it's GCP Cloud Run only, AWS and Azure are on the way.

### Is there any SDK?
Unfortunately not for now. However, it should be straightforward to implement a basic client.

### Should I use it in production?
No.
