#ifndef _SMART_ENVIRONMENT_H_
#define _SMART_ENVIRONMENT_H_

// Nanocoap unterstützt nur 64bit Zahlen
#define IR_TEMP_SENSOR		0x0000000000000001
#define HUMID_SENSOR		0x0000000000000002
#define MAG_SENSOR			0x0000000000000004
#define RGB_LIGHT_SENDSOR	0x0000000000000008
#define PRESS_SENSOR		0x0000000000000010
#define ACC_SENSOR			0x0000000000000020

#define SERVER_CONN_PORT	2017


// statische globals, nach außen hin nicht sichtbar
static const char app_id[] = "riot-swp-2017-se";
#define APP_ID_LEN			sizeof(app_id)

typedef struct PingMsg {
	char app_id[APP_ID_LEN];
	unsigned long sensors;
} PingMsg;

#define APP_PING_MSG_LEN	sizeof(PingMsg)


#endif
