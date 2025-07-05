package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/miekg/dns"
)

var verbose bool = false
var recursion bool = false

// IPScanner handles concurrent IP scanning with worker pools
type IPScanner struct {
	Workers   int
	Timeout   time.Duration
	DNSServer string
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewIPScanner creates a new scanner with specified workers, timeout, and DNS server
func NewIPScanner(workers int, timeout time.Duration, dnsServer string) *IPScanner {
	ctx, cancel := context.WithCancel(context.Background())
	return &IPScanner{
		Workers:   workers,
		Timeout:   timeout,
		DNSServer: dnsServer,
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Stop cancels the current scan
func (s *IPScanner) Stop() {
	s.cancel()
}

// generateIPs generates all IPs in a CIDR range
func (s *IPScanner) generateIPs(cidr string) ([]string, error) {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, fmt.Errorf("invalid CIDR: %v", err)
	}

	var ips []string
	var cpt = 0
	for ip := ipnet.IP.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
		cpt++
	}

	fmt.Printf("\nCIDR %s (%d IPs)\n", cidr, cpt)

	return ips, nil
}

// inc increments an IP address
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

// reverseResolve performs reverse DNS lookup on an IP using custom DNS server
func (s *IPScanner) reverseResolve(ip string) error {

	// Check if scan was cancelled
	select {
	case <-s.ctx.Done():
		return nil
	default:
	}

	c := new(dns.Client)
	m := new(dns.Msg)
	arpa, _ := dns.ReverseAddr(ip)
	m.SetQuestion(arpa, dns.TypePTR)
	m.RecursionDesired = recursion

	r, _, err := c.Exchange(m, s.DNSServer)
	if err != nil {
		return err
	}

	if r.Rcode == dns.RcodeNameError {
		return fmt.Errorf("NXDOMAIN - no local record")
	}

	for _, ans := range r.Answer {
		fmt.Printf("%s %v\n", ip, ans)
	}

	return nil
}

// worker processes IPs from the jobs channel
func (s *IPScanner) worker(jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case ip, ok := <-jobs:
			if !ok {
				return
			}
			err := s.reverseResolve(ip)
			if err != nil && verbose {
				fmt.Printf("%s %v\n", ip, err)
			}
		case <-s.ctx.Done():
			return
		}
	}
}

// ScanRange scans an IP range using worker pools
func (s *IPScanner) ScanRange(cidr string) error {
	ips, err := s.generateIPs(cidr)
	if err != nil {
		return err
	}

	jobs := make(chan string, len(ips))
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < s.Workers; i++ {
		wg.Add(1)
		go s.worker(jobs, &wg)
	}

	// Send jobs
	go func() {
		for _, ip := range ips {
			select {
			case jobs <- ip:
			case <-s.ctx.Done():
				break
			}
		}
		close(jobs)
	}()

	// Wait for all workers to complete
	wg.Wait()
	return nil
}

func main() {
	// Command line flags
	var (
		flag_cidr      = flag.String("cidr", "", "CIDR range to scan e.g: 192.168.0.0/16")
		flag_workers   = flag.Int("workers", 64, "Number of concurrent workers")
		flag_timeout   = flag.Duration("timeout", 2*time.Second, "Timeout for DNS lookups")
		flag_dnsServer = flag.String("dns", "", "DNS server (IP:port or hostname:port)")
		flag_verbose   = flag.Bool("verbose", false, "Show verbose output including failed lookups")
		flag_recursion = flag.Bool("recursion", false, "Ask srv to perform recursion")
	)

	flag.Parse()

	// DNS
	if *flag_dnsServer == "" {
		fmt.Printf("Missing DNS server\n")
		os.Exit(1)
	}

	// Create scanner
	scanner := NewIPScanner(*flag_workers, *flag_timeout, *flag_dnsServer)

	// Setup Enter key handling for graceful shutdown
	go func() {
		reader := bufio.NewReader(os.Stdin)
		reader.ReadLine()
		fmt.Println("Enter pressed, stopping scan...")
		scanner.Stop()
	}()

	fmt.Printf("Workers: %d, Timeout: %v, DNS: %s", *flag_workers, *flag_timeout, *flag_dnsServer)

	verbose = *flag_verbose
	if verbose {
		fmt.Printf(", Verbose")
	}

	recursion = *flag_recursion
	if recursion {
		fmt.Printf(", Recursive")
	}

	fmt.Printf("\nPress Enter to stop scanning...\n")

	// Start timing
	start := time.Now()

	// Scan the range
	if *flag_cidr == "" {
		for _, cidr := range []string{"192.168.0.0/16", "172.16.0.0/12", "10.0.0.0/8"} {
			err := scanner.ScanRange(cidr)
			if err != nil {
				fmt.Printf("Error scanning range: %v \n", err)
				os.Exit(1)
			}
		}
	} else if *flag_cidr == "k8s" {
		for _, cidr := range []string{"10.96.0.0/12", "10.100.0.0/16", "10.0.0.0/16", "172.20.0.0/16"} {
			err := scanner.ScanRange(cidr)
			if err != nil {
				fmt.Printf("Error scanning range: %v \n", err)
				os.Exit(1)
			}
		}
	} else {
		err := scanner.ScanRange(*flag_cidr)
		if err != nil {
			fmt.Printf("Error scanning range: %v \n", err)
			os.Exit(1)
		}
	}

	// Print timing
	elapsed := time.Since(start)
	fmt.Printf("\nScan completed in %v \n", elapsed)
}
