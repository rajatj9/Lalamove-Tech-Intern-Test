package main

import (
	"bufio"
	"strings"
	"os"
	"sort"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"github.com/google/go-github/github"
	"github.com/coreos/go-semver/semver"

)
type descendingVersions []*semver.Version

func (x descendingVersions) Len() int { return len(x)}
func (x descendingVersions) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x descendingVersions) Less(i, j int) bool { return x[j].LessThan(*x[i])}

// LatestVersions returns a sorted slice with the highest version as its first element and the highest version of the smaller minor versions in a descending order
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {
	var versionSlice []*semver.Version
	// This is just an example structure of the code, if you implement this interface, the test cases in main_test.go are very easy to run
	for _, version := range releases {
		if !(version).LessThan(*minVersion) { 
			appendToSlice := true
			for _, retVersion := range versionSlice {
				if (*version).Major == (*retVersion).Major && (*version).Minor == (*retVersion).Minor {
					appendToSlice = false
					if (*version).Patch > (*retVersion).Patch { 
						*retVersion = *version
						break
					} else if (*version).Patch == (*retVersion).Patch && (*version).PreRelease > (*retVersion).PreRelease { //handling the case for Pre releases
						*retVersion = *version
						break
					}
				}
			}
			if appendToSlice {
				versionSlice = append(versionSlice, version)
			}
		}
	}
	sort.Sort(descendingVersions(versionSlice))
	return versionSlice
}

func CSVFileRead(fileLocation string) ([][]string, error) {
	inputFile, err := os.Open(fileLocation)
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()
	reader := csv.NewReader(bufio.NewReader(inputFile))
	entries, err := reader.ReadAll()
	return entries, nil
}

// Here we implement the basics of communicating with github through the library as well as printing the version
// You will need to implement LatestVersions function as well as make this application support the file format outlined in the README
// Please use the format defined by the fmt.Printf line at the bottom, as we will define a passing coding challenge as one that outputs
// the correct information, including this line
func main() {
	// Github
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{PerPage: 10}
	
	fileLocation := os.Args[1]
	csv, err := CSVFileRead(fileLocation)

	if err != nil {
		log.Fatal(err)
	}
	for i, entry := range csv { 
		if i == 0 {
			continue
		}
		repo := strings.Split(entry[0], "/")
		releases, _, err := client.Repositories.ListReleases(ctx, repo[0], repo[1], opt)
		if err != nil {
			log.Fatal(err)
		}

	minVersion := semver.New("1.8.0")
	allReleases := make([]*semver.Version, len(releases))
	for i, release := range releases {
		versionString := *release.TagName
		if versionString[0] == 'v' {
			versionString = versionString[1:]
		}
		allReleases[i] = semver.New(versionString)
	}
	versionSlice := LatestVersions(allReleases, minVersion)

	fmt.Printf("latest versions of %s/%s: %s", repo[0], repo[1], versionSlice)
	fmt.Println()
}
}