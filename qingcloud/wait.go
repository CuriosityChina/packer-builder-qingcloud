package qingcloud

import (
	"fmt"
	"log"
	"time"

	qc "github.com/yunify/qingcloud-sdk-go/service"
)

func waitForInstanceState(
	desiredState string, id string,
	client *qc.InstanceService, timeout time.Duration) error {
	done := make(chan struct{})
	defer close(done)

	result := make(chan error, 1)
	go func() {
		attempts := 0
		for {
			attempts++

			log.Printf("Checking instance status... (attempt: %d)", attempts)
			results, err := client.DescribeInstances(&qc.DescribeInstancesInput{
				Instances: []*string{qc.String(id)},
			})

			if err != nil {
				log.Printf("%s", err)
				result <- err
				return
			}

			if *results.InstanceSet[0].Status == desiredState {
				result <- nil
				return
			}

			// Wait 3 seconds in between
			time.Sleep(3 * time.Second)

			// Verify we shouldn't exit
			select {
			case <-done:
				// We finished, so just exit the goroutine
				return
			default:
				// Keep going
			}
		}
	}()

	log.Printf("Waiting for up to %d seconds for eip to become %s", timeout/time.Second, desiredState)
	select {
	case err := <-result:
		return err
	case <-time.After(timeout):
		err := fmt.Errorf("Timeout while waiting to for eip to become '%s'", desiredState)
		return err
	}
}

func waitForEIPState(
	desiredState string, id string,
	client *qc.EIPService, timeout time.Duration) error {
	done := make(chan struct{})
	defer close(done)

	result := make(chan error, 1)
	go func() {
		attempts := 0
		for {
			attempts++

			log.Printf("Checking eip status... (attempt: %d)", attempts)
			results, err := client.DescribeEIPs(&qc.DescribeEIPsInput{
				EIPs: []*string{qc.String(id)},
			})

			if err != nil {
				log.Printf("%s", err)
				result <- err
				return
			}

			if *results.EIPSet[0].Status == desiredState {
				result <- nil
				return
			}

			// Wait 3 seconds in between
			time.Sleep(3 * time.Second)

			// Verify we shouldn't exit
			select {
			case <-done:
				// We finished, so just exit the goroutine
				return
			default:
				// Keep going
			}
		}
	}()

	log.Printf("Waiting for up to %d seconds for eip to become %s", timeout/time.Second, desiredState)
	select {
	case err := <-result:
		return err
	case <-time.After(timeout):
		err := fmt.Errorf("Timeout while waiting to for eip to become '%s'", desiredState)
		return err
	}
}
