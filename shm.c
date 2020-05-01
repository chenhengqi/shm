#include "shm.h"

#include <unistd.h>
#include <sys/mman.h>
#include <fcntl.h>

int posix_create_shm(const char* name, int flag, mode_t mode, off_t length, void** addr) {
    int fd = shm_open(name, flag, mode);
    if (fd == -1) {
        return -1;
    }

    int code = ftruncate(fd, length);
    if (code == -1) {
        return -1;
    }

    void* shm_addr = mmap(NULL, length, PROT_READ|PROT_WRITE, MAP_SHARED, fd, 0);
    if (shm_addr == MAP_FAILED) {
        return -1;
    }
    *addr = shm_addr;
    return fd;
}

int posix_destroy_shm(const char* name, void *addr, size_t length) {
    int code = munmap(addr, length);
    if (code == -1) {
        return -1;
    }

    return shm_unlink(name);
}
