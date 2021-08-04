FROM golang:1.16-alpine

WORKDIR /app

COPY . .

RUN go get github.com/prometheus/client_golang/prometheus
RUN go get github.com/prometheus/client_golang/prometheus/promauto
RUN go get github.com/prometheus/client_golang/prometheus/promhttp
RUN go build q1.go
RUN go build q2.go
RUN go build q3.go
RUN go build q4.go
RUN go build q5.go
RUN go build q6.go

EXPOSE 2112

CMD [ "/app/server" ]
