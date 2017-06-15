#ifndef _SMART_ENVIRONMENT_H_
#define _SMART_ENVIRONMENT_H_

#define TEMP_SENSOR		0x00000001
#define HUMID_SENSOR	0x00000002
#define MAG_SENSOR		0x00000004


#define SERVER_CONN_PORT	2017
#define APP_ID_LEN			sizeof(app_id)
#define APP_PING_MSG_LEN	APP_ID_LEN + 1 + 4

// statische globals, nach außen hin nicht sichtbar
const char app_id[] = "riot-swp-2017-se";

#endif
