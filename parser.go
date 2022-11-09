package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	aeroconf "github.com/rglonek/aerospike-config-file-parser"
)

func main() {
	if len(os.Args) < 4 {
		help()
	}
	command := os.Args[1]
	path := os.Args[2]
	filename := os.Args[len(os.Args)-1]
	setValues := []string{""}

	switch command {
	case "delete":
		if len(os.Args) != 4 {
			help()
		}
	case "set":
		if len(os.Args) > 4 {
			setValues = os.Args[3 : len(os.Args)-1]
		}
	case "create":
		if len(os.Args) != 4 {
			help()
		}
	}

	s, err := aeroconf.ParseFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	sa := s
	pathn := strings.Split(path, ".")
	if pathn[0] == "" && len(pathn) > 1 {
		pathn = pathn[1:]
	}
	switch command {
	case "delete":
		for _, i := range pathn[0 : len(pathn)-1] {
			sa = sa.Stanza(i)
			if sa == nil {
				log.Fatal("Stanza not found")
			}
		}
		err = sa.Delete(pathn[len(pathn)-1])
		if err != nil {
			log.Fatal(err)
		}
	case "set":
		for _, i := range pathn[0 : len(pathn)-1] {
			sa = sa.Stanza(i)
			if sa == nil {
				log.Fatal("Stanza not found")
			}
		}
		err = sa.SetValues(pathn[len(pathn)-1], aeroconf.SliceToValues(setValues))
		if err != nil {
			log.Fatal(err)
		}
	case "create":
		for _, i := range pathn {
			if sa.Stanza(i) == nil {
				err = sa.NewStanza(i)
				if err != nil {
					log.Fatal(err)
				}
			}
			sa = sa.Stanza(i)
		}
	}

	err = s.WriteFile(filename, "", "    ", true)
	if err != nil {
		log.Fatal(err)
	}
}

func help() {
	fmt.Printf("\nUsage: %s command path [set-value1] [set-value2] [...set-valueX] filename\n", os.Args[0])
	fmt.Println("\n" + `Commands:
	delete - delete configuration/stanza
	set    - set configuration parameter
	create - create a new stanza`)
	fmt.Println("\n" + `Path: .path.to.item or .path.to.stanza, e.g. .network.heartbeat`)
	fmt.Println("\n" + `Set-value: for the 'set' command - used to specify value of parameter; leave empty to crete no-value param`)
	fmt.Printf("\n"+`Example:
	touch new.conf
	%s create network.heartbeat new.conf
	%s set network.heartbeat.mode mesh new.conf
	%s set network.heartbeat.mesh-seed-address-port "172.17.0.2 3000" "172.17.0.3 3000" new.conf
	%s create service new.conf
	%s set service.proto-fd-max 3000 new.conf
	`+"\n", os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
	os.Exit(1)
}
