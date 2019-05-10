package main

import (
	"crypto/md5"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"github.com/mattn/go-encoding"
)

// Node is HTML node traversing to convert to text.
type Node struct {
	Name      xml.Name
	Attr      []xml.Attr
	Children  []interface{}
	RootSpace map[string]string
}

// MarshalXML encode the Node as XML emement. This handle only Comment and
// CharData. Any other values are not converted.
func (n *Node) MarshalXML(e *xml.Encoder, s xml.StartElement) error {
	s.Name = n.Name
	s.Attr = n.Attr

	if ns, ok := n.RootSpace[s.Name.Space]; ok {
		s.Name.Local = ns + ":" + s.Name.Local
	}
	s.Name.Space = ""

	var newattr []xml.Attr
	for _, attr := range s.Attr {
		if ns, ok := n.RootSpace[attr.Name.Space]; ok {
			attr.Name.Local = ns + ":" + attr.Name.Local
		} else if attr.Name.Space != "" {
			attr.Name.Local = attr.Name.Space + ":" + attr.Name.Local
		}
		attr.Name.Space = ""
		newattr = append(newattr, attr)
	}
	s.Attr = newattr

	e.EncodeToken(s)
	for _, v := range n.Children {
		switch v.(type) {
		case xml.Comment:
			e.EncodeToken(v.(xml.Comment))
		case xml.CharData:
			e.EncodeToken(v.(xml.CharData))
		case *Node:
			v.(*Node).RootSpace = n.RootSpace
			if err := e.Encode(v.(*Node)); err != nil {
				return err
			}
		}
	}
	e.EncodeToken(s.End())
	return nil
}

// UnmarshalXML decodes as single XML element to read content.
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
}

func hasAttr(node *Node, name, value string) bool {
	for _, attr := range node.Attr {
		if attr.Name.Local == name && attr.Value == value {
			return true
		}
	}
	return false
}

func toCamelCase(name string) string {
	out := ""
	rs := []rune(name)
	out += string(unicode.ToUpper(rs[0]))
	for i := 1; i < len(rs); i++ {
		if i < len(rs)-1 && (rs[i] == '.' || rs[i] == '-') {
			i++
			out += string(unicode.ToUpper(rs[i]))
		} else {
			out += string(rs[i])
		}
	}
	return out
}

func updateElementAttr(node *Node, name, attr, oldValue, newValue string) bool {
	nodes := findElementAll(node, name)
	var updated bool
	for _, p := range nodes {
		for i := range p.Attr {
			if p.Attr[i].Name.Local == attr && p.Attr[i].Value == oldValue {
				p.Attr[i].Value = newValue
				updated = true
			}
		}
	}
	return updated
}

func findElementAll(node *Node, name string) []*Node {
	var nodes []*Node
	for _, n := range node.Children {
		p, ok := n.(*Node)
		if !ok {
			continue
		}
		if p.Name.Local == name {
			nodes = append(nodes, p)
		}
		nodes = append(nodes, findElementAll(p, name)...)
	}
	return nodes
}

