#include <ffi.h>
#include <stdio.h>

void puts_binding(ffi_cif *cif, unsigned int *ret, void* args[],
    FILE *stream);

void binding(ffi_cif *cif, unsigned int *ret, void* args[], void *data);
void call(void * fn);
