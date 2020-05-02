#pragma once

#include <sys/stat.h>
#include <sys/types.h>
#include <sys/ipc.h>
#include <sys/shm.h>

int posix_shm_create(const char* name, int flag, mode_t mode, off_t length, void** addr);
int posix_shm_remove(const char* name, void *addr, size_t length);
int64_t posix_shm_seek(int fd, off_t offset, int whence);
int posix_shm_read(int fd, void* buf, size_t count);
int posix_shm_write(int fd, void* buf, size_t count);

int sysv_shm_create(const char *pathname, int proj_id, size_t size, int flag, int mode, void** addr);
int sysv_shm_remove(int shmid, void* addr);
int sysv_shm_read(void* dest, void* src, int offset, size_t count);
int sysv_shm_write(void* dest, int offset, void* src, size_t count);
