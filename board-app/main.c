#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <inttypes.h>

#include "net/ipv6/addr.h"
#include "net/sock/udp.h"
#include "net/gnrc/ipv6/netif.h"
#include "shell.h"

int main(void) {
    printf("Smart Envoirement App on %s\n", RIOT_BOARD);

	// ff02::1 -> addr für link-local broadcast
	ipv6_addr_t addr;
	ipv6_addr_from_str(&addr, "fe80::");

	char addr_str[IPV6_ADDR_MAX_STR_LEN];
	ipv6_addr_to_str(addr_str, &addr, IPV6_ADDR_MAX_STR_LEN);
	printf("prefix: %s\n", addr_str);

	ipv6_addr_t* out = NULL;
	printf("DEBUG\n");
	gnrc_ipv6_netif_find_by_prefix(&out, &addr);
	
	ipv6_addr_to_str(addr_str, out, IPV6_ADDR_MAX_STR_LEN);
	printf("prefix: %s\n", addr_str);
	

	/*
	sock_udp_ep_t remote = SOCK_IPV6_EP_ANY;
	remote.port = 2017;
	ipv6_addr_from_str((ipv6_addr_t *)&remote.addr.ipv6, "ff02::1");
	
	size_t intro_msg_len = sizeof("swp-se-app") + IPV6_ADDR_MAX_STR_LEN; 
	char* intro_msg[intro_msg_len] = "swp-se-app";
	
	if(sock_udp_send(NULL, intro_msg, intro_msg_len, &remote) < 0) {
		puts("Error sending intro message!\n");
		return EXIT_FAILURE;
	}
	*/

    return EXIT_SUCCESS;
}
