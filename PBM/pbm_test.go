package Netpbm

import (
	"os"
	"testing"
)

const imageWidth = 15
const imageHeight = 15

var imageDataP1 = []bool{
	false, false, false, false, false, false, false, true, true, true, true, false, false, false, false,
	false, false, false, false, false, false, true, false, false, false, false, true, false, false, false,
	false, false, false, false, false, true, false, false, false, false, false, false, true, false, false,
	false, false, false, false, false, true, false, false, false, false, true, false, true, true, true,
	false, false, false, false, false, true, false, false, false, false, false, false, false, false, true,
	false, false, false, false, false, true, true, false, false, false, false, false, false, true, true,
	true, true, false, false, false, false, true, true, false, false, false, true, true, true, false,
	true, false, true, true, false, false, false, true, false, false, false, true, false, false, false,
	true, false, false, false, true, true, true, false, false, false, false, true, false, false, false,
	true, false, false, false, false, false, false, false, false, false, false, false, true, false, false,
	true, false, false, false, false, false, false, false, false, false, false, false, true, false, false,
	false, true, false, false, false, false, false, false, false, false, false, false, true, false, false,
	false, true, true, false, false, false, false, false, false, false, false, true, false, false, false,
	false, false, true, true, false, false, false, false, false, true, true, false, false, false, false,
	false, false, false, false, true, true, true, true, true, true, false, false, false, false, false,
}

var imageDataInvert = []bool{
	true, true, true, true, true, true, true, false, false, false, false, true, true, true, true,
	true, true, true, true, true, true, false, true, true, true, true, false, true, true, true,
	true, true, true, true, true, false, true, true, true, true, true, true, false, true, true,
	true, true, true, true, true, false, true, true, true, true, false, true, false, false, false,
	true, true, true, true, true, false, true, true, true, true, true, true, true, true, false,
	true, true, true, true, true, false, false, true, true, true, true, true, true, false, false,
	false, false, true, true, true, true, false, false, true, true, true, false, false, false, true,
	false, true, false, false, true, true, true, false, true, true, true, false, true, true, true,
	false, true, true, true, false, false, false, true, true, true, true, false, true, true, true,
	false, true, true, true, true, true, true, true, true, true, true, true, false, true, true,
	false, true, true, true, true, true, true, true, true, true, true, true, false, true, true,
	true, false, true, true, true, true, true, true, true, true, true, true, false, true, true,
	true, false, false, true, true, true, true, true, true, true, true, false, true, true, true,
	true, true, false, false, true, true, true, true, true, false, false, true, true, true, true,
	true, true, true, true, false, false, false, false, false, false, true, true, true, true, true,
}

var imageDataFlip = []bool{
	false, false, false, false, true, true, true, true, false, false, false, false, false, false, false,
	false, false, false, true, false, false, false, false, true, false, false, false, false, false, false,
	false, false, true, false, false, false, false, false, false, true, false, false, false, false, false,
	true, true, true, false, true, false, false, false, false, true, false, false, false, false, false,
	true, false, false, false, false, false, false, false, false, true, false, false, false, false, false,
	true, true, false, false, false, false, false, false, true, true, false, false, false, false, false,
	false, true, true, true, false, false, false, true, true, false, false, false, false, true, true,
	false, false, false, true, false, false, false, true, false, false, false, true, true, false, true,
	false, false, false, true, false, false, false, false, true, true, true, false, false, false, true,
	false, false, true, false, false, false, false, false, false, false, false, false, false, false, true,
	false, false, true, false, false, false, false, false, false, false, false, false, false, false, true,
	false, false, true, false, false, false, false, false, false, false, false, false, false, true, false,
	false, false, false, true, false, false, false, false, false, false, false, false, true, true, false,
	false, false, false, false, true, true, false, false, false, false, false, true, true, false, false,
	false, false, false, false, false, true, true, true, true, true, true, false, false, false, false,
}

var imageDataFlop = []bool{
	false, false, false, false, true, true, true, true, true, true, false, false, false, false, false,
	false, false, true, true, false, false, false, false, false, true, true, false, false, false, false,
	false, true, true, false, false, false, false, false, false, false, false, true, false, false, false,
	false, true, false, false, false, false, false, false, false, false, false, false, true, false, false,
	true, false, false, false, false, false, false, false, false, false, false, false, true, false, false,
	true, false, false, false, false, false, false, false, false, false, false, false, true, false, false,
	true, false, false, false, true, true, true, false, false, false, false, true, false, false, false,
	true, false, true, true, false, false, false, true, false, false, false, true, false, false, false,
	true, true, false, false, false, false, true, true, false, false, false, true, true, true, false,
	false, false, false, false, false, true, true, false, false, false, false, false, false, true, true,
	false, false, false, false, false, true, false, false, false, false, false, false, false, false, true,
	false, false, false, false, false, true, false, false, false, false, true, false, true, true, true,
	false, false, false, false, false, true, false, false, false, false, false, false, true, false, false,
	false, false, false, false, false, false, true, false, false, false, false, true, false, false, false,
	false, false, false, false, false, false, false, true, true, true, true, false, false, false, false,
}

