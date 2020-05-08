FROM node:7-alpine

RUN apk update
RUN apk add openssh-client

EXPOSE 3000
