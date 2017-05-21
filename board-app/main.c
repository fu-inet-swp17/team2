#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "net/ipv6/addr.h"
#include "net/sock/udp.h"
#include "shell.h"

int main(void)
{
    puts("This is the board app.");

	sock_udp_ep_t remote = SOCK_IPV6_EP_ANY;
	remote.port = 2017;
	// ff02::1 -> addr für link-local broadcast
	ipv6_addr_from_str((ipv6_addr_t *)&remote.addr.ipv6, "ff02::1");
	
	
	if(sock_udp_send(NULL, "Hello world!", sizeof("Hello world!"), &remote) < 0) {
		puts("Error sending \"Hello world!\"");
		return EXIT_FAILURE;
	}
	
	puts("Successfully send \"Hello world\"");


    return EXIT_SUCCESS;
}
