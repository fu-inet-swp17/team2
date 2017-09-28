/**
 * Copyright (C) 2017 Freie Universität Berlin
 * Authors: Adrian Herrmann <adrian.herrmann@fu-berlin.de>
 *          Thomas Tegethoff <cone@zedat.fu-berlin.de>
 *
 *
 * This file is subject to the terms and conditions of the GNU Lesser
 * General Public License v2.1. See the file LICENSE in the top level
 * directory for more details.
 */

#ifndef __SMART_ENVIRONMENT_SWP_2017_H__
#define __SMART_ENVIRONMENT_SWP_2017_H__

/**
 * @brief Port at which the ping thread is sending
 */
#define SERVER_CONN_PORT	2017

/**
 * @brief ID of the SenPhy application
 */
static const char app_id[] = "riot-swp-2017-se";

/**
 * @brief Length of app_id
 */
#define APP_ID_LEN			sizeof(app_id)

#endif