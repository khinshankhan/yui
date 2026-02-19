package cli

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

type RunFunc func(ctx *Context, args []string) error

type Flag struct {
	Name        string
	Short       string
	Value       string
	Description string
}

type Arg struct {
	Name     string
	Optional bool
	Variadic bool
}

type Section struct {
	Title string
	Lines []string
}

type Example struct {
	Args []string
	Note string
}

type Command struct {
	Name              string
	Aliases           []string
	Description       string
	UsageLines        []string
	Args              []Arg
	SubcommandName    string
	DefaultSubcommand string
	Flags             []Flag
	Sections          []Section
	ExampleSpecs      []Example
	Examples          []string
	Subcommands       []*Command
	Run               RunFunc
}

type Context struct {
	Command *Command
	Path    []string
	Stdout  io.Writer
	Stderr  io.Writer
}

func New(name, description string) *Command {
	return &Command{
		Name:           name,
		Description:    description,
		SubcommandName: "subcommand",
	}
}

func (c *Command) WithAliases(aliases ...string) *Command {
	c.Aliases = append(c.Aliases, aliases...)
	return c
}

func (c *Command) WithUsage(lines ...string) *Command {
	c.UsageLines = append(c.UsageLines, lines...)
	return c
}

func (c *Command) WithArgs(args ...Arg) *Command {
	c.Args = append(c.Args, args...)
	return c
}

func (c *Command) WithSubcommandName(name string) *Command {
	c.SubcommandName = name
	return c
}

func (c *Command) WithSections(sections ...Section) *Command {
	c.Sections = append(c.Sections, sections...)
	return c
}

func (c *Command) WithExamples(examples ...string) *Command {
	c.Examples = append(c.Examples, examples...)
	return c
}

func (c *Command) WithExample(note string, args ...string) *Command {
	c.ExampleSpecs = append(c.ExampleSpecs, Example{
		Args: args,
		Note: note,
	})
	return c
}

func (c *Command) WithRun(run RunFunc) *Command {
	c.Run = run
	return c
}

func (c *Command) WithDefaultSubcommand(name string) *Command {
	c.DefaultSubcommand = name
	return c
}

func (c *Command) Register(subcommands ...*Command) *Command {
	c.Subcommands = append(c.Subcommands, subcommands...)
	return c
}

func (c *Command) RegisterFlags(flags ...Flag) *Command {
	c.Flags = append(c.Flags, flags...)
	return c
}

func RequiredArg(name string) Arg {
	return Arg{Name: name}
}

func OptionalArg(name string) Arg {
	return Arg{Name: name, Optional: true}
}

func VariadicArg(name string) Arg {
	return Arg{Name: name, Variadic: true}
}

func Validate(root *Command) error {
	var errs []string
	validateCommand(root, []string{root.Name}, &errs)
	if len(errs) == 0 {
		return nil
	}
	return errors.New(strings.Join(errs, "\n"))
}

func Execute(root *Command, args []string, stdout, stderr io.Writer) int {
	current := root
	path := []string{root.Name}
	remaining := args

	for len(remaining) > 0 {
		if isHelpArg(remaining[0]) {
			fmt.Fprintln(stdout, current.Help(path))
			return 0
		}

		next := current.findSubcommand(remaining[0])
		if next == nil {
			break
		}

		current = next
		path = append(path, current.Name)
		remaining = remaining[1:]
	}

	if len(remaining) > 0 && len(current.Subcommands) > 0 && current.findSubcommand(remaining[0]) == nil {
		fmt.Fprintf(stderr, "Error: unknown %s: %s\n\n", unknownLabel(path), remaining[0])
		fmt.Fprintln(stderr, current.Help(path))
		return 1
	}

	if len(remaining) == 0 && len(current.Subcommands) > 0 && current.Run == nil && current.DefaultSubcommand != "" {
		next := current.findSubcommand(current.DefaultSubcommand)
		if next == nil {
			fmt.Fprintf(stderr, "Error: unknown default %s: %s\n\n", unknownLabel(path), current.DefaultSubcommand)
			fmt.Fprintln(stderr, current.Help(path))
			return 1
		}
		current = next
		path = append(path, current.Name)
	}

	if current.Run != nil {
		ctx := &Context{
			Command: current,
			Path:    path,
			Stdout:  stdout,
			Stderr:  stderr,
		}
		if err := current.Run(ctx, remaining); err != nil {
			fmt.Fprintf(stderr, "Error: %v\n\n", err)
			fmt.Fprintln(stderr, current.Help(path))
			return 1
		}
		return 0
	}

	if len(current.Subcommands) > 0 {
		fmt.Fprintf(stderr, "Error: %s required\n\n", unknownLabel(path))
		fmt.Fprintln(stderr, current.Help(path))
		return 1
	}

	if len(remaining) > 0 {
		fmt.Fprintf(stderr, "Error: unexpected argument: %s\n\n", remaining[0])
		fmt.Fprintln(stderr, current.Help(path))
		return 1
	}

	return 0
}

