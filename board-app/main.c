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
	ipv6_addr_from_str((ipv6_addr_t *)&remote.addr.ipv6, "ff01::1");
	if(sock_udp_send(NULL, "Hello world!", sizeof("Hello world!"), &remote) < 0) {
		puts("Error sending\"Hello world!\"");
		return EXIT_FAILURE;
	}
	
	puts("Successfully send\"Hello world\"");

	// char line_buf[SHELL_DEFAULT_BUFSIZE];
    // shell_run(NULL, line_buf, SHELL_DEFAULT_BUFSIZE);

    return EXIT_SUCCESS;
}
