#include "csignal.h"

/* #include <sys/stat.h> */
#include <signal.h>

#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>
#include <fcntl.h>
#include <errno.h>
#include <sys/stat.h>
#include <sys/socket.h>
#include <sys/un.h>

int signal_ignore(int signum) {
	struct sigaction action = {.sa_handler = SIG_IGN};
	return sigaction(signum, &action, 0);
}

/* int open_socket(struct sockaddr_un * addr, size_t path_size) */
/* { */
/*     int fd; */
/*     socklen_t size = offsetof(typeof(*addr), sun_path) + path_size + 1; */

/*     if ((fd = socket(PF_LOCAL, SOCK_STREAM | SOCK_CLOEXEC, 0)) < 0) */
/*         goto error0; */

/*     /\* Unlink the socket location in case it was being used by a process which */
/*      * left around a stale lockfile. *\/ */
/*     unlink(addr->sun_path); */

/*     if (bind(fd, (struct sockaddr *) addr, size) < 0) */
/*         goto error1; */

/*     if (listen(fd, 1) < 0) */
/*         goto error2; */

/*     return fd; */

/*   error2: */
/*     if (addr->sun_path[0]) */
/*         unlink(addr->sun_path); */
/*   error1: */
/*     close(fd); */
/*   error0: */
/*     return -1; */
/* } */
