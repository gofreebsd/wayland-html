#include "_cgo_export.h"

#include "cfn.h"

void puts_binding(ffi_cif *cif, unsigned int *ret, void* args[],
    FILE *stream) {
  *ret = fputs(*(char **)args[0], stream);
}

void binding(ffi_cif *cif, unsigned int *ret, void **args, void *data) {
    /* printf("%s\n", *((char **)(args[0]))); */
    /* printf("%d\n", *((int *)(args[1]))); */
    cfn_go_callback(data, args, ret);
}


typedef float(*callback )(char *, int, int);

void call(void * fn) {
   callback cb = (callback)fn;
   float i = cb("helloworld", 123, 456);
   printf("%f\n", i);
}
