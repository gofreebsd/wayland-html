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
