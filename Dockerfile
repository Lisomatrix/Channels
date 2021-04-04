FROM golang:stretch
WORKDIR /
COPY .  .
RUN CGO_ENABLED=0
RUN mkdir build
RUN go build -o build/
ENV PORT=8090
EXPOSE 8090
CMD [ "build/channels" ]