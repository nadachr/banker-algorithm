package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	process []string
	allo    []int
	max     []int
	need    []int
	res     []int
)

func initialized() {
	process = make([]string, 10)
	allo = make([]int, 30)
	max = make([]int, 30)
	need = make([]int, 30)
	res = append(res, 10)
	res = append(res, 10)
	res = append(res, 10)
}

func showProcess() {
	fmt.Printf("\n-----------------------------------------------\n")
	fmt.Printf(" Process |Allocate|  Need |  Max  | Available ")
	fmt.Printf("\n         | a b c  | a b c | a b c | ")
	fmt.Printf("\n-----------------------------------------------\n")

	for n := range need {
		need[n] = calNeed(max[n], allo[n])
	}

	if process[0] == "" {
		fmt.Printf("    -    | - - -  | - - - | - - - | %d %d %d\n", res[0], res[1], res[2])
	} else {
		for i := range process {
			if process[i] == "" {
				continue
			}
			if i == 0 {
				fmt.Printf("    %s   | %d %d %d  | %d %d %d | %d %d %d | %d %d %d\n", process[i], allo[0], allo[1], allo[2], need[0], need[1], need[2], max[0], max[1], max[2], res[0], res[1], res[2])
			} else {
				fmt.Printf("    %s   | %d %d %d  | %d %d %d | %d %d %d |\n", process[i], allo[0+(3*i)], allo[1+(3*i)], allo[2+(3*i)], need[0+(3*i)], need[1+(3*i)], need[2+(3*i)], max[0+(3*i)], max[1+(3*i)], max[2+(3*i)])
			}
		}
	}
	fmt.Printf("\n")
	fmt.Printf("\nCommand > ")
}

func getCommand() string {
	reader := bufio.NewReader(os.Stdin)
	data, _ := reader.ReadString('\n')
	data = strings.Trim(data, "\n")
	return data
}

func command_create(p string, m1, m2, m3 int) {
	if len(process) > 0 {
		for i := range process {
			if process[i] == "" {
				process[i] = p
				max[0+(i*3)] = m1
				max[1+(i*3)] = m2
				max[2+(i*3)] = m3
				break
			}
		}
	}
}

func isSafe(p string, pt int, req int) bool {
	x := 0  // number of process that have created
	pc := 0 // address of process that request
	nsafe := 0 // not safe
	for n := range process {
		if process[n] != "" {
			x++
		}
		if process[n] == p {
			pc = n
		}
	}
	for j := range need {
		if j >= x*3 {
			break
		}
		fmt.Println("here1")
		fail := 0
		for i := range res {
			if i+(j*3) >= x*3 {
				continue
			}
			if i == pt {
				check := res[pt] - req
				if (i+(j*3)) == (pt+(pc*3)) {
					new_need := need[pt+(pc*3)] - req
					//fmt.Printf("RES[%d] = %d |||| (new)need[%d] = %d\n", pt, check, i+(j*3), new_need)
					if check < new_need {	// || (check == 0 && check == new_need)
						fail++
						fmt.Println("fail1 = ", fail)
					}
					continue
				}
				//fmt.Printf("RES[%d] = %d |||| need[%d] = %d\n", pt, check, i+(j*3), need[i+(j*3)])
				if check < need[i+(j*3)] {		//|| (check == 0 && check == need[i+(j*3)]) 
					fail++
					fmt.Println("fail2 = ", fail)
				}
			} else {
				// fmt.Printf("res[%d] = %d |||| need[%d] = %d\n", i, res[i], i+(j*3), need[i+(j*3)])
				if res[i] < need[i+(j*3)] {		//|| (res[i] == 0 && res[i] == need[i+(j*3)])
					fmt.Println("here3")
					fail++
					fmt.Println("fail3 = ", fail)
				}
			}
		}
		if fail > 0 {
			nsafe++
			fmt.Println("can't end the process", nsafe, "times")
			if nsafe == x {
				return false
			}
		}
	}
	return true
}

func command_request(p string, t string, req int) {
	pt := 0
	switch t {
	case "a":
		pt = 0
	case "b":
		pt = 1
	case "c":
		pt = 2
	default:
		fmt.Printf("\nSyntax error, please input weither a, b or c.")
	}

	if len(process) > 0 {
		for i := range process {
			if process[i] == p {
				if need[pt+(3*i)] >= req {
					if isSafe(p, pt, req) {
						if (res[pt] - req) >= 0 {
							fmt.Printf("\nSafe!\n")
							res[pt] = res[pt] - req
							allo[pt+(3*i)] = allo[pt+(3*i)] + req
							for n := range need {
								need[n] = calNeed(max[n], allo[n])
							}
							if need[0+(3*i)] == 0 && need[1+(3*i)] == 0 && need[2+(3*i)] == 0 {
								terminate(process[i])
								break
							}
							break
						}
					} else {
						fmt.Println("\nNot Safe")
					}
				} else {
					fmt.Printf("\nError, You request resource more than process need.\n")
				}
			}
		}
	} else {
		fmt.Println("Error, There's no process to requested.")
	}
}

func terminate(p string) {
	for i := range process {
		if process[i] == p {
			copy(process[i:], process[i+1:])
			for j := range res {
				res[j] = res[j] + allo[j+(3*i)]
				allo[j+(3*i)] = 0
				max[j+(3*i)] = 0
			}
		}
	}
}

func calNeed(m int, a int) int {
	need := m - a
	return need
}

func main() {
	initialized()
	for {
		showProcess()
		command := getCommand()
		commandx := strings.Split(command, " ")

		switch commandx[0] {
		case "create":
			if len(commandx) > 2 {
				mx := strings.Split(commandx[2], ",")
				mx1, _ := strconv.Atoi(mx[0])
				mx2, _ := strconv.Atoi(mx[1])
				mx3, _ := strconv.Atoi(mx[2])
				command_create(commandx[1], mx1, mx2, mx3)
			} else {
				fmt.Println("\nSyntax Error.")
			}
		case "req":
			sub := strings.Split(commandx[1], "-")
			if len(sub) > 2 {
				req, _ := strconv.Atoi(sub[2])
				command_request(sub[0], sub[1], req)
			} else {
				fmt.Println("\nSyntax Error!")
			}
		case "exit":
			return
		default:
			fmt.Println("\nSyntax Error!")
		}
	}
}