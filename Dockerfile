FROM golang:1.6.2

ADD . /go/src/github.com/theothertomelliott/github-watchlists

ADD ./scripts/ /scripts/

RUN chmod 755 /scripts/startup.sh

# Install revel and the revel CLI.
RUN go get github.com/revel/revel
RUN go get github.com/revel/cmd/revel

RUN go get github.com/theothertomelliott/github-watchlists/...

# Use the revel CLI to start up our application.
ENTRYPOINT exec /scripts/startup.sh

# Open up the port where the app is running.
EXPOSE 9000
