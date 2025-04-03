curl --http1.1 -v http://localhost:4221/ http://localhost:4221/
* Host localhost:4221 was resolved.
* IPv6: ::1
* IPv4: 127.0.0.1
*   Trying [::1]:4221...
* Connected to localhost (::1) port 4221
> GET / HTTP/1.1
> Host: localhost:4221
> User-Agent: curl/8.7.1
> Accept: */*
> 
* Request completely sent off
< HTTP/1.1 200 OK
< Server: SimpleHTTP/0.6 Python/3.12.9
< Date: Thu, 03 Apr 2025 07:31:54 GMT
< Content-type: text/html; charset=utf-8
< Content-Length: 340
< 
<!DOCTYPE HTML>
* Connection #0 to host localhost left intact
* Found bundle for host: 0x148b04500 [serially]
* Re-using existing connection with host localhost
> GET / HTTP/1.1
> Host: localhost:4221
> User-Agent: curl/8.7.1
> Accept: */*
> 
* Request completely sent off
< HTTP/1.1 200 OK
< Server: SimpleHTTP/0.6 Python/3.12.9
< Date: Thu, 03 Apr 2025 07:31:54 GMT
< Content-type: text/html; charset=utf-8
< Content-Length: 340
< 
<!DOCTYPE HTML>
* Connection #0 to host localhost left intact