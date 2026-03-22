package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// BatteryInfo holds battery status information
type BatteryInfo struct {
	ChargePercent     int     `json:"charge_percent"`
	HealthPercent     int     `json:"health_percent"`
	CycleCount        int     `json:"cycle_count"`
	MaxCycles         int     `json:"max_cycles"`
	TemperatureCelsius float64 `json:"temperature_celsius"`
	TimeRemaining     string  `json:"time_remaining"`
	Status            string  `json:"status"`
	PowerSource       string  `json:"power_source"`
	EstimatedLifeYears float64 `json:"estimated_life_years"`
}

var (
	version = "1.0.0"
	theme   = "default"
	noColor = false
)

// Color codes for terminal output
var colors = map[string]map[string]string{
	"default": {
		"red":     "\033[31m",
		"green":   "\033[32m",
		"yellow":  "\033[33m",
		"blue":    "\033[34m",
		"magenta": "\033[35m",
		"cyan":    "\033[36m",
		"white":   "\033[37m",
		"reset":   "\033[0m",
		"bold":    "\033[1m",
	},
	"neon": {
		"red":     "\033[38;5;197m",
		"green":   "\033[38;5;46m",
		"yellow":  "\033[38;5;226m",
		"blue":    "\033[38;5;51m",
		"magenta": "\033[38;5;201m",
		"cyan":    "\033[38;5;87m",
		"white":   "\033[38;5;15m",
		"reset":   "\033[0m",
		"bold":    "\033[1m",
	},
	"dark": {
		"red":     "",
		"green":   "",
		"yellow":  "",
		"blue":    "",
		"magenta": "",
		"cyan":    "",
		"white":   "",
		"reset":   "",
		"bold":    "",
	},
	"minimal": {
		"red":     "",
		"green":   "",
		"yellow":  "",
		"blue":    "",
		"magenta": "",
		"cyan":    "",
		"white":   "",
		"reset":   "",
		"bold":    "",
	},
}

func color(name string) string {
	if noColor {
		return ""
	}
	if c, ok := colors[theme][name]; ok {
		return c
	}
	return ""
}

func getBatteryInfoMacOS() (*BatteryInfo, error) {
	info := &BatteryInfo{
		MaxCycles: 1000, // Default for modern MacBooks
	}

	// Get battery info using pmset
	out, err := exec.Command("pmset", "-g", "batt").Output()
	if err != nil {
		return nil, err
	}
	output := string(out)

	// Parse charge percentage
	if strings.Contains(output, "%") {
		parts := strings.Split(output, "%")
		if len(parts) > 0 {
			percentStr := parts[0]
			lastSpace := strings.LastIndex(percentStr, " ")
			if lastSpace != -1 {
				percentStr = percentStr[lastSpace+1:]
				if p, err := strconv.Atoi(percentStr); err == nil {
					info.ChargePercent = p
				}
			}
		}
	}

	// Determine status
	if strings.Contains(output, "charging") {
		info.Status = "charging"
		info.PowerSource = "USB-C"
	} else if strings.Contains(output, "discharging") {
		info.Status = "discharging"
		info.PowerSource = "Battery"
	} else if strings.Contains(output, "charged") || info.ChargePercent == 100 {
		info.Status = "charged"
		info.PowerSource = "AC Power"
	}

	// Get cycle count and health using system_profiler
	out, err = exec.Command("system_profiler", "SPPowerDataType").Output()
	if err == nil {
		profilerOutput := string(out)
		lines := strings.Split(profilerOutput, "\n")
		
		for i, line := range lines {
			line = strings.TrimSpace(line)
			
			if strings.Contains(line, "Cycle Count") {
				parts := strings.Split(line, ":")
				if len(parts) == 2 {
					countStr := strings.TrimSpace(parts[1])
					if c, err := strconv.Atoi(countStr); err == nil {
						info.CycleCount = c
					}
				}
			}
			
			if strings.Contains(line, "Maximum Capacity") {
				parts := strings.Split(line, ":")
				if len(parts) == 2 {
					capStr := strings.TrimSpace(parts[1])
					capStr = strings.TrimSuffix(capStr, "%")
					if h, err := strconv.Atoi(capStr); err == nil {
						info.HealthPercent = h
					}
				}
			}
			
			if strings.Contains(line, "Temperature") {
				parts := strings.Split(line, ":")
				if len(parts) == 2 {
					tempStr := strings.TrimSpace(parts[1])
					// Parse temperature (usually in Celsius)
					if strings.Contains(tempStr, "°C") {
						tempStr = strings.TrimSuffix(tempStr, "°C")
						tempStr = strings.TrimSpace(tempStr)
						if t, err := strconv.ParseFloat(tempStr, 64); err == nil {
							info.TemperatureCelsius = t
						}
					}
				}
			}
			
			// Estimate remaining life based on cycle count
			if info.CycleCount > 0 && i == len(lines)-1 {
				// Rough estimate: cycles used / max cycles
				usedRatio := float64(info.CycleCount) / float64(info.MaxCycles)
				info.EstimatedLifeYears = (1.0 - usedRatio) * 5.0 // Assume ~5 years max
			}
		}
	}

	// Get time estimate
	out, err = exec.Command("pmset", "-g", "ps").Output()
	if err == nil {
		psOutput := string(out)
		if strings.Contains(psOutput, "remaining") {
			parts := strings.Split(psOutput, "remaining")
			if len(parts) > 1 {
				timeStr := strings.TrimSpace(parts[0])
				lastPart := timeStr[strings.LastIndex(timeStr, " ")+1:]
				info.TimeRemaining = strings.TrimSpace(lastPart + " remaining")
			}
		}
	}

	return info, nil
}

