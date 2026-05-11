# Example Dockerfile to test the parser

FROM alpine:latest

USER hello
RUN echo hello \
    && echo world