func TestReadPBM(t *testing.T) {

	// read the image with P1 magic number
	pbm, err := ReadPBM("./testImages/pbm/testP1.pbm")
	if err != nil {
		t.Error(err)
	}
	// check the magic number
	if pbm.magicNumber != "P1" {
		t.Error("Wrong magic number")
	}

	if pbm.width != 15 {
		t.Error("Wrong width")
	}
	if pbm.height != 15 {
		t.Error("Wrong height")
	}

	// compare the data
	for i := 0; i < imageWidth*imageHeight; i++ {
		var x = i % imageWidth
		var y = i / imageWidth
		if pbm.data[y][x] != imageDataP1[i] {
			t.Error("Wrong data")
		}
	}

	// read the image with P4 magic number
	pbm, err = ReadPBM("./testImages/pbm/testP4.pbm")
	if err != nil {
		t.Error(err)
	}
	// check the magic number
	if pbm.magicNumber != "P4" {
		t.Error("Wrong magic number")
	}
	if pbm.width != 15 {
		t.Error("Wrong width")
	}
	if pbm.height != 15 {
		t.Error("Wrong height")
	}

	// compare the data
	for i := 0; i < imageWidth*imageHeight; i++ {
		var x = i % imageWidth
		var y = i / imageWidth
		if pbm.data[y][x] != imageDataP1[i] {
			t.Error("Wrong data")
		}
	}
}

func TestSize(t *testing.T) {
	pbm, err := ReadPBM("./testImages/pbm/testP1.pbm")
	if err != nil {
		t.Error(err)
	}
	w, h := pbm.Size()
	if w != imageWidth || h != imageHeight {
		t.Error("Wrong size")
	}
}

func TestAt(t *testing.T) {
	pbm, err := ReadPBM("./testImages/pbm/testP1.pbm")
	if err != nil {
		t.Error(err)
	}
	if pbm.At(0, 8) != true {
		t.Error("Wrong value")
	}
}

func TestSet(t *testing.T) {
	pbm, err := ReadPBM("./testImages/pbm/testP1.pbm")
	if err != nil {
		t.Error(err)
	}
	pbm.Set(1, 3, true)
	if pbm.At(1, 3) != true {
		t.Error("Wrong value")
	}
}

func TestSave(t *testing.T) {
	pbm, err := ReadPBM("./testImages/pbm/testP1.pbm")
	if err != nil {
		t.Error(err)
	}
	pbm.SetMagicNumber("P1")
	err = pbm.Save("./testImages/pbm/testP1Save.pbm")
	if err != nil {
		t.Error(err)
	}
	pbm2, err := ReadPBM("./testImages/pbm/testP1Save.pbm")
	if err != nil {
		t.Error(err)
	}
	if pbm2.magicNumber != "P1" {
		t.Error("Wrong magic number")
	}
	if pbm2.width != 15 {
		t.Error("Wrong width")
	}
	if pbm2.height != 15 {
		t.Error("Wrong height")
	}
	// compare the data
	for i := 0; i < imageWidth*imageHeight; i++ {
		var x = i % imageWidth
		var y = i / imageWidth
		if pbm2.data[y][x] != imageDataP1[i] {
			t.Error("Wrong data")
		}
	}

	pbm, err = ReadPBM("./testImages/pbm/testP4.pbm")
	if err != nil {
		t.Error(err)
	}
	pbm.SetMagicNumber("P4")
	err = pbm.Save("./testImages/pbm/testP4Save.pbm")
	if err != nil {
		t.Error(err)
	}
	pbm2, err = ReadPBM("./testImages/pbm/testP4Save.pbm")
	if err != nil {
		t.Error(err)
	}
	if pbm2.magicNumber != "P4" {
		t.Error("Wrong magic number")
	}
	if pbm2.width != 15 {
		t.Error("Wrong width")
	}
	if pbm2.height != 15 {
		t.Error("Wrong height")
	}
	// compare the data
	for i := 0; i < imageWidth*imageHeight; i++ {
		var x = i % imageWidth
		var y = i / imageWidth
		if pbm2.data[y][x] != imageDataP1[i] {
			t.Error("Wrong data")
		}
	}
	// remove the test files
	err = os.Remove("./testImages/pbm/testP1Save.pbm")
	if err != nil {
		t.Error(err)
	}
	err = os.Remove("./testImages/pbm/testP4Save.pbm")
	if err != nil {
		t.Error(err)
	}
}

func TestInvert(t *testing.T) {
	pbm, err := ReadPBM("./testImages/pbm/testP1.pbm")
	if err != nil {
		t.Error(err)
	}
	pbm.Invert()
	// compare the data
	for i := 0; i < imageWidth*imageHeight; i++ {
		var x = i % imageWidth
		var y = i / imageWidth
		if pbm.data[y][x] != imageDataInvert[i] {
			t.Error("Wrong data")
		}
	}
}

func TestFlip(t *testing.T) {
	pbm, err := ReadPBM("./testImages/pbm/testP1.pbm")
	if err != nil {
		t.Error(err)
	}
	pbm.Flip()
	// compare the data
	for i := 0; i < imageWidth*imageHeight; i++ {
		var x = i % imageWidth
		var y = i / imageWidth
		if pbm.data[y][x] != imageDataFlip[i] {
			t.Error("Wrong data")
		}
	}
}

func TestFlop(t *testing.T) {
	pbm, err := ReadPBM("./testImages/pbm/testP1.pbm")
	if err != nil {
		t.Error(err)
	}
	pbm.Flop()
	// compare the data
	for i := 0; i < imageWidth*imageHeight; i++ {
		var x = i % imageWidth
		var y = i / imageWidth
		if pbm.data[y][x] != imageDataFlop[i] {
			t.Error("Wrong data")
		}
	}
}

func TestSetMagicNumber(t *testing.T) {
	pbm, err := ReadPBM("./testImages/pbm/testP1.pbm")
	if err != nil {
		t.Error(err)
	}
	pbm.SetMagicNumber("P4")
	if pbm.magicNumber != "P4" {
		t.Error("Wrong magic number")
	}
}
