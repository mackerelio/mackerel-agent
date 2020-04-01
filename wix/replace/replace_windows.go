package main

/*
 * 1. Fresh installation
 * => replace.exe <sample.conf> <conf> ___YOUR_API_KEY___ <apikey>
 *	Substitute ___YOUR_API_KEY___ by <apikey> in <sample.conf> then it write to <conf>.
 *
 * 2. Upgrade to x64 from x86
 * => replace.exe <sample.conf> <conf> ___YOUR_API_KEY___ '' (empty string)
 *	Just copy <conf> from under ProgramFiles(x86)\Mackerel folder to <conf> if it exists.
 *	Otherwise completely same as the case of fresh installation.
 *
 * 3. Upgrade from same architecture
 * => replace.exe <sample.conf> <conf> ___YOUR_API_KEY___ '' (empty string)
 *	Do nothing.
 *
 * Currently it is impossible to distinguish between the case 2 and 3 from the arguments.
 * Therefore replace.exe looks <conf> on the location of mackerel-agent.exe,
 * if exists, it will judge the state is in the case of 3.
 */

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 5 {
		log.Fatal("Usage: replace [in file] [out file] [old string] [new string]")
	}
	inFile := os.Args[1]
	outFile := os.Args[2]
	oldStr := os.Args[3]
	newStr := os.Args[4]

	if newStr == "" { // upgrade
		if !filepath.IsAbs(outFile) {
			// This program is usually called by the installer; safe.
			log.Fatalf("%q: it must be an absolute path", outFile)
		}
		dir := filepath.Dir(outFile)
		name := filepath.Base(outFile)
		if oldDir := FallbackConfigDir(dir); oldDir != "" {
			oldFile := filepath.Join(oldDir, name)
			if err := migrateFile(outFile, oldFile); err != nil {
				if !os.IsNotExist(err) {
					log.Fatalf("migrate %q to %q: %v", oldFile, outFile, err)
				}
				// Don't fallback; continue to the case 1.
				goto out
			}
			idFile := filepath.Join(dir, "id")
			oldIDFile := filepath.Join(oldDir, "id")
			if err := migrateFile(idFile, oldIDFile); err != nil {
				if !os.IsNotExist(err) {
					log.Fatalf("migrate %q to %q: %v", oldIDFile, idFile, err)
				}
			}
			os.Exit(0)
		}
	}

out:
	content, err := ioutil.ReadFile(inFile)
	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Stat(outFile)
	outFileIsExists := err == nil
	if !(outFileIsExists) {
		err = ioutil.WriteFile(outFile, []byte(strings.Replace(string(content), oldStr, newStr, -1)), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// migrateFile copies inFile to outFile if needed.
// If inFile is not exist, this will return os.ErrNotExist or its variants.
func migrateFile(outFile, inFile string) error {
	w, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return err
	}
	r, err := os.Open(inFile)
	if os.IsNotExist(err) {
		return err
	}
	_, err = io.Copy(w, r)
	return err
}

func isExist(file string) (bool, error) {
	f, err := os.Open(file)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	f.Close()
	return true, nil
}
