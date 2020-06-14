package tenantipmanager

import (
    "context"
    "fmt"
	"errors"
	"sync"
    //goipam "github.com/metal-stack/go-ipam"
	goipam "github.com/VeenaSL/go-ipam"
)

const (
	MT_IP_PREFIX = "10.0.0.0/8"
	//MT_IP_PREFIX = "10.10.10.0/24"
)

/*
type DeviceIntf interface {

	DeviceID() string
	TenantID() string
	DeviceIP() string
	AllocatedIP() string
}
*/

type Device struct {

	tntid string	//Tenant ID
	devid string	//Device ID
	extip string	//External IP
	intip *goipam.IP //Internal IP
}

func (d *Device)DeviceID() string {
	return d.devid
}

func (d *Device)TenantID() string {
	return d.tntid
}

func (d *Device)DeviceIP() string {
	return d.extip
}

func (d *Device)AllocatedIP () string {
	return d.intip.IP.String()
}

type Tenant struct {

	tntid string
	devMap map[string]*Device
}

type TenantIpManager struct {

	ipam goipam.Ipamer
	prefixVal string
	prefix *goipam.Prefix
	tntMap map[string]*Tenant		//Key is tntid
}

var timObj *TenantIpManager
var once sync.Once

func GetTenantIpManager(ctx context.Context, prefix string) *TenantIpManager {

	once.Do(func() {
		timObj = newTenantIpManager(ctx, prefix)
	})
	return timObj

}

func newTenantIpManager(ctx context.Context, prefix string) *TenantIpManager {

	ipam := goipam.New()
    ph, err := ipam.NewPrefix(prefix)
    if err != nil {
		fmt.Println(err)
        return nil
    }

	ret := &TenantIpManager {

		ipam: ipam,
		prefix: ph,
		tntMap:  make(map[string]*Tenant),
	}
	if prefix == "" {
		ret.prefixVal = MT_IP_PREFIX
	} else {
		ret.prefixVal = prefix
	}


	return ret
}

//AllocateIP - Allocate IP for tenant id and device id
func (tim *TenantIpManager)AllocateIP(ctx context.Context,
	tntid string, devid string, devip string) (*Device, error) {

	var tnt *Tenant
	var tntOk bool

	tnt, tntOk = tim.tntMap[tntid]
	if !tntOk {

		//New Tenant
		tnt = &Tenant {

			tntid: tntid,
			devMap: make(map[string]*Device),
		}

		tim.tntMap[tntid] = tnt
	}

	var dev *Device
	var devOk bool

	dev, devOk = tnt.devMap[devid]
	if !devOk {

		//New Device
		dev = &Device {

			tntid: tntid,
			devid: devid,
			extip: devip,
		}

		tnt.devMap[devid] = dev
	} else {

		return dev, nil
	}

    ip, err := tim.ipam.AcquireIP(tim.prefix.Cidr)
    if err != nil {
        return nil, err
    }

    dev.intip = ip

	return dev, nil
}

//GetAllocatedIP -
func (tim *TenantIpManager)GetAllocatedIP(ctx context.Context,
	tntid string, devid string) (*Device, error) {

	var tnt *Tenant
	var tntOk bool
	tnt, tntOk = tim.tntMap[tntid]
	if !tntOk {
		return nil, errors.New("tenant not found")
	}

	var dev *Device
	var devOk bool
	dev, devOk = tnt.devMap[devid]
	if !devOk {
		return nil, errors.New("device not found")
	}

	return dev, nil
}

//ReleaseIP -
func (tim *TenantIpManager)ReleaseIP(ctx context.Context,
	tntid string, devid string) (*Device, error) {

	var tnt *Tenant
	var tntOk bool
	tnt, tntOk = tim.tntMap[tntid]
	if !tntOk {
		return nil, errors.New("tenant not found")
	}

	var dev *Device
	var devOk bool
	dev, devOk = tnt.devMap[devid]
	if !devOk {
		return nil, errors.New("device not found")
	}

    _, err := tim.ipam.ReleaseIP(dev.intip)
    if err != nil {
        return nil, err
    }
	delete(tnt.devMap, devid)

	if len(tnt.devMap) == 0 {

		delete(tim.tntMap, tntid)
	}

	return dev, nil
}

//AllocatedIPs -
func (tim *TenantIpManager)AllocatedIPs(ctx context.Context) uint64 {

	var count uint64
	for _, tnt := range tim.tntMap {

		count += uint64(len(tnt.devMap))
	}

	return count
}

//AvailableIPs -
func (tim *TenantIpManager)AvailableIPs(ctx context.Context) uint64 {

	var count uint64
	for _, tnt := range tim.tntMap {

		count += uint64(len(tnt.devMap))
	}

	usg := tim.prefix.Usage()

	return (usg.AvailableIPs - count)
}

func PrintDevice (dev *Device) {

	fmt.Printf("%s %s %s %s\n", dev.TenantID(), dev.DeviceID(), dev.DeviceIP(),
		dev.AllocatedIP())
}

