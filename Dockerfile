FROM golang:1.21-alpine as builder
WORKDIR /
COPY . ./
RUN go mod download


RUN go build -o /vancouver-conditions-lambda


FROM alpine
COPY --from=builder /vancouver-conditions-lambda .


EXPOSE 80
CMD [ "/vancouver-conditions-lambda" ]