func (c *Command) Help(path []string) string {
	cmdPath := strings.Join(path, " ")
	var b strings.Builder

	fmt.Fprintf(&b, "%s - %s\n\n", cmdPath, c.Description)

	b.WriteString("USAGE:\n")
	if len(c.UsageLines) == 0 {
		fmt.Fprintf(&b, "    %s%s\n", cmdPath, c.autoUsageSuffix())
	} else {
		for _, line := range c.UsageLines {
			fmt.Fprintf(&b, "    %s\n", strings.ReplaceAll(line, "%cmd%", cmdPath))
		}
	}

	if len(c.Subcommands) > 0 {
		b.WriteString("\nCOMMANDS:\n")
		width := maxCommandWidth(c.Subcommands)
		for _, sub := range c.Subcommands {
			name := sub.Name
			if len(sub.Aliases) > 0 {
				name = name + ", " + strings.Join(sub.Aliases, ", ")
			}
			fmt.Fprintf(&b, "    %-*s  %s\n", width, name, sub.Description)
		}
		if c.DefaultSubcommand != "" {
			fmt.Fprintf(&b, "\nDEFAULT:\n    %s\n", c.DefaultSubcommand)
		}
	}

	if len(c.Flags) > 0 {
		b.WriteString("\nFLAGS:\n")
		width := maxFlagWidth(c.Flags)
		for _, flag := range c.Flags {
			fmt.Fprintf(&b, "    %-*s  %s\n", width, formatFlag(flag), flag.Description)
		}
	}

	for _, section := range c.Sections {
		if len(section.Lines) == 0 {
			continue
		}
		fmt.Fprintf(&b, "\n%s:\n", section.Title)
		for _, line := range section.Lines {
			fmt.Fprintf(&b, "    %s\n", strings.ReplaceAll(line, "%cmd%", cmdPath))
		}
	}

	if len(c.ExampleSpecs) > 0 || len(c.Examples) > 0 {
		b.WriteString("\nEXAMPLES:\n")
		for _, example := range c.ExampleSpecs {
			fmt.Fprintf(&b, "    %s\n", formatExampleLine(cmdPath, example))
		}
		for _, example := range c.Examples {
			fmt.Fprintf(&b, "    %s\n", strings.ReplaceAll(example, "%cmd%", cmdPath))
		}
	}

	if len(c.Subcommands) > 0 {
		fmt.Fprintf(&b, "\nUse \"%s help\" for more information.", cmdPath)
	}

	return b.String()
}

func (c *Command) findSubcommand(token string) *Command {
	token = strings.ToLower(token)
	for _, sub := range c.Subcommands {
		if token == strings.ToLower(sub.Name) {
			return sub
		}
		for _, alias := range sub.Aliases {
			if token == strings.ToLower(alias) {
				return sub
			}
		}
	}
	return nil
}

func (c *Command) autoUsageSuffix() string {
	var parts []string

	if len(c.Subcommands) > 0 {
		label := c.SubcommandName
		if label == "" {
			label = "subcommand"
		}
		if c.DefaultSubcommand != "" {
			parts = append(parts, "["+label+"]", "[arguments]")
		} else {
			parts = append(parts, "<"+label+">", "[arguments]")
		}
	}

	for _, arg := range c.Args {
		parts = append(parts, formatArg(arg))
	}

	if len(parts) == 0 {
		return ""
	}
	return " " + strings.Join(parts, " ")
}

func formatArg(arg Arg) string {
	name := arg.Name
	if arg.Variadic {
		name = name + "..."
	}
	if arg.Optional {
		return "[" + name + "]"
	}
	return "<" + name + ">"
}

func isHelpArg(arg string) bool {
	switch arg {
	case "help", "-h", "--help":
		return true
	default:
		return false
	}
}

func unknownLabel(path []string) string {
	if len(path) <= 1 {
		return "command"
	}
	return path[len(path)-1] + " subcommand"
}

func maxCommandWidth(commands []*Command) int {
	max := 0
	for _, command := range commands {
		name := command.Name
		if len(command.Aliases) > 0 {
			name = name + ", " + strings.Join(command.Aliases, ", ")
		}
		if len(name) > max {
			max = len(name)
		}
	}
	return max
}

