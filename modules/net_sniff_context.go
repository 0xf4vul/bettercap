package modules

import (
	"os"
	"regexp"

	"github.com/evilsocket/bettercap-ng/core"
	"github.com/evilsocket/bettercap-ng/log"
	"github.com/evilsocket/bettercap-ng/session"

	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
)

type SnifferContext struct {
	Handle       *pcap.Handle
	DumpLocal    bool
	Verbose      bool
	Filter       string
	Expression   string
	Compiled     *regexp.Regexp
	Output       string
	OutputFile   *os.File
	OutputWriter *pcapgo.Writer
}

func (s *Sniffer) GetContext() (error, *SnifferContext) {
	var err error

	ctx := NewSnifferContext()

	if ctx.Handle, err = pcap.OpenLive(s.Session.Interface.Name(), 65536, true, pcap.BlockForever); err != nil {
		return err, ctx
	}

	if err, ctx.Verbose = s.BoolParam("net.sniff.verbose"); err != nil {
		return err, ctx
	}

	if err, ctx.DumpLocal = s.BoolParam("net.sniff.local"); err != nil {
		return err, ctx
	}

	if err, ctx.Filter = s.StringParam("net.sniff.filter"); err != nil {
		return err, ctx
	} else if ctx.Filter != "" {
		err = ctx.Handle.SetBPFFilter(ctx.Filter)
		if err != nil {
			return err, ctx
		}
	}

	if err, ctx.Expression = s.StringParam("net.sniff.regexp"); err != nil {
		return err, ctx
	} else if ctx.Expression != "" {
		if ctx.Compiled, err = regexp.Compile(ctx.Expression); err != nil {
			return err, ctx
		}
	}

	if err, ctx.Output = s.StringParam("net.sniff.output"); err != nil {
		return err, ctx
	} else if ctx.Output != "" {
		if ctx.OutputFile, err = os.Create(ctx.Output); err != nil {
			return err, ctx
		}

		ctx.OutputWriter = pcapgo.NewWriter(ctx.OutputFile)
		ctx.OutputWriter.WriteFileHeader(65536, layers.LinkTypeEthernet)
	}

	return nil, ctx
}

func NewSnifferContext() *SnifferContext {
	return &SnifferContext{
		Handle:       nil,
		DumpLocal:    false,
		Verbose:      true,
		Filter:       "",
		Expression:   "",
		Compiled:     nil,
		Output:       "",
		OutputFile:   nil,
		OutputWriter: nil,
	}
}

var (
	no  = core.Red("no")
	yes = core.Green("yes")
)

func (c *SnifferContext) Log(sess *session.Session) {
	if c.DumpLocal {
		log.Info("Skip local packets : %s", no)
	} else {
		log.Info("Skip local packets : %s", yes)
	}

	if c.Verbose {
		log.Info("Verbose            : %s", yes)
	} else {
		log.Info("Verbose            : %s", no)
	}

	if c.Filter != "" {
		log.Info("BPF Filter         : '%s'", core.Yellow(c.Filter))
	}

	if c.Expression != "" {
		log.Info("Regular expression : '%s'", core.Yellow(c.Expression))
	}

	if c.Output != "" {
		log.Info("File output        : '%s'", core.Yellow(c.Output))
	}
}

func (c *SnifferContext) Close() {
	if c.Handle != nil {
		c.Handle.Close()
		c.Handle = nil
	}

	if c.OutputFile != nil {
		c.OutputFile.Close()
		c.OutputFile = nil
	}
}
