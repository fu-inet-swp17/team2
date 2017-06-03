/*
 * Copyright (C) 2017 Freie Universität Berlin
 *
 * This file is subject to the terms and conditions of the GNU Lesser
 * General Public License v2.1. See the file LICENSE in the top level
 * directory for more details.
 */

/**
 * @ingroup     driver_hdc1000
 * @{
 *
 * @file
 * @brief       HDC1000 adaption to the RIOT actuator/sensor interface
 *
 * @author      Hauke Petersen <hauke.petersen@fu-berlin.de>
 *
 * @}
 */

#include "hdc1000.h"

char* hdc1000_name = "hdc1000";
char* hdc1000_saul_name = "humidity sensor";

int read_temp(void *dev, phydat_t *res);

int read_hum(void *dev, phydat_t *res);

const saul_driver_t hdc1000_saul_temp_driver;

const saul_driver_t hdc1000_saul_hum_driver;