func findInstallDir(node *Node, id string) *Node {
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
				if !ok || node.Name.Local != "Directory" || !hasAttr(node, "Id", id) {
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

// Platform represents architecture configuration.
type Platform string

// ProgramFiles returns ID for Program Files folder.
func (a *Platform) ProgramFiles() string {
	switch *a {
	case "x64":
		return "ProgramFiles64Folder"
	case "x86":
		return "ProgramFilesFolder"
	default:
		log.Fatalln("unknown platform:", *a)
		return ""
	}
}

// Win64 returns yes if target platform is 64bit machine. otherwise no.
func (a *Platform) Win64() string {
	switch *a {
	case "x64":
		return "yes"
	case "x86":
		return "no"
	default:
		log.Fatalln("unknown platform:", *a)
		return ""
	}
}

// String implements flag.Value interface.
func (a *Platform) String() string {
	return string(*a)
}

// Set implements flag.Value interface.
func (a *Platform) Set(s string) error {
	switch s {
	case "x86", "x64":
		*a = Platform(s)
		return nil
	default:
		return fmt.Errorf("unknown platform: %s", s)
	}
}

var (
	buildDir       = flag.String("buildDir", "", "build directory")
	productVersion = flag.String("productVersion", "", "product version")
	templateFile   = flag.String("templateFile", "", "path to template file")
	outputFile     = flag.String("outputFile", "", "path to output file")

	platform = Platform("x86")
)

func init() {
	flag.Var(&platform, "platform", "target platform")
}

func main() {
	flag.Parse()
	if *buildDir == "" || *productVersion == "" {
		flag.Usage()
		os.Exit(1)
	}

	names, err := fileNames(*buildDir)
	if err != nil {
		log.Fatal(err)
	}

	v, err := parseTemplate(*templateFile)
	if err != nil {
		log.Fatal(err)
	}

	ok := updateElementAttr(v, "Product", "Version", "___VERSION___", *productVersion)
	if !ok {
		log.Fatal("___VERSION___ not found")
	}
	ok = updateElementAttr(v, "Package", "Platform", "___PLATFORM___", platform.String())
	if !ok {
		log.Fatal("___PLATFORM___ not found")
	}
	ok = updateElementAttr(v, "Directory", "Id", "___PROGRAMFILES___", platform.ProgramFiles())
	if !ok {
		log.Fatal("___PROGRAMFILES___ not found")
	}
	ok = updateElementAttr(v, "Component", "Win64", "___WIN64___", platform.Win64())
	if !ok {
		log.Fatal("___WIN64___ not found")
	}

	// generate Component/File(s) from listing plugins on *buildDir.
	installDir := findInstallDir(v, platform.ProgramFiles())
	if installDir == nil {
		log.Fatal("INSTALLDIR not found")
	}
	installDir.Children = append(installDir.Children, xml.CharData("  "))

	component := new(Node)
	component.Name = xml.Name{Local: "Component", Space: ""}
	idlist := []string{}
	installDir.Children = append(installDir.Children, component)
	for _, name := range names {
		if !strings.HasPrefix(name, "check-") && !strings.HasPrefix(name, "mackerel-plugin-") && name != "mkr.exe" {
			continue
		}
		id := toCamelCase(name)
		fname := filepath.Join(*buildDir, name)
		file := new(Node)
		file.Name = xml.Name{Local: "File", Space: ""}
		file.Attr = []xml.Attr{
			{Name: xml.Name{Local: "Id", Space: ""}, Value: id},
			{Name: xml.Name{Local: "Name", Space: ""}, Value: name},
			{Name: xml.Name{Local: "DiskId", Space: ""}, Value: "1"},
			{Name: xml.Name{Local: "Source", Space: ""}, Value: fname},
		}
		component.Children = append(component.Children, xml.CharData("\n              "))
		component.Children = append(component.Children, file)
		idlist = append(idlist, id)
	}
	sort.Strings(idlist)
	b := md5.Sum([]byte(strings.Join(idlist, "\n")))
	component.Attr = []xml.Attr{
		{Name: xml.Name{Local: "Id", Space: ""}, Value: "Plugins"},
		{Name: xml.Name{Local: "Guid", Space: ""}, Value: fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])},
	}
	component.Children = append(component.Children, xml.CharData("\n            "))
	installDir.Children = append(installDir.Children, xml.CharData("\n          "))

	f, err := os.Create(*outputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	v.RootSpace = make(map[string]string)
	for _, attr := range v.Attr {
		if attr.Name.Local != "xmlns" {
			v.RootSpace[attr.Value] = attr.Name.Local
		}
	}

	f.Write([]byte("<?xml version=\"1.0\" encoding=\"windows-1252\"?>\n"))
	err = xml.NewEncoder(f).Encode(v)
	if err != nil {
		log.Fatal(err)
	}
}
