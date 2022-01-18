# miniroot

> docker but simpler

In contrast to Docker, miniroot **does not**:
- have a concept of images and layers
- isolate networks
- handle volumes or bind mounts
- have a central daemon

It's just a simple tool that utilises namespaces to run applications in a semi-isolated environment.

### Usage

- get rootfs (you can use `docker export container_name > export.tar` for that)
- extract it somewhere
- run miniroot; example:
```
miniroot -root "/opt/something" -workdir "/usr/src/app" -init "/usr/local/bin/npm start"
```
