#include <stdio.h>
#include <stddef.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>

#include "periph/gpio.h"
#include "net/ipv6/addr.h"
#include "net/sock/udp.h"
#include "net/gnrc/ipv6/netif.h"
#include "shell.h"
#include "board.h"
#include "thread.h"

#include "smart_environment.h"

#define MAX_BOARD_NUM 5
#define TIMEOUT (uint32_t)2000000

typedef struct SensorBoard {
	bool init;
	char addr[IPV6_ADDR_MAX_STR_LEN];
	sock_udp_ep_t ep;
} SensorBoard;

// gefährlich! sollte eigentlich gemutext werden
SensorBoard boards[MAX_BOARD_NUM];

char own_addr[IPV6_ADDR_MAX_STR_LEN];
char connect_stack[THREAD_STACKSIZE_MAIN];

uint8_t conn_buf[CLIENT_INIT_MSG_LEN];

uint8_t req_buf[128];

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

/*
 * #description: strncpy in Vernünftig, garantiert Nullterminator
 * #param[char* dst]: Speicher in den rein kopiert wird
 * #param[const char* src]: String der in 'dst' kopiert werden soll
 * #param[size_t num]: max. Anzahl Zeichen die in dst geschrieben werden darf 
 * #return[size_t]: Anzahl geschriebener Zeichen
 */
static inline size_t my_strncpy(char* dst, const char* src, size_t num) {
	size_t i;
	for(i=0; src[i] && i<(num-1); i++) {
		dst[i] = src[i];
	}
	
	dst[i] = 0;
	return i;
}

void button_handler(void* args) {
	// puts("Requesting data from boards...");
	
	for(size_t i=0; i<MAX_BOARD_NUM; i++) {
		if(!boards[i].init) {
			printf("Slot %u not used.\n", i);
			continue;
		}
		
		char request[] = "request";
		ssize_t res = sock_udp_send(
			NULL,
			request, 
			sizeof("request"),
			&boards[i].ep
		);
		if(res < 0) {
			fprintf(
				stderr,
				"Error sending server request to board nr. %u: %d!\n",
				i,
				res
			);
		} else {
			printf("Reqest send to board %u.\n", i);
			
			sock_udp_ep_t resp_ep = SOCK_IPV6_EP_ANY;
			resp_ep.port = SERVER_POLL_PORT;
			
			sock_udp_t resp_sock;
			sock_udp_create(&resp_sock, &resp_ep, NULL, 0);
			// printf("beg: %p, size: %u\n", (void*)req_buf, sizeof(req_buf));
			res = sock_udp_recv(
				&resp_sock,
				req_buf,
				sizeof(req_buf),
				(uint32_t)2000000,
				&boards[i].ep
			);
			
			if(res == -ETIMEDOUT) {
				fputs("-ETIMEDOUT\n", stderr);
			} else if(res == -EADDRNOTAVAIL) {
				fputs("-EADDRNOTAVAIL\n", stderr);
			} else if(res == -EAGAIN) {
				fputs("-EAGAIN\n", stderr);
			} else if(res == -ENOBUFS) {
				fputs("-ENOBUFS\n", stderr);
			} else if(res == -ENOMEM) {
				fputs("-EPROTO\n", stderr);
			} else if(res == -EPROTO) {
				fputs("-ENOMEM\n", stderr);
			} else if(res < 0) {
				fprintf(stderr, "An error occoured: %d!\n", res);
			} else {
				printf("%s\n", req_buf);
			}
		}
		
	}
}


/*
 * Deamon der im Hintergrund läuft und neue Boards registriert
 */
