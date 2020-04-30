package main

import (
	"fmt"
	logpkg "log"
	"os"
	"strconv"

	"github.com/andrewn3wman7/statuscake"
	"strings"
)

var log *logpkg.Logger

type command func(*statuscake.Client, ...string) error

var commands map[string]command

func init() {
	log = logpkg.New(os.Stderr, "", 0)
	commands = map[string]command{
		"list":   cmdList,
		"listpagespeed": cmdListPg,
		"listssl": cmdListSsl,
		"detail": cmdDetail,
		"detailpg": cmdDetailPg,
		"delete": cmdDelete,
		"deletepg": cmdDeletePg,
		"create": cmdCreate,
		"createpg": cmdCreatePg,
		"update": cmdUpdate,
		"updatepg": cmdUpdatePg,
	}
}

func colouredStatus(s string) string {
	switch s {
	case "Up":
		return fmt.Sprintf("\033[0;32m%s\033[0m", s)
	case "Down":
		return fmt.Sprintf("\033[0;31m%s\033[0m", s)
	default:
		return s
	}
}

func getEnv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		log.Fatalf("`%s` env variable is required", name)
	}

	return v
}

func cmdList(c *statuscake.Client, args ...string) error {
	tt := c.Tests()
	tests, err := tt.All()
	if err != nil {
		return err
	}

	for _, t := range tests {
		var paused string
		if t.Paused {
			paused = "yes"
		} else {
			paused = "no"
		}

		fmt.Printf("* %d: %s\n", t.TestID, colouredStatus(t.Status))
		fmt.Printf("  WebsiteName: %s\n", t.WebsiteName)
		fmt.Printf("  TestType: %s\n", t.TestType)
		fmt.Printf("  Paused: %s\n", paused)
		fmt.Printf("  ContactGroup: %s\n", fmt.Sprint(t.ContactGroup))
		fmt.Printf("  Uptime: %f\n", t.Uptime)
	}

	return nil
}

func cmdListSsl(c *statuscake.Client, args ...string) error {
	tt := c.Ssls()
	tests, err := tt.All()
	if err != nil {
		return err
	}

	for _, t := range tests {
		fmt.Printf("* %d: %s\n", t.ID)
		fmt.Printf("  Domain: %s\n", t.Domain)
	}

	return nil
}

func cmdListPg(c *statuscake.Client, args ...string) error {
	tt := c.PageSpeeds()
	tests, err := tt.All()
	if err != nil {
		return err
	}

	for _, t := range tests.Data {
		fmt.Printf("* %d \n", t.ID)
		fmt.Printf("  URL: %s\n", t.URL)
		fmt.Printf("  Name: %d \n", t.Title)
		fmt.Println(t.LatestStats.Requests)
	}
	return nil
}

func cmdDetailPg(c *statuscake.Client, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("command `detail` requires a single argument `TestID`")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	tt := c.PageSpeeds()
	t, err := tt.Detail(id)
	if err != nil {
		return err
	}

	fmt.Printf("* %d \n", t.ID)
	fmt.Printf("* %d \n", t.Website_url)
	fmt.Printf("* %d \n", t.Checkrate)
	fmt.Printf("* %d \n", t.AlertSmaller)
	fmt.Printf("* %d \n", t.Location_iso)
	return nil
}

func cmdDetail(c *statuscake.Client, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("command `detail` requires a single argument `TestID`")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	tt := c.Tests()
	t, err := tt.Detail(id)
	if err != nil {
		return err
	}

	var paused string
	if t.Paused {
		paused = "yes"
	} else {
		paused = "no"
	}

	fmt.Printf("* %d: %s\n", t.TestID, colouredStatus(t.Status))
	fmt.Printf("  WebsiteName: %s\n", t.WebsiteName)
	fmt.Printf("  WebsiteURL: %s\n", t.WebsiteURL)
	fmt.Printf("  PingURL: %s\n", t.PingURL)
	fmt.Printf("  TestType: %s\n", t.TestType)
	fmt.Printf("  Paused: %s\n", paused)
	fmt.Printf("  ContactGroup: %s\n", fmt.Sprint(t.ContactGroup))
	fmt.Printf("  Uptime: %f\n", t.Uptime)
	fmt.Printf("  NodeLocations: %s\n", fmt.Sprint(t.NodeLocations))

	return nil
}

func cmdDelete(c *statuscake.Client, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("command `delete` requires a single argument `TestID`")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	return c.Tests().Delete(id)
}

