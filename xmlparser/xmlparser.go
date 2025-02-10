package xmlparser

import (
	"encoding/xml"
	"fmt"
	"strings"
	"text/tabwriter"
)

type Protocol struct {
	XMLName     xml.Name    `xml:"protocol"`
	Name        string      `xml:"name,attr"`
	Copyright   string      `xml:"copyright"`
	Description Description `xml:"description"`
	Interfaces  []struct {
		XMLName     xml.Name    `xml:"interface"`
		Name        string      `xml:"name,attr"`
		Version     string      `xml:"version,attr"`
		Description Description `xml:"description"`
		Requests    []struct {
			XMLName     xml.Name    `xml:"request"`
			Name        string      `xml:"name,attr"`
			Type        string      `xml:"type,attr"`
			Since       string      `xml:"since,attr"`
			Description Description `xml:"description"`
			Arguments   []Argument  `xml:"arg"`
		} `xml:"request"`
		Events []struct {
			XMLName         xml.Name    `xml:"event"`
			Name            string      `xml:"name,attr"`
			Type            string      `xml:"type,attr"`
			Since           string      `xml:"since,attr"`
			DeprecatedSince string      `xml:"deprecated-since,attr"`
			Description     Description `xml:"description"`
			Arguments       []Argument  `xml:"arg"`
		} `xml:"event"`
		Enums []struct {
			XMLName     xml.Name    `xml:"enum"`
			Name        string      `xml:"name,attr"`
			Bitfield    string      `xml:"bitfield,attr"`
			Since       string      `xml:"since,attr"`
			Description Description `xml:"description"`

			Entries []Entry `xml:"entry"`
		} `xml:"enum"`
	} `xml:"interface"`
}

type Description struct {
	XMLName xml.Name `xml:"description"`
	Summary string   `xml:"summary,attr"`
	Content string   `xml:",chardata"`
}

type Argument struct {
	XMLName     xml.Name    `xml:"arg"`
	Name        string      `xml:"name,attr"`
	Type        string      `xml:"type,attr"`
	Interface   string      `xml:"interface,attr"`
	Enum        string      `xml:"enum,attr"`
	AllowNull   string      `xml:"allow-null,attr"`
	Summary     string      `xml:"summary,attr"`
	Since       string      `xml:"since,attr"`
	Description Description `xml:"description"`
}

type Entry struct {
	XMLName     xml.Name    `xml:"entry"`
	Name        string      `xml:"name,attr"`
	Value       string      `xml:"value,attr"`
	Summary     string      `xml:"summary,attr"`
	Since       string      `xml:"since,attr"`
	Description Description `xml:"description"`
}

func (a Argument) render() string {
	var argSb strings.Builder

	if a.AllowNull == "true" {
		argSb.WriteString(fmt.Sprintf("%s: ?%s", a.Name, a.Type))
	} else {
		argSb.WriteString(fmt.Sprintf("%s: %s", a.Name, a.Type))
	}

	if a.Interface != "" {
		argSb.WriteString(fmt.Sprintf("<%s>", a.Interface))
	} else if a.Enum != "" {
		argSb.WriteString(fmt.Sprintf("<%s>", a.Enum))
	}

	return argSb.String()
}

func renderArgumentSignature(sb *strings.Builder, args []Argument) {
	arglen := len(args)
	if arglen > 1 {
		sb.WriteByte('(')
		sb.WriteString(args[0].render())

		for i := 1; i < arglen; i++ {
			sb.WriteString(fmt.Sprintf(", %s", args[i].render()))
		}

		sb.WriteString(")\n")
	} else if arglen == 1 {
		sb.WriteString(fmt.Sprintf("(%s)\n", args[0].render()))
	}
}

func renderArgumentList(sb *strings.Builder, args []Argument) {
	sb.WriteByte('\n')
	w := tabwriter.NewWriter(sb, 0, 0, 6, ' ', 0)

	w.Write([]byte("[name]\t[type]\t"))

	var hasSummaryDecl bool = false
	var hasDescSummaryDecl bool = false
	var hasDescContentDecl bool = false

	for _, arg := range args {
		if arg.Summary != "" {
			hasSummaryDecl = true
		}
		if arg.Description.Summary != "" {
			hasDescSummaryDecl = true
		}
		if arg.Description.Content != "" {
			hasDescContentDecl = true
		}
	}

	if hasSummaryDecl {
		w.Write([]byte("[summary]\t"))
	}

	if hasDescSummaryDecl {
		w.Write([]byte("[description summary]\t"))
	}

	if hasDescContentDecl {
		w.Write([]byte("[description]\t"))
	}

	w.Write([]byte("\n"))

	for _, arg := range args {
		w.Write([]byte(arg.Name))

		if arg.Since != "" {
			fmt.Fprintf(w, " (since version: %s)", arg.Since)
		}

		fmt.Fprintf(w, "\t%s", arg.Type)

		if arg.Interface != "" {
			fmt.Fprintf(w, "<%s>", arg.Interface)
		} else if arg.Enum != "" {
			fmt.Fprintf(w, "<%s>", arg.Enum)
		}

		if arg.AllowNull == "true" {
			fmt.Fprint(w, " (nullable)")
		}

		if arg.Summary != "" {
			fmt.Fprintf(w, "\t'%s'", arg.Summary)
		}

		if arg.Description.Summary != "" {
			fmt.Fprintf(w, "\t'%s'", arg.Description.Summary)
		}

		if arg.Description.Summary != "" {
			fmt.Fprintf(w, "\t'%s'", arg.Description.Content)
		}

		w.Write([]byte("\t\n"))
	}

	w.Flush()
	sb.WriteByte('\n')
}

