# Persistence
HTTP/1.1 defaults to the use of "persistent connections", allowing multiple requests and responses to be carried over a single connection. HTTP implementations SHOULD support persistent connections.

- If the `close` connection option is present, the connection will **not persist** after the current response; else,
- If the received protocol is `HTTP/1.1` (or later), the connection **will persist** after the current response; else,
- If the received protocol is `HTTP/1.0`, the `keep-alive` connection option is present, either the recipient is not a proxy or the message is a response, and the recipient wishes to honor the HTTP/1.0 "keep-alive" mechanism, the connection **will persist** after the current response; otherwise,
- The connection will **close** after the current response.

We don't support HTTP/1.0 so ideally we shouldn't send the `keep-alive` header. 
All requests implicitly default to `persistent`. (But MDN says: *persistence is the default, and the header is no longer needed (but it is often added as a defensive measure against cases requiring a fallback to HTTP/1.0).*), so we can try sending the `Connection: keep-alive` header.

## Stage breakdown
1. Send multiple sequential requests on single connection (*implicit persistence*)
```
Conn1
R1 -> 
 <- R1

... wait 2s

R2 -> 
 <- R2
```
2. Send multiple sequential requests on individual parallel connections (*implicit persistence*)
```
Conn1   x   Conn2
R1 ->       R1 ->

 <- R1       <- R1

... wait 2s

R2 ->       R2 ->

 <- R2       <- R2
```
3. Send `close` header from client: Server also should send `close` and close conn
```
Conn1
R1 (/get-token, (Connection: close))
 <- R2 (no-token) (Connection: close)
```
4. Send `close` header from client: Server also should send `close` and close conn, 4XX shouldn't affect open session (Connection should be treated as session, test using some session token)
```
Conn1
R1 (token) -> 
 <- R1 (processed)
R2 (/error) -> 
 <- R2 (404)
R3 (/get-token, (Connection: close))
 <- R3 (token) (Connection: close)

Conn2
R1 (/get-token, (Connection: close))
 <- R1 (no-token) (Connection: close)
```
5. Backward compatibility (`Connection: keep-alive`) [1](https://www.rfc-editor.org/rfc/rfc2068#section-19.7.1.1) (We send HTTP/1.1)
```
Conn1
R1 /set?seed=foo (Connection: keep-alive) ->
 <- R1
R2 (Get /random) (Connection: close)
 <- R2 (Connection: close)
```
6. Timeout conn from server side (5s) (can use something like `Keep-Alive: timeout=5`) (We send HTTP/1.1)
```
Conn1
R1 /set?seed=foo (Connection: keep-alive) (Keep-Alive: timeout=5) ->
 <- R1
R2 (Get /random) ->
 <- R2 

... sleep 5s

R3 ->
connection should be closed
```

# Pipelining

Possibly can add a couple of stages for Pipelining. 
Pipelining should ideally be done with Idempotent methods.
We would first need to add support for `PUT`, `HEAD`, `DELETE`. [2](https://www.rfc-editor.org/rfc/rfc2068#section-9.1.1) `GET` is not technically idempotent, but it is *safe*. We could use `GET` if we were so inclined. (wiki says `GET`s are always usable in pipelines, RFCs disagree and make it very clear) [3](https://www.rfc-editor.org/rfc/rfc9112.html#section-9.3.2)

## Stage breakdown

1. Basic Pipelining (Multiple requests sent without waiting for responses)
```
Conn1
R1,R2,R3 -> 
 <- R1
 <- R2
 <- R3
```

2. Mixed Pipelined and Non-Pipelined Requests
```
Conn1
R1 ->
 <- R1
R2,R3,R4 ->
 <- R2
 <- R3
 <- R4
R5 ->
 <- R5
```

3. Error Handling in Pipelined Requests (Preserving response order)
```
Conn1
R1,R2(invalid),R3 ->
 <- R1
 <- R2(400 Bad Request)
 <- R3
```

4. Partial Response Processing (Client processes responses as they arrive)
```
Conn1
R1,R2,R3 ->
 <- R1 (client processes)
 <- R2 (client processes)
 <- R3 (client processes)
```

5. Head-of-Line Blocking Test (Slow request impacts pipeline)
```
Conn1
R1(fast),R2(slow),R3(fast) ->
 <- R1
 ... wait for R2 processing ...
 <- R2
 <- R3
```

6. Pipeline with Connection Close
```
Conn1
R1,R2,R3(Connection: close) ->
 <- R1
 <- R2
 <- R3 (Connection: close)
connection closed
```

7. Timeout During Pipelined Processing
```
Conn1
R1,R2(timeout-trigger),R3 ->
<- R1
... server timeout on R2 ...
<- R2(408 Request Timeout, Connection: close)
connection closed (R3 is never processed)
```


Ref: https://en.wikipedia.org/wiki/HTTP_persistent_connection  
https://www.rfc-editor.org/rfc/rfc9112.html  
https://www.rfc-editor.org/rfc/rfc2068  
https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Connection_management_in_HTTP_1.x