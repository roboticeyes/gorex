package rex

import (
	"bytes"
	b64 "encoding/base64"
	"image"
	_ "image/png"
	"strings"
	"testing"
)

func TestWriteImage(t *testing.T) {

	b, err := b64.StdEncoding.DecodeString(testImage)
	if err != nil {
		panic(err)
	}
	img := Image{
		ID:          11,
		Compression: png,
		Data:        b,
	}

	var buf []byte
	w := bytes.NewBuffer(buf)
	err = img.Write(w)
	if err != nil {
		t.Fatal("Error: ", err)
	}
}

func TestReadImage(t *testing.T) {

	r := b64.NewDecoder(b64.StdEncoding, strings.NewReader(string(rexImageBlock)))
	// read header
	hdr, err := ReadDataBlockHeader(r)
	if err != nil {
		t.Fatal("Cannot read header")
	}
	if hdr.ID != 11 || hdr.Type != typeImage || hdr.Version != 1 || hdr.Size != 182 {
		t.Fatalf("Header has unexpected data: %v", hdr)
	}

	img, err := ReadImage(r, hdr)

	if err != nil {
		t.Fatal("Error: ", err)
	}
	if img.ID != hdr.ID {
		t.Fatal("ID does not match")
	}
	if img.Compression != png {
		t.Fatal("Compression does not match")
	}

	if len(img.Data) != 178 {
		t.Fatalf("Size does not match expected=178 actual=%d\n", len(img.Data))
	}

	imgReader := bytes.NewReader(img.Data)
	m, _, err := image.Decode(imgReader)
	if err != nil {
		t.Fatal(err)
	}
	if m.Bounds().Dx() != 32 {
		t.Fatal("Image size does  not match")
	}
}

const testImage = `
iVBORw0KGgoAAAANSUhEUgAAACAAAAAgCAIAAAD8GO2jAAAACXBIWXMAAC4jAAAuIwF4pT92AAAA
B3RJTUUH4wMdCB4gWOt4sQAAABl0RVh0Q29tbWVudABDcmVhdGVkIHdpdGggR0lNUFeBDhcAAAAs
SURBVEjH7c1BAQBABAAwrn8gLTSQ5Urw2wospysuvTgmEAgEAoFAIBAItnyV3gKlvdJsuwAAAABJ
RU5ErkJggg==
`

const rexImageBlock = `
BAABALYAAAALAAAAAAAAAAIAAACJUE5HDQoaCgAAAA1JSERSAAAAIAAAACAIAgAAAPwY7aMAAAAJ
cEhZcwAALiMAAC4jAXilP3YAAAAHdElNRQfjAx0IHiBY63ixAAAAGXRFWHRDb21tZW50AENyZWF0
ZWQgd2l0aCBHSU1QV4EOFwAAACxJREFUSMftzUEBAEAEADCufyAtNJDlSvDbCiynKy69OCYQCAQC
gUAgEAi2fJXeAqW90my7AAAAAElFTkSuQmCC
`
