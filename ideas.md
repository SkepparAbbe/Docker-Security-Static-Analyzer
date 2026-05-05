## Limitations
* Ignore docker compose and swarm. Focus on dockerfiles and host docker configuration.

## User not defined
* (DockerFile) Check if USER is defined, if the deployment is run as root (if rootless mode isn't enabled)
* (Host configuration) Check the host's daemon.json to see if rootless mode or userns-remap (user namespaces) is enabled.

## breaking least privilege capabilities 
(Host configuration) Check if capabilities are more than they need to be.

## Check for volume exposure
* (DockerFile) Are bound mounts resonable? 
* (DockerFile) Is root on host bound to the container?
* (DockerFile) Is the docker socket (/var/run/docker.sock) mounted to the container? Gives control over the daemon to the container.

## DCT - image signatures
* (DockerFile) Check if used images are signed.
* (Host configuration) Check if DOCKER_CONTENT_TRUST is enabled on host

## Security leaks through add/copy
(DockerFile) Files added through  add/copy that are deleted through `COMMAND ... rm file` is deleted from the container but remains in the produced image.

## ":Latest"
(DockerFile) Using the latest version may seem good for security updates but if the upstream image dependency gets malicious code included in it, a new build of the container will include said code because the latest dependency image is pulled. Potentially warn about depenency images without hash digest, ensures that the image can't be tampered with after it's been published to the trusted registry.

## reading the configuration:
(Host configuration) `docker info` displays the current effective configuration:
- So for example check if rootless can be done through checking for rootless in the output
- Check that cgroupns is enabled
