package main

import (
	"flag"
	"log"
	"syscall"
	"unsafe"
)

/*
#include <sys/sem.h>
typedef struct sembuf sembuf;
typedef union semun semun;

*/
import "C"

func semget(key int) int {
	r1, r2, err := syscall.Syscall(syscall.SYS_SEMGET, uintptr(key),
		uintptr(1), uintptr(00666))
	if int(r1) < 0 {
		r1, r2, err = syscall.Syscall(syscall.SYS_SEMGET, uintptr(key),
			uintptr(1), uintptr(C.IPC_CREAT|C.IPC_EXCL|00666))
		if int(r1) < 0 {
			log.Printf("error:semget error is %v\n", err)
		}
	} else {
		log.Printf("success :semget is %v,%v,%v\n", r1, r2, err)
	}
	return int(r1)
}

func semLock(semid int) int {

	stSemBuf := C.sembuf{
		sem_num: 0,
		sem_op:  -1,
		sem_flg: C.IPC_NOWAIT | C.SEM_UNDO,
	}

	r1, r2, err := syscall.Syscall(syscall.SYS_SEMOP, uintptr(semid), uintptr(unsafe.Pointer(&stSemBuf)), 1)
	if int(r1) < 0 {
		log.Printf("error:semget error is %v,%v,%v\n", r1, r2, err)
	}
	return int(r1)
}

//int  semctl(int _semid  ,int _semnum,int _cmd  ……);
var view = flag.Bool("v", false, "get the current value of a sem")
var key = flag.Int("k", 0, "Set the Sem Key")

//flag.StringVar(&operate,"o", "add", "operation for calc")

func main() {
	flag.Parse()
	if *key == 0 {
		flag.Usage()
		log.Fatal("Must have a key")
	}

	sem := semget(*key)
	if sem < 0 {
		log.Printf("Open Sem Failed with %d!\n", sem)
		return
	}
	ret := semLock(sem)
	if ret != 0 {
		log.Fatal("Wait Sem Failed!\n")
	}
	log.Printf("Wait for Sem %d success and go on!\n", *key)

}
