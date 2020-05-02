#pragma once

#include <sys/stat.h>
#include <sys/types.h>

int posix_shm_create(const char* name, int flag, mode_t mode, off_t length, void** addr);
int posix_shm_destroy(const char* name, void *addr, size_t length);
int64_t posix_shm_seek(int fd, off_t offset, int whence);
int posix_shm_read(int fd, void* buf, size_t count);
int posix_shm_write(int fd, void* buf, size_t count);
