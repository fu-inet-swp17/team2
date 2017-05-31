#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <inttypes.h>

#include "net/ipv6/addr.h"
#include "net/sock/udp.h"
#include "net/gnrc/ipv6/netif.h"
#include "shell.h"
#include "saul.h"
#include "saul_reg.h"

#include "mag3110_params.h"
#include "mag3110_saul.h"
#include "mma8x5x_params.h"
#include "mma8x5x_saul.h"

#include "smart_environment.h"

static inline void initialize_sensors(void) {
	mag3110_t mag3110;
	memset(&mag3110, 0, sizeof(mag3110));
	if (!mag3110_init(&mag3110, mag3110_params)) {
		printf("mag3110 (magnetometer) successfully initialized.\n");
	}

	saul_reg_t mag3110_saul = {
		NULL,
		&mag3110,
		"mag3110 (magnetometer)",
		&mag3110_saul_driver
	};

	int mag3110_saul_add = saul_reg_add(&mag3110_saul);

	if (!mag3110_saul_add) {
		printf("mag3110 (magnetometer) successfully added to SAUL registry.\n");
	}


	mma8x5x_t mma8x5x;
	memset(&mma8x5x, 0, sizeof(mma8x5x));
	if (!mma8x5x_init(&mma8x5x, mma8x5x_params)) {
		printf("mma8652 (accelerometer) successfully initialized.\n");
	}

	saul_reg_t mma8x5x_saul = {
		NULL,
		&mma8x5x,
		"mma8652 (accelerometer)",
		&mma8x5x_saul_driver
	};

	int mma8x5x_saul_add = saul_reg_add(&mma8x5x_saul);

	if (!mma8x5x_saul_add) {
		printf("mma8652 (accelerometer) successfully added to SAUL registry.\n");
	}
}

int main(void) {
    printf("Smart environment app on %s\n", RIOT_BOARD);

	initialize_sensors();

	printf("The SAUL registry contains the following devices: ");
	do {
		printf("%s", saul_reg->name);
		if (saul_reg->next != NULL) {
			printf(", ");
			saul_reg = saul_reg->next;
		}
	} while (saul_reg->next != NULL);
	printf("\n");

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
