# Coalescence

This project aims to compare and analyze request coalescing behavior in Varnish and NGINX. It consists of a simple Go backend server, and configurations for both Varnish and NGINX as reverse proxies.

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
