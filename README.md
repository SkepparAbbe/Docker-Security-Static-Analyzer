# Docker-Security-Static-Analyzer

This is a project in the Chalmers course TDA602 - Language-based Security by Albin Skeppstedt and [Isak Söderlind @executem](https://github.com/executem)

We built a prototype Dockerfile and Docker daemon configuration static analyzer that finds some cases where best practices aren't followed and prints helpful educational issues to stdout.

## Requirements:
- go
- docker
## Build
```sh
cd Docker-Security-Static-Analyzer/
go build
chmod +x docker-ssa # Optional  
```

## Run the program
```sh
./docker-ssa # Can be used in the same dir as a Dockerfile or
./docker-ssa path/to/Dockerfile # Take a Dockerfile as optional argument
```
