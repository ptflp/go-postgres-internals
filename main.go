package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

// PageHeader is based on pg's PageHeaderData struct
// from src/include/storage/bufpage.h.
type PageHeader struct {
	XLogID            uint32
	XRecOff           uint32
	PdChecksum        uint16
	PdFlags           uint16
	PdLower           uint16
	PdUpper           uint16
	PdSpecial         uint16
	PdPagesizeVersion uint16
	PdPruneXID        uint32
}

type RawItemIDData struct {
	LpOff uint16
	LpLen uint16
}

type ItemIDData struct {
	LpOff   uint16
	LpLen   uint16
	LpFlags byte
}

const PageHeaderByteSize int = 24
const ItemIDByteSize int = 4

func main() {
	ReadTableHeader()
}

func ReadTableHeader() {
	path := "./data/base/16384/16388"

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Could not open '%s': %v", path, err)
	}
	defer file.Close() // ignore possible close error

	header := &PageHeader{}
	err = binary.Read(file, binary.LittleEndian, header)
	if err != nil {
		log.Fatalf("Could not read page header: %v", err)
	}

	fmt.Printf("page header: %+v\n", header)

	itemIDDataSize := (int(header.PdLower) - PageHeaderByteSize) / ItemIDByteSize

	fmt.Printf("Number of ItemIDData structs: %d\n", itemIDDataSize)

	itemIDDataPointers := make([]ItemIDData, itemIDDataSize)

	for i := 0; i < itemIDDataSize; i++ {
		raw := &RawItemIDData{}
		err = binary.Read(file, binary.LittleEndian, raw)
		if err != nil {
			log.Fatalf("Could not read %dth ItemIDData: %v", i, err)
		}
		newItem := ItemIDData{
			LpOff: raw.LpOff & 0b_01111111_11111111,
			LpLen: raw.LpLen >> 1,
			LpFlags: byte(((raw.LpOff >> 15) & 0b_00000000_00000001) |
				((raw.LpLen << 1) & 0b_00000000_00000010)),
		}
		itemIDDataPointers[i] = newItem
	}

	for _, p := range itemIDDataPointers {
		fmt.Printf("item ID data pointer: %+v\n", p)
	}
}
