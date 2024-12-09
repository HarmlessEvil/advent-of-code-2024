package main

import (
	"container/heap"
	"fmt"
	"os"
)

func main() {
	if res, err := part1(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}

	if res, err := part2(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}

func part1() (int, error) {
	diskMap, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	disk := mapToDisk(diskMap)

	fragmentDisk(diskMap, disk)

	checksum := calculateChecksum(disk)

	return checksum, nil
}

func part2() (int, error) {
	diskMap, err := readInput()
	if err != nil {
		return 0, fmt.Errorf("readInput: %w", err)
	}

	sizeMap := offsetToSizeMap(diskMap)
	h := mapToFreeSpaceHeap(diskMap)

	disk := mapToDisk(diskMap)

	deFragmentDisk(diskMap, disk, sizeMap, h)

	checksum := calculateChecksum(disk)

	return checksum, nil
}

type SizeMap struct {
	Free  map[int]int
	Files map[int]int
}

func (m SizeMap) Reverse() SizeMap {
	sizeMap := SizeMap{
		Free:  make(map[int]int, len(m.Free)),
		Files: make(map[int]int, len(m.Files)),
	}

	for start, size := range m.Free {
		sizeMap.Free[start+size] = size
	}

	for start, size := range m.Files {
		sizeMap.Files[start+size] = size
	}

	return sizeMap
}

type FreeSpaceHeap []int

func (h FreeSpaceHeap) Len() int {
	return len(h)
}

func (h FreeSpaceHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h FreeSpaceHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *FreeSpaceHeap) Push(x any) {
	*h = append(*h, x.(int))
}

func (h *FreeSpaceHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

const maxFileSize = 9

func mapToFreeSpaceHeap(diskMap []byte) [maxFileSize]FreeSpaceHeap {
	var h [maxFileSize]FreeSpaceHeap

	offset := 0
	for i := 0; i < len(diskMap); i += 2 {
		offset += int(diskMap[i])

		if i+1 < len(diskMap) {
			size := diskMap[i+1]
			if size > 0 {
				heap.Push(&h[size-1], offset)
			}

			offset += int(size)
		}
	}

	return h
}

func offsetToSizeMap(diskMap []byte) SizeMap {
	sizeMap := SizeMap{
		Free:  map[int]int{},
		Files: map[int]int{},
	}

	offset := 0
	for i := 0; i < len(diskMap); i += 2 {
		size := int(diskMap[i])
		sizeMap.Files[offset] = size
		offset += size

		if i+1 < len(diskMap) {
			size := int(diskMap[i+1])
			sizeMap.Free[offset] = size

			offset += size
		}
	}

	return sizeMap
}

func calculateChecksum(disk Disk) int {
	checksum := 0

	for i, id := range disk.Storage {
		if id == -1 {
			continue
		}

		checksum += i * id
	}

	return checksum
}

func deFragmentDisk(diskMap []byte, disk Disk, sizeMap SizeMap, h [maxFileSize]FreeSpaceHeap) {
	endOffsets := sizeMap.Reverse()

	for fileOffset := len(disk.Storage) - int(diskMap[len(diskMap)-1]); fileOffset > 0; {
		fileSize := sizeMap.Files[fileOffset]

		offset := len(disk.Storage)
		size := -1

		for i := fileSize; i <= maxFileSize; i++ {
			if len(h[i-1]) == 0 {
				continue
			}

			freeBlockOffset := h[i-1][0]
			if freeBlockOffset < offset {
				offset = freeBlockOffset
				size = i
			}
		}

		if offset < fileOffset {
			diff := size - fileSize

			for i := 0; i < fileSize; i++ {
				disk.Storage[fileOffset+i], disk.Storage[offset+i] = disk.Storage[offset+i], disk.Storage[fileOffset+i]
			}

			heap.Pop(&h[size-1])
			if diff > 0 {
				heap.Push(&h[diff-1], offset+fileSize)
			}
		}

		fileOffset -= endOffsets.Free[fileOffset]
		fileOffset -= endOffsets.Files[fileOffset]
	}
}

func fragmentDisk(diskMap []byte, disk Disk) {
	freeBlockIndex := 0
	freeBlockSize := int(diskMap[1])
	freeBlockShift := 0

	lastFreeBlockIndex := len(diskMap) - 2
	for diskMap[lastFreeBlockIndex] == 0 {
		lastFreeBlockIndex -= 2
	}

	lastFreeBlockSize := int(diskMap[lastFreeBlockIndex])

	start := int(diskMap[0])
	end := len(disk.Storage) - 1

	for start+freeBlockShift < len(disk.Storage)-disk.FreeSpace {
		if disk.Storage[end] == -1 {
			end -= lastFreeBlockSize

			lastFreeBlockIndex -= 2
			for diskMap[lastFreeBlockIndex] == 0 {
				lastFreeBlockIndex -= 2
			}

			lastFreeBlockSize = int(diskMap[lastFreeBlockIndex])
		}

		disk.Storage[start+freeBlockShift], disk.Storage[end] = disk.Storage[end], disk.Storage[start+freeBlockShift]

		freeBlockShift++
		if freeBlockShift >= freeBlockSize {
			freeBlockShift = 0

			freeBlockIndex++
			if freeBlockIndex*2+1 >= len(diskMap) {
				return
			}

			start += freeBlockSize + int(diskMap[freeBlockIndex*2])

			freeBlockSize = int(diskMap[freeBlockIndex*2+1])
			for freeBlockSize == 0 {
				freeBlockIndex++
				if freeBlockIndex*2+1 >= len(diskMap) {
					return
				}

				start += int(diskMap[freeBlockIndex*2])

				freeBlockSize = int(diskMap[freeBlockIndex*2+1])
			}
		}

		end--
	}

	return
}

type Disk struct {
	Storage   []int
	FreeSpace int
}

func mapToDisk(diskMap []byte) Disk {
	id := 0

	freeSpace := 0

	var disk []int
	for i := 0; i < len(diskMap); i += 2 {
		for range diskMap[i] {
			disk = append(disk, id)
		}

		id++

		if i+1 < len(diskMap) {
			freeBlockSize := diskMap[i+1]
			freeSpace += int(freeBlockSize)

			for range freeBlockSize {
				disk = append(disk, -1)
			}
		}
	}

	return Disk{
		Storage:   disk,
		FreeSpace: freeSpace,
	}
}

func readInput() ([]byte, error) {
	bytes, err := os.ReadFile("input/day-09.txt")
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile: %w", err)
	}

	for i := range bytes {
		bytes[i] -= '0'
	}

	return bytes, err
}
