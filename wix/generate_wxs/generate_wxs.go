package main

import (
	"encoding/xml"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattn/go-encoding"
)

type Node struct {
	Name     xml.Name
	Attr     []xml.Attr
	Children []interface{}
}

func (n *Node) MarshalXML(e *xml.Encoder, s xml.StartElement) error {
	s.Name = n.Name
	s.Name.Space = ""
	s.Attr = n.Attr
	e.EncodeToken(s)
	for _, v := range n.Children {
		switch v.(type) {
		case xml.Comment:
			e.EncodeToken(v.(xml.Comment))
		case xml.CharData:
			e.EncodeToken(v.(xml.CharData))
		case *Node:
			if err := e.Encode(v.(*Node)); err != nil {
				return err
			}
		}
	}
	e.EncodeToken(s.End())
	return nil
}

func (n *Node) UnmarshalXML(d *xml.Decoder, s xml.StartElement) error {
	n.Name = s.Name
	n.Attr = s.Attr
	for {
		t, err := d.Token()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		switch v := t.(type) {
		case xml.CharData:
			n.Children = append(n.Children, t.(xml.CharData).Copy())
		case xml.Comment:
			n.Children = append(n.Children, t.(xml.Comment).Copy())
		case xml.StartElement:
			var c *Node
			if err := d.DecodeElement(&c, &v); err != nil {
				return err
			}
			n.Children = append(n.Children, c)
		}
	}
	return nil
}

func hasAttr(node *Node, name, value string) bool {
	for _, attr := range node.Attr {
		if attr.Name.Local == name && attr.Value == value {
			return true
		}
	}
	return false
}

func findPackage(node *Node) *Node {
	var ok bool
	for _, n := range node.Children {
		node, ok = n.(*Node)
		if !ok || node.Name.Local != "Product" {
			continue
		}
		for _, n := range node.Children {
			node, ok = n.(*Node)
			if !ok || node.Name.Local != "Package" {
				continue
			}
			return node
		}
	}
	return nil
}

func findInstallDir(node *Node) *Node {
	var ok bool
	for _, n := range node.Children {
		node, ok = n.(*Node)
		if !ok || node.Name.Local != "Product" {
			continue
		}
		for _, n := range node.Children {
			node, ok = n.(*Node)
			if !ok || node.Name.Local != "Directory" || !hasAttr(node, "Id", "TARGETDIR") {
				continue
			}
			for _, n := range node.Children {
				node, ok = n.(*Node)
				if !ok || node.Name.Local != "Directory" || !hasAttr(node, "Id", "ProgramFilesFolder") {
					continue
				}
				for _, n := range node.Children {
					node, ok = n.(*Node)
					if !ok || node.Name.Local != "Directory" || !hasAttr(node, "Id", "Mackerel") {
						continue
					}
					for _, n := range node.Children {
						node, ok = n.(*Node)
						if !ok || node.Name.Local != "Directory" || !hasAttr(node, "Id", "INSTALLDIR") {
							continue
						}
						return node
					}
				}
			}
		}
	}
	return nil
}

func fileNames(p string) ([]string, error) {
	dir, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	return dir.Readdirnames(-1)
}

func parseTemplate(n string) (*Node, error) {
	f, err := os.Open(n)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	v := &Node{}
	dec := xml.NewDecoder(f)
	dec.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		enc := encoding.GetEncoding(charset)
		if err == nil {
			return input, err
		}
		return enc.NewDecoder().Reader(input), nil
	}
	err = dec.Decode(&v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

var (
	pluginDir      = flag.String("pluginDir", "", "plugin directory")
	productVersion = flag.String("productVersion", "", "product version")
	templateFile   = flag.String("templateFile", "", "path to template file")
	outputFile     = flag.String("outputFile", "", "path to output file")
)

func main() {
	flag.Parse()
	if *pluginDir == "" || *productVersion == "" {
		flag.Usage()
		os.Exit(1)
	}

	names, err := fileNames(*pluginDir)
	if err != nil {
		log.Fatal(err)
	}

	v, err := parseTemplate(*templateFile)
	if err != nil {
		log.Fatal(err)
	}

	// update __VERSION__ to specified *productVersion
	pkg := findPackage(v)
	if pkg == nil {
		log.Fatal("__VERSION__ not found")
	}
	for _, attr := range pkg.Attr {
		if attr.Name.Local == "Version" {
			attr.Value = *productVersion
		}
	}

	// generate Component/File(s) from listing plugins on *pluginDir.
	installDir := findInstallDir(v)
	if installDir == nil {
		log.Fatal("INSTALLDIR not found")
	}
	installDir.Children = append(installDir.Children, xml.CharData("  "))
	component := new(Node)
	component.Name = xml.Name{Local: "Component", Space: ""}
	component.Attr = []xml.Attr{
		{Name: xml.Name{Local: "Id", Space: ""}, Value: "Plugins"},
	}
	component.Children = append(component.Children, xml.CharData("\n            "))
	for _, name := range names {
		if !strings.HasPrefix(name, "check-") {
			continue
		}
		fname := filepath.Join(*pluginDir, name)
		component.Children = append(component.Children, xml.CharData("  "))
		file := new(Node)
		file.Name = xml.Name{Local: "File", Space: ""}
		file.Attr = []xml.Attr{
			{Name: xml.Name{Local: "Id", Space: ""}, Value: name},
			{Name: xml.Name{Local: "Name", Space: ""}, Value: name},
			{Name: xml.Name{Local: "DiskId", Space: ""}, Value: "1"},
			{Name: xml.Name{Local: "Source", Space: ""}, Value: fname},
			{Name: xml.Name{Local: "KeyPath", Space: ""}, Value: "yes"},
		}
		component.Children = append(component.Children, file)
		component.Children = append(component.Children, xml.CharData("\n            "))
	}
	installDir.Children = append(installDir.Children, component)
	installDir.Children = append(installDir.Children, xml.CharData("\n          "))

	f, err := os.Create(*outputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	err = xml.NewEncoder(f).Encode(v)
	if err != nil {
		log.Fatal(err)
	}
}
