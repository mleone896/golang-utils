package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

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

func iterateResToMap(resp *ec2.DescribeInstancesOutput) map[string]string {
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

func main() {
	// Create an EC2 service object in the "us-west-2" region
	// Note that you can also configure your region globally by
	// exporting the AWS_REGION environment variable
	svc := ec2.New(&aws.Config{Region: "us-west-2"})

	// Call the DescribeInstances Operation
	resp, err := svc.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}
	// re-format to method call
	newMap := iterateResToMap(resp)

	r := recieveStatus(newMap)

	for {
		select {
		case result := <-r:
			if result == "stopped" {
				fmt.Printf("Trannsisioned to %v taking some actions \n", result)
			}
		}
	}

}
