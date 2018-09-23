# iron/go is the alpine image with only ca-certificates added
FROM iron/go

WORKDIR /app

# copy env file
COPY .env /app/

# Now just add the binary
ADD main /app/
ENTRYPOINT ["./main"]
