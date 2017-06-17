#include <stdio.h>

#include "saul_reg.h"

void sensors_test(void) {
    puts("\n##############################################");
    puts("############## SAUL sensor test ##############\n");

    saul_reg_t* dev = saul_reg;

    while(1) {
        phydat_t res;
        int dim = saul_reg_read(dev, &res);
        printf("%s of type %s returns the following data: \n", dev->name, saul_class_to_str(dev->driver->type));
        phydat_dump(&res, dim);
        if (dev->next == NULL) {
            break;
        }
        dev = dev->next;
    };

    puts("##############################################\n");
}
