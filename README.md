# Learn Request Coalescing

This project aims to compare and analyze request coalescing behavior in [Varnish](https://varnish-cache.org/) and [NGINX](https://nginx.org/). It consists of a simple Go backend server, and configurations for both Varnish and NGINX as reverse proxies.

---

What is request coalescing? It is a technique used to reduce the number of requests sent to the backend server. Have you heard about the [thundering herd problem](https://en.wikipedia.org/wiki/Thundering_herd_problem)?

When multiple requests are sent to the same resource at the same time, the webserver is instructed to send only one of the requests to the backend. It reduces the load on the backend service and improves the overall performance.

**It is important to note that I created this project for learning purposes, and the information here might be misleading or incorrect. If you find any mistakes, please let me know by opening an issue or pull request.**

## Table of Contents

1. [System Architecture](#system-architecture)
2. [Components](#components)
3. [Setup](#setup)
4. [Results Analysis](#results-analysis)

## System Architecture

Our toy system consists of three main components:
* An HTTP server writen in Go, designed to reply slow and streaming responses.
* Reverse Proxies: Varnish and NGINX configured to handle HTTP requests.

## Components

### Backend

The backend is a simple Go HTTP server that implements a `/stream` and `slow` endpoints. This endpoint simulates slow responses by streaming data with deliberate delays.

It listens on port 8080 and has two endpoints:
* `/stream`: Streams data to the client using chunked encoding.
* `/slow`: Delays the response for a fixed amount of time.

### Varnish

Varnish is a caching HTTP reverse proxy.

* Enable *streaming* by setting `do_stream`to `true`.
* Custom headers to track coalescing and caching behavior.

Varnish listens on port 9090 and forwards requests to the backend server.

### NGINX

NGINX is configured as an alternative reverse proxy to Varnish. It uses cache locking for request coalescing.

NGINX listens on port 9191 and forwards requests to the backend server.

## Setup

You need to have Docker and [just](https://github.com/casey/just) installed.

### Building

```sh
just build
```

### Running

```sh
just run
```

Running the command above will start the backend server, Varnish, and NGINX.

Now we can play with the `/stream` route, either by requesting the backend directly or through the reverse proxies.


## Results Analysis

A bit of explanation about the possible values of the `X-Cache-Status` header

* **HIT:** response was served from the cache.
* **MISS:** response was not found in the cache and was fetched from the backend/upstream.
* **UPDATING:** response is being updated in the cache (NGINX only).

For both Varnish and NGINX, any request sent to the `/stream` endpoint will be coalesced. The difference between them is how they handle the response to all clients waiting. Requests to Varnish will start streaming as soon as the backend starts sending data, while NGINX will wait for the backend to respond, then NGINX will write the data to the cache, and finally the response will be sent to all clients, all served from the cache.

Varnish sends `hit` in the `X-Cache-Status` header as soon as the first receives the response headers of the backend while NGINX sends `hit` only when the response is fully received and written to the cache.

For NGINX, I think we can have a similar behavior to Varnish by setting `proxy_buffering` to `off`. When `proxy_buffering` is off, NGINX will start sending the response to the client as soon as it receives the response headers from the backend. However, features like rate limiting and caching don't work.
