package main

import (
	"fmt"
	"strings"
)

var Disk *Disk_t
var root_dir *Inode_t	// '/'
var cur_dir *Inode_t	// '$'



func initDisk(){
	Disk = Open()
	//Create root inode
	_, root := CreateInode(Disk,0x800)
	Disk.Root_inode = root
	//Write changes to disk
	Disk.WriteDisk()
	root_dir = &Disk.Root_inode
	cur_dir = root_dir
}

func cd(inode * Inode_t){
	cur_dir = inode
}

func ls(inode * Inode_t){
	if inode == nil{
		inode = cur_dir
	}
	if inode.I_mode != 0x800{
		fmt.Printf("ls: File is not a directory")
		return
	}
	entries := inode.ReadDirInode(Disk)

	for _,e := range entries{
		name :=string(e.Name[:])
		fmt.Printf("%s ",name)
	}
	fmt.Println()
	
}

func cat(inode * Inode_t) string{
	return string(inode.ReadInodeData(Disk))
}

func touch(parent * Inode_t, name, contents string) Inode_t{
	if parent == nil{
		parent = cur_dir
	}
	//malloc new entry 
	 _, entry := parent.CreateDirEntry(Disk, name,REGULAR_FILE_INODE)
	 entry.WriteInodeData(Disk,[]byte(contents))
	return entry
}

func mkdir(parent * Inode_t, name string) Inode_t{
	if parent == nil{
		parent = cur_dir
	}
	_ , ind := parent.CreateDirEntry(Disk,name,0x400)
	return ind
}

func tree(inode Inode_t, h int){
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
			tree(ind, h+1)
		}else{
			val := cat(&ind)
			fmt.Printf("│%s└─ %s: '%s'\n",path, name,val)
		}
	}
}


func main() {
	initDisk()
	//Add home dir to root

	//Add /usr
	usrDir := mkdir(root_dir,"usr")
	//Add /dev
	mkdir(root_dir,"dev")
	
	//Add /usr
	touch(&usrDir,"notes.txt", "hello")
	cd(&usrDir)
	mkdir(nil,"bin")
	//Disk.Root_inode.CreateDirEntry(Disk,"bin",0x400)

	tree(Disk.Root_inode, 0)
	ls(root_dir)

	//readFile(Disk.Root_inode)
	Disk.Close()
}

