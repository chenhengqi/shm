#include "shm.h"

#include <unistd.h>
#include <sys/mman.h>
#include <fcntl.h>

int posix_shm_create(const char* name, int flag, mode_t mode, off_t length, void** addr) {
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

int posix_shm_destroy(const char* name, void *addr, size_t length) {
    int code = munmap(addr, length);
    if (code == -1) {
        return -1;
    }

    return shm_unlink(name);
}

int64_t posix_shm_seek(int fd, off_t offset, int whence) {
    off_t new_offset = lseek(fd, offset, whence);
    if (new_offset == (off_t)-1) {
        return -1;
    }
    return (int64_t)new_offset;
}

int posix_shm_read(int fd, void* buf, size_t count) {
    return read(fd, buf, count);
}

int posix_shm_write(int fd, void* buf, size_t count) {
    return write(fd, buf, count);
}
