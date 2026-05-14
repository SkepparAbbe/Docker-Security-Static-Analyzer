# Example Dockerfile to test the parser

FROM alpine:latest

ADD testfile container_testfile
USER hello
RUN echo hello \
    && echo world