func getBatteryInfoLinux() (*BatteryInfo, error) {
	info := &BatteryInfo{
		MaxCycles: 1000,
	}

	// Try to read from /sys/class/power_supply
	batteryPath := "/sys/class/power_supply/BAT0"
	if _, err := os.Stat(batteryPath); os.IsNotExist(err) {
		batteryPath = "/sys/class/power_supply/BAT1"
	}

	// Read charge percentage
	if data, err := os.ReadFile(batteryPath + "/capacity"); err == nil {
		if p, err := strconv.Atoi(strings.TrimSpace(string(data))); err == nil {
			info.ChargePercent = p
		}
	}

	// Read status
	if data, err := os.ReadFile(batteryPath + "/status"); err == nil {
		status := strings.TrimSpace(string(data))
		info.Status = strings.ToLower(status)
		if status == "Charging" {
			info.PowerSource = "AC Adapter"
		} else if status == "Discharging" {
			info.PowerSource = "Battery"
		} else {
			info.PowerSource = "AC Power"
		}
	}

	// Read cycle count
	if data, err := os.ReadFile(batteryPath + "/cycle_count"); err == nil {
		if c, err := strconv.Atoi(strings.TrimSpace(string(data))); err == nil {
			info.CycleCount = c
		}
	}

	// Read health (capacity)
	if data, err := os.ReadFile(batteryPath + "/health"); err == nil {
		health := strings.TrimSpace(string(data))
		// Parse "Good", "Fair", etc.
		if health == "Good" {
			info.HealthPercent = 90
		} else if health == "Fair" {
			info.HealthPercent = 70
		}
	}

	// Calculate estimated life
	if info.CycleCount > 0 {
		usedRatio := float64(info.CycleCount) / float64(info.MaxCycles)
		info.EstimatedLifeYears = (1.0 - usedRatio) * 5.0
	}

	return info, nil
}

