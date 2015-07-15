package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type Conn struct {
	aw2 *ec2.EC2
}

// create a struct to map toml config file

type AwsConfig struct {
	AwsSecretKey string `toml:"AWS_ACCESS_KEY_ID"`
	AwsAccessKey string `toml:"AWS_SECRET_ACCESS_KEY"`
	Region       string `toml:"AWS_REGION"`
}

func recieveStatus(dataMap map[string]string) <-chan string {

	c := make(chan string)
	go func() {
		for _, v := range dataMap {
			c <- fmt.Sprintf("%s", v)
			//			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()
	return c

}

func (c *Conn) iterateResToMap(resp *ec2.DescribeInstancesOutput) map[string]string {
	insMap := make(map[string]string)
	fmt.Println("> Number of reservation sets: ", len(resp.Reservations))
	for idx, _ := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			// fmt.Printf("   Instance State: %v InstanceID: %v \n", *inst.State.Name, *inst.InstanceID)
			// dereference pointer
			var id, state string
			id = *inst.PrivateDNSName
			state = *inst.State.Name
			statstr := state + ":" + id

			insMap[id] = statstr
		}
	}
	return insMap
}

func (c *Conn) startLoop() bool {
	// Call the DescribeInstances Operation
	resp, err := c.aw2.DescribeInstances(nil)
	if err != nil {
		log.Fatal(err)
	}
	// re-format to method call
	newMap := c.iterateResToMap(resp)

	r := recieveStatus(newMap)

	// set a timeout for the channel
	timeout := time.After(5 * time.Second)

	// begin channel operations
	// TODO: this should have interfaces and structs and a poller
	for {
		select {
		case result := <-r:
			matched, _ := regexp.MatchString("stopped", result)
			if matched {
				res := strings.Split(result, ":")
				status := res[0]
				ip := res[1]
				fmt.Printf("Trannsisioned to %v taking some actions for %v  \n", status, ip)
			}
		case <-timeout:
			fmt.Println("You took too long")
			return false
		}
	}

}

func NewEc2() *Conn {
	c := new(Conn)
	c.aw2 = ec2.New(&aws.Config{Region: "us-west-2"})

	return c
}

func main() {
	// Create an EC2 service object in the "us-west-2" region
	// Note that you can also configure your region globally by
	// exporting the AWS_REGION environment variable
	c := NewEc2()

	/* here is where the crazyness happens
	   we sit in a for loop waiting to hit timetout of the channel
	   if we do we will start again (hoopefully getting new data)
	*/

	for {
		if c.startLoop() {
			fmt.Println("starting up again")
			c.startLoop()
		}
	}

}
