#include <stdio.h>
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

#define MAX_BOARD_NUM 10
#define BUFF_SIZE_NUM (sizeof(app_id) + IPV6_ADDR_MAX_STR_LEN + 2)

const char app_id[] = "krasse-riot-se-app";

char connect_stack[THREAD_STACKSIZE_DEFAULT];

typedef struct SensorBoard {
	ipv6_addr_t addr;
} SensorBoard;

SensorBoard boards[MAX_BOARD_NUM];

/*
 * #description: checks if string prefix is a prefix of string str
 * #requires: stdbool.h, string.h
 * #param[char* str]: string that should start with something
 * #param[char* prefix]: the prefix to check for
 * #return[bool]: true wenn prefix is a prefix of str, false else
 */
bool startsWith(const char* str, const char* prefix) {
	char* substr = strstr(str, prefix);
	return substr == str;
}


void button_handler(void* args) {
	printf("hipp hipp hurra\n");
}


/*
 * Deamon der im Hintergrund läuft und neue Boards registriert
 * Frage: kann der das Array 'boards' sehen? 
 */
void* connect_thread_handler(void* args) {
	(void)args;
	char buf[BUFF_SIZE_NUM];
	
	sock_udp_ep_t local = SOCK_IPV6_EP_ANY;
	
	sock_udp_t sock;
	local.port = 2017;
	bool sock_ready = false;
	if(sock_udp_create(&sock, &local, NULL, 0) < 0) {
		puts("Error creating UDP sock");
		sock_ready = false;
	} else {
		puts("Successfull createt UDP socket.");
		puts("Waiting for boards.");
		sock_ready = true;
	}
	
	while(sock_ready) {
		ssize_t res = sock_udp_recv(
			&sock,
			buf,
			sizeof(buf),
			SOCK_NO_TIMEOUT,
			NULL
		);
		
		if(res < 0) {
			printf("Error during \"sock_udp_recv\": %d!\nAboarting.", res);
			sock_ready = false;
		} else if((res) >= 0) {
			if(startsWith(buf, app_id)) {
				char* addr_str = strchr(buf, ' ');
				printf("Board found on: %s\n", addr_str);
			}
		}
	}


	return NULL;
}

int main(void) {
	printf("This is the server-emulator on %s with ", RIOT_BOARD);
	
	// zu suchendes Prefix in ipv6 umwandeln
	ipv6_addr_t addr;
	ipv6_addr_from_str(&addr, "fe80::");

	// eigene link-local adresse anhand prefix suchen 
	ipv6_addr_t* ll_addr = NULL;
	gnrc_ipv6_netif_find_by_prefix(&ll_addr, &addr);
	
	char addr_str[IPV6_ADDR_MAX_STR_LEN];
	ipv6_addr_to_str(addr_str, ll_addr, IPV6_ADDR_MAX_STR_LEN);
	printf("\"%s\" as link-local address.\n", addr_str);
	
	// funktioniert nur mit GPIO_IN_PU
	// mit GPIO_RISING wird der Handler 2 mal bei Boardstart aufgerufen,
	// Grund unbekannt
	gpio_init_int(BUTTON_GPIO, GPIO_IN_PU, GPIO_FALLING, button_handler, NULL);
	printf("Button activated.\n");
	
	/*kernel_pid_t conn_t_pid =*/ thread_create(
		connect_stack,
		THREAD_STACKSIZE_DEFAULT,
		THREAD_PRIORITY_MAIN-1,
		THREAD_CREATE_STACKTEST,
		connect_thread_handler,
		NULL,
		"connection thread"
	);
	
	// TODO shell-commands bauen?
	char line_buf[SHELL_DEFAULT_BUFSIZE];
	shell_run(NULL, line_buf, SHELL_DEFAULT_BUFSIZE);
    
    return EXIT_SUCCESS;
}