func formatFlag(flag Flag) string {
	long := "--" + flag.Name
	if flag.Value != "" {
		long = long + " <" + flag.Value + ">"
	}

	if flag.Short == "" {
		return long
	}

	short := "-" + flag.Short
	if flag.Value != "" {
		short = short + " <" + flag.Value + ">"
	}
	return long + ", " + short
}

func maxFlagWidth(flags []Flag) int {
	max := 0
	for _, flag := range flags {
		w := len(formatFlag(flag))
		if w > max {
			max = w
		}
	}
	return max
}

func validateCommand(c *Command, path []string, errs *[]string) {
	cmdPath := strings.Join(path, " ")
	if c.Name == "" {
		*errs = append(*errs, fmt.Sprintf("command at %q has empty name", cmdPath))
	}

	if c.DefaultSubcommand != "" && c.findSubcommand(c.DefaultSubcommand) == nil {
		*errs = append(*errs, fmt.Sprintf("command %q has unknown default subcommand %q", cmdPath, c.DefaultSubcommand))
	}

	seen := make(map[string]string)
	addToken := func(token, kind string) {
		if token == "" {
			*errs = append(*errs, fmt.Sprintf("command %q has empty %s", cmdPath, kind))
			return
		}
		key := strings.ToLower(token)
		if prev, ok := seen[key]; ok {
			*errs = append(*errs, fmt.Sprintf("command %q has duplicate token %q (%s conflicts with %s)", cmdPath, token, kind, prev))
			return
		}
		seen[key] = kind
	}

	for _, sub := range c.Subcommands {
		addToken(sub.Name, "subcommand")
		for _, alias := range sub.Aliases {
			addToken(alias, "alias")
		}
	}

	for _, ex := range c.Examples {
		validateExample(c, path, ex, errs)
	}
	for _, ex := range c.ExampleSpecs {
		validateExampleSpec(c, path, ex, errs)
	}

	for _, sub := range c.Subcommands {
		validateCommand(sub, append(path, sub.Name), errs)
	}
}

func validateExample(c *Command, path []string, example string, errs *[]string) {
	cmdPath := strings.Join(path, " ")
	line := strings.TrimSpace(strings.ReplaceAll(example, "%cmd%", cmdPath))
	if line == "" {
		return
	}
	if i := strings.Index(line, " #"); i >= 0 {
		line = strings.TrimSpace(line[:i])
	}
	if line == "" {
		return
	}

	tokens := strings.Fields(line)
	start := findPathStart(tokens, path)
	if start < 0 {
		*errs = append(*errs, fmt.Sprintf("example %q for %q does not include full command path", example, cmdPath))
		return
	}

	current := c
	currentPath := append([]string{}, path...)
	remaining := tokens[start+len(path):]
	for len(remaining) > 0 {
		if isHelpArg(remaining[0]) {
			return
		}

		next := current.findSubcommand(remaining[0])
		if next == nil {
			break
		}
		current = next
		currentPath = append(currentPath, next.Name)
		remaining = remaining[1:]
	}

	if len(remaining) > 0 && len(current.Subcommands) > 0 && current.findSubcommand(remaining[0]) == nil {
		*errs = append(*errs, fmt.Sprintf("example %q for %q has unknown %s %q", example, cmdPath, unknownLabel(currentPath), remaining[0]))
	}
}

func findPathStart(tokens []string, path []string) int {
	if len(tokens) < len(path) {
		return -1
	}
	for i := 0; i <= len(tokens)-len(path); i++ {
		ok := true
		for j := 0; j < len(path); j++ {
			if tokens[i+j] != path[j] {
				ok = false
				break
			}
		}
		if ok {
			return i
		}
	}
	return -1
}

func formatExampleLine(cmdPath string, example Example) string {
	line := cmdPath
	if len(example.Args) > 0 {
		line = line + " " + strings.Join(example.Args, " ")
	}
	if example.Note != "" {
		line = line + "  # " + example.Note
	}
	return line
}

func validateExampleSpec(c *Command, path []string, example Example, errs *[]string) {
	cmdPath := strings.Join(path, " ")
	current := c
	currentPath := append([]string{}, path...)
	remaining := append([]string{}, example.Args...)

	for len(remaining) > 0 {
		if isHelpArg(remaining[0]) {
			return
		}
		next := current.findSubcommand(remaining[0])
		if next == nil {
			break
		}
		current = next
		currentPath = append(currentPath, next.Name)
		remaining = remaining[1:]
	}

	if len(remaining) > 0 && len(current.Subcommands) > 0 && current.findSubcommand(remaining[0]) == nil {
		*errs = append(*errs, fmt.Sprintf("structured example for %q has unknown %s %q", cmdPath, unknownLabel(currentPath), remaining[0]))
	}
}
