package main

import (
	"fmt"
	"flag"
        "log"
        "os"
	"strings"
	"strconv"
        "text/template"
    	"net"
	"github.com/praserx/ipconv"
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

func getEIPs(jobiteration int, numaddr int, addrSlice []string) []string{
	eips := []string{}

	for i := 0; i < numaddr; i++ {
		eips = append(eips, addrSlice[(jobiteration * numaddr) + i])
	}
	return eips
}

// kube-burner will add EIP addresses in EIP object template.
// Each job iteration (i.e a namespace) will have one EIP object
func kube_burner() {
	// Export environment variables for kube-burner
    	numJobIterations, err := strconv.Atoi(os.Getenv("JOB_ITERATIONS"))
    	if err != nil {
        	fmt.Println("Error converting to integer:", err)
	        return
    	}	
	addressesPerIteration, err := strconv.Atoi(os.Getenv("ADDRESSES_PER_ITERATION"))
    	if err != nil {
        	fmt.Println("Error converting to integer:", err)
	        return
    	}	
	eipAddresses := os.Getenv("EIP_ADDRESSES")

	addrSlice := strings.Split(eipAddresses, " ")

    	// render template for each job
	for i := 0; i < numJobIterations; i++ {
		renderTemplate(getEIPs(i, addressesPerIteration, addrSlice))
	}
}

// generate addresses and export to kube-burner
func kube_burner_ocp(numJobIterations int, addressesPerIteration int, excludeAddresses string) {

	// kube-burner-ocp calculates egressCIDR and nodeips from ocp cluster. Hard coded here for testing
	egressCIDR := "192.168.1.0/24"
	nodeips := []string{"192.168.1.2", "192.168.1.3"}

	if excludeAddresses != "" {
		nodeips = append(nodeips, strings.Split(excludeAddresses, " ")...)
	}

	addrSlice := make([]string, 0, (numJobIterations * addressesPerIteration))
        baseAddr, _, err := net.ParseCIDR(egressCIDR)
	if err != nil {
                 log.Fatal("Error: ", err)
	 }
        baseAddrInt, err := ipconv.IPv4ToInt(baseAddr)
	if err != nil {
                 log.Fatal("Error: ", err)
	 }

	// map to store nodeips
	nodeMap := make(map[uint32]bool)
	for _, nodeip := range nodeips {
		nodeipuint32, err := ipconv.IPv4ToInt(net.ParseIP(nodeip))
		if err != nil {
                	log.Fatal("Error: ", err)
	 	}
		nodeMap[nodeipuint32] = true	
	}

	// Generate ip addresses from CIDR by excluding nodeips
	var newAddr uint32
	for i := 0; i < ((numJobIterations * addressesPerIteration) + len(nodeips) ); i++ {
		newAddr = baseAddrInt + uint32(i)
		if !nodeMap[newAddr] {
			addrSlice = append(addrSlice, ipconv.IntToIPv4(newAddr).String())
		}
	}

	// Export environment variables for kube-burner
	os.Setenv("JOB_ITERATIONS", fmt.Sprint(numJobIterations))
	os.Setenv("ADDRESSES_PER_ITERATION", fmt.Sprint(addressesPerIteration))

    	// combine all addresses to a string and export as an environment variable
	envVarName := "EIP_ADDRESSES"
	os.Setenv(envVarName, strings.Join(addrSlice, " "))
}

func main() {
	iterations := flag.Int("iterations", 1, "Number of Job iterations")
	addressesPerIteration := flag.Int("addresses-per-iteration", 1, "EIP address per Iteration")
	excludeAddresses := flag.String("exclude-addresses", "", "List of addresses to exclude for EIP. Example '192.168.1.0 192.168.1.1'")
	flag.Parse()

	// generate addresses and export to kube-burner
	kube_burner_ocp(*iterations, *addressesPerIteration, *excludeAddresses)

	// render template from provided addresses
	kube_burner()
}
