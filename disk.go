package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
)

const DISK_SIZE = 2048
const BLOCK_SIZE  = 128

const BLOCK_GROUP_0 int16 = 0x100

const BLOCK_GROUP_1 int16 = 0x200


const FILE_NAME = "disk.img"
var fd *os.File

type Disk_t struct {
	Size       int16
	InodeCount int16
	Root_inode Inode_t

	//inode_table [DISK_SIZE]byte;
}


/*  @Desc: If a disk img exists open the file and return the contents as a *Disk_t
 *  	   If a disk img does not exist create a new img file and return a new *Disk_t
 *  		Set the fileDescripter variable for future use
*/
func Open() *Disk_t{
	disk  := new(Disk_t)
	_, err := os.Stat(FILE_NAME)
	if os.IsNotExist(err) {
		fd, err = os.Create(FILE_NAME)
		fd.Write(make([]byte,DISK_SIZE))
	}else{
		fd, err = os.OpenFile(FILE_NAME,os.O_RDWR,os.ModePerm)
		binary.Read(fd,binary.LittleEndian,disk)
	}

	if err!=nil{
		return nil
	}
	return disk
}

/*  @Desc: Re-Writes data to disk and then closes fd
 */
func (disk *Disk_t)Close()  {
	disk.WriteDisk()
	fd.Close()
}


/*  @Desc: Reads n bytes from address offset and returns data in []byte
 */
func (disk * Disk_t)ReadData(n, offset int16) []byte{
	buff := make([]byte,n)
	wroteN, err := fd.ReadAt(buff,int64(offset))
	if err != nil{
		fmt.Println(err, wroteN)
	}
	return buff
}

/*  @Desc: Writes len(data) bytes to disk at address offset and returns number of bytes written
 */
func (disk * Disk_t)WriteData(data []byte, offset int16) int16{
	var buff bytes.Buffer

	binary.Write(&buff, binary.LittleEndian, data)
	n, err := fd.WriteAt(buff.Bytes(),int64(offset))
	if err!=nil{
		fmt.Printf("ERROR: %e\n",err)
		return -1
	}
	return int16(n)

}


/*  @Desc: Writes Whole Disk to img file
 */
func (disk * Disk_t)WriteDisk(){

	var buff bytes.Buffer

	 binary.Write(&buff, binary.LittleEndian, disk)
	 _, err := fd.Write(buff.Bytes())
	 if err!=nil{
	 	fmt.Printf("ERROR: %e\n",err)
	 }
}


func (disk * Disk_t)mallocBlock() int16{
	address := disk.Size+BLOCK_GROUP_0;
	disk.Size+=BLOCK_SIZE
	return address
}
