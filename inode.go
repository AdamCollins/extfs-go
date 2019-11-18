package main

import (
	"bytes"
	"encoding/binary"
)



type Block = [BLOCK_SIZE] byte
type Inode_t struct {
	I_address	   int16
	I_mode         int16
	I_size         int16
	I_directBlocks [3] int16
}

const DIRECTORY_INODE = 0x400
const REGULAR_FILE_INODE = 0x800

type DirEntry_t struct{
	Inode_P int16	//Inode pointer
	Name [10]byte
}

/*
 * Writes a blank Inode to disk of type mode. and returns the address as int16 and the new Inode
*/
func CreateInode(disk *Disk_t, mode int16) (int16, Inode_t){
	//create inode
	newInode := Inode_t{I_mode:mode,I_size:0}
	//Malloc direct Blocks
	newInode.I_directBlocks[0] = disk.mallocBlock()
	newInode.I_directBlocks[1] = disk.mallocBlock()
	newInode.I_directBlocks[2] = disk.mallocBlock()

	//Malloc a block for writting the inode
	address := disk.mallocBlock()
	newInode.I_address = address

	WriteInode(disk, newInode)
	//update disk.InodeCount
	disk.InodeCount++
	//return saved location
	return address, newInode
}

/** Accepts Directory inodes and adds an inode(file) with name and type mode to the directory
 *
*/
func (inode *Inode_t) CreateDirEntry(disk * Disk_t, name string, mode int16) (int16,Inode_t){
	//create Inode
	address ,newInode := CreateInode(disk,mode)
	//create dir_entry w/ Inode and name
	charArr := [10]byte{}
	copy(charArr[:],name)
	dir_entry := DirEntry_t{Inode_P:address, Name:charArr}

	//convert dir_entry to []byte
	buff := &bytes.Buffer{}

	binary.Write(buff, binary.LittleEndian, dir_entry)

	//write dir_entry to given inode
	inode.WriteInodeData(disk,buff.Bytes())

	return address, newInode

}

/**
*	@Desc: Read inode at mem address pInode
*/
func ReadInode(disk *Disk_t, pInode int16) Inode_t{
	inode := Inode_t{}
	b := disk.ReadData(int16(binary.Size(inode)),pInode)
	buff := bytes.NewBuffer(b)
	binary.Read(buff,binary.LittleEndian, &inode)
	return inode
}

func  WriteInode(disk * Disk_t,inode Inode_t){
	//cast to bytes
	buff := &bytes.Buffer{}
	binary.Write(buff,binary.LittleEndian, inode)
	//Write new Inode
	disk.WriteData(buff.Bytes(),inode.I_address)
}

func (inode *Inode_t) ReadDirInode(disk *Disk_t) []DirEntry_t{
	data := inode.ReadInodeData(disk)
	entrySize := binary.Size(DirEntry_t{})
	numEntries := len(data)/entrySize
	entries := make([]DirEntry_t, numEntries)

	//For each entry found in inode blocks load into slice
	for i:=0; i<numEntries; i++{
		b := data[i*entrySize:(i+1)*entrySize]
		buff := bytes.NewBuffer(b)
		binary.Read(buff,binary.LittleEndian,&entries[i])
	}

	return entries;
}



/*  @Desc: Reads all data stored inside of an inode
 */

func (inode *Inode_t) ReadInodeData(disk *Disk_t) []byte{
	size := inode.I_size
	data := make([]byte,0)

	bytesRead := int16(0)
	block_index := int16(0)
	//Read each block's data (or remaining data) and append to data
	for bytesRead < size {
		readSize := min(size-bytesRead,BLOCK_SIZE)
		newData := disk.ReadData(readSize,inode.I_directBlocks[block_index])
		data = append(data, newData...)
		bytesRead += readSize
		block_index++
	}
	return data
}



/**
	Returns the block_num and offset of next availible space can be written
	Aswell as how many more bytes can be written to the block
	*/
func (inode *Inode_t) getWriteLocation() (int16, int16,int16){

	block_num := inode.I_size/BLOCK_SIZE
	offset := inode.I_size%BLOCK_SIZE;
	//if remaining data cannot fit into the remainder of the old block skip to next

	return block_num, offset, BLOCK_SIZE-offset
}

func (inode *Inode_t) WriteInodeData(disk *Disk_t, data []byte){

	bytesLeft := int16(len(data))
	bytesWriten := int16(0)
	for bytesLeft>0 {
		block_num, off, bytes_avail := inode.getWriteLocation()
		//how many bytes should I write
		writeSize := min(bytesLeft,bytes_avail)
		//where to write?
		startAddress := inode.I_directBlocks[block_num] + off
		//Which range of bytes to write

		writeData := data[bytesWriten:bytesWriten+writeSize]

		//write data
		w := disk.WriteData(writeData,startAddress)
		if w!=writeSize {
			panic("wrong number of bytes written")
			return
		}
		bytesWriten += w
		bytesLeft -= w
		inode.I_size +=w
	}
	WriteInode(disk, *inode)
}


func min(a, b int16) int16 {
	if a<b{
		return a
	}
	return b
}