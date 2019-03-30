package sys

func CurrentBatteryCapacity(batteries ...string) (uint64, error) {
	var capacity uint64
	for _, b := range batteries {
		cap, err := Uint64(SysPath + batteryPath + b + "/energy_now")
		if err != nil {
			return capacity, err
		}
		capacity += cap
	}
	return capacity, nil

}

func MaxBatteryCapacity(batteries ...string) (uint64, error) {
	var capacity uint64
	for _, b := range batteries {
		cap, err := Uint64(SysPath + batteryPath + b + "/energy_full")
		if err != nil {
			return capacity, err
		}
		capacity += cap
	}
	return capacity, nil
}
