package main


import (
    "context"
    "fmt"
    "time"
	"tenantipmanager"
)

func main() {

	ctx := context.Background()
	tim := tenantipmanager.GetTenantIpManager(ctx, tenantipmanager.MT_IP_PREFIX)

	devip := "200.1.1.2"

	startTime := time.Now()
	for j := 1; j <= 1; j++ {

		tntid := fmt.Sprintf("%d", j)

		for i := 1; i <= 3; i++ {
			devid := fmt.Sprintf("%d", i)
			if _, err := tim.AllocateIP(ctx, tntid, devid, devip); err != nil {
				fmt.Println(err)
			}
		}

		//fmt.Printf("Available IPs: %d\n", tim.AvailableIPs(ctx))
		//fmt.Printf("Allocated IPs: %d\n", tim.AllocatedIPs(ctx))
	}
	endTime := time.Now()
	fmt.Println(endTime.Sub(startTime))

	//startTime = time.Now()
	//if _, err := tim.AllocateIP(ctx, "1", "256", ""); err != nil {

	//	fmt.Println(err)
	//}
	//endTime = time.Now()
	//fmt.Println(endTime.Sub(startTime))

	//for i := 1; i < 10; i++ {
	//	devid := fmt.Sprintf("%d", i)
	//	if _, err := tim.ReleaseIP(ctx, tntid, devid); err != nil {
	//		fmt.Println(err)
	//	}
	//}

	//for i := 1; i < 10; i++ {
	//	devid := fmt.Sprintf("%d", i)

	//	if dev, err := tim.GetAllocatedIP(ctx, tntid, devid); err != nil {
	//		fmt.Println(err)
	//	} else {
	//		tenantipmanager.PrintDevice(dev)
	//	}
	//}
}
