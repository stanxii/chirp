# Chirp
Chirp is a simplified Twitter clone backed by a REST API written in Go. 

Features:
- tweet and retweet posts
- like tweets and view what other users liked
- add tags in your tweets
- follow other users
- user signup and authenticaiton

## Setup Locally
### Requirements
- Golang v1.10 or later
- PostgreSQL 9.6 or later

### Installation
```shell
 git clone https://github.com/xiao-vincent/chirp.git
 cd chirp 
```
### Configuration
The appliation config can be changed in the [.config](./.config) JSON file
```json
{
  "port": 3000,
  "env": "dev",
  "pepper": "super-secret-pepper-string",
  "hmac_key": "super-secret-hmac-key",
  "database": {
    "host": "localhost",
    "port": 5432,
    "user": "vince",
    "password": "your-password",
    "name": "chirp_dev"
  }
}
```
## Running the application
In the command line, enter
```shell
go run *.go
```
or install the auto rebuild/rerun tool [Refresh](https://github.com/markbates/refresh) and enter
```shell
refresh
```
to rerun the applicatoin when files change. The refresh config file is defined at
[refresh.yml](./refresh.yml)

### Test Connection
Test the api with
```shell 
curl -i localhost:3000/ping
```
and you should get response
```
HTTP/1.1 200 OK
Pinging the server...Success!
```

