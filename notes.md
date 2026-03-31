## From Docker.com/engine/security
4 major areas of docker security
    - Instrinsic security of the used kernel.
    - The attack surface of the Docker Daemon itself
    - Loopholes in the container ocnfiguration profile, either by default or when customizeed bu users.
    - The "hardeing" secuirt features of the kernel and how theu interact with containters.

## How Docker (Containers) works briefly
Containers are isolated environments for code to run. Containers are a more lightweight alternative to fullscale VMs. They work by letting the container use the same kernel as the host system, but virtualizing the user space. 

The main mechanism for isolation are the unix concepts of cgroups and namespaces. Cgroups specifies the resource allocation to a container (CPU, memory, network bandwidth etc...) And namespaces isolates the container from processes outside the namespace. They cannot see them nor affect them (both the host system and other containers.)

Namespaces are an old feature of the kernel, introduced in 2008 and is considered safe since it has had a long time for verification and bug finding.

## Networking
Each container gets their own network stack, meaning that they do not by default have priviledged access to other container's sockets. The host can setup and configure the traffic between containers aswell as external hosts. Docker comes pre configured with a network between all containers called the bridge, which is an abstraction that makes the interface between containers behave just as they were regular machines connected on an Ethernet switch.

## Control groups
Provides limiting and metrics of system resources. While they do not limit containers from interacting with eachothers resources. They do provide fairness of system resources and the limits it ensures mitigates denial of service attacks so that a single container cannot bring down the host system.

## The docker daemon uses Unix sockets
Since Docker 0.5.2, the Daemon API uses unix sockets instead of TCP sockets on local host. This is to mitigate cross site forgery attacks that TCP traffic is more prone to. It's also possible to use the unix permissions system for unix sockets which enables access control to the API. 

Only trusted users should be able to run the Daemon, since it has the possibility of creating new containers. The daemon usually runs in root mode (root-less mode is a less common possibility.) Containers can share file system's with the host system and there is nothing hindering sharing the host's / to the container, giving an attacker control of the hosts whole file system. 

## Kernel priviledges and root access
Containers are designed to have a limited capability set. Instead of giving a container root access, the host only gives what the container needs through configuration. And many aspects, such as networking, SSH access, log management is handled by the host system through the infrastructure around Docker. Basically, docker utilizes an allowlist of capabilities instead of a denylist. Meaning the focus is on what is a container allowed to do instead of plugging security holes. This doesn't usually affect regular applications, but considerably hinders adversaries were they to gain root access inside a container. 

The defaults capabilities of a container may be insufficient to the application suppsoed to run there. Or too much. The recommended best practice is to analyze what capabilities the application needs, and only allow those. 

## Rootless mode
Rootless mode allows docker to be installed and run without root priviledges. More over, containers are also ran without root. 

## Docker Trusted Content
Provides the possibility of using signatures for images published to registries. This feature allows verification of pulled software. A registry can contain images that are signed and those that are unsigned, for example images with different tags. Image consumers can enable DCT to only be able to pull signed images, which hardens their security profile. A producer of images creates a key set which consists of:
    - A root key which is stored offline, used to create...
    - Signing and verification keys that are actually used. They can be recreated from the root key which makes that key extra necessary to store in a secure backedup environment. A root key cannot simly be recovered. 
Commands that are effected by DCT are:
    - pull
    - build
    - create
    - push
    - run
