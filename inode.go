package main

import (
	"bytes"
	"encoding/binary"
)



type Block = [BLOCK_SIZE] byte
type Inode_t struct {
	I_mode         int16
	I_size         int16
	I_directBlocks [3] int16
}

const DIRECTORY_INODE = 0x400
const REGULAR_FILE_INODE = 0x800

type DirEntry_t struct{
	Indode Inode_t
	Name [8]byte
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

	//cast to bytes
	buff := &bytes.Buffer{}
	binary.Write(buff,binary.LittleEndian,newInode)
	//Malloc a block for writting the inode
	address := disk.mallocBlock()
	//Write new Inode
	disk.WriteData(buff.Bytes(),address)
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
	charArr := [8]byte{}
	copy(charArr[:],name)
	dir_entry := DirEntry_t{Indode:newInode, Name:charArr}

	//convert dir_entry to []byte
	buff := &bytes.Buffer{}

	binary.Write(buff, binary.LittleEndian, dir_entry)

	//write dir_entry to given inode
	inode.WriteInodeData(disk,buff.Bytes())

	return address, newInode

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
	//Read one block at a time
	for i:=0; i<int(size); i+=BLOCK_SIZE {
		address :=inode.I_directBlocks[i/BLOCK_SIZE]
		newData := disk.ReadData(size%BLOCK_SIZE, address)
		data = append(data, newData...)
	}
	inode.I_size+=int16(len(data))
	return data
}

/**
	Returns the block and offset where n bytes can be written
 */
func (inode *Inode_t) getWriteLocation(n int16) (int16, int16){

	block_num := inode.I_size/BLOCK_SIZE
	offset := inode.I_size%BLOCK_SIZE;
	//if remaining data cannot fit into the remainder of the old block skip to next
	if offset + n > BLOCK_SIZE {
		return block_num+1, 0
	}

	return block_num, offset
}

func (inode *Inode_t) WriteInodeData(disk *Disk_t, data []byte){

	for i:=0; i<len(data); i+=BLOCK_SIZE {
		block_num, offset := inode.getWriteLocation(BLOCK_SIZE-1);//write each one to a new block
		dataRange := [2]int16{offset,min(int16(len(data)),BLOCK_SIZE)}

		//Write to directBlock address + offset
		address :=inode.I_directBlocks[block_num] + offset
		//subset of data that can fit into one block
		dataBlock := data[dataRange[0] : dataRange[1]]
		inode.I_size += disk.WriteData(dataBlock,address)
	}
}


//func (inode *Inode_t) WriteInodeData(disk *Disk_t, data []byte) int16{
//	var byteCount int16 = 0
//	//Write to one block at a time
//	//If set read range
//	addressRange := [2]int16{inode.I_size,min(int16(len(data)),BLOCK_SIZE)}
//	for i:=0; i<len(data); i+=BLOCK_SIZE {
//		address :=inode.I_directBlocks[i/BLOCK_SIZE]
//		//subset of data that can fit into one block
//		dataBlock := data[addressRange[0] : addressRange[1]]
//		bytesRead := disk.WriteData(dataBlock,address)
//		addressRange[0] += bytesRead;
//		addressRange[1] += bytesRead;
//		byteCount += bytesRead
//
//	}
//	inode.I_size += byteCount
//	return byteCount
//}

func min(a, b int16) int16 {
	if a<b{
		return a
	}
	return b
}