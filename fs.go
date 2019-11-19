package main

import (
	"fmt"
	"strings"
)

var Disk *Disk_t
var root_dir *Inode_t

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
	root_dir = &Disk.Root_inode
}

func readFile(inode Inode_t) {

	str := string(inode.ReadInodeData(Disk))
	fmt.Println(str)
}


/**
	* @Desc Recursively prints directories
*/
func printDirectory(inode Inode_t, h int){
	entries := inode.ReadDirInode(Disk)
	if len(entries)>0 && h==0{
		fmt.Println(".")
	}
	for _,e := range entries{
		name := string(e.Name[:])
		path :=strings.Repeat("   ",h)
		ind := ReadInode(Disk, e.Inode_P)
		if ind.I_mode == DIRECTORY_INODE{
			if h>0 {
				fmt.Printf("│")
			}
			fmt.Printf("%s└─ %s/\n",path, name)
			printDirectory(ind, h+1)
		}else{
			fmt.Printf("│%s└─ %s\n",path, name)
		}
	}
}


func createTextFile(parent * Inode_t, name, contents string) Inode_t{
	//malloc new entry 
	 _, entry := parent.CreateDirEntry(Disk, name,REGULAR_FILE_INODE)
	 entry.WriteInodeData(Disk,[]byte(contents))
	return entry
}

func mkdir(inode * Inode_t, name string) Inode_t{
	_ , ind := inode.CreateDirEntry(Disk,name,0x400)
	return ind
}

func main() {
	initDisk()
	//Add home dir to root

	//Add /usr
	usrDir := mkdir(root_dir,"usr")
	//Add /dev
	mkdir(root_dir,"dev")
	
	
	//Add /usr
	createTextFile(&usrDir,"notes.txt", "hello")

	mkdir(&usrDir,"bin")

	//Disk.Root_inode.CreateDirEntry(Disk,"bin",0x400)

	printDirectory(Disk.Root_inode, 0)

	//readFile(Disk.Root_inode)
	Disk.Close()
}

