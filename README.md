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

### NGINX

NGINX is configured as an alternative reverse proxy to Varnish. It uses cache locking for request coalescing.
