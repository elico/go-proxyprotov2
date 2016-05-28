# proxyproto

This library provides the `proxyproto` package which can be used for servers
listening behind HAProxy of Amazon ELB load balancers. Those load balancers
support the use of a proxy protocol (http://haproxy.1wt.eu/download/1.5/doc/proxy-protocol.txt),
which provides a simple mechansim for the server to get the address of the client
instead of the load balancer.

This library provides both a net.Listener and net.Conn implementation that
can be used to handle situation in which you may be using the proxy protocol.
Only proxy protocol version 1, the human-readable form, is understood.

The only caveat is that we check for the "PROXY " prefix to determine if the protocol
is being used. If that string may occur as part of your input, then it is ambiguous
if the protocol is being used and you may have problems.

# Trying to replicated the code at:
http://bazaar.launchpad.net/~squid/squid/3.5/revision/13570
- ConnStateData::parseProxyProtocolHeader()
- ConnStateData::parseProxy1p0()
- ConnStateData::parseProxy2p0()

# Examples
