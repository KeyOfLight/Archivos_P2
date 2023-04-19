package Estructuras

import (
	"fmt"
)

type ParamStruct struct {
	Nombre    string
	Direccion string
	Tam       string
	Unit      string
	Fit       string
	Status    string
	Tipo      string
	Start     string
	Size      string
	Delete    string
	Add       string
	Fs        string
	User      string
	Pwd       string
}

type MBR struct {
	Mbr_tamano         int64
	Mbr_fecha_creacion [19]byte
	Mbr_dsk_signature  int64
	Dsk_fit            byte
	Mbr_partition      [4]Particion
}

type Particion struct {
	Part_status byte
	Part_type   byte
	Part_fit    byte
	Part_start  int64
	Part_size   int64
	Part_name   [16]byte
}

type EBR struct {
	Part_status byte
	Part_fit    byte
	Part_start  int64
	Part_size   int64
	Part_next   int64
	Part_name   [16]byte
}

type Sblock struct {
	S_filesystem_type   int64
	S_inodes_count      int64
	S_blocks_count      int64
	S_free_blocks_count int64
	S_free_inodes_count int64
	S_mtime             [19]byte
	S_mnt_count         int64
	S_magic             int64
	S_inode_size        int64
	S_block_size        int64
	S_firts_ino         int64
	S_first_blo         int64
	S_bm_inode_start    int64
	S_bm_block_start    int64
	S_inode_start       int64
	S_block_start       int64
}

type I_node struct {
	I_uid   int64
	I_gid   int64
	I_size  int64
	I_atime [19]byte
	I_ctime [19]byte
	I_mtime [19]byte
	I_block [16]byte
	I_type  byte
	I_perm  int64
}

var Listado = ListaMontados{}

var Logueado = User{}

// Bloque Carpetas
type Bcontent struct {
	B_name  [15]byte
	B_inodo byte
}

type BloqueCarpetas struct {
	B_content [4]Bcontent
}

// Bloque Archivos
type BloqueArchivos struct {
	B_content [64]byte
}

// Bloque Apuntador
type Bapuntador struct {
	B_pointers [16]int64
}

type PartMounted struct {
	Id         string
	Name       string
	Path       string
	StartPoint int64
	Size       int64
	Tipo       byte
}

type ListaMontados struct {
	Montado []PartMounted
}

func (X *ListaMontados) Montar(Nombre string, Path string, StartPoint int64, Size int64, Tipo byte) {
	Montar := PartMounted{}
	lens := len(X.Montado)
	Montar.Id = "27" + fmt.Sprint(lens) + Nombre
	fmt.Println(Montar.Id)
	Montar.Path = Path
	Montar.Size = Size
	Montar.StartPoint = StartPoint
	Montar.Tipo = Tipo
	Listado.Montado = append(X.Montado, Montar)
}

type User struct {
	UID        byte
	Tipo       byte
	Grupo      [11]byte
	Pass       [11]byte
	Uss        [11]byte
	Path       string
	Startpoint int64
}

func (X *ListaMontados) GetLista() ListaMontados {

	return Listado
}

func (X *User) Loguear(Usuario User) {
	Logueado = Usuario
}

func (X *User) Getusser() User {

	return Logueado
}

func (X *User) Logout() {

	Logueado = User{}
}
