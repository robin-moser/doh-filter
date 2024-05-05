# DOH filter API

This is a simple API that queries and caches the DOH blocklist from [DoH-IP-blocklists](https://github.com/dibdot/DoH-IP-blocklists) and filters out localhost IP addresses. This is needed until the localhosts are removed from the blocklist, see [pull request](https://github.com/dibdot/DoH-IP-blocklists/pull/10).

The server will start listening on port 8080. You can access the filtered IPv4 list at `/ipv4` and the IPv6 list at `/ipv6`.
