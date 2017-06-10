#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <inttypes.h>
#include <stdbool.h>

#include "net/ipv6/addr.h"
#include "net/sock/udp.h"
#include "net/gnrc/ipv6/netif.h"
#include "shell.h"
#include "saul.h"
#include "saul_reg.h"
#include "net/gcoap.h"
#include "thread.h"
#include "xtimer.h"

#include "smart_environment.h"

#define PING_TIMEOUT 10

char ping_stack[THREAD_STACKSIZE_DEFAULT];
char coap_stack[THREAD_STACKSIZE_DEFAULT];

char intro_msg[CLIENT_INIT_MSG_LEN];
char own_addr[IPV6_ADDR_MAX_STR_LEN];

/*
 * #description: checks if string prefix is a prefix of string str
 * #requires: stdbool.h, string.h
 * #param[char* str]: string that should start with something
 * #param[char* prefix]: the prefix to check for
 * #return[bool]: true if prefix is a prefix of str, false else
 */
static inline bool startsWith(const char* str, const char* prefix) {
    char* substr = strstr(str, prefix);
    return substr == str;
}

/**********************************COAP STUFF**********************************/
coap_pkt_t pdu;
uint8_t coap_buff[GCOAP_PDU_BUF_SIZE];

// static ssize_t sensors_handler(coap_pkt_t* pdu, uint8_t *buf, size_t len);
// static ssize_t temp_handler(coap_pkt_t* pdu, uint8_t *buf, size_t len);
// static ssize_t humid_handler(coap_pkt_t* pdu, uint8_t *buf, size_t len);

// TODO Ressourcen, Listener und Handler für alle Sensoren anlegen

// Array mit den CoAP Ressourcen
// Die Ressourcen muessen in alphabetischer Reihnfolge eingefuegt werden!!!
/*static const coap_resource_t se_resources[] = {
    // ressource-path, ressource-type, response-handler
    {"/se-app/sensors", COAP_GET, sensors_handler},
    // {"/se-app/humid", COAP_GET, humid_handler},
    // {"/se-app/temp", COAP_GET, temp_handler},
};*/

// TODO sind die listener so korrekt initialisiert?
// listener für gcoap
/*
static gcoap_listener_t sensors_listener = {
    (coap_resource_t*)&se_resources[0],
    sizeof(se_resources) / sizeof(se_resources[0]),
    NULL
};

static gcoap_listener_t humid_listener = {
    (coap_resource_t*)&se_resources[1],
    sizeof(se_resources) / sizeof(se_resources[1]),
    NULL
};

static gcoap_listener_t temp_listener = {
    (coap_resource_t*)&se_resources[2],
    sizeof(se_resources) / sizeof(se_resources[2]),
    NULL
};*/

// handler der angesprochen werden sollte wenn, eine anfrage nach den Sensoren
// kommmt
/*
static ssize_t sensors_handler(coap_pkt_t* pdu, uint8_t *buf, size_t len) {
    gcoap_resp_init(pdu, buf, len, COAP_CODE_CONTENT);
    size_t payload_size = (size_t)snprintf(
        (char*)pdu->payload, 
        len,
        "%s",
        "temp-humid"
    );
    puts("responed");
    return gcoap_finish(pdu, payload_size, COAP_FORMAT_TEXT);
}

static ssize_t temp_handler(coap_pkt_t* pdu, uint8_t *buf, size_t len) {
    gcoap_resp_init(pdu, buf, len, COAP_CODE_CONTENT);
    // Temperatur-Sensorwert mit SAUL abfragen und in pdu->payload schreiben
    // Anzahl der Bytes speichern und an gcoap_finish übergeben
    size_t payload_size = 0;
    return gcoap_finish(pdu, payload_size, COAP_FORMAT_NONE);
}

static ssize_t humid_handler(coap_pkt_t* pdu, uint8_t *buf, size_t len) {
    gcoap_resp_init(pdu, buf, len, COAP_CODE_CONTENT);
    // Luftfeuchtigkeit-Sensorwert mit SAUL abfragen und in pdu->payload schreiben
    // Anzahl der Bytes speichern und an gcoap_finish übergeben
    size_t payload_size = 0;
    return gcoap_finish(pdu, payload_size, COAP_FORMAT_NONE);
}
*/
/**********************************COAP STUFF**********************************/

