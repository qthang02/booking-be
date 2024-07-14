FROM ubuntu:latest
LABEL authors="thang"

ENTRYPOINT ["top", "-b"]