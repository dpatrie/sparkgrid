# Running mongo
`docker run --name spark-mongo -d -p 27017:27017 mongo:4.1-xenial`

# Running S1
`go run s1/main.go`

# Running S2
`go run s2/main.go`

# What's missing
- Authentication in Mongo
- Docker compose for all of it
- Unique index for UUID in Mongo
- Better error messages
- Security for DELETE request in S1
- Concurency when importing records in PUT method

