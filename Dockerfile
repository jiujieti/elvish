FROM golang:alpine as builder
RUN apk update && \
    apk add --virtual build-deps make git
# Build Elvish
COPY . /go/src/src.elv.sh
RUN make -C /go/src/src.elv.sh get

FROM alpine
COPY --from=builder /go/bin/elvish /bin/elvish
RUN adduser -D elf
RUN apk update && apk add tmux man man-pages vim curl git
USER elf
WORKDIR /home/elf
CMD ["/bin/elvish"]
