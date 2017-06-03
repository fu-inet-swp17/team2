/*
 * Copyright (C) 2016 Freie Universität Berlin
 *
 * This file is subject to the terms and conditions of the GNU Lesser
 * General Public License v2.1. See the file LICENSE in the top level
 * directory for more details.
 */

/**
 * @ingroup     driver_tmp006
 * @{
 *
 * @file
 * @brief       TMP006 adaption to the RIOT actuator/sensor interface
 *
 * @author      Adrian Herrmann <adrian.herrmann@fu-berlin.de>
 *
 * @}
 */

#include "tmp006.h"

int read_temp(void *dev, phydat_t *res);

const saul_driver_t tmp006_saul_driver;