void* connect_thread_handler(void* args) {
	(void)args;
	
	sock_udp_ep_t local = SOCK_IPV6_EP_ANY;
	
	sock_udp_t sock;
	local.port = SERVER_CONN_PORT;
	bool sock_ready = false;
	if(sock_udp_create(&sock, &local, NULL, 0) < 0) {
		fputs("Error creating UDP sock!\n", stderr);
		sock_ready = false;
	} else {
		puts("Successfull createt UDP socket.");
		puts("Waiting for boards...");
		sock_ready = true;
	}
	
	while(sock_ready) {
		ssize_t res = sock_udp_recv(
			&sock,
			conn_buf,
			CLIENT_INIT_MSG_LEN,
			SOCK_NO_TIMEOUT,
			NULL
		);
		
		if(res < 0) {
			fprintf(
				stderr,
				"Error during \"sock_udp_recv\": %d!\nAboarting.",
				res
			);
			sock_ready = false;
		} else if((res) >= 0) {
			if(!startsWith((const char*)conn_buf, client_id)) {
				continue;
			}
		
			char* addr_str = strchr((const char*)conn_buf, ' ');
			addr_str++; // addr_str zeigt sonst auf das Leerzeichen!
				
			size_t board_index;
			for(board_index=0;
				board_index<MAX_BOARD_NUM && boards[board_index].init;
				board_index++
			);
				
			if(board_index >= MAX_BOARD_NUM) {
				fputs("Max nr. of boards registered!\n", stderr);
			} else {
				boards[board_index].init = true;
				my_strncpy(
					boards[board_index].addr,
					addr_str,
					IPV6_ADDR_MAX_STR_LEN
				);
				
				boards[board_index].ep.family = AF_INET6;
				boards[board_index].ep.port = CLIENT_PORT;
				ipv6_addr_from_str(
					(ipv6_addr_t *)&boards[board_index].ep.addr.ipv6,
					boards[board_index].addr
				);
					
				char resp_msg[SERVER_RESP_MSG_LEN];
				snprintf(
					resp_msg,
					SERVER_RESP_MSG_LEN,
					"%s %s",
					server_id,
					own_addr
				);
					
				res = sock_udp_send(
					NULL,
					resp_msg,
					SERVER_RESP_MSG_LEN,  
					&boards[board_index].ep
				);
				
				if(res == -EAFNOSUPPORT) {
					fputs("-EAFNOSUPPORT\n", stderr);
				} else if(res == -EHOSTUNREACH) {
					fputs("-EHOSTUNREACH\n", stderr);
				} else if(res == -EINVAL) {
					fputs("-EINVAL\n", stderr);
				} else if(res == -ENOMEM) {
					fputs("-ENOMEM\n", stderr);
				} else if(res == -ENOTCONN) {
					fputs("-ENOTCONN\n", stderr);
				} else if(res < 0) {
					fputs("error\n", stderr);
				} else {
					printf("Board found on: %s\n", addr_str);
				}
			}
		}
	}

	return NULL;
}

int main(void) {
	printf("This is the server-emulator on \"%s\".\n", RIOT_BOARD);
	
	for(size_t i=0; i<MAX_BOARD_NUM; i++) {
		boards[i].init = false;
	}
	
	// zu suchendes Prefix in ipv6 umwandeln
	ipv6_addr_t addr;
	ipv6_addr_from_str(&addr, "fe80::");

	// eigene link-local adresse anhand prefix suchen 
	ipv6_addr_t* ll_addr = NULL;
	gnrc_ipv6_netif_find_by_prefix(&ll_addr, &addr);
	
	ipv6_addr_to_str(own_addr, ll_addr, IPV6_ADDR_MAX_STR_LEN);
	printf("Link-local address: \"%s\".\n", own_addr);
	
	// funktioniert nur mit GPIO_IN_PU
	// mit GPIO_RISING wird der Handler 2 mal bei Boardstart aufgerufen,
	// Grund unbekannt
	gpio_init_int(BUTTON_GPIO, GPIO_IN_PU, GPIO_FALLING, button_handler, NULL);
	puts("Button activated.");
	
	thread_create(
		connect_stack,
		sizeof(connect_stack),
		THREAD_PRIORITY_MAIN-1,
		THREAD_CREATE_STACKTEST,
		connect_thread_handler,
		NULL,
		"connect thread"
	);
		
	// TODO shell-commands bauen?
	char line_buf[SHELL_DEFAULT_BUFSIZE];
	shell_run(NULL, line_buf, SHELL_DEFAULT_BUFSIZE);
    
    return EXIT_SUCCESS;
}