func renderEntryList(sb *strings.Builder, entries []Entry) {
	sb.WriteByte('\n')
	w := tabwriter.NewWriter(sb, 0, 0, 6, ' ', 0)

	w.Write([]byte("[name]\t[value]\t"))

	var hasSummaryDecl bool = false
	var hasDescSummaryDecl bool = false
	var hasDescContentDecl bool = false

	for _, entry := range entries {
		if entry.Summary != "" {
			hasSummaryDecl = true
		}
		if entry.Description.Summary != "" {
			hasDescSummaryDecl = true
		}
		if entry.Description.Content != "" {
			hasDescContentDecl = true
		}
	}

	if hasSummaryDecl {
		w.Write([]byte("[summary]\t"))
	}

	if hasDescSummaryDecl {
		w.Write([]byte("[description summary]\t"))
	}

	if hasDescContentDecl {
		w.Write([]byte("[description]\t"))
	}

	w.Write([]byte("\n"))

	for _, entry := range entries {
		w.Write([]byte(entry.Name))

		if entry.Since != "" {
			fmt.Fprintf(w, " (since version: %s)", entry.Since)
		}

		fmt.Fprintf(w, "\t%s", entry.Value)

		if entry.Summary != "" {
			fmt.Fprintf(w, "\t'%s'", entry.Summary)
		}

		if entry.Description.Summary != "" {
			fmt.Fprintf(w, "\t'%s'", entry.Description.Summary)
		}

		if entry.Description.Summary != "" {
			fmt.Fprintf(w, "\t'%s'", entry.Description.Content)
		}

		w.Write([]byte("\t\n"))
	}

	w.Flush()
	sb.WriteByte('\n')
}

func (d Description) render(sb *strings.Builder) {
	if d.Summary != "" {
		sb.WriteString(fmt.Sprintf("%s\n", d.Summary))
	}

	if d.Content != "" {
		sb.WriteString(fmt.Sprintf("%s\n", d.Content))
	}
}

func (p Protocol) Render() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%s\n\n", p.Name))

	p.Description.render(&sb)

	for _, iface := range p.Interfaces {
		sb.WriteString(fmt.Sprintf("interface: %s version: %s\n", iface.Name, iface.Version))
		iface.Description.render(&sb)

		for _, request := range iface.Requests {
			sb.WriteString(fmt.Sprintf("request: %s.%s", iface.Name, request.Name))

			if request.Type != "" {
				sb.WriteString(fmt.Sprintf(" type: %s", request.Type))
			}

			if request.Since != "" {
				sb.WriteString(fmt.Sprintf(" since: version %s", request.Since))
			}

			if len(request.Arguments) > 0 {
				renderArgumentSignature(&sb, request.Arguments)
				renderArgumentList(&sb, request.Arguments)
			}

			request.Description.render(&sb)
		}

		for _, event := range iface.Events {
			sb.WriteString(fmt.Sprintf("event: %s.%s", iface.Name, event.Name))

			if event.Type != "" {
				sb.WriteString(fmt.Sprintf(" type: %s", event.Type))
			}

			if event.Since != "" {
				sb.WriteString(fmt.Sprintf(" since: version %s", event.Since))
			}

			if event.DeprecatedSince != "" {
				sb.WriteString(fmt.Sprintf(" deprecated-since: version %s", event.DeprecatedSince))
			}

			if len(event.Arguments) > 0 {
				renderArgumentSignature(&sb, event.Arguments)
				renderArgumentList(&sb, event.Arguments)
			}

			event.Description.render(&sb)
		}

		for _, enum := range iface.Enums {
			sb.WriteString(fmt.Sprintf("enum: %s.%s", iface.Name, enum.Name))

			if enum.Bitfield == "true" {
				sb.WriteString(" (bitfield)")
			}
			if enum.Since != "" {
				sb.WriteString(fmt.Sprintf(" (since version: %s)", enum.Since))
			}

			sb.WriteByte('\n')

			renderEntryList(&sb, enum.Entries)

			enum.Description.render(&sb)
		}
	}

	sb.WriteString("copyright:\n")
	sb.WriteString(p.Copyright)

	return sb.String()
}

func ParseProtocol(p []byte) Protocol {
	var protocol Protocol
	xml.Unmarshal(p, &protocol)
	return protocol
}
