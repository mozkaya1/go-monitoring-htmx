package hardware

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/docker"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/sensors"
)

const megabyteDiv uint64 = 1024 * 1024
const gigabyteDiv uint64 = megabyteDiv * 1024

func GetDocker() (string, error) {
	dock, err := docker.GetDockerStat()
	if err != nil {
		return "", err
	}

	html := "<div class='system-data'><table class='table table-striped table-hover table-sm'><thead><tr><th>Name</th><th>Status</th><th>Running</th></tr></thead><tbody>"
	for i := range dock {
		html = html + "<tr><td>" + dock[i].Name + "</td><td>" + dock[i].Status + "</td><td>" + strconv.FormatBool(dock[i].Running) + "</td></tr>"
	}
	// fmt.Println(dock[0].Running, dock[0].Name, dock[0].Image, dock[0].Status, err)
	html = html + "</tbody></table></div>"
	return html, nil
}

func GetSensors() (string, error) {
	temp, _ := sensors.SensorsTemperatures()

	html := "<div class='system-data'><table class='table table-striped table-hover table-sm'><thead><tr><th>Device</th><th>Temp</th><th>Threshold</th></tr></thead><tbody>"
	for i := range temp {
		html = html + "<tr><td>" + temp[i].SensorKey + "</td><td>" + strconv.FormatFloat(temp[i].Temperature, 'f', 2, 64) + "</td><td>" + strconv.FormatFloat(temp[i].High, 'f', 2, 64) + "</td></tr>"
	}
	html = html + "</tbody></table></div>"
	return html, nil
}

func GetLoad() (string, error) {
	load, err := load.Avg()
	if err != nil {
		return "", err
	}

	html := "<div class='load'><table class='table table-striped table-hover table-sm'><tbody>"
	html = html + "<tr><td>Average Load: </td><td>" + strconv.FormatFloat(load.Load1, 'f', 2, 64) + "</td><td>" + strconv.FormatFloat(load.Load5, 'f', 2, 64) + "</td><td>" + strconv.FormatFloat(load.Load15, 'f', 2, 64) + "</td></tr></tbody></table></div>"
	return html, nil
}

func GetSystemSection() (string, error) {
	runTimeOS := runtime.GOOS
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return "", err
	}

	hostStat, err := host.Info()
	if err != nil {
		return "", err
	}

	html := "<div class='system-data'><table class='table table-striped table-hover table-sm'><tbody>"
	html = html + "<tr><td>Operating System:</td> <td><i class='fa fa-brands fa-linux'></i> " + runTimeOS + "</td></tr>"
	html = html + "<tr><td>Platform:</td><td> <i class='fa fa-brands fa-fed'></i> " + hostStat.Platform + "</td></tr>"
	html = html + "<tr><td>Hostname:</td><td>" + hostStat.Hostname + "</td></tr>"
	html = html + "<tr><td>Number of processes running:</td><td>" + strconv.FormatUint(hostStat.Procs, 10) + "</td></tr>"
	html = html + "<tr><td>Total memory:</td><td>" + strconv.FormatUint(vmStat.Total/megabyteDiv, 10) + " MB</td></tr>"
	html = html + "<tr><td>Free memory:</td><td>" + strconv.FormatUint(vmStat.Free/megabyteDiv, 10) + " MB</td></tr>"
	html = html + "<tr><td>Percentage used memory:</td><td>" + strconv.FormatFloat(vmStat.UsedPercent, 'f', 2, 64) + "%</td></tr></tbody></table>"

	html = html + "</div>"

	return html, nil
}

func GetDiskSection() (string, error) {
	diskStat, err := disk.Usage("/")
	if err != nil {
		return "", err
	}

	html := "<div class='disk-data'><table class='table table-striped table-hover table-sm'><tbody>"
	html = html + "<tr><td>Total disk space:</td><td>" + strconv.FormatUint(diskStat.Total/gigabyteDiv, 10) + " GB</td></tr>"
	html = html + "<tr><td>Used disk space:</td><td>" + strconv.FormatUint(diskStat.Used/gigabyteDiv, 10) + " GB</td></tr>"
	html = html + "<tr><td>Free disk space:</td><td>" + strconv.FormatUint(diskStat.Free/gigabyteDiv, 10) + " GB</td></tr>"
	html = html + "<tr><td>Percentage disk space usage:</td><td>" + strconv.FormatFloat(diskStat.UsedPercent, 'f', 2, 64) + "%</td></tr></tbody></table></div>"
	return html, nil
}

func GetCpuSection() (string, error) {
	cpuStat, err := cpu.Info()

	if err != nil {
		fmt.Println("Error getting CPU info", err)

	}
	percentage, err := cpu.Percent(0, true)
	if err != nil {
		return "", err
	}

	html := "<div class='cpu-data'><table class='table table-striped table-hover table-sm'><tbody>"

	if len(cpuStat) != 0 {
		html = html + "<tr><td>Model Name:</td><td>" + cpuStat[0].ModelName + "</td></tr>"
		html = html + "<tr><td>Family:</td><td>" + cpuStat[0].Family + "</td></tr>"
		html = html + "<tr><td>Speed:</td><td>" + strconv.FormatFloat(cpuStat[0].Mhz, 'f', 2, 64) + " MHz</td></tr>"
	}

	firstCpus := percentage[:len(percentage)/2]
	secondCpus := percentage[len(percentage)/2:]
	html = html + "<tr><td>Cores: </td><td><div class='row mb-4'><div class='col-md-6'><table class='table table-sm'><tbody>"
	for idx, cpupercent := range firstCpus {
		html = html + "<tr><td>CPU [" + strconv.Itoa(idx) + "]: " + strconv.FormatFloat(cpupercent, 'f', 2, 64) + "%</td></tr>"
	}
	html = html + "</tbody></table></div><div class='col-md-6'><table class='table table-sm'><tbody>"
	for idx, cpupercent := range secondCpus {
		html = html + "<tr><td>CPU [" + strconv.Itoa(idx+8) + "]: " + strconv.FormatFloat(cpupercent, 'f', 2, 64) + "%</td></tr>"
	}
	html = html + "</tbody></table></div></div></td></tr></tbody></table></div>"
	return html, nil

}
