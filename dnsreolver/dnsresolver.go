package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

// vars
var hostFile = flag.String("hostFile", "./hostfile", "The path to the cloudtrax host file")
var outFile = flag.String("outFile", "./outfile", "The path to the cloudtrax host file")

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(string(path))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func ReturnIpMap(hh []string) map[string]string {
	hostMap := make(map[string]string)

	// loop over the hosts to resolve to IP
	for _, host := range hh {
		DOMAINcname, err := net.LookupIP(host)
		if err != nil {
			log.Fatal(err)
		}

		ipformat := strings.Trim(fmt.Sprintf("%s", DOMAINcname), "[]")

		hostMap[host] = ipformat
	}

	return hostMap

}

func WriteToFile(mm map[string]string, ofile string) {

	fo, err := os.Create(ofile)
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	for k, v := range mm {
		fo.WriteString(k + "\n")
		v = strings.Replace(v, " ", "\n", -1)
		fo.WriteString(v + "\n")

	}

}

func main() {
	flag.Parse() // parse command line flags
	var convertHosts string
	var convertOutFile string
	convertOutFile = *outFile
	convertHosts = *hostFile // dereference pointer
	hosts, _ := ReadLines(convertHosts)
	ipset := ReturnIpMap(hosts)

	WriteToFile(ipset, convertOutFile)

}