func getBatteryInfo() (*BatteryInfo, error) {
	switch runtime.GOOS {
	case "darwin":
		return getBatteryInfoMacOS()
	case "linux":
		return getBatteryInfoLinux()
	default:
		return nil, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func progressBar(percent int, width int) string {
	filled := percent * width / 100
	empty := width - filled
	
	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	return bar
}

func getColorByPercent(percent int) string {
	if percent >= 80 {
		return "green"
	} else if percent >= 50 {
		return "yellow"
	}
	return "red"
}

func printBox(info *BatteryInfo) {
	width := 57
	
	fmt.Printf("┌%s┐\n", strings.Repeat("─", width-2))
	fmt.Printf("│  🔋 Battery Status%s│\n", strings.Repeat(" ", width-21))
	fmt.Printf("├%s┤\n", strings.Repeat("─", width-2))
	
	// Charge
	chargeColor := getColorByPercent(info.ChargePercent)
	fmt.Printf("│  Charge      %s%s%s  %s%d%%%s        │\n",
		color(chargeColor), progressBar(info.ChargePercent, 24), color("reset"),
		color(chargeColor), info.ChargePercent, color("reset"))
	
	// Health
	healthColor := getColorByPercent(info.HealthPercent)
	fmt.Printf("│  Health      %s%s%s  %s%d%%%s        │\n",
		color(healthColor), progressBar(info.HealthPercent, 24), color("reset"),
		color(healthColor), info.HealthPercent, color("reset"))
	
	// Cycles
	cyclePercent := info.CycleCount * 100 / info.MaxCycles
	cycleColor := "green"
	if cyclePercent > 70 {
		cycleColor = "yellow"
	}
	if cyclePercent > 90 {
		cycleColor = "red"
	}
	fmt.Printf("│  Cycles      %s%s%s  %s%d / %d%s │\n",
		color(cycleColor), progressBar(cyclePercent, 24), color("reset"),
		color(cycleColor), info.CycleCount, info.MaxCycles, color("reset"))
	
	// Temperature
	tempColor := "green"
	if info.TemperatureCelsius > 40 {
		tempColor = "yellow"
	}
	if info.TemperatureCelsius > 45 {
		tempColor = "red"
	}
	fmt.Printf("│  Temp        %s%.1f°C%s  %s    │\n",
		color(tempColor), info.TemperatureCelsius, color("reset"),
		strings.Repeat(" ", 20))
	
	fmt.Printf("├%s┤\n", strings.Repeat("─", width-2))
	
	// Time estimate
	if info.TimeRemaining != "" {
		fmt.Printf("│  ⏱️  %s%s│\n", info.TimeRemaining, strings.Repeat(" ", width-7-len(info.TimeRemaining)))
	}
	
	// Status
	statusIcon := "🔋"
	if info.Status == "charging" {
		statusIcon = "⚡"
	}
	statusText := fmt.Sprintf("%s Status: %s", statusIcon, info.Status)
	if info.PowerSource != "" {
		statusText += fmt.Sprintf(" (%s)", info.PowerSource)
	}
	fmt.Printf("│  %s%s│\n", statusText, strings.Repeat(" ", width-3-len(statusText)))
	
	// Estimated life
	if info.EstimatedLifeYears > 0 {
		lifeText := fmt.Sprintf("📅 Estimated Life: ~%.1f years remaining", info.EstimatedLifeYears)
		fmt.Printf("│  %s%s│\n", lifeText, strings.Repeat(" ", width-3-len(lifeText)))
	}
	
	fmt.Printf("└%s┘\n", strings.Repeat("─", width-2))
}

func printCompact(info *BatteryInfo) {
	now := time.Now().Format("15:04:05")
	
	statusIcon := "🔋"
	if info.Status == "charging" {
		statusIcon = "⚡"
	}
	
	tempStr := ""
	if info.TemperatureCelsius > 0 {
		tempStr = fmt.Sprintf(" | %.0f°C", info.TemperatureCelsius)
	}
	
	fmt.Printf("[%s] %s %d%% | %s%s\n",
		now, statusIcon, info.ChargePercent, info.Status, tempStr)
}

func main() {
	watch := flag.Bool("watch", false, "Watch mode, update every 5 seconds")
	w := flag.Bool("w", false, "Watch mode (shorthand)")
	interval := flag.Int("interval", 5, "Watch interval in seconds")
	i := flag.Int("i", 5, "Watch interval in seconds (shorthand)")
	jsonOutput := flag.Bool("json", false, "Output as JSON")
	j := flag.Bool("j", false, "Output as JSON (shorthand)")
	detail := flag.Bool("detail", false, "Show detailed battery information")
	d := flag.Bool("d", false, "Show detailed info (shorthand)")
	themeFlag := flag.String("theme", "default", "Color theme (default, dark, neon, minimal)")
	t := flag.String("t", "default", "Color theme (shorthand)")
	noColorFlag := flag.Bool("no-color", false, "Disable colored output")
	n := flag.Bool("n", false, "Disable colored output (shorthand)")
	showVersion := flag.Bool("version", false, "Show version")
	v := flag.Bool("v", false, "Show version (shorthand)")
	
	flag.Parse()
	
	if *showVersion || *v {
		fmt.Printf("cli-battery v%s\n", version)
		os.Exit(0)
	}
	
	theme = *themeFlag
	if *t != "default" {
		theme = *t
	}
	noColor = *noColorFlag || *n
	
	isWatch := *watch || *w
	isJSON := *jsonOutput || *j
	isDetail := *detail || *d
	watchInterval := *interval
	if *i != 5 {
		watchInterval = *i
	}
	
	for {
		info, err := getBatteryInfo()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		
		if isJSON {
			jsonData, _ := json.MarshalIndent(info, "", "  ")
			fmt.Println(string(jsonData))
		} else if isDetail || !isWatch {
			printBox(info)
		} else {
			printCompact(info)
		}
		
		if !isWatch {
			break
		}
		
		time.Sleep(time.Duration(watchInterval) * time.Second)
	}
}
