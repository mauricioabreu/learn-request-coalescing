vcl 4.1;

probe healthcheck {
    .url = "/healthcheck";
    .interval = 10s;
    .timeout = 3s;
}

backend default {
    .host = "backend";
    .port = "8080";
    .probe = healthcheck;
}

sub vcl_backend_response {
    set beresp.do_stream = true;
    set beresp.ttl = 10s;
}

sub vcl_deliver {
    if (obj.hits > 0) {
        set resp.http.X-Cache = "HIT";
    } else {
        set resp.http.X-Cache = "MISS";
    }
}
