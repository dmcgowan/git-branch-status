package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
)

var sortFlag string
var showAll bool

func init() {
	flag.Usage = func() {
		name := os.Args[0]
		if name == "git-branch-status" {
			name = "git branch-status"
		}
		fmt.Fprintf(os.Stderr, "usage: %s [<options>] [<branch>]\n\noptions:\n", name)
		flag.PrintDefaults()
	}
	flag.StringVar(&sortFlag, "sort", "name", `How to sort branches (name | left | right)
		"name": Sort by name of branch (left side branch)
		"left": Sort by number of commits in left branch and not right
		"right": Sort by number of commits in right branch and not left`)
	flag.BoolVar(&showAll, "a", false, "Show all branches including branches with no commit differences")
}

func main() {
	flag.Parse()
	if flag.NArg() > 1 {
		flag.Usage()
		log.Fatalf("Too many arguments")
	}

	upstream := "%(upstream:short)"
	if flag.NArg() == 1 {
		upstream = flag.Arg(0)
	}

	sortT := SortName
	if sortFlag != "" {
		switch sortFlag {
		case "left":
			sortT = SortLeft
		case "right":
			sortT = SortRight
		case "name":
		default:
			flag.Usage()
		}
	}

	cmd := exec.Command("git", "for-each-ref", fmt.Sprintf("--format=%%(refname:short) %s", upstream), "refs/heads")
	cmd.Stderr = os.Stderr
	foreachOut, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Error getting pipe: %s", err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("Error starting process: %s", err)
	}

	// Prepare output records
	var statuses []BranchStatus
	var maxLeftLen int

	scanner := bufio.NewScanner(foreachOut)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(scanner.Text())
		if len(fields) > 2 {
			log.Fatalf("Unexpected output line: %q", line)
		}
		if len(fields) == 1 {
			// No branch to compare, skip
			continue
		}

		status, err := GetBranchStatus(fields[0], fields[1])
		if err != nil {
			log.Fatalf("Error getting branch status: %s", err)
		}
		if len(status.Left) > maxLeftLen {
			maxLeftLen = len(status.Left)
		}
		if !showAll && status.LeftCount == 0 && status.RightCount == 0 {
			continue
		}

		statuses = append(statuses, status)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error starting process: %s", err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatalf("Error running process: %s", err)
	}

	l := NewBranchStatusList(statuses, sortT)
	sort.Sort(l)
	fmtStr := fmt.Sprintf("%%-%ds %%5d|%%-5d %%s\n", maxLeftLen)
	for _, status := range l.Statuses() {
		fmt.Printf(fmtStr, status.Left, status.LeftCount, status.RightCount, status.Right)
	}
}

type BranchStatus struct {
	Left       string
	Right      string
	LeftCount  int
	RightCount int
}

func GetBranchStatus(left, right string) (BranchStatus, error) {
	args := []string{
		"rev-list",
		"--left-right",
		left + "..." + right,
		"--",
	}
	cmd := exec.Command("git", args...)
	cmd.Stderr = os.Stderr

	p, err := cmd.StdoutPipe()
	if err != nil {
		return BranchStatus{}, err
	}
	if err := cmd.Start(); err != nil {
		return BranchStatus{}, err
	}

	status := BranchStatus{
		Left:  left,
		Right: right,
	}
	scanner := bufio.NewScanner(p)
	for scanner.Scan() {
		b := scanner.Bytes()
		switch b[0] {
		case '<':
			status.LeftCount++
		case '>':
			status.RightCount++
		default:
			return BranchStatus{}, fmt.Errorf("Unexpected revision line: %q", scanner.Text())
		}
	}
	if err := scanner.Err(); err != nil {
		return BranchStatus{}, err
	}

	if err := cmd.Wait(); err != nil {
		return BranchStatus{}, err
	}

	return status, nil
}

type SortType int

const (
	SortName SortType = iota
	SortLeft
	SortRight
)

func LeftLess(s1, s2 BranchStatus) bool {
	if s1.LeftCount == s2.LeftCount {
		if s1.RightCount == s2.RightCount {
			return s1.Left < s2.Left
		}
		return s1.RightCount < s2.RightCount
	}
	return s1.LeftCount < s2.LeftCount
}

func RightLess(s1, s2 BranchStatus) bool {
	if s1.RightCount == s2.RightCount {
		if s1.LeftCount == s2.LeftCount {
			return s1.Left < s2.Left
		}
		return s1.LeftCount < s2.LeftCount
	}
	return s1.RightCount < s2.RightCount
}

func NameLess(s1, s2 BranchStatus) bool {
	if s1.Left == s2.Left {
		if s1.LeftCount == s2.LeftCount {
			return s1.RightCount < s2.RightCount
		}
		return s1.LeftCount < s2.LeftCount
	}
	return s1.Left < s2.Left

}

type BranchStatusList struct {
	statuses []BranchStatus
	sort     SortType
}

func NewBranchStatusList(statuses []BranchStatus, sort SortType) *BranchStatusList {
	return &BranchStatusList{
		statuses: statuses,
		sort:     sort,
	}
}

func (l *BranchStatusList) Statuses() []BranchStatus {
	return l.statuses
}

func (l *BranchStatusList) Len() int {
	return len(l.statuses)
}

func (l *BranchStatusList) Less(i, j int) bool {
	switch l.sort {
	case SortLeft:
		return LeftLess(l.statuses[i], l.statuses[j])
	case SortRight:
		return RightLess(l.statuses[i], l.statuses[j])
	default:
		return NameLess(l.statuses[i], l.statuses[j])
	}
}

func (l *BranchStatusList) Swap(i, j int) {
	l.statuses[i], l.statuses[j] = l.statuses[j], l.statuses[i]
}
