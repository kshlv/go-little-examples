# TCP-UDP ping-pong

**Huge disclaimer here**: The presented code is half-working, really. It works on half of scenarios but badly need refactoring and tons of debugging and fixing.
- TCP seem to drop the connection after the first dial. So if you send segments with `nc` you'll have to restart it to send another one.
- Somehow `nc -u` is not working with this program. For sending UDP datagrams, you'll need to use something else. *(Or write a simple UDP client.)*

## What about the same port?

Usually, when you try to run a program that implies listening on a port that is already occupied, you'll get an error like this:
```
2020/06/23 22:58:48 listen udp4 :8080: bind: address already in use
exit status 1
```

But it's won't happen if you're trying to run a TCP listener and a UDP listener on the same port.

At first it looks very wrong, so let's sort it out:
- Any working application has to have a socket.
- All sockets are created with [`socket()`](https://man7.org/linux/man-pages/man2/socket.2.html).
- [`socket()`](https://man7.org/linux/man-pages/man2/socket.2.html) takes three parameters: `domain`, `type` and `protocol`, and returns a socket file descriptor.
    - For all IPv4 interactions the `domain` is `AF_INET`, let's skip the others.
    - Now, `type` specifies the *communication semantics*, which and also represents actual *protocol* we're going to use. Thus, for TCP the `SOCK_STREAM` is used, and for UDP - `SOCK_DGRAM`.
    - In all examples we found `protocol` is left `0`. *([*Assigned Internet Protocol Numbers*](https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml) looks quite legit though.)*
- After that the socket is binded to a certain port (*address* actually but whatever) with [`bind()`](https://man7.org/linux/man-pages/man2/bind.2.html). If the address is already binded then `EADDRINUSE` is returned. (*Ok, fine, we haven't find exactly **why** this happens but at least have found some [pieces of evidence](https://github.com/golang/go/blob/master/src/net/listen_test.go#L158) that it **shoud** work like this.*)

Anyways, to check that, we wrote a simple TCP-UDP ping-pong program. TCP and UDP listeners run in two goroutines and send messages to each other. Very thrilling.

