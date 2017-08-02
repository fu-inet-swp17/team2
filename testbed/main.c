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

#define LED_NUM 10

ws2812b_stripe_t test;

int main(void) {
	gpio_t pin = GPIO_PIN(PORT_A, 2);
	gpio_init(pin, GPIO_OUT);
	gpio_clear(pin);
	ws2812b_init(&test, LED_NUM, pin);
		

	for(size_t i=0; i<LED_NUM; i++) {
		color_rgb_t led;
		
		led.r = 0xAA;
		led.g = 0xAA;
		led.b = 0xAA;
		/*
		if(i % 6 == 0) { // yellow
			led.r = 0xFF;
			led.g = 0xFF;
			led.b = 0x00;
		} else if(i % 5 == 0) { // orange
			led.r = 0xFF;
			led.g = 0x80;
			led.b = 0x00;
		} else if(i % 4 == 0) { // red
			led.r = 0xFF;
			led.g = 0x80;
			led.b = 0x00;
		} else if(i % 3 == 0) { // purple
			led.r = 0x80;
			led.g = 0x00;
			led.b = 0x80;
		} else if(i % 2 == 0) { // blue
			led.r = 0x00;
			led.g = 0x00;
			led.b = 0xFF;
		} else { // green
			led.r = 0x00;
			led.g = 0xFF;
			led.b = 0x00;
		}
		*/
		ws2812b_set_color(&test, i, led);
	}
	ws2812b_show(&test);
	    
    // shell starten damit Application nicht terminiert
    char line_buf[SHELL_DEFAULT_BUFSIZE];
    shell_run(NULL, line_buf, SHELL_DEFAULT_BUFSIZE);
    
    return EXIT_SUCCESS;
}

