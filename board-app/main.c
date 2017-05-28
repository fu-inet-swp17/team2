#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <inttypes.h>

#include "net/ipv6/addr.h"
#include "net/sock/udp.h"
#include "net/gnrc/ipv6/netif.h"
#include "shell.h"

#include "./../smart_environment.h"

int main(void) {
    printf("Smart environment app on %s\n", RIOT_BOARD);

	// ff02::1 -> addr für link-local broadcast
	ipv6_addr_t addr;
	// char prefix[IPV6_ADDR_MAX_STR_LEN];
	// my_strncpy(prefix, , IPV6_ADDR_MAX_STR_LEN);
	ipv6_addr_from_str(&addr, "fe80::");

	char addr_str[IPV6_ADDR_MAX_STR_LEN];
	ipv6_addr_to_str(addr_str, &addr, IPV6_ADDR_MAX_STR_LEN);

	ipv6_addr_t* out = NULL;
	gnrc_ipv6_netif_find_by_prefix(&out, &addr);
	
	ipv6_addr_to_str(addr_str, out, IPV6_ADDR_MAX_STR_LEN);
	printf("own ipv6 addr: %s\n", addr_str);

	sock_udp_ep_t remote = SOCK_IPV6_EP_ANY;
	remote.port = SERVER_CONN_PORT;
	ipv6_addr_from_str((ipv6_addr_t *)&remote.addr.ipv6, "ff02::1");
	
	char intro_msg[CLIENT_INIT_MSG_LEN];
	snprintf(intro_msg, CLIENT_INIT_MSG_LEN, "%s %s", client_id, addr_str);
	
	if(sock_udp_send(NULL, intro_msg, CLIENT_INIT_MSG_LEN, &remote) < 0) {
		puts("Error sending intro message!\n");
		return EXIT_FAILURE;
	}
	
	uint8_t buf[SERVER_RESP_MSG_LEN];
	sock_udp_ep_t local = SOCK_IPV6_EP_ANY;
	local.port = CLIENT_PORT; // empfangssocket also auf eigenen Port stellen
	sock_udp_t sock;
	sock_udp_create(&sock, &local, NULL, 0);
	ssize_t res = sock_udp_recv(
			&sock,
			buf,
			sizeof(buf),
			SOCK_NO_TIMEOUT,
			NULL
	);
	
	if(res < 0) {
		printf("Error during receive!\n");
	} else {
		printf("%s\n", buf);
	}
	
    return EXIT_SUCCESS;
}
