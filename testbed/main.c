#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <inttypes.h>
#include <stdbool.h>

#include "shell.h"
#include "xtimer.h"
#include "periph/gpio.h"
#include "board.h"

#include "ws2812b_driver.h"

#include "smart_environment.h"

#define LED_NUM 50

ws2812b_stripe_t test;

int main(void) {
	
	gpio_t pin = GPIO_PIN(0, 7);
	ws2812b_init(&test, LED_NUM, pin);
	
	for(size_t i=0; i<LED_NUM; i++) {
		color_rgb_t led;
	
		led.r = 0xAA;
		led.g = 0xAA;
		led.b = 0xAA;
		
		ws2812b_set_color(&test, i, led);
	}
	
	ws2812b_show(&test);
    
    // shell starten damit Application nicht terminiert
    char line_buf[SHELL_DEFAULT_BUFSIZE];
    shell_run(NULL, line_buf, SHELL_DEFAULT_BUFSIZE);
    
    return EXIT_SUCCESS;
}

