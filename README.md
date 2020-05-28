# Proxlight


### the simpler "single node http reverse proxy"

[![](https://img.shields.io/badge/docker-glutamatt/proxlight-green.svg?logo=docker&longCache=true&style=flat-square)](https://hub.docker.com/r/glutamatt/proxlight/)

You can proxify the service http://example.com on proxy port `1234` with:

`docker run --network=host -it --rm glutamatt/proxlight http://example.com 0.0.0.0:1234`

You can proxify the service http://192.168.10.1:7890 on proxy port `1234` with:

`docker run --network=host -it --rm glutamatt/proxlight http://192.168.10.1:7890 0.0.0.0:1234`

with options ; example : throttling 10 rps + cors headers
`docker run --network=host -it --rm glutamatt/proxlight -cors -throttle 10 http://192.168.10.1:7890 0.0.0.0:1234`
