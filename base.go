package main

import (
	"log"
	"time"
	"github.com/codegangsta/cli"
	"github.com/mqu/openldap"
)

type Result struct {
	wid int
	count int
	success int
	startTime time.Time
	endTime time.Time
	elapsedTime float64
}

type Job interface {
	Init(int, *cli.Context) bool
	Finish()
	Request() bool
	GetVerbose() int
	IncCount()
	GetCount() int
	IncSuccess()
	GetSuccess() int
}

type BaseJob struct {
	ldap *openldap.Ldap
	wid int
	count int
	success int
	verbose int
}

func (job *BaseJob) Request() bool {
	if job.verbose >= 3 {
		log.Printf("[%d]: %d\n", job.wid, job.count)
	}
	time.Sleep(1000 * time.Millisecond)
	return true
}

func (job *BaseJob) GetVerbose() int {
	return job.verbose
}

func (job *BaseJob) IncCount() {
	job.count++
}

func (job *BaseJob) GetCount() int {
	return job.count
}

func (job *BaseJob) IncSuccess() {
	job.success++
}

func (job *BaseJob) GetSuccess() int {
	return job.success
}

func (job *BaseJob) Init(wid int, c *cli.Context) bool {
	job.wid = wid
	job.verbose = c.Int("verbose")
	url := c.Args().First()
	if job.verbose >= 2 {
		log.Printf("worker[%d]: initialize %s\n", job.wid, url)
	}
	var err error
	job.ldap, err = openldap.Initialize(url)
	if err != nil {
		log.Fatal("initialize err: ", err)
		return false
	}
	job.ldap.SetOption(openldap.LDAP_OPT_PROTOCOL_VERSION, openldap.LDAP_VERSION3)
	err = job.ldap.Bind(c.String("D"), c.String("w"))
	if err != nil {
		log.Fatal("bind err: ", err)
		return false
	}
	return true
}

func (job *BaseJob) Finish() {
	if job.verbose >= 2 {
		log.Printf("worker[%d]: finalize\n", job.wid)
	}
	job.ldap.Close()
}
