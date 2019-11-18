package main

import (
	"fmt"
)

var Disk *Disk_t


/*
 *  \
 *    --> file1: "hello"
 *
 */

func initDisk(){
	Disk = Open()
	//Create root inode
	_, root := CreateInode(Disk,0x400)
	Disk.Root_inode = root

	//Write changes to disk
	Disk.WriteDisk()
}

func readFile(inode Inode_t) {

	str := string(inode.ReadInodeData(Disk))
	fmt.Println(str)
}

func printDirectory(inode Inode_t){
	entries := inode.ReadDirInode(Disk)
	for _,e := range entries{
		name := string(e.Name[:])
		if e.Indode.I_mode == DIRECTORY_INODE{
			fmt.Printf("/%s\n",name)
		}

	}
}



func main() {
	initDisk()
	//writeFile("hello!")
	//Add home dir to root

	Disk.Root_inode.CreateDirEntry(Disk,"usr",0x400)
	Disk.Root_inode.CreateDirEntry(Disk,"dev",0x400)
	//Disk.Root_inode.CreateDirEntry(Disk,"bin",0x400)

	printDirectory(Disk.Root_inode)

	//readFile(Disk.Root_inode)
	Disk.Close()
}

