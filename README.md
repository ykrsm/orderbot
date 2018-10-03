# orderbot

# Set up interactive components for localhost
```
ngrok http 3000
```
and set it on slack interavtive

# Compile for linux
```
dep ensure
env GOOS=linux GOARCH=386 go build -o main
```
# Start docker
```
docker-compose build
docker-compose up
```

# Feeling lazy (one liner)
```
env GOOS=linux GOARCH=386 go build -o main && docker-compose up --build
```