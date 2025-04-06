package smartgo

import (
	"fmt"

	"github.com/anatol/smart.go"
)

/* Extract SATA drive SMART data */
func getSATASMART(diskpath string) (*SMARTData, error) {
	dev, err := smart.OpenSata(diskpath)
	if err != nil {
		return nil, err
	}
	defer dev.Close()

	c, err := dev.Identify()
	if err != nil {
		return nil, err
	}

	sm, err := dev.ReadSMARTData()
	if err != nil {
		return nil, err
	}

	_, capacity, _, _, _ := c.Capacity()
	sataAttrs := []*SATAAttrData{}
	for _, attr := range sm.Attrs {
		thisAttrID := attr.Id
		thisAttrName := attr.Name
		thisAttrType := attr.Type
		thisAttrRawVal := attr.ValueRaw
		val, low, high, _, err := attr.ParseAsTemperature()
		if err != nil {
			continue
		}

		fmt.Println("Temperature: ", val, low, high)
		thisAttrData := SATAAttrData{
			Id:      thisAttrID,
			Name:    thisAttrName,
			Type:    thisAttrType,
			RawVal:  thisAttrRawVal,
			Current: attr.Current,
			Worst:   attr.Worst,
		}

		sataAttrs = append(sataAttrs, &thisAttrData)
	}

	smartData := SMARTData{
		ModelNumber:  c.ModelNumber(),
		SerialNumber: c.SerialNumber(),
		Size:         capacity,

		SATAAttrs: sataAttrs,
	}

	return &smartData, nil
}

/* Extract NVMe drive SMART data */
func getNVMESMART(diskpath string) (*SMARTData, error) {
	dev, err := smart.OpenNVMe(diskpath)
	if err != nil {
		return nil, err
	}
	defer dev.Close()

	c, nss, err := dev.Identify()
	if err != nil {
		return nil, err
	}

	NameSpaceUtilizations := []uint64{}
	for _, ns := range nss {
		NameSpaceUtilizations = append(NameSpaceUtilizations, ns.Nuse*ns.LbaSize())
	}

	sm, _ := dev.ReadSMART()
	smartData := SMARTData{
		ModelNumber:           c.ModelNumber(),
		SerialNumber:          c.SerialNumber(),
		Size:                  c.Tnvmcap.Val[0],
		NameSpaceUtilizations: NameSpaceUtilizations,
		Temperature:           int(sm.Temperature),
		PowerOnHours:          sm.PowerOnHours.Val[0],
		PowerCycles:           sm.PowerCycles.Val[0],
		UnsafeShutdowns:       sm.UnsafeShutdowns.Val[0],
		MediaErrors:           sm.MediaErrors.Val[0],
	}

	return &smartData, nil
}
