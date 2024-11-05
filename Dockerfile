FROM golang:1.21 as builder

ENV ROOT=/go/src/app
RUN mkdir /built
WORKDIR ${ROOT}
RUN apt update
COPY go.mod go.sum ./
RUN go mod download

COPY . ${ROOT}
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /built/khi cmd/kubernetes-history-inspector/main.go
RUN mkdir /built/data
COPY ./resources /built/resources
COPY ./dist /built/web

FROM alpine
ENV ROOT=/go/src/app
WORKDIR ${ROOT}
COPY --from=builder /built ${ROOT}
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV GOMEMLIMIT=10000MiB
EXPOSE 8080
CMD /go/src/app/khi --host=0.0.0.0 --temporary-folder=/ --data-destination-folder=/ --iam-token=${IAM_TOKEN} --access-token=${GCP_ACCESS_TOKEN} --fixed-project-id=${KHI_FIXED_PROJECT_ID} --ga-labels=${KHI_GA_LABELS}
