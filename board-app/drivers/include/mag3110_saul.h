/*
 * Copyright (C) 2017 HAW Hamburg
 *
 * This file is subject to the terms and conditions of the GNU Lesser
 * General Public License v2.1. See the file LICENSE in the top level
 * directory for more details.
 */

/**
 * @ingroup     driver_mag3110
 * @{
 *
 * @file
 * @brief       MAG3110 adaption to the RIOT actuator/sensor interface
 *
 * @author      Sebastian Meiling <s@mlng.net>
 *
 * @}
 */

#include "mag3110.h"

int read_mag(void *dev, phydat_t *res);

const saul_driver_t mag3110_saul_driver;
