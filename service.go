// services.go - fig commend support
package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"unicode"
)

type service struct {
	name        string
	namespace   string
	Build       string            `yaml:"build"`      // not supported as a command
	Command     string            `yaml:"command"`    // not supported as a command
	Image       string            `yaml:"image"`
	Ports       []string          `yaml:"ports"`
	Links       []string          `yaml:"links"`
	Environment map[string]string `yaml:"environment"`
	Volumes     []string          `yaml:"volumes"`
	Net         string            `yaml:"net"`
	Entrypoint  string            `yaml:"entrypoint"`
	Hostname    string            `yaml:"hostname"`
	User        string            `yaml:"user"`
	MemLimit    string            `yaml:"mem_limit"`
	Privileged  string            `yaml:"privileged"`
	Workdir     string            `yaml:"working_dir"`
	RunFlags    string            `yaml:"run_flags"`

	containerRe    *regexp.Regexp
	containers     []*container
	linkedServices []*link
}

type link struct {
	alias   string
	service *service
}

// init fields not found in the yaml file
func (s *service) init(name string, serviceMap map[string]*service) {
	s.name = name
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	base := path.Base(dir)
	// escape namespace
	for _, r := range base {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			s.namespace += string(r)
		}
	}
	s.containerRe = regexp.MustCompile(fmt.Sprintf("%s_%s_\\d+", s.namespace, s.name))
	s.containers = []*container{}

	s.linkedServices = []*link{}
	for _, l := range s.Links {
		fields := strings.Split(l, ":")
		link := &link{}
		link.alias = fields[len(fields)-1]
		linkedService, ok := serviceMap[fields[0]]
		if ok {
			link.service = linkedService
			s.linkedServices = append(s.linkedServices, link)
		}
	}
}

func (s *service) run(logsCh chan<- string, daemon, verbose bool) error {
	logfmt := "recreating %s...\n"
	if len(s.containers) == 0 {
		c := newContainer(s, 1)
		s.containers = append(s.containers, c)
		logfmt = "creating %s...\n"
	}
	for _, c := range s.containers {
		fmt.Printf(logfmt, c.name)
		err := c.run(logsCh, daemon, verbose)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) start(verbose bool) error {
	for _, c := range s.containers {
		fmt.Printf("starting %s...\n", c.name)
		if err := c.start(verbose); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) stop(verbose bool) error {
	for _, c := range s.containers {
		fmt.Printf("stopping %s...\n", c.name)
		if err := c.stop(verbose); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) kill(verbose bool) error {
	for _, c := range s.containers {
		fmt.Printf("killiing %s...\n", c.name)
		if err := c.kill(verbose); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) logs(ch chan<- string, timestamps, verbose bool) (int, error) {
	count := 0
	for _, c := range s.containers {
		err := c.logs(ch, timestamps, verbose)
		if err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

func (s *service) rmf(verbose bool) {
	for _, c := range s.containers {
		fmt.Printf("removing %s...\n", c.name)
		c.rmf(verbose)
	}
}

// returns true if container belongs to this service
func (s *service) matchContainer(container string) bool {
	return s.containerRe.MatchString(container)
}

func (s *service) String() string {
	return s.namespace + "/" + s.name
}

func (s *service) isLinked(s2 *service) bool {
	for _, l := range s.linkedServices {
		if l.service == s2 {
			return true
		}
	}
	return false
}

// sort services by linked dependencies
//type ByServiceDependency []*service

//func (s ByServiceDependency) Len() int {
//	return len(s)
//}

//func (s ByServiceDependency) Less(i, j int) bool {
//	return !s[i].isLinked(s[j])
//}

//func (s ByServiceDependency) Swap(i, j int) {
//	s[i], s[j] = s[j], s[i]
//}
