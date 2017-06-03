/*
 * Copyright (C) 2017 Freie Universität Berlin
 *
 * This file is subject to the terms and conditions of the GNU Lesser
 * General Public License v2.1. See the file LICENSE in the top level
 * directory for more details.
 */

/**
 * @ingroup     driver_tcs37727
 * @{
 *
 * @file
 * @brief       TCS37727 adaption to the RIOT actuator/sensor interface
 *
 * @author      Hauke Petersen <hauke.petersen@fu-berlin.de>
 *
 * @}
 */

#include "tcs37727.h"

char* tcs37727_name = "tcs37727";
char* tcs37727_saul_name = "RGBC light sensor";

int read(void *dev, phydat_t *res);

const saul_driver_t tcs37727_saul_driver;