// Thread-Routine zum senden regelmäßiger Pings
static void* ping_handler(void* args) {
	(void)args;
	sock_udp_ep_t remote = SOCK_IPV6_EP_ANY;
    remote.port = SERVER_CONN_PORT;
    // ff02::1 -> addr für link-local broadcast
    ipv6_addr_from_str((ipv6_addr_t *)&remote.addr.ipv6, "ff02::1");
    
    
    snprintf(intro_msg, CLIENT_INIT_MSG_LEN, "%s %s", client_id, own_addr);
    
    while(true) {
		ssize_t res = sock_udp_send(
			NULL,
			intro_msg,
			CLIENT_INIT_MSG_LEN,
			&remote
		);
		
		if(res == -EAFNOSUPPORT) {
			fputs("Ping: EAFNOSUPPORT", stderr);
		} else if(res == -EHOSTUNREACH) {
			fputs("Ping: EHOSTUNREACH", stderr);
		} else if(res == -EINVAL) {
			fputs("Ping: EINVAL", stderr);
		} else if(res == -ENOMEM) {
			fputs("Ping: ENOMEM", stderr);
		} else if(res == -ENOTCONN) {
			fputs("Ping: ENOTCONN", stderr);
		} else if(res < 0) {
			fprintf(stderr, "Ping error: %d\n", res);
		}
		
		if(res < 0) {
			break;
		} else {
			xtimer_sleep(PING_TIMEOUT);
		}
	}
	
    return NULL;
}

// Thread-Routine für einen CoAP-Server
static void* coap_handler(void* args) {
	(void)args;
	sock_udp_ep_t local = SOCK_IPV6_EP_ANY;
	local.port = GCOAP_PORT;
	
	sock_udp_t sock;
	// TODO detalliertere Fehlerabfrage
	if(sock_udp_create(&sock, &local, NULL, 0) < 0) {
		puts("coap: error creating socket!");
		return NULL;
	}
	
	while(true) {
		sock_udp_ep_t remote;
		ssize_t msg_len = sock_udp_recv(
			&sock,
			coap_buff,
			GCOAP_PDU_BUF_SIZE,
			SOCK_NO_TIMEOUT, // ggf. ändern
			&remote
		);
		
		if(msg_len == -EADDRNOTAVAIL) {
			puts("coap: EADDRNOTAVAIL");
		} else if(msg_len == -EAGAIN) {
			puts("coap: EAGAIN");
		} else if(msg_len == -ENOBUFS) {
			puts("coap: ENOBUFS");
		} else if(msg_len == -ENOMEM) {
			puts("coap: ENOMEM");
		} else if(msg_len == -EPROTO) {
			puts("coap: EPROTO");
		/*} else if(msg_len == -ETIMEDOUT) {
			puts("coap: ETIMEDOUT");*/
		} else if(msg_len < 0) {
			printf("coap error: %d\n", msg_len);
		}

		if(msg_len < 0) {
			break;
		}
		
		ssize_t res = coap_parse(&pdu, coap_buff, msg_len);
		if(res < 0) {
			puts("coap: parse error!");
			break;
		}
	}	
	
	return NULL;
}

static inline void sensors_test(void) {
    puts("\n##############################################");
    puts("############## SAUL sensor test ##############\n");

    saul_reg_t* dev = saul_reg;

    while(1) {
        phydat_t res;
        int dim = saul_reg_read(dev, &res);
        printf("%s of type %s returns the following data: \n", dev->name, saul_class_to_str(dev->driver->type));
        phydat_dump(&res, dim);
        if (dev->next == NULL) {
            break;
        }
        dev = dev->next;
    };

    puts("##############################################\n");
}

int main(void) {
    printf("Smart environment app on %s\n", RIOT_BOARD);

    saul_reg_t* dev = saul_reg;

    if (dev == NULL) {
        puts("The SAUL registry does not contain any devices.");
    } else {
        printf("The SAUL registry contains the following devices: ");
        while(1) {
            printf("%s", dev->name);
            if (dev->next == NULL) {
                break;
            }
            printf(", ");
            dev = dev->next;
        };
        printf("\n");
    }

    sensors_test();

    ipv6_addr_t addr;
    ipv6_addr_from_str(&addr, "fe80::");

    ipv6_addr_to_str(own_addr, &addr, IPV6_ADDR_MAX_STR_LEN);

    ipv6_addr_t* out = NULL;
    gnrc_ipv6_netif_find_by_prefix(&out, &addr);
    
    ipv6_addr_to_str(own_addr, out, IPV6_ADDR_MAX_STR_LEN);
    printf("own ipv6 addr: %s\n", own_addr);

	// thread für pings
	thread_create(
		ping_stack,
		THREAD_STACKSIZE_DEFAULT,
		THREAD_PRIORITY_MAIN-1,
		0,
		ping_handler,
		NULL,
		"ping"
	);
	
	// thread für coap
	thread_create(
		coap_stack,
		THREAD_STACKSIZE_DEFAULT,
		THREAD_PRIORITY_MAIN-1,
		0,
		coap_handler,
		NULL,
		"coap"
	);
    
    // shell starten damit Application nicht terminiert
    char line_buf[SHELL_DEFAULT_BUFSIZE];
    shell_run(NULL, line_buf, SHELL_DEFAULT_BUFSIZE);
        
    return EXIT_SUCCESS;
}
