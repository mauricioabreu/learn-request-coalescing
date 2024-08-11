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

sub vcl_recv {
	unset req.http.X-Cache-Status;
}

sub vcl_hit {
	set req.http.X-Cache-Status = "hit";
	if (obj.ttl <= 0s && obj.grace > 0s) {
		set req.http.X-Cache-Status = "hit graced";
	}
}

sub vcl_miss {
	set req.http.X-Cache-Status = "miss";
}

sub vcl_pass {
	set req.http.X-Cache-Status = "pass";
}

sub vcl_pipe {
	set req.http.X-Cache-Status = "pipe uncacheable";
}

sub vcl_synth {
	set req.http.X-Cache-Status = "synth synth";
	set resp.http.X-Cache-Status = req.http.X-Cache-Status;
}

sub vcl_deliver {
	if (obj.uncacheable) {
		set req.http.X-Cache-Status = req.http.X-Cache-Status + " uncacheable" ;
	} else {
		set req.http.X-Cache-Status = req.http.X-Cache-Status + " cached" ;
	}

	set resp.http.X-Cache-Status = req.http.X-Cache-Status;
}
