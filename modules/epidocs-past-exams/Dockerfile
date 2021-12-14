FROM alpine:latest
RUN apk add git

WORKDIR /output

ENTRYPOINT ["git", "clone", "https://github.com/Epidocs/Past-Exams.git"]