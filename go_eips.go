package main

import (
        "log"
        "os"
        "text/template"
    "net"
    "github.com/PraserX/ipconv"
)

type EgressIPs struct {
        IPAddrs []string
}

const tmpl = `EgressIPs:
{{range .IPAddrs}} - {{.}}
{{end}}
`

func renderTemplate(eips []string) {
        person := EgressIPs{
                IPAddrs: eips,
        }

        t := template.New("EgressIPs template")

        t, err := t.Parse(tmpl)
        if err != nil {
                log.Fatal("Parse: ", err)
                return
        }

        err = t.Execute(os.Stdout, person)
        if err != nil {
                log.Fatal("Execute: ", err)
                return
        }
}

func getEIPs(jobiteration uint32, numaddr uint32, addrSlice []uint32) []string{
	var i uint32 
	eips := []string{}

    	for i = 0; i < numaddr; i++ {
		eip := addrSlice[(jobiteration * numaddr) + i]
        	eips = append(eips, ipconv.IntToIPv4(eip).String())
	}
	return eips
}

func main() {
	eip_cidr := "192.168.1.1/24"
	var numIterations, addrPerIteration, i uint32 = 3, 2, 0
	addrSlice := make([]uint32, 0, (numIterations * addrPerIteration))
        baseAddr, _, err := net.ParseCIDR(eip_cidr)
	if err != nil {
                 log.Fatal("Error: ", err)
	 }
        baseAddrInt := ipconv.IPv4ToInt(baseAddr)

	nodeMap := make(map[uint32]bool)

	nodeips := []string{"192.168.1.2", "192.168.1.3"}
	for _, nodeip := range nodeips {
		nodeipuint32 := ipconv.IPv4ToInt(net.ParseIP(nodeip))
		nodeMap[nodeipuint32] = true	
	}
	for i = 0; i < ((numIterations * addrPerIteration) + uint32(len(nodeips)) ); i++ {
		if !nodeMap[baseAddrInt + i] {
			addrSlice = append(addrSlice, baseAddrInt + i)
		}
	}

	for i = 1; i < numIterations; i++ {
		eips := getEIPs(i, addrPerIteration, addrSlice)
		renderTemplate(eips)
	}
}
