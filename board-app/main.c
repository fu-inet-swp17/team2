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
#include "net/gcoap.h"

#include "mag3110_params.h"
#include "mag3110_saul.h"
#include "mma8x5x_params.h"
#include "mma8x5x_saul.h"
#include "tmp006_saul.h"
#include "mpl3115a2_saul.h"
#include "hdc1000_params.h"
#include "hdc1000_saul.h"
#include "tcs37727_params.h"
#include "tcs37727_saul.h"

#include "smart_environment.h"

/*
 * #description: checks if string prefix is a prefix of string str
 * #requires: stdbool.h, string.h
 * #param[char* str]: string that should start with something
 * #param[char* prefix]: the prefix to check for
 * #return[bool]: true wenn prefix is a prefix of str, false else
 */
static inline bool startsWith(const char* str, const char* prefix) {
	char* substr = strstr(str, prefix);
	return substr == str;
}


/**********************************COAP STUFF**********************************/
static ssize_t temp_handler(coap_pkt_t* pdu, uint8_t *buf, size_t len);
static ssize_t humid_handler(coap_pkt_t* pdu, uint8_t *buf, size_t len);

char own_addr[IPV6_ADDR_MAX_STR_LEN];
sock_udp_ep_t local = SOCK_IPV6_EP_ANY;
sock_udp_ep_t server_ep;

// TODO Ressourcen, Listener und Handler für alle Sensoren anlegen

// Die Ressourcen muessen in alphabetischer Reihnfolge eingefuegt werden!!!
static const coap_resource_t se_resources[] = {
	{"/se-app/humid", COAP_GET, humid_handler},
	{"/se-app/temp", COAP_GET, temp_handler},
};

static gcoap_listener_t humid_listener = {
	(coap_resource_t*)&se_resources[1],
	sizeof(se_resources) / sizeof(se_resources[1]),
	NULL
};

static gcoap_listener_t temp_listener = {
	(coap_resource_t*)&se_resources[0],
	sizeof(se_resources) / sizeof(se_resources[0]),
	&temp_listener
};


static ssize_t temp_handler(coap_pkt_t* pdu, uint8_t *buf, size_t len) {
	gcoap_resp_init(pdu, buf, len, COAP_CODE_CONTENT);
	// Temperatur Sensorwert mit SAUL abfragen und in pdu->payload schreiben
	// Anzahl der Bytes speichern und an gcoap_finish übergeben
	size_t payload_size = 0;
	return gcoap_finish(pdu, payload_size, COAP_FORMAT_NONE);
}

static ssize_t humid_handler(coap_pkt_t* pdu, uint8_t *buf, size_t len) {
	gcoap_resp_init(pdu, buf, len, COAP_CODE_CONTENT);
	// Luftfeutigkeit Sensorwert mit SAUL abfragen und in pdu->payload schreiben
	// Anzahl der Bytes speichern und an gcoap_finish übergeben
	size_t payload_size = 0;
	return gcoap_finish(pdu, payload_size, COAP_FORMAT_NONE);
}

/**********************************COAP STUFF**********************************/

static inline void sensors_initialize(void) {
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


	hdc1000_t hdc1000;
	memset(&hdc1000, 0, sizeof(hdc1000));
	uint8_t hdc1000_init_res = hdc1000_init(&hdc1000, hdc1000_params);
	if (hdc1000_init_res) {
		printf("Initialization of %s (%s) failed.\n", hdc1000_name, hdc1000_saul_name);
	} else {
		printf("%s (%s) successfully initialized.\n", hdc1000_name, hdc1000_saul_name);
	}

	/**
	* This needs to be outside the else clause because the SAUL registry requires it to be inside a persistent memory location.
	* If the device is not added to the SAUL registry or the addition fails, the memory taken up by the struct is freed later.
	*/
	saul_reg_t hdc1000_saul = {
		NULL,
		&hdc1000,
		hdc1000_saul_name,
		&hdc1000_saul_hum_driver
	};

	if (!hdc1000_init_res && !saul_reg_add(&hdc1000_saul)) {
		printf("%s successfully added to SAUL registry.\n", hdc1000_name);
	} else {
		free(&hdc1000_saul);
	}


	tcs37727_t tcs37727;
	memset(&tcs37727, 0, sizeof(tcs37727));
	uint8_t tcs37727_init_res = tcs37727_init(&tcs37727, tcs37727_params);
	if (tcs37727_init_res) {
		printf("Initialization of %s (%s) failed.\n", tcs37727_name, tcs37727_saul_name);
	} else {
		printf("%s (%s) successfully initialized.\n", tcs37727_name, tcs37727_saul_name);
	}

	/**
	* This needs to be outside the else clause because the SAUL registry requires it to be inside a persistent memory location.
	* If the device is not added to the SAUL registry or the addition fails, the memory taken up by the struct is freed later.
	*/
	saul_reg_t tcs37727_saul = {
		NULL,
		&tcs37727,
		tcs37727_saul_name,
		&tcs37727_saul_driver
	};

	if (!tcs37727_init_res && !saul_reg_add(&tcs37727_saul)) {
		printf("%s successfully added to SAUL registry.\n", tcs37727_name);
	} else {
		free(&tcs37727_saul);
	}
}

