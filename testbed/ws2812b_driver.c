#include "ws2812b_driver.h"
#include "periph/gpio.h"

#define ENABLE_DEBUG    (0)
#include "debug.h"
#include "xtimer.h"

static inline void send_zero(gpio_t pin) {
	// complet transmission: 1.25us +/- 150ns
	gpio_set(pin);
	// wait 250ns +/- 150ns
	__asm__ volatile(
		"nop \n\t"
	);
	
	gpio_clear(pin);
	// wait 900ns +/- 150ns
	__asm__ volatile(
		"nop\n"
	);
}

static inline void send_one(gpio_t pin) {
	// complet transmission: 1.25us +/- 150ns
	gpio_set(pin);
	// wait 900ns +/- 150ns
	__asm__ volatile(
		"nop\n"
		"nop\n"
		"nop\n"
		"nop\n"	
		"nop\n"
		"nop\n"
		"nop\n"
		"nop\n"	
		"nop\n"
		"nop\n"
		"nop\n"
		"nop\n"	
		"nop\n"
		"nop\n"
		"nop\n"
		"nop\n"
		"nop\n"
	);
	
	gpio_clear(pin);
	// wait 350ns +/- 150ns
}

int ws2812b_init(ws2812b_stripe_t *dev, uint32_t led_count, gpio_t pin) {
	if (gpio_init(pin, GPIO_OUT) < 0) {
		DEBUG("ws2812b_init: failed to init gpio %d\n", pin);
		return -1;
	}

	gpio_clear(pin);

	dev->led_count = led_count;
	dev->pin = pin;

	return 0;
}

int ws2812b_set_color(ws2812b_stripe_t *dev, uint32_t index, color_rgb_t led) {
	if(index >= dev->led_count) {
		DEBUG("ws2812b_set_color: index out of bound.\n");
		return -1;
	}

	dev->leds[index].r = led.r;
	dev->leds[index].g = led.g;
	dev->leds[index].b = led.b;

	return 0;
}

// maybe a double Pointer for the led?
int ws2812b_get_color(ws2812b_stripe_t *dev, uint32_t index, color_rgb_t *led) {
	if (index >= dev->led_count) {
		DEBUG("ws2812b_get_color: index ouf of bound.\n");
		return -1;
	}

	/* NOTE: For color order see neopixel_set_pixel_color */
	led->r = dev->leds[index].r;
	led->g = dev->leds[index].g;
	led->b = dev->leds[index].b;

	return 0;
}

void ws2812b_show(ws2812b_stripe_t *dev) {
	int irq_state;
	uint32_t i;	// Timing!
	uint8_t bit;
	color_rgb_t *pixel; // Timing!
	gpio_t pin = dev->pin; // Timing!

	gpio_clear(dev->pin);
	/*for(i=0; i<dev->led_count; i++) {
		printf(
			"LED %"PRIu32": r: %d g: %d b: %d\n",
			i,
			dev->leds[i].r,
			dev->leds[i].g,
			dev->leds[i].b
		);
	}*/

	/* Disable IRQ's, we have hard time bounds here.
	 * Restore it afterwards */
	// gpio_set(dev->pin);
	// gpio_clear(dev->pin);
	
	irq_state = irq_disable();
	
	/* Loop over all connected leds
	 * For the WS2812B only the HIGH Signal has to be distinctive,
	 * for LOW we don't need any NOP's.
	 * We could also loop over all three colors for each LED,
	 * but since we need hard timings for each 24bit we unroled the loop.
	 */
	for(i=0; i<dev->led_count; i++) {
		pixel = &(dev->leds[i]);
		// TODO Ausrollen!
		// green
		for(bit=0x80; bit; bit >>= 1) {
			if(pixel->g & bit) {
				send_one(pin);					
			} else {
				send_zero(pin);
			}
		}
		
		// red
		for(bit=0x80; bit; bit >>= 1) {
			if(pixel->r & bit) {
				send_one(pin);		
			} else {
				send_zero(pin);
			}
		}
		
		for(bit=0x80; bit; bit >>= 1) {
			if(pixel->b & bit) {
				send_one(pin);
			} else {
				send_zero(pin);
			}
		}
	}
	
	/* Restore IRQ State */
	irq_restore(irq_state);
	xtimer_usleep(50);
}

