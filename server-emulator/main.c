#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>

#include "net/sock/udp.h"
#include "shell.h"

uint8_t buf[128];

int main(void) {
	printf("This is the server-emulator on %s.\n", RIOT_BOARD);

	sock_udp_ep_t local = SOCK_IPV6_EP_ANY;
	
	/*{
		.family = AF_INET6,
		.netif = SOCK_ADDR_ANY_NETIF
	};*/
	
	sock_udp_t sock;
	local.port = 2017;
	if(sock_udp_create(&sock, &local, NULL, 0) < 0) {
		puts("Error creating UDP sock");
		return EXIT_FAILURE;
	} else {
		puts("Successfull createt UDP socket");
	}
	
	while(true) {
		ssize_t res;
		if((res = sock_udp_recv(&sock, buf, sizeof(buf), SOCK_NO_TIMEOUT, NULL)) >= 0) {
			printf("%.128s\n", buf);
		}
	}
	
	// TODO herausfinden was genau 'shell_run' tut
	char line_buf[SHELL_DEFAULT_BUFSIZE];
	shell_run(NULL, line_buf, SHELL_DEFAULT_BUFSIZE);
	
    return EXIT_SUCCESS;
}
