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
#include "tmp006_saul.h"
#include "mpl3115a2_saul.h"

#include "smart_environment.h"

static inline void initialize_sensors(void) {
	mag3110_t mag3110;
	memset(&mag3110, 0, sizeof(mag3110));
	uint8_t mag3110_init_res = mag3110_init(&mag3110, mag3110_params);
	if (mag3110_init_res) {
		printf("Initialization of %s (%s) failed.\n", mag3110_name, mag3110_saul_name);
	} else {
		printf("%s (%s) successfully initialized.\n", mag3110_name, mag3110_saul_name);
	}

	/**
	* This needs to be outside the else clause because the SAUL registry requires it to be inside a persistent memory location.
	* If the device is not added to the SAUL registry or the addition fails, the memory taken up by the struct is freed later.
	*/
	saul_reg_t mag3110_saul = {
		NULL,
		&mag3110,
		mag3110_saul_name,
		&mag3110_saul_driver
	};

	if (!mag3110_init_res && !saul_reg_add(&mag3110_saul)) {
		printf("%s successfully added to SAUL registry.\n", mag3110_name);
	} else {
		free(&mag3110_saul);
	}


	mma8x5x_t mma8x5x;
	memset(&mma8x5x, 0, sizeof(mma8x5x));
	uint8_t mma8x5x_init_res = mma8x5x_init(&mma8x5x, mma8x5x_params);
	if (mma8x5x_init_res) {
		printf("Initialization of %s (%s) failed.\n", mma8x5x_name, mma8x5x_saul_name);
	} else {
		printf("%s (%s) successfully initialized.\n", mma8x5x_name, mma8x5x_saul_name);
	}

	/**
	* This needs to be outside the else clause because the SAUL registry requires it to be inside a persistent memory location.
	* If the device is not added to the SAUL registry or the addition fails, the memory taken up by the struct is freed later.
	*/
	saul_reg_t mma8x5x_saul = {
		NULL,
		&mma8x5x,
		mma8x5x_saul_name,
		&mma8x5x_saul_driver
	};

	if (!mma8x5x_init_res && !saul_reg_add(&mma8x5x_saul)) {
		printf("%s successfully added to SAUL registry.\n", mma8x5x_name);
	} else {
		free(&mma8x5x_saul);
	}


	tmp006_t tmp006;
	memset(&tmp006, 0, sizeof(tmp006));
	uint8_t tmp006_init_res = tmp006_init(&tmp006, 0, TMP006_I2C_ADDRESS, TMP006_CONFIG_CR_DEF);
	if (tmp006_init_res) {
		printf("Initialization of %s (%s) failed.\n", tmp006_name, tmp006_saul_name);
	} else {
		printf("%s (%s) successfully initialized.\n", tmp006_name, tmp006_saul_name);
	}

	/**
	* This needs to be outside the else clause because the SAUL registry requires it to be inside a persistent memory location.
	* If the device is not added to the SAUL registry or the addition fails, the memory taken up by the struct is freed later.
	*/
	saul_reg_t tmp006_saul = {
		NULL,
		&tmp006,
		tmp006_saul_name,
		&tmp006_saul_driver
	};

	if (!tmp006_init_res && !saul_reg_add(&tmp006_saul)) {
		printf("%s successfully added to SAUL registry.\n", tmp006_name);
	} else {
		free(&tmp006_saul);
	}


	mpl3115a2_t mpl3115a2;
	memset(&mpl3115a2, 0, sizeof(mpl3115a2));
	uint8_t mpl3115a2_init_res = mpl3115a2_init(&mpl3115a2, 0, MPL3115A2_I2C_ADDRESS, MPL3115A2_OS_RATIO_DEFAULT);
	if (mpl3115a2_init_res) {
		printf("Initialization of %s (%s) failed.\n", mpl3115a2_name, mpl3115a2_saul_name);
	} else {
		printf("%s (%s) successfully initialized.\n", mpl3115a2_name, mpl3115a2_saul_name);
	}

	/**
	* This needs to be outside the else clause because the SAUL registry requires it to be inside a persistent memory location.
	* If the device is not added to the SAUL registry or the addition fails, the memory taken up by the struct is freed later.
	*/
	saul_reg_t mpl3115a2_saul = {
		NULL,
		&mpl3115a2,
		mpl3115a2_saul_name,
		&mpl3115a2_saul_driver
	};

	if (!mpl3115a2_init_res && !saul_reg_add(&mpl3115a2_saul)) {
		printf("%s successfully added to SAUL registry.\n", mpl3115a2_name);
	} else {
		free(&mpl3115a2_saul);
	}
}

int main(void) {
    printf("Smart environment app on %s\n", RIOT_BOARD);

	initialize_sensors();

	if (saul_reg == NULL) {
		puts("The SAUL registry does not contain any devices.");
	} else {
		printf("The SAUL registry contains the following devices: ");
		while(1) {
			printf("%s", saul_reg->name);
			if (saul_reg->next == NULL) {
				break;
			}
			printf(", ");
			saul_reg = saul_reg->next;
		};
		printf("\n");
	}

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
		puts("Error sending intro message!");
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
		puts("Error during receive!");
	} else {
		printf("%s\n", buf);
	}
	
    return EXIT_SUCCESS;
}