int main(void) {
    printf("Smart environment app on %s\n", RIOT_BOARD);

	sensors_initialize();

	saul_reg_t* saul_reg_cur = saul_reg;

	if (saul_reg_cur == NULL) {
		puts("The SAUL registry does not contain any devices.");
	} else {
		printf("The SAUL registry contains the following devices: ");
		while(1) {
			printf("%s", saul_reg_cur->name);
			if (saul_reg_cur->next == NULL) {
				break;
			}
			printf(", ");
			saul_reg_cur = saul_reg_cur->next;
		};
		printf("\n");
	}

	// TODO listener registrieren
	gcoap_register_listener(&temp_listener);
	gcoap_register_listener(&humid_listener);

	// ff02::1 -> addr für link-local broadcast
	ipv6_addr_t addr;
	// char prefix[IPV6_ADDR_MAX_STR_LEN];
	// my_strncpy(prefix, , IPV6_ADDR_MAX_STR_LEN);
	ipv6_addr_from_str(&addr, "fe80::");

	ipv6_addr_to_str(own_addr, &addr, IPV6_ADDR_MAX_STR_LEN);

	ipv6_addr_t* out = NULL;
	gnrc_ipv6_netif_find_by_prefix(&out, &addr);
	
	ipv6_addr_to_str(own_addr, out, IPV6_ADDR_MAX_STR_LEN);
	printf("own ipv6 addr: %s\n", own_addr);

	sock_udp_ep_t remote = SOCK_IPV6_EP_ANY;
	remote.port = SERVER_CONN_PORT;
	ipv6_addr_from_str((ipv6_addr_t *)&remote.addr.ipv6, "ff02::1");
	
	char intro_msg[CLIENT_INIT_MSG_LEN];
	snprintf(intro_msg, CLIENT_INIT_MSG_LEN, "%s %s", client_id, own_addr);
	
	if(sock_udp_send(NULL, intro_msg, CLIENT_INIT_MSG_LEN, &remote) < 0) {
		puts("Error sending intro message!");
		return EXIT_FAILURE;
	}
	
	uint8_t resp_buf[SERVER_RESP_MSG_LEN];
	sock_udp_ep_t local = SOCK_IPV6_EP_ANY;
	local.port = CLIENT_PORT; // empfangssocket also auf eigenen Port stellen
	sock_udp_t sock;
	sock_udp_create(&sock, &local, NULL, 0);
	do {
		ssize_t res = sock_udp_recv(
			&sock,
			resp_buf,
			sizeof(resp_buf),
			SOCK_NO_TIMEOUT,
			NULL
		);
		
		if(res < 0) {
			puts("Error during receive!");
			continue;
		}
		
	} while(!startsWith((const char*)resp_buf, server_id));
	
	char* se_addr_str = strchr((const char*)resp_buf, ' ');
	se_addr_str++; // addr_str zeigt sonst auf das Leerzeichen!
	printf("server found on: %s\n", se_addr_str);
	
	server_ep.family =  AF_INET6;
	server_ep.port = SERVER_POLL_PORT;
	ipv6_addr_from_str(
		(ipv6_addr_t *)&server_ep.addr.ipv6,
		se_addr_str
	);
	
	
	
    return EXIT_SUCCESS;
}