func cmdDeletePg(c *statuscake.Client, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("command `delete` requires a single argument `ID`")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	return c.PageSpeeds().Delete(id)
}

func askString(name string) string {
	var v string

	fmt.Printf("%s: ", name)
	_, err := fmt.Scanln(&v)
	if err != nil {
		log.Fatal(err)
	}

	return v
}

func askInt(name string) int {
	v := askString(name)
	i, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("Invalid number `%s`", v)
	}

	return i
}

func cmdCreate(c *statuscake.Client, args ...string) error {
	websiteName := askString("WebsiteName")
	websiteURL := askString("WebsiteURL")
	testType := askString("TestType")
	checkRate := askInt("CheckRate")
	contactGroupString := askString("ContactGroup (comma separated list)")
	contactGroup := strings.Split(contactGroupString, ",")
	nodeLocationsString := askString("NodeLocations (comma separated list)")
	nodeLocations := strings.Split(nodeLocationsString, ",")

	t := &statuscake.Test{
		WebsiteName:   websiteName,
		WebsiteURL:    websiteURL,
		TestType:      testType,
		CheckRate:     checkRate,
		NodeLocations: nodeLocations,
		ContactGroup:  contactGroup,
	}

	t2, err := c.Tests().Update(t)
	if err != nil {
		return err
	}

	fmt.Printf("CREATED: \n%+v\n", t2)

	return nil
}


func cmdCreatePg(c *statuscake.Client, args ...string) error {
	name := askString("name")
	websiteURL := askString("WebsiteURL")
	checkRate := askInt("CheckRate")
	locationIso :=askString("LocationIso")

	t := &statuscake.PageSpeed{
		Name:   name,
		Website_url:    websiteURL,
		Checkrate:      checkRate,
		Location_iso:	locationIso,
	}

	t2, err := c.PageSpeeds().Create(t)
	if err != nil {
		return err
	}

	fmt.Printf("CREATED: \n%+v\n", t2.ID)

	return nil
}

func cmdUpdatePg(c *statuscake.Client, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("command `update` requires a single argument `TestID`")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	tt := c.PageSpeeds()
	t, err := tt.Detail(id)
	if err != nil {
		return err
	}

	t.ID = id
	t.Website_url = askString(fmt.Sprintf("WebsiteName [%s]", t.Website_url))

	t2, err := c.PageSpeeds().Update(t)
	if err != nil {
		return err
	}

	fmt.Printf("UPDATED: \n%+v\n", t2)

	return nil
}

func cmdUpdate(c *statuscake.Client, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("command `update` requires a single argument `TestID`")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	tt := c.Tests()
	t, err := tt.Detail(id)
	if err != nil {
		return err
	}

	t.TestID = id
	t.WebsiteName = askString(fmt.Sprintf("WebsiteName [%s]", t.WebsiteName))
	t.WebsiteURL = askString(fmt.Sprintf("WebsiteURL [%s]", t.WebsiteURL))
	t.TestType = askString(fmt.Sprintf("TestType [%s]", t.TestType))
	t.CheckRate = askInt(fmt.Sprintf("CheckRate [%d]", t.CheckRate))
	contactGroupString := askString("ContactGroup (comma separated list)")
	t.ContactGroup = strings.Split(contactGroupString, ",")
	nodeLocationsString := askString("NodeLocations (comma separated list)")
	t.NodeLocations = strings.Split(nodeLocationsString, ",")

	t2, err := c.Tests().Update(t)
	if err != nil {
		return err
	}

	fmt.Printf("UPDATED: \n%+v\n", t2)

	return nil
}

func usage() {
	fmt.Printf("Usage:\n")
	fmt.Printf("  %s COMMAND\n", os.Args[0])
	fmt.Printf("Available commands:\n")
	for k := range commands {
		fmt.Printf("  %+v\n", k)
	}
}

func main() {
	username := getEnv("STATUSCAKE_USERNAME")
	apikey := getEnv("STATUSCAKE_APIKEY")

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	var err error

	c, err := statuscake.New(statuscake.Auth{Username: username, Apikey: apikey})
	if err != nil {
		log.Fatal(err)
	}

	if cmd, ok := commands[os.Args[1]]; ok {
		err = cmd(c, os.Args[2:]...)
	} else {
		err = fmt.Errorf("Unknown command `%s`", os.Args[1])
	}

	if err != nil {
		log.Fatalf("Error running command `%s`: %s", os.Args[1], err.Error())
	}